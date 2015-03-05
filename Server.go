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
	_ "github.com/go-sql-driver/mysql"
	"net"
	"os"
	"sync"
)

var databaseG *sql.DB //The G means its a global var
var eventQueuMutexG sync.Mutex

func main() {

	//databaseTest()
	//GobTest()
	//LogInTest()
	intializeDatabaseConnection()

	listener := setUpServer()

	
	for{
		conn, err := listener.Accept()
		checkError(err)
		fmt.Println("Connection established")
	
		go handleClient(conn)
		//handleClient(conn)	
	}
	
	
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func intializeDatabaseConnection(){
	var err error
	databaseG, err = sql.Open("mysql",
		"admin1:admin@tcp(127.0.0.1:3306)/mud-database")
	checkError(err)

	err = databaseG.Ping()
	checkError(err)
}

func isGoodLogin(name string, pw string) bool{
	rows, err := databaseG.Query("select Password from Login where CharacterNameLI = ?", name)
	
	checkError(err)
	defer rows.Close()
	
	var DBpassword string
	
	if( rows.Next()){
		err := rows.Scan(&DBpassword)
		checkError(err)
		if(DBpassword == pw){
			return true
		}
	}
	
	return false
}

func handleClient(client net.Conn){
	//encoder := gob.NewEncoder(client)
	decoder := gob.NewDecoder(client)
	
	var clientsMessage ClientMessage
	decoder.Decode(&clientsMessage)

	fmt.Println("clients message: " + clientsMessage.Value)
	
	if(isGoodLogin(clientsMessage.getUsername(), clientsMessage.getPassword())){
		fmt.Println("Good Login!")
	}else{
		fmt.Println("Bad Login!")
	}
	databaseG.Close() //TODO remove these closes 
	client.Close()
}

func setUpServer() *net.TCPListener{
	service := "127.0.0.1:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	return listener
}
