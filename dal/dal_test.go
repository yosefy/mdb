package dal

import (
	"github.com/Bnei-Baruch/mdb/migrations"
	"github.com/Bnei-Baruch/mdb/rest"
	"github.com/Bnei-Baruch/mdb/utils"

    "regexp"
    "fmt"
    "math/rand"
    "os"
    "strings"
    "testing"
    "time"
    "path/filepath"

	"github.com/spf13/viper"
    "github.com/jinzhu/gorm"
)

func RunMigrations(tmpDb *gorm.DB) {

    var visit = func(path string, f os.FileInfo, err error) error {
        match, _ := regexp.MatchString(".*\\.sql$", path);
        if !match {
            fmt.Printf("Did not match sql file %s\n", path)
            return nil
        }

        fmt.Printf("Migrating %s\n", path)
        m, err := migrations.NewMigration(path)
        if err != nil {
            fmt.Printf("Error migrating %s, %s", path, err.Error())
            return err
        }

        for _, statement := range m.Up() {
            if _, err := tmpDb.CommonDB().Exec(statement); err != nil {
                return fmt.Errorf("Unable to apply migration %s: %s\nStatement: %s\n", m.Name, err, statement)
            }
        }

        return nil
    }

    err := filepath.Walk("../migrations", visit)
    if err != nil {
        panic(fmt.Sprintf("Could not load and run all migrations. %s", err.Error()))
    }

}

func InitTestConfig() {
	viper.SetConfigName("config")
	viper.AddConfigPath("../")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Could not read config, using: ", viper.ConfigFileUsed(), err)
	}
}

func SwitchToTmpDb() (*gorm.DB, *gorm.DB, string) {
    InitTestConfig()

    baseDb, err := Init()
    if err != nil {
        panic(fmt.Sprintf("Could not connect to database. %s", err))
    }

    name := strings.ToLower(utils.GenerateName(10))
    if err := baseDb.Exec(fmt.Sprintf("CREATE DATABASE %s", name)).Error; err != nil {
        panic(fmt.Sprintf("Could not create tmp database %s due to %s.", name, err))
    }

	url := viper.GetString("test.url-template")
    var tmpDb *gorm.DB
    tmpDb, err = InitByUrl(fmt.Sprintf(url, name))

    RunMigrations(tmpDb)

    return baseDb, tmpDb, name
}

func DropTmpDB(baseDb *gorm.DB, tmpDb *gorm.DB, name string) {
    tmpDb.Close()
    if err := baseDb.Exec(fmt.Sprintf("DROP DATABASE %s", name)).Error; err != nil {
        panic(fmt.Sprintf("Could not drop test database. %s", err))
    }
}

func TestInit(t *testing.T) {
    InitTestConfig()

    if _, err := InitByUrl("bad://database-connection-url"); err == nil {
		t.Error("Expected not nil, got nil.")
	}
	if _, err := Init(); err != nil {
		t.Error("Expected nil, got ", err.Error())
	}
}

func TestCaptureStart(t *testing.T) {
    SwitchToTmpDb()
    // baseDb, tmpDb, name := SwitchToTmpDb()
    // defer DropTmpDB(baseDb, tmpDb, name)

    cs := rest.CaptureStart{
        Type: "type",
        Station: "a station",
        User: "operator@dev.com",
        FileName: "some.file.name",
        CaptureID: "this.is.capture.id",
    }

    if err := CaptureStart(cs); err != nil {
        t.Error("CaptureStart should succeed.", err)
    }
}

func TestMain(m *testing.M) {
    rand.Seed(time.Now().UTC().UnixNano())
    os.Exit(m.Run())
}
