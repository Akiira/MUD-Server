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

		go HandleLoginClient(conn)
	}

}

func HandleLoginClient(myConn net.Conn) {

	//waiting for login msg
	//then validate and send connection for world server back
	//or return error validation fail

	var clientResponse ClientMessage
	myDecoder := gob.NewDecoder(myConn)

	for {

		err := myDecoder.Decode(&clientResponse)
		checkError(err)

		if err == nil {

			if clientResponse.CommandType == CommandLogin {
				value := clientResponse.Value
				data := strings.Fields(value)
				username := data[0]
				password := data[1]

				fielName := "./Characters/" + username + ".xml"

				if _, err := os.Stat(fielName); err == nil {

					xmlFile, err := os.Open(fielName)
					checkError(err)
					defer xmlFile.Close()

					//should send charData to world server
					//as an notification that a character is going to join

					var charData CharacterXML
					XMLdata, _ := ioutil.ReadAll(xmlFile)
					xml.Unmarshal(XMLdata, &charData)
					if password == charData.Password {

						//find the world's addr for character to respawn
						lookupWorld := charData.CurrentWorld
						found := false
						var addr string
						for i := 0; i < serverNum; i++ {
							if serverNames[i] == lookupWorld {
								found = true
								addr = serverAddrs[i]
								break
							}
						}

						if found {
							//send world addr back to client
							var svMsg ServerMessage
							svMsg.MsgType = CommandRedirectServer
							svMsg.MsgDetail = addr
							gob.NewEncoder(myConn).Encode(svMsg)
						} else {
							var svMsg ServerMessage
							svMsg.MsgType = ErrorWorldIsNotFound
							svMsg.MsgDetail = "error world is not found"
							gob.NewEncoder(myConn).Encode(svMsg)
						}
					} else {
						var svMsg ServerMessage
						svMsg.MsgType = ErrorAuthorizationFail
						svMsg.MsgDetail = "error password is not correct"
						gob.NewEncoder(myConn).Encode(svMsg)
					}
				} else {
					var svMsg ServerMessage
					svMsg.MsgType = ErrorAuthorizationFail
					svMsg.MsgDetail = "error cannot find user data"
					gob.NewEncoder(myConn).Encode(svMsg)
				}

			} else {
				var svMsg ServerMessage
				svMsg.MsgType = ErrorUnexpectedCommand
				svMsg.MsgDetail = "error unexpected command"
				gob.NewEncoder(myConn).Encode(svMsg)
			}

		} else {
			break
		}
	}
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
