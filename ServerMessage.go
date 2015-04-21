// ServerMessage
package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"strings"
)

const (
	REDIRECT  = 1
	GETFILE   = 2
	SAVEFILE  = 3
	GAMEPLAY  = 4
	PING      = 5
	REPLYPING = 6
)

type ServerMessage struct {
	Value    []FormattedString
	MsgType  int
	CharInfo CharacterInfo
}

type CharacterInfo struct {
	Name      string
	CurrentHP int
	MaxHP     int
	Exp       int
}

func newServerMessageFS(msgs []FormattedString) ServerMessage {
	return ServerMessage{MsgType: GAMEPLAY, Value: msgs}
}

func newServerMessageS(msg string) ServerMessage {
	return ServerMessage{MsgType: GAMEPLAY, Value: newFormattedStringSplice(msg)}
}

func newServerMessageTypeFS(typeOfMsg int, msgs []FormattedString) ServerMessage {
	return ServerMessage{MsgType: typeOfMsg, Value: msgs}
}

func newServerMessageTypeS(typeOfMsg int, msg string) ServerMessage {
	return ServerMessage{MsgType: typeOfMsg, Value: newFormattedStringSplice(msg)}
}

func (msg *ServerMessage) addCharInfo(char *Character) {
	msg.CharInfo = CharacterInfo{CurrentHP: char.currentHP, MaxHP: char.MaxHitPoints, Name: char.Name, Exp: char.experience}
}

func (msg *ServerMessage) getFormattedCharInfo() []FormattedString {
	fs := newFormattedStringCollection()
	fs.addMessage(ct.Red, fmt.Sprintf("\n%d/%d", msg.getCurrentHP(), msg.getMaxHP()))
	fs.addMessage2("|")
	fs.addMessage(ct.Green, fmt.Sprintf("%d", msg.CharInfo.Exp))
	fs.addMessage2("|")
	fs.addMessage(ct.Blue, msg.CharInfo.Name)
	return fs.fmtedStrings
}

func (msg *ServerMessage) getCurrentHP() int {
	return msg.CharInfo.CurrentHP
}

func (msg *ServerMessage) getMaxHP() int {
	return msg.CharInfo.MaxHP
}

func (msg *ServerMessage) getMessage() string {
	if len(msg.Value) <= 0 {
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
