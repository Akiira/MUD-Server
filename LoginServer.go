package main

import (
	"bufio"
	"encoding/gob"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strings"
)

var serverNames [10]string
var serverAddrs [10]string
var serverNum int

func main() {

	//files, _ := ioutil.ReadDir("./")

	/*
		for _, f := range files {
			fmt.Println(f.Name())
		}*/
	readServerList()

	listener := setUpServer()

	for {
		conn, err := listener.Accept()
		checkError(err)
		fmt.Println("Connection established")

		go HandleClient(conn)
	}

}

func HandleClient(myConn net.Conn) {

	//waiting for login msg
	//then validate and send connection for world server back
	//or return error validation fail

	var clientResponse ClientMessage
	myDecoder := gob.NewDecoder(myConn)

	for {

		err := myDecoder.Decode(&clientResponse)
		checkError(err)

		if err == nil {

			if clientResponse.MsgType == CommandLogin {
				value := clientResponse.Value
				data := strings.Fields(value)
				username := data[0]
				password := data[1]

				fielName := username + ".xml"

				if _, err := os.Stat(fielName); err == nil {
					fmt.Printf("file exists; processing...")

					xmlFile, err := os.Open(fielName)
					checkError(err)
					defer xmlFile.Close()

					//should send charData to world server
					//as an notification that a character is going to join

					var charData CharacterXML
					XMLdata, _ := ioutil.ReadAll(xmlFile)
					xml.Unmarshal(XMLdata, &charData)

					//find the world's addr for character to respawn
					lookupWorld := charData.CurrentWorld
					found := false
					var addr string
					for i := 0; i < serverNum; i++ {
						if worldsName[i] == lookupWorld {
							found = true
							addr = worldAddrs[i]
							break
						}
					}

					if found {
						//send world addr back to client
						var svMsg ServerMessage
						svMsg.MsgType = CommandRedirectServer
						svMsg.Value = addr
						gob.NewDecoder(myConn).Encode(svMsg)
					}
					else {
						var svMsg ServerMessage
						svMsg.MsgType = ErrorWorldIsNotFound
						svMsg.Value = "error world is not found"
						gob.NewDecoder(myConn).Encode(svMsg)		
					}
				}

			} else {
				var svMsg ServerMessage
				svMsg.MsgType = ErrorUnexpectedCommand
				svMsg.Value = "error unexpected command"
				gob.NewDecoder(myConn).Encode(svMsg)
			}

		} else {
			break
		}
	}
}

func setUpServer() *net.TCPListener {
	service := "127.0.0.1:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	return listener
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
		serverNames[serverNum] = readData[0]
		serverAddrs[serverNum] = readData[1]
		serverNum++
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
