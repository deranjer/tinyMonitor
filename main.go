package main

import (
	"os"
	"time"

	"github.com/deranjer/tinyMonitor/config"
	"github.com/deranjer/tinyMonitor/messaging"

	"github.com/asdine/storm"
	"github.com/rs/zerolog"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pair"

	// register transports
	_ "nanomsg.org/go/mangos/v2/transport/tcp" //TODO change to /all to register all transports
)

var (
	err  error
	msg  []byte
	sock mangos.Socket
	//Logger is the global Logger variable
	Logger zerolog.Logger
	//StormDB is the global bolt database variable
	StormDB        *storm.DB
	serverSettings config.ServerConfig
)

func initializeDatabase() {

}

func pairSocket(serverSettings config.ServerConfig, channel chan string) {
	sock, err = pair.NewSocket()
	if err != nil {
		Logger.Fatal().Err(err).Msg("Can't get new pub socket")
	}
	err = sock.Listen(serverSettings.ListenAddr)
	if err != nil {
		Logger.Fatal().Err(err).Str("Connection String", serverSettings.ListenAddr).Msg("Failed to create listen socket")
	}
	Logger.Info().Str("Address", serverSettings.ListenAddr).Msg("Listen socket created with no errors")
	for {
		msg, err = sock.Recv()
		if err != nil {
			Logger.Error().Err(err).Msg("Server failed to receive pair message from agent")
		} else {
			messaging.MessageDecode(msg) //Send the message over to the messaging package for decoding and handling
			sock.Send([]byte("Server sending response"))
		}

	}
}

func main() {
	serverSettings, Logger, StormDB = config.SetupServer() //receiving the server settings, logger and bbolt database from the config package
	defer StormDB.Close()                                  //Will close the db if the main function exits
	initializeDatabase()
	messaging.Logger = Logger   //injecting the Logger into the messaging package
	messaging.StormDB = StormDB //injecting the Database into the messaging package
	Logger.Info().Msg("Server and Logger configuration complete")
	agentListenCh := make(chan string)
	go pairSocket(serverSettings, agentListenCh) //setting up a permanent pair socket to listen for messages from agents
	for {
		select {
		case <-agentListenCh:
			{
				os.Exit(0)
			}
		default:
			{
				time.Sleep(5 * time.Second)
				Logger.Info().Msg("Just sitting here waiting in main routine... go coroutine running in background")
			}
		}
	}
	//
}
