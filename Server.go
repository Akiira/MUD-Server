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
	//_ "github.com/go-sql-driver/mysql"
	"bufio"
	//"log"
	"net"
	"os"
)

var databaseG *sql.DB //The G means its a global var
var eventManager *EventManager

func main() {
	//populateTestData()
	//MovementAndCombatTest()
	runServer()
}

func runServer() {
	eventManager = newEventManager()

	listener := setUpServer()

	go createDummyMsg()

	for {
		conn, err := listener.Accept()
		checkError(err)
		fmt.Println("Connection established")

		go HandleClient(conn)
	}
}

func createDummyMsg() {

	for {
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		eventManager.sendMessageToRoom(text)
	}
}

func HandleClient(myConn net.Conn) {

	clientConnection := newClientConnection(myConn, eventManager)
	_ = clientConnection

	clientConnection.receiveMsgFromClient()
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
