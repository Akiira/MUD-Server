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
var onlinePlayers map[string]*Character
var eventQueuMutexG sync.Mutex
var worldRoomsG [4]*Room

func main() {

//	ct.ChangeColor(ct.Red, true, ct.White, false)
//	fmt.Println("Test")
	onlinePlayers = make(map[string]*Character)
	onlinePlayers["Ragnar"] = new(Character)
	foo := onlinePlayers["Ragnar"]
	foo.Name = "Ragnar"
	foo.RoomIN = 0
	foo.HitPoints = 30
	worldRoomsG = loadRooms()
	
	worldRoomsG[0].populateRoomWithMonsters()
//	fmt.Println(worldRoomsG[0].MonstersInRoom["Rabbit"])


//	fmt.Println("length: ", len(worldRoomsG))
	//fmt.Println(rooms[1].Description)
	//fmt.Println(rooms[0].ExitLinksToRooms[1].Description)
	
	MovementAndCombatTest()
	//combatTest()
	//roomAndMoveTest2()
	//roomAndMoveTest(worldRoomsG)

	//databaseTest()
	//GobTest()
	//LogInTest()
	
//	var m map[string]Character
//	m = make(map[string]Character)
	
//	m["Ragnar"] = Character{ name: "Ragnar"}
	
//	fmt.Println(m["Ragnar"])
	
//	intializeDatabaseConnection()

//	listener := setUpServer()

	
//	for{
//		conn, err := listener.Accept()
//		checkError(err)
//		fmt.Println("Connection established")
	
//		go handleClient(conn)
//		//handleClient(conn)	
//	}
	
	
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


	//get clients character name
	// load info from database
}

func loadCharacterFromDB(characterName string){
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

func setUpServer() *net.TCPListener{
	service := "127.0.0.1:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	return listener
}
