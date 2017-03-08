package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/Bnei-Baruch/mdb/migrations"
	"github.com/spf13/viper"
	"github.com/vattle/sqlboiler/boil"

	// List all test dependencies here until this bug is fixed in godeps
	// which should allow us to use `godeps save -t`
	// https://github.com/tools/godep/issues/405
	_ "github.com/adams-sarah/test2doc/test"
	_ "github.com/stretchr/testify"

	// used in models tests generated by sqlboiler
	_ "github.com/kat-co/vala"
	_ "github.com/vattle/sqlboiler/bdb"
	_ "github.com/vattle/sqlboiler/bdb/drivers"
	_ "github.com/vattle/sqlboiler/randomize"
	_ "github.com/vattle/sqlboiler/strmangle"
)

type TestDBManager struct {
	testDB string
}

func (m *TestDBManager) InitTestDB() error {
	m.testDB = fmt.Sprintf("test_%s", strings.ToLower(GenerateName(10)))
	fmt.Println("Initializing test DB: ", m.testDB)

	m.initConfig()

	// Open connection to RDBMS
	db, err := sql.Open("postgres", viper.GetString("mdb.url"))
	if err != nil {
		return err
	}

	// Create a new temporary test database
	if _, err := db.Exec("CREATE DATABASE " + m.testDB); err != nil {
		return err
	}

	// Close first connection and connect to temp database
	db.Close()
	db, err = sql.Open("postgres", fmt.Sprintf(viper.GetString("test.url-template"), m.testDB))
	if err != nil {
		return err
	}

	// Run migrations
	m.runMigrations(db)

	// Setup SQLBoiler
	boil.SetDB(db)
	boil.DebugMode = viper.GetBool("test.debug-sql")

	return nil
}

func (m *TestDBManager) DestroyTestDB() error {
	fmt.Println("Destroying testDB: ", m.testDB)

	// Close temp DB
	err := boil.GetDB().(*sql.DB).Close()
	if err != nil {
		return err
	}

	// Connect to MDB
	db, err := sql.Open("postgres", viper.GetString("mdb.url"))
	if err != nil {
		return err
	}

	// Drop test DB
	_, err = db.Exec("DROP DATABASE " + m.testDB)
	if err != nil {
		return err
	}

	return nil
}

func (m *TestDBManager) initConfig() {
	viper.SetDefault("test", map[string]interface{}{
		"url-template": "postgres://localhost/%s?sslmode=disable&?user=postgres",
		"debug-sql":    true,
	})

	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Could not read config, using: ", viper.ConfigFileUsed(), err.Error())
	}
}

func (m *TestDBManager) runMigrations(db *sql.DB) error {
	var visit = func(path string, f os.FileInfo, err error) error {
		match, _ := regexp.MatchString(".*\\.sql$", path)
		if !match {
			return nil
		}

		fmt.Printf("Applying migration %s\n", path)
		m, err := migrations.NewMigration(path)
		if err != nil {
			fmt.Printf("Error migrating %s, %s", path, err.Error())
			return err
		}

		for _, statement := range m.Up() {
			if _, err := db.Exec(statement); err != nil {
				return fmt.Errorf("Unable to apply migration %s: %s\nStatement: %s\n", m.Name, err, statement)
			}
		}

		return nil
	}

	return filepath.Walk("../migrations", visit)
}
