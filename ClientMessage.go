// ClientMessage
package main

import (
	"strings"
)

//command for system
const CommandLogin = 101
const CommandLogout = 102

//command for create user
const CommandRegister = 111

//command in a room
const CommandAttack = 11
const CommandItem = 12
const CommandLeave = 13 // leave occur the same time with enter the room??

//command between room?
const CommandJoinWorld = 21 // will change the room occur the same time with leave?
// probably use after authenticate with login server and move to the first world as well

//this is suppose to be an event
type ClientMessage struct {
	Command int
	Value   string
}

func ClientMessageConstructor(cmd int, val string) ClientMessage {
	return ClientMessage{Command: cmd, Value: val}
}

func (message *ClientMessage) getPassword() string {
	split := strings.Split(message.Value, " ")
	return split[1]
}

func (message *ClientMessage) getUsername() string {
	split := strings.Split(message.Value, " ")
	return split[0]
}
