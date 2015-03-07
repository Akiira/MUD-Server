// ClientMessage
package main

import (
	"strings"
)

type ClientMessage struct {
	Command int
	Value   string
}


func ClientMessageConstructor(cmd int, val string) ClientMessage{
	return ClientMessage{ Command: cmd, Value: val}
}

func (message *ClientMessage) getPassword() string{
	split := strings.Split(message.Value, " ")
	return split[1]
}

func (message *ClientMessage) getUsername() string{
	split := strings.Split(message.Value, " ")
	return split[0]
}