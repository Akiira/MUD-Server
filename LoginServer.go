package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

var servers map[string]string

//func main() {
func notMain() {
	servers = make(map[string]string)
	//files, _ := ioutil.ReadDir("./")

	/*
		for _, f := range files {
			fmt.Println(f.Name())
		}*/
	readServerList()

	go runCharacterServer()
	runClientServer()
}

func runCharacterServer() {
	listener := setUpServerWithPort(1301)
	for {
		fmt.Println("Character Server: i'm waiting")
		conn, err := listener.Accept()
		//checkError(err)
		if err == nil {
			fmt.Println("Character Server:Connection established")

		}
	}
}

func runClientServer() {
	listener := setUpServer()

	for {
		fmt.Println("Client Server: i'm waiting")
		conn, err := listener.Accept()
		//checkError(err)
		if err == nil {
			fmt.Println("Client Server:Connection established")
			go HandleLoginClient(conn)
		}
	}
}

func HandleLoginClient(myConn net.Conn) {

	//waiting for login msg
	//then validate and send connection for world server back
	//or return error validation fail

	var clientResponse ClientMessage
	myDecoder := gob.NewDecoder(myConn)

	err := myDecoder.Decode(&clientResponse)

	c := getCharacterFromFile(clientResponse.getUsername(), clientResponse.getPassword())

	//TODO
	if c == nil {
		//incorrect password or character name
	} else {
		// Correct!
	}
}
func replyWorldIsNotFound(myConn net.Conn) {
	var svMsg ServerMessage
	//svMsg.MsgType = ErrorWorldIsNotFound
	//svMsg.MsgDetail = "error world is not found"
	gob.NewEncoder(myConn).Encode(svMsg)
}

func replyFailAuthorizationCommand(myConn net.Conn, msgDetail string) {
	var svMsg ServerMessage
	//svMsg.MsgType = ErrorAuthorizationFail
	//svMsg.MsgDetail = msgDetail
	gob.NewEncoder(myConn).Encode(svMsg)
}

func replyUnexpectedCommand(myConn net.Conn) {
	var svMsg ServerMessage
	//svMsg.MsgType = ErrorUnexpectedCommand
	//svMsg.MsgDetail = "error unexpected command"
	gob.NewEncoder(myConn).Encode(svMsg)
}

func readServerList() {
	//this should be the one that read list of servers, including central server
	serverNum = 0
	file, err := os.Open("serverConfig/serverList.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		readData := strings.Fields(scanner.Text())
		fmt.Println(readData)
		servers[data[0]] = data[1]
	}

	for i := 0; i < serverNum; i++ {
		fmt.Println(serverNames[i], " ", serverAddrs[i])
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	//Pattanapoom Hand
	//start model
}
