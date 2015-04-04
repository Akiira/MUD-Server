// ServerMessage
package main

import (
	"strings"
)

const (
	REDIRECT = 1
	GETFILE  = 2
	SAVEFILE = 3
)

type ServerMessage struct {
	MsgType int
	Value   []FormattedString
}

func newServerMessage(msgs []FormattedString) ServerMessage {
	return ServerMessage{Value: msgs}
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
