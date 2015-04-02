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
	_ "github.com/daviddengcn/go-colortext"
	"io"
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
	//runServer()
	getCharactersFile("Ragnar")
}

func runServer() {
	eventManager = newEventManager()

	listener := setUpServerWithPort(1300)

	for {
		conn, err := listener.Accept()
		//checkError(err)
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
	var clientResponse ClientMessage

	err := gob.NewDecoder(myConn).Decode(&clientResponse)
	checkError(err)

	getCharactersFile(clientResponse.getUsername())
	playerChar := getCharacterFromFile(clientResponse.getUsername())

	clientConnection := newClientConnection(myConn, eventManager, playerChar)
	clientConnection.receiveMsgFromClient()
}

func getCharactersFile(name string) {
	conn, err := net.Dial("tcp", servers["characterStorage"])
	checkError(err)
	defer conn.Close()

	gob.NewEncoder(conn).Encode(&ServerMessage{Value: newFormattedString(name)})

	file, err := os.Create("Characters/" + name + ".xml")
	checkError(err)
	defer file.Close()

	_, err = io.Copy(file, conn)
	checkError(err)
}
