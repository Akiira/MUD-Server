package main

import (
	"fmt"
	"net"
	"os"
)

//command for error
const ErrorUnexpectedCommand = 201
const ErrorWorldIsNotFound = 202
const ErrorAuthorizationFail = 203

//command for system
const CommandLogin = 101
const CommandLogout = 102
const CommandRedirectServer = 103
const CommandEnterWorld = 104
const CommandQueryCharacter = 105
const CommandCharacterDetail = 106

//command for create user
const CommandRegister = 111

//command in a room
const CommandAttack = 11
const CommandItem = 12
const CommandLeave = 13 // leave occur the same time with enter the room??

//command between room?
const CommandJoinWorld = 21 // will change the room occur the same time with leave?
// probably use after authenticate with login server and move to the first world as well

func setUpServerWithAddress(addr string) *net.TCPListener {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err, true)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, true)
	return listener
}

func checkError(err error, exitIfError bool) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		if exitIfError {
			os.Exit(1)
		}
	}
}

func checkErrorWithMessage(err error, exitIfError bool, messageIfError string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		fmt.Fprintf(os.Stderr, "Additional Message: %s", messageIfError)
		if exitIfError {
			os.Exit(1)
		}
	}
}
