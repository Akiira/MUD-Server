package main

//"database/sql"
//_ "github.com/go-sql-driver/mysql"
// port 3306, tcp
// user: admin1
// pw: admin
import (
	"database/sql"
	"encoding/gob"
	"fmt"
	"strings"
	//_ "github.com/go-sql-driver/mysql"
	"bufio"
	"log"
	"net"
	"os"
	"strconv"
)

const centralServer = "central"

var serverNames [10]string
var serverAddrs [10]string
var serverNum int
var databaseG *sql.DB //The G means its a global var
var eventManager *EventManager

func main() {

	//populateTestData()
	//MovementAndCombatTest()
	readServerList()

	eventManager = newEventManager()

	listener := setUpServerWithPort(1300)

	go createDummyMsg()

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

func createDummyMsg() {

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		eventManager.sendMessageToRoom(text)
	}
}

func HandleClientLogin(myConn net.Conn) {

	var playerChar *Character

	//need to check for authentication first
	for {
		var clientResponse ClientMessage
		myDecoder := gob.NewDecoder(myConn)
		err := myDecoder.Decode(&clientResponse)

		if err == nil {
			if clientResponse.CommandType == CommandLogin {
				value := clientResponse.Value
				data := strings.Fields(value)
				username := data[0]
				password := data[1]
				playerChar = queryCharacterFromCentral(username, password)
			} else {
				var svMsg ServerMessage
				svMsg.MsgType = ErrorUnexpectedCommand
				svMsg.MsgDetail = "error unexpected command"
				gob.NewEncoder(myConn).Encode(svMsg)
			}
		}
	}
	//once the authentication is good player can continue on gameplay

	clientConnection := newClientConnection(myConn, eventManager, playerChar)
	_ = clientConnection

	//TODO this should actually only be called once at the servers start up
	go eventManager.waitForTick()

	clientConnection.receiveMsgFromClient()
}

func queryCharacterFromCentral(username string, password string) *Character {
	var address string
	for i := 0; i < serverNum; i++ {
		if serverNames[i] == centralServer {
			address = serverAddrs[i]
			break
		}
	}

	conn, err := net.Dial("tcp", address)
	//checkError(err)
	if err == nil {

		clientResponse := ClientMessage{CommandType: CommandQueryCharacter, Command: "queryCharacter", Value: username + " " + password}
		gob.NewEncoder(conn).Encode(clientResponse)
		var serversResponse ServerMessage
		gob.NewDecoder(conn).Decode(&serversResponse)

		if serversResponse.MsgType == CommandCharacterDetail {
			data := strings.Fields(serversResponse.MsgDetail)
			//name string, room int, hp int, def int
			room, err := strconv.Atoi(data[1])
			checkError(err)
			hp, err := strconv.Atoi(data[2])
			checkError(err)
			def, err := strconv.Atoi(data[3])
			checkError(err)
			queriedChar := newCharacter(data[0], room, hp, def)
			return queriedChar
		} else {
			fmt.Println(serversResponse.MsgDetail)
		}

	} else {
		return nil
	}

	return nil
}

func handleClient(client net.Conn) {
	//encoder := gob.NewEncoder(client)
	decoder := gob.NewDecoder(client)

	var clientsMessage ClientMessage
	decoder.Decode(&clientsMessage)

	fmt.Println("clients message: " + clientsMessage.Value)

	if isGoodLogin(clientsMessage.getUsername(), clientsMessage.getPassword()) {
		fmt.Println("Good Login!")
	} else {
		fmt.Println("Bad Login!")
	}
	databaseG.Close() //TODO remove these closes
	client.Close()

	//get clients character name
	// load info from database
}

func intializeDatabaseConnection() {
	var err error
	databaseG, err = sql.Open("mysql",
		"admin1:admin@tcp(127.0.0.1:3306)/mud-database")
	checkError(err)

	err = databaseG.Ping()
	checkError(err)
}

func isGoodLogin(name string, pw string) bool {
	rows, err := databaseG.Query("select Password from Login where CharacterNameLI = ?", name)

	checkError(err)
	defer rows.Close()

	var DBpassword string

	if rows.Next() {
		err := rows.Scan(&DBpassword)
		checkError(err)
		if DBpassword == pw {
			return true
		}
	}

	return false
}
