package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var servers map[string]string
var eventManager *EventManager

func main() {

	if len(os.Args) < 2 {
		fmt.Println(os.Args[0] + " requires 1 arguments, worldname")
		os.Exit(1)
	}

	readServerList()
	runServer()

	//getCharactersFile("Ragnar")
	//sendCharactersFile("Tiefling")
}

func runServer() {
	loadMonsterData()
	readServerList()
	eventManager = newEventManager(os.Args[1])

	listener := setUpServerWithAddress(servers[os.Args[1]])

	for {
		conn, err := listener.Accept()
		checkError(err, false)
		if err == nil {
			fmt.Println("Connection established")

			go HandleClientConnection(conn)
		}
	}
}

func readServerList() {
	//this should be the one that read list of servers, including central server
	servers = make(map[string]string)
	file, err := os.Open("serverConfig/serverList.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		readData := strings.Fields(scanner.Text())
		fmt.Println(readData)
		servers[readData[0]] = readData[1]
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func HandleClientConnection(myConn net.Conn) {

	clientConnection := newClientConnection(myConn, eventManager)
	clientConnection.receiveMsgFromClient()
}

func sendCharactersFile(name string) {
	conn, err := net.Dial("tcp", servers["characterStorage"])
	checkError(err, true)
	defer conn.Close()

	err = gob.NewEncoder(conn).Encode(&ServerMessage{MsgType: SAVEFILE, Value: newFormattedStringSplice(name)})
	checkError(err, true)

	file, err := os.Open("Characters/" + name + ".xml")
	checkError(err, true)
	written, err := io.Copy(conn, file)
	checkError(err, true)
	fmt.Println("Amount Send: ", written)
	file.Close()
}
