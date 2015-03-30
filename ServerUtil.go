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

//command for create user
const CommandRegister = 111

//command in a room
const CommandAttack = 11
const CommandItem = 12
const CommandLeave = 13 // leave occur the same time with enter the room??

//command between room?
const CommandJoinWorld = 21 // will change the room occur the same time with leave?
// probably use after authenticate with login server and move to the first world as well

func setUpServer() *net.TCPListener {
	service := "127.0.0.1:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	return listener
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
