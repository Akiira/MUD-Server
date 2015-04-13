// ClientMessage
package main

import (
	"strconv"
	"strings"
	"time"
)

//this is suppose to be an event
type ClientMessage struct {
	CombatAction bool
	Command      string
	Value        string
}

func newClientMessage(cmd string, val string) *ClientMessage {
	return &ClientMessage{CombatAction: false, Command: cmd, Value: val}
}
func (msg *ClientMessage) getCommand() string {
	if strings.Contains(msg.Command, ";") {
		return strings.Split(msg.Command, ";")[0]
	} else {
		return msg.Command
	}
}
func (msg *ClientMessage) setCommand(cmd string) {
	msg.CombatAction = false
	msg.Command = cmd
	msg.Value = ""
}

func (msg *ClientMessage) setCommandWithTimestamp(cmd string) {
	msg.CombatAction = false
	msg.Command = cmd + ";" + time.Now().String()
	msg.Value = ""
}

func (msg *ClientMessage) setMsgWithTimestamp(cmd string, value string) {
	msg.CombatAction = false
	msg.Command = cmd + ";" + time.Now().String()
	msg.Value = value
}

func (msg *ClientMessage) getTimeStamp() string {

	peices := strings.Split(msg.Command, ";")
	if len(peices) == 2 {
		return peices[1]
	} else {
		return ""
	}
}

func (msg *ClientMessage) setToMovementMessage(direction string) {
	msg.CombatAction = false
	msg.Command = "move"
	msg.Value = direction
}

func (msg *ClientMessage) setToSayMessage(thingToSay string) {
	msg.CombatAction = false
	msg.Command = "say"
	msg.Value = thingToSay
}

func (msg *ClientMessage) setToGetMessage(item string) {
	msg.CombatAction = false
	msg.Command = "get"
	msg.Value = item
}

func (msg *ClientMessage) setToLookMessage(target string) {
	msg.CombatAction = false
	msg.Command = "look"
	msg.Value = target
}

func (msg *ClientMessage) setToAttackMessage(target string) {
	msg.CombatAction = true
	msg.Command = "attack"
	msg.Value = target
}

func (msg *ClientMessage) setToExitMessage() {
	msg.CombatAction = false
	msg.Command = "exit"
	msg.Value = ""
}

func (msg *ClientMessage) setAll(combatAction bool, cmd string, val string) {
	msg.CombatAction = combatAction
	msg.Command = cmd
	msg.Value = val
}

func (msg *ClientMessage) setAllNonCombat(cmd string, val string) {
	msg.CombatAction = false
	msg.Command = cmd
	msg.Value = val
}

func (message *ClientMessage) getPassword() string {
	split := strings.Split(message.Value, " ")
	return split[1]
}

func (message *ClientMessage) getUsername() string {
	split := strings.Split(message.Value, " ")
	return split[0]
}

func (msg *ClientMessage) getBid() int {
	val, err := strconv.Atoi(msg.Value)
	checkError(err, false)
	return val
}
