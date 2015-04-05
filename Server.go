package main

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"fmt"
	_ "github.com/daviddengcn/go-colortext"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var servers map[string]string
var eventManager *EventManager

func main() {

	readServerList()
	runServer()

	//getCharactersFile("Ragnar")
	//sendCharactersFile("Tiefling")
}

func runServer() {
	loadMonsterData()
	readServerList()
	eventManager = newEventManager()

	listener := setUpServerWithAddress(os.Args[1])

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

func getCharactersFile(name string) {
	conn, err := net.Dial("tcp", servers["characterStorage"])
	checkError(err, true)
	defer conn.Close()

	err = gob.NewEncoder(conn).Encode(&ServerMessage{MsgType: GETFILE, Value: newFormattedStringSplice(name)})
	checkError(err, true)

	file, err := os.Create("Characters/" + name + ".xml")
	checkError(err, true)
	defer file.Close()

	sent, err := io.Copy(file, conn)
	checkError(err, true)
	fmt.Println("Amount Written: ", sent)
}

func sendCharactersFile(name string) {
	conn, err := net.Dial("tcp", servers["characterStorage"])
	checkError(err, true)
	defer conn.Close()

	err = gob.NewEncoder(conn).Encode(&ServerMessage{MsgType: SAVEFILE, Value: newFormattedStringSplice(name)})
	checkError(err, true)

	file, err := os.Open("Characters/" + name + ".xml")
	checkError(err, true)
	defer file.Close()

	buf := new(bytes.Buffer)
	io.Copy(buf, file)
	written, err := conn.Write(buf.Bytes())

	//written, err := io.Copy(conn, buf)
	checkError(err, true)
	fmt.Println("Amount Written: ", written)
}
