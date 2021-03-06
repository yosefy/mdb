package cmd

import (
	"context"
	"database/sql"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/coreos/go-oidc"
	"github.com/nats-io/go-nats-streaming"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/stvp/rollbar"
	"github.com/volatiletech/sqlboiler/boil"
	"gopkg.in/gin-contrib/cors.v1"
	"gopkg.in/gin-gonic/gin.v1"

	"github.com/Bnei-Baruch/mdb/api"
	"github.com/Bnei-Baruch/mdb/events"
	"github.com/Bnei-Baruch/mdb/permissions"
	"github.com/Bnei-Baruch/mdb/utils"
	"github.com/Bnei-Baruch/mdb/version"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "MDB api server",
	Run:   serverFn,
}

func init() {
	RootCmd.AddCommand(serverCmd)
}

func serverFn(cmd *cobra.Command, args []string) {
	log.SetFormatter(&log.TextFormatter{FullTimestamp: true})
	log.Infof("Starting MDB API server version %s", version.Version)

	log.Info("Setting up connection to MDB")
	db, err := sql.Open("postgres", viper.GetString("mdb.url"))
	utils.Must(err)
	defer db.Close()
	//boil.SetDB(db)
	boil.DebugMode = viper.GetString("server.mode") == "debug"

	log.Info("Initializing type registries")
	utils.Must(api.InitTypeRegistries(db))

	// Setup events handlers
	eventHandlers := make([]events.EventHandler, 0)
	hNames := viper.GetStringSlice("events.handlers")
	if len(hNames) > 0 {
		for i := range hNames {
			switch hNames[i] {
			case "logger":
				eventHandlers = append(eventHandlers, new(events.LoggerEventHandler))
			case "nats":
				log.Info("Initializing nats streaming event handler")
				h, err := events.NewNatsStreamingEventHandler(
					viper.GetString("nats.subject"),
					viper.GetString("nats.cluster-id"),
					viper.GetString("nats.client-id"),
					stan.NatsURL(viper.GetString("nats.url")),
					stan.PubAckWait(viper.GetDuration("nats.pub-ack-wait")),
				)
				if err != nil {
					log.Errorf("Error connecting to nats streaming server: %s", err)
				} else {
					eventHandlers = append(eventHandlers, h)
				}
			default:
				log.Warnf("Unknown event handler: %s", hNames[i])
			}
		}
	}

	emitter, err := events.NewBufferedEmitter(viper.GetInt("events.emitter-size"), eventHandlers...)
	utils.Must(err)

	// Setup Rollbar
	rollbar.Token = viper.GetString("server.rollbar-token")
	rollbar.Environment = viper.GetString("server.rollbar-environment")
	rollbar.CodeVersion = version.Version

	// cors
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowMethods = append(corsConfig.AllowMethods, http.MethodDelete)
	corsConfig.AllowHeaders = append(corsConfig.AllowHeaders, "Authorization")
	corsConfig.AllowAllOrigins = true

	// Authentication
	var oidcIDTokenVerifier *oidc.IDTokenVerifier
	if viper.GetBool("authentication.enable") {
		oidcProvider, err := oidc.NewProvider(context.TODO(), viper.GetString("authentication.issuer"))
		utils.Must(err)
		oidcIDTokenVerifier = oidcProvider.Verifier(&oidc.Config{
			SkipClientIDCheck: true,
		})
	}

	// casbin
	enforcer, err := permissions.NewEnforcer()
	utils.Must(err)
	enforcer.EnableEnforce(viper.GetBool("permissions.enable"))
	enforcer.EnableLog(viper.GetBool("permissions.log"))

	// Setup gin
	gin.SetMode(viper.GetString("server.mode"))
	router := gin.New()
	router.Use(
		utils.MdbLoggerMiddleware(),
		utils.EnvMiddleware(db, emitter, enforcer, oidcIDTokenVerifier),
		utils.ErrorHandlingMiddleware(),
		permissions.AuthenticationMiddleware(),
		cors.New(corsConfig),
		utils.RecoveryMiddleware())

	api.SetupRoutes(router)

	srv := &http.Server{
		Addr:    viper.GetString("server.bind-address"),
		Handler: router,
	}

	go func() {
		// service connections
		log.Infoln("Running application")
		if err := srv.ListenAndServe(); err != nil {
			log.Infof("Server listen: %s", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Infof("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Infof("Server exiting")

	log.Infof("Closing event handlers")
	for i := range eventHandlers {
		if h, ok := eventHandlers[i].(io.Closer); ok {
			if err := h.Close(); err != nil {
				log.Fatal("Close event handler:", err)
			}
		}
	}

	if len(rollbar.Token) > 0 {
		log.Infof("Wait for rollbar")
		rollbar.Wait()
	}

}
