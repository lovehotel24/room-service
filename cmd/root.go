package cmd

import (
	"net"
	"os"
	"strings"

	"github.com/labstack/echo/v4"
	echomiddleware "github.com/labstack/echo/v4/middleware"
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
	rootCmd.Flags().String("grpc-host", ":50051", "grpc server to connect")
	rootCmd.Flags().String("s3-url", "http://localhost:9000", "AWS S3 endpoint url")
	rootCmd.Flags().String("region", "us-east-1", "AWS Region for S3")
	rootCmd.Flags().String("bucket", "room", "AWS S3 bucket name for room photos")
	rootCmd.Flags().String("access-key", "minio", "AWS access key id")
	rootCmd.Flags().String("secret-key", "minioTopSecret", "AWS secret key id")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.SetEnvPrefix("room")
	viper.BindPFlag("pg-user", rootCmd.Flags().Lookup("pg-user"))
	viper.BindPFlag("pg-pass", rootCmd.Flags().Lookup("pg-pass"))
	viper.BindPFlag("pg-host", rootCmd.Flags().Lookup("pg-host"))
	viper.BindPFlag("pg-port", rootCmd.Flags().Lookup("pg-port"))
	viper.BindPFlag("pg-db", rootCmd.Flags().Lookup("pg-db"))
	viper.BindPFlag("pg-ssl", rootCmd.Flags().Lookup("pg-ssl"))
	viper.BindPFlag("port", rootCmd.Flags().Lookup("port"))
	viper.BindPFlag("grpc-port", rootCmd.Flags().Lookup("grpc-port"))
	viper.BindPFlag("s3-url", rootCmd.Flags().Lookup("s3-url"))
	viper.BindPFlag("region", rootCmd.Flags().Lookup("region"))
	viper.BindPFlag("bucket", rootCmd.Flags().Lookup("bucket"))
	viper.BindPFlag("access-key", rootCmd.Flags().Lookup("access-key"))
	viper.BindPFlag("secret-key", rootCmd.Flags().Lookup("secret-key"))
	viper.BindPFlag("grpc-host", rootCmd.Flags().Lookup("grpc-host"))
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

	//swagger, err := routers.GetSwagger()
	//if err != nil {
	//	fmt.Fprintf(os.Stderr, "Error loading swagger spec\n: %s", err)
	//	os.Exit(1)
	//}

	var log = logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)
	log.SetLevel(logrus.InfoLevel)

	//swagger.Servers = nil
	//openapi3filter.RegisterBodyDecoder("multipart/form-data", openapi3filter.FileBodyDecoder)

	db, err := configs.NewDB(dbConf)
	if err != nil {
		log.WithError(err).Error("failed to connect db")
		os.Exit(1)
	}

	err = configs.Migrate(db)
	if err != nil {
		log.WithError(err).Error("failed to migrate db schema")
	}

	s3Client := configs.NewS3Client(viper.GetString("s3-url"),
		viper.GetString("region"),
		viper.GetString("access-key"),
		viper.GetString("secret-key"))

	userClient, err := configs.NewGrpcUserService(viper.GetString("grpc-host"))
	if err != nil {
		log.WithError(err).Error("failed to connect GRPC user service")
		os.Exit(1)
	}

	api := controllers.NewApp(db, log, s3Client, userClient)

	e := echo.New()

	e.Use(echomiddleware.Logger())

	//e.Use(middleware.OapiRequestValidator(swagger))

	routers.RouterV1(e, api)

	log.Fatalln(e.Start(net.JoinHostPort("0.0.0.0", viper.GetString("port"))))
}
