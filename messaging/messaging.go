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
	AgentID       int `storm:"id,unique"`
	AgentHostName string
	AgentIPAddr   string
	AgentJoinDate time.Time
}

//MessageDecode handles all incoming messages
func MessageDecode(msg []byte) {
	decodedMsg := BaseMessage{}
	b := bytes.NewBuffer(msg)
	gob.NewDecoder(b).Decode(&decodedMsg)
	Logger.Debug().Msg("Message Received from Agent")
	fmt.Printf("%+v\n", decodedMsg)
	var messageBody map[string]interface{}
	if decodedMsg.MessageBody != nil && decodedMsg.MessageBody != "" {
		messageBody = decodedMsg.MessageBody.(map[string]interface{})
	}
	switch decodedMsg.MessageType {
	case "RegisterAgent":
		registerAgentHandler(messageBody)
	}
}

//MessageEncode will encode and test the message (and maybe send it?)
func MessageEncode(msg *BaseMessage) *bytes.Buffer {
	b := new(bytes.Buffer)
	gob.NewEncoder(b).Encode(&msg)
	return b
}

func registerAgentHandler(messageBody map[string]interface{}) {
	newAgent := RegisterAgent{
		messageBody["AgentID"].(int),
		messageBody["AgentHostName"].(string),
		messageBody["AgentIPAddr"].(string),
		messageBody["AgentJoinDate"].(time.Time),
	}
	StormDB.Save(&newAgent)
}
