// ServerMessage
package main

import (
	"github.com/daviddengcn/go-colortext"
	"strings"
)

const (
	REDIRECT = 1
	GETFILE  = 2
	SAVEFILE = 3
	GAMEPLAY = 4
)

type ServerMessage struct {
	Value   []FormattedString
	MsgType int
}

func newServerMessage(msgs []FormattedString) ServerMessage {
	return ServerMessage{Value: msgs}
}

func newServerMessageWithType(typeOfMsg int, msgs []FormattedString) ServerMessage {
	return ServerMessage{MsgType: typeOfMsg, Value: msgs}
}

func newSimpleServerMessage(typeOfMsg int, msg string) ServerMessage {
	output := make([]FormattedString, 1, 1)
	output[0].Color = ct.White
	output[0].Value = msg
	return ServerMessage{MsgType: typeOfMsg, Value: output}
}

func (msg *ServerMessage) getMessage() string {
	if len(msg.Value) == 0 {
		return ""
	}
	return msg.Value[0].Value
}

func (msg *ServerMessage) isError() bool {
	if len(msg.Value) == 0 {
		return false
	}

	return (strings.Split(msg.Value[0].Value, " ")[0] == "error")
}
