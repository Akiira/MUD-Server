package main

import (
	"bufio"
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"io/ioutil"
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
	//runServer()

	getCharacterFromCentral("Ragnar")
	sendCharactersFile("Tiefling")
}

func runServer() {
	loadMonsterData()

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

	encdr := gob.NewEncoder(conn)
	err = encdr.Encode(&ServerMessage{MsgType: SAVEFILE, Value: newFormattedStringSplice(name)})
	checkError(err, true)

	file, err := os.Open("Characters/" + name + ".xml")
	checkError(err, true)
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	checkError(err, false)

	var c CharacterXML
	err = xml.Unmarshal(data, &c)
	checkError(err, false)

	err = encdr.Encode(c)
	checkError(err, false)
}
