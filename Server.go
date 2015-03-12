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
	"sync"
)

var databaseG *sql.DB //The G means its a global var
var onlinePlayers map[string]*Character
var eventQueuMutexG sync.Mutex

//might be easier to read from file and take command from console to update list by read from file again
//can shutdown by command from console
var worldRoomsG []*Room
var numEventManagerG int

func main() {

	//this should be the one that read list of servers, including central server
	/*
		file, err := os.Open("serverList.txt")
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
	*/

	populateTestData()

	//Pattanapoom Hand
	//start model

	//	eventManagersG[0] = new(EventManager)
	//	//might change to init() later
	//	(*eventManagersG[0]).numListener = 0

	//	listener := setUpServer()

	//	go createDummyMsg()

	//	for {
	//		conn, err := listener.Accept()
	//		checkError(err)
	//		fmt.Println("Connection established")

	//		go dummyHandleClient(conn)
	//	}

	/*for {
		time.Sleep(1 * time.Microsecond)
	}*/
}

func createDummyMsg() {

	for {

		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter text: ")
		text, _ := reader.ReadString('\n')
		(*eventManagersG[0]).dummySentMsg(text)

	}
}

func dummyHandleClient(myConn net.Conn) {

	//*(eventManagersG[0]).subscribeListener()
	newChar := new(Character)
	(*newChar).init(myConn, "name", eventManagersG[0])
	(eventManagersG[0]).subscribeListener(newChar)
	go (*newChar).receiveMessage()
	//character should start routine to get msg from client here then return to main loop in server
	// however eventmanager, in its own routine, can call method getEventmessage on character instantly without extra routine nor channel
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

func loadCharacterFromDB(characterName string) {
	//	rows, err := databaseG.Query("select * from Character where CharacterName = ?", characterName)
	//	checkError(err)
	//	defer rows.Close()

	//	var char Character

	//	if( rows.Next()){
	//		err := rows.Scan(&char.Name, ....)
	//		checkError(err)
	//		if(DBpassword == pw){
	//			return true
	//		}
	//	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
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

func setUpServer() *net.TCPListener {
	service := "127.0.0.1:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	return listener
}
