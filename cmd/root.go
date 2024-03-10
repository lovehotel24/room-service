package cmd

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
	middleware "github.com/oapi-codegen/echo-middleware"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/lovehotel24/room-service/pkg/configs"
	"github.com/lovehotel24/room-service/pkg/controllers"
	"github.com/lovehotel24/room-service/pkg/routers"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "room-service",
	Short: "A brief description of your application",
	Run:   runCommand,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().String("pg-user", "postgres", "user name for postgres database")
	rootCmd.Flags().String("pg-pass", "postgres", "password for postgres database")
	rootCmd.Flags().String("pg-host", "localhost", "postgres server address")
	rootCmd.Flags().String("pg-port", "5432", "postgres server port")
	rootCmd.Flags().String("pg-db", "postgres", "postgres database name")
	rootCmd.Flags().Bool("pg-ssl", false, "postgres server ssl mode on or not")
	rootCmd.Flags().String("port", "8082", "port for test HTTP server")
	rootCmd.Flags().String("grpc-port", "50051", "booking service grpc port")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("book")
	viper.BindPFlag("pg-user", rootCmd.Flags().Lookup("pg-user"))
	viper.BindPFlag("pg-pass", rootCmd.Flags().Lookup("pg-pass"))
	viper.BindPFlag("pg-host", rootCmd.Flags().Lookup("pg-host"))
	viper.BindPFlag("pg-port", rootCmd.Flags().Lookup("pg-port"))
	viper.BindPFlag("pg-db", rootCmd.Flags().Lookup("pg-db"))
	viper.BindPFlag("pg-ssl", rootCmd.Flags().Lookup("pg-ssl"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("grpc-port", rootCmd.Flags().Lookup("grpc-port"))
	viper.BindEnv("gin_mode", "GIN_MODE")
	viper.AutomaticEnv()
}

func runCommand(cmd *cobra.Command, args []string) {

	dbConf := configs.NewDBConfig().
		WithHost(viper.GetString("pg-host")).
		WithPort(viper.GetString("pg-port")).
		WithUser(viper.GetString("pg-user")).
		WithPass(viper.GetString("pg-pass")).
		WithName(viper.GetString("pg-db")).
		WithSecure(viper.GetBool("pg-ssl"))

	swagger, err := routers.GetSwagger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
		os.Exit(1)
	}

	var log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetLevel(logrus.InfoLevel)

	// Clear out the servers array in the swagger spec, that skips validating
	// that server names match. We don't know how this thing will be run.
	swagger.Servers = nil

	db, err := configs.NewDB(dbConf)
	if err != nil {
		log.WithError(err).Error("failed to connect db")
	}

	err = configs.Migrate(db)
	if err != nil {
		log.WithError(err).Error("failed to migrate db schema")
	}

	// Create an instance of our handler which satisfies the generated interface
	api := controllers.NewApp(db, log)

	// This is how you set up a basic Echo router
	e := echo.New()
	// Log all requests
	e.Use(echomiddleware.Logger())
	// Use our validation middleware to check all requests against the
	// OpenAPI schema.
	e.Use(middleware.OapiRequestValidator(swagger))

	wapper := routers.ServerInterfaceWrapper{Handler: api}
	e.GET("/v1/roomtype", wapper.GetAllRoomType, testMiddleware)

	// And we serve HTTP until the world ends.
	log.Fatalln(e.Start(net.JoinHostPort("0.0.0.0", viper.GetString("port"))))
}

func testMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		logrus.Info("test middleware")
		return next(c)
	}
}
