package main

import (
	"strings"
)

//command for error
const ErrorUnexpectedCommand = 201
const ErrorWorldIsNotFound = 202

//command for system
const CommandLogin = 101
const CommandLogout = 102
const CommandRedirectServer = 103
const CommandEnterWorld = 104
const CommandQueryCharacter = 105

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
	MsgType      int
	CombatAction bool
	Command      string
	Value        string
}

func ClientMessageConstructor(cmd string, val string) ClientMessage {
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
