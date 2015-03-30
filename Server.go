package main

//"database/sql"
//_ "github.com/go-sql-driver/mysql"
// port 3306, tcp
// user: admin1
// pw: admin
import (
	"bufio"
	"encoding/gob"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"log"
	"net"
	"os"
	"strings"
)

const centralServer = "central"

var servers map[string]string

var eventManager *EventManager

func main() {
	//populateTestData()
	//MovementAndCombatTest()

	readServerList()
	runServer()
}

func runServer() {
	eventManager = newEventManager()

	listener := setUpServerWithPort(1300)

	for {
		conn, err := listener.Accept()
		//checkError(err)
		if err == nil {
			fmt.Println("Connection established")

			go HandleClientLogin(conn)
		}
	}
}

func readServerList() {
	//this should be the one that read list of servers, including central server

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

	//Pattanapoom Hand
	//start model
}

func HandleClientLogin(myConn net.Conn) {

	var playerChar *Character

	//need to check for authentication first

	var clientResponse ClientMessage
	myDecoder := gob.NewDecoder(myConn)
	err := myDecoder.Decode(&clientResponse)

	if err == nil {
		if clientResponse.CommandType == CommandLogin {
			username := clientResponse.getUsername()
			password := clientResponse.getPassword()
			playerChar = queryCharacterFromCentral(username, password)
		} else {
			svMsg := newServerMessage(newFormattedString2(ct.Red, "error unexpected command"))
			gob.NewEncoder(myConn).Encode(svMsg)
		}
	}

	//once the authentication is good player can continue on gameplay

	clientConnection := newClientConnection(myConn, eventManager, playerChar)

	clientConnection.receiveMsgFromClient()
}

func queryCharacterFromCentral(username string, password string) *Character {

	address := servers["central"]

	conn, err := net.Dial("tcp", address)
	checkError(err)

	clientResponse := ClientMessage{Value: username}
	err = gob.NewEncoder(conn).Encode(clientResponse)
	checkError(err)

	var queriedChar Character
	err = gob.NewDecoder(conn).Decode(&queriedChar)
	checkError(err)

	return &queriedChar
}
