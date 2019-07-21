package main

import (
	"fmt"
	"os"
	"time"

	"github.com/deranjer/tinyMonitor/config"

	"github.com/rs/zerolog"
	"nanomsg.org/go/mangos/v2"
	"nanomsg.org/go/mangos/v2/protocol/pair"

	// register transports
	_ "nanomsg.org/go/mangos/v2/transport/tcp" //TODO change to /all to register all transports
)

func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}

var (
	err  error
	msg  []byte
	sock mangos.Socket
	//Logger is the global Logger variable
	Logger zerolog.Logger
)

func date() string {
	return time.Now().Format(time.ANSIC)
}

func pairSocket(serverSettings config.ServerConfig, Logger zerolog.Logger, channel chan string) {
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
		//sock.SetOption(mangos.OptionRecvDeadline, 100*time.Millisecond)
		msg, err = sock.Recv()
		if err != nil {
			Logger.Error().Err(err).Msg("Server failed to receive pair message from agent")
		} else {
			Logger.Debug().Str("Message Body", string(msg)).Msg("Message Received from Agent")
			sock.Send([]byte("Server sending response"))
		}

	}
}

func main() {
	serverSettings, Logger := config.SetupServer() //setup logging and all server settings
	Logger.Info().Msg("Server and Logger configuration complete")
	agentListenCh := make(chan string)
	go pairSocket(serverSettings, Logger, agentListenCh) //setting up a permanent pair socket to listen for messages from agents
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
