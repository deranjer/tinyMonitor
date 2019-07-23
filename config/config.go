package config

import (
	"fmt"
	"os"

	"github.com/asdine/storm"
	"github.com/rs/zerolog"
	"github.com/spf13/viper"
)

var (
	//Logger is the global zap logger
	Logger zerolog.Logger
	//StormDB is the global bolt database variable
	StormDB *storm.DB
)

//ServerConfig contains all of the server settings defined in the TOML file
type ServerConfig struct {
	ListenAddr string
}

//SetupServer does the initial configuration
func SetupServer() (ServerConfig, zerolog.Logger, *storm.DB) {
	viper.AddConfigPath("config/")
	viper.AddConfigPath(".")
	viper.SetConfigName("serverConfig")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("fatal error config file: %s \n", err))
	}
	setupLogging()
	StormDB = setupDatabase()
	Logger.Info().Msg("Logger is setup")
	serverSettings := ServerConfig{}
	serverPort := viper.GetString("serverConfig.ServerPort")
	serverAddr := viper.GetString("serverConfig.ServerAddr")
	serverSettings.ListenAddr = "tcp://" + serverAddr + ":" + serverPort //TODO in the future support more than just TCP

	return serverSettings, Logger, StormDB
}

func setupLogging() {
	logLevelString := viper.GetString("logging.Level")
	switch logLevelString { //Options = Debug 0, Info 1, Warn 2, Error 3, Fatal 4, Panic 5
	case "Panic", "panic":
		zerolog.SetGlobalLevel(zerolog.PanicLevel)
	case "Fatal", "fatal":
		zerolog.SetGlobalLevel(zerolog.FatalLevel)
	case "Error", "error":
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
	case "Warn", "warn":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	case "Info", "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "Debug", "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	}
	//zapConfig.Encoding = viper.GetString("logging.Encoding")
	//zapConfig.OutputPaths = viper.GetStringSlice("logging.OutputPaths")
	//zapConfig.ErrorOutputPaths = viper.GetStringSlice("logging.ErrorOutputPaths")
	logger := zerolog.New(os.Stdout).With().Timestamp().Logger() //TODO let user specifiy output options
	Logger = logger
}

func setupDatabase() (db *storm.DB) {
	db, err := storm.Open("tinyMonitorDB.db")
	if err != nil {
		Logger.Fatal().Err(err).Msg("Unable to create/open database!")
	}
	return db
}
