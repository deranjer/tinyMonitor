package messaging

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/asdine/storm"
	"github.com/rs/zerolog"
)

var (
	//Logger is the global zap logger
	Logger zerolog.Logger
	//StormDB is the global bolt database variable
	StormDB *storm.DB
)

//BaseMessage is the base message for all messages sent/received
type BaseMessage struct {
	MessageType string
	MessageBody interface{}
}

//RegisterAgent is for agents to register with the server
type RegisterAgent struct {
	AgentID       int `storm:"id,increment"`
	AgentHostName string
	AgentIPAddr   string
	AgentJoinDate time.Time
}

//MessageDecode handles all incoming messages
func MessageDecode(msg []byte) {
	gob.Register(RegisterAgent{})
	decodedMsg := BaseMessage{}
	b := bytes.NewBuffer(msg)
	err := gob.NewDecoder(b).Decode(&decodedMsg)
	if err != nil {
		Logger.Error().Err(err).Msg("Unable to decode message!")
	}
	Logger.Debug().Msg("Message Received from Agent")
	fmt.Printf("%+v\n", decodedMsg)
	if decodedMsg.MessageBody == nil || decodedMsg.MessageBody == "" { //If the message body is missing or blank, return
		Logger.Error().Msg("Unable to decode message from agent")
		return
	}
	switch decodedMsg.MessageType {
	case "RegisterAgent":
		messageBody := decodedMsg.MessageBody.(RegisterAgent)
		registerAgentHandler(messageBody)
	}
}

//MessageEncode will encode and test the message (and maybe send it?)
func MessageEncode(msg BaseMessage) []byte {
	gob.Register(RegisterAgent{})
	fmt.Printf("%+v\n", msg)
	b := new(bytes.Buffer)
	err := gob.NewEncoder(b).Encode(msg)
	if err != nil {
		Logger.Error().Err(err).Msg("Unable to encode message!")
	}
	Logger.Info().Msg("Message encoded")
	return b.Bytes()
}

func registerAgentHandler(newAgent RegisterAgent) {
	fmt.Println("AgentHostName", newAgent.AgentHostName, "AgentIPAddr", newAgent.AgentIPAddr, "AgentJoinDate", newAgent.AgentJoinDate)
	err := StormDB.Save(&newAgent)
	if err != nil {
		Logger.Error().Err(err).Msg("Unable to write to database!")
	}
}
