// test
package main

import (
	//"database/sql"
	//"encoding/gob"
	"fmt"
	//"log"
	//"net"
	"github.com/daviddengcn/go-colortext"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func populateTestData() {
	monsterTemplatesG = make(map[string]*Monster)
	onlinePlayers = make(map[string]*Character)
	loadMonsterData()
	worldRoomsG = loadRooms()
	loadCharacterData("Ragnar")
	onlinePlayers["Ragnar"].addItemToInventory(Item{name: "Brown-Hat", description: "A moldy old hat with holes in it"})
	onlinePlayers["Ragnar"].addItemToInventory(Item{name: "Walking-Stick", description: "A sturdy walking stick made of oak"})

	worldRoomsG[0].populateRoomWithMonsters()
}

func printFormatedOutput(output []FormattedString) {
	for _, element := range output {
		ct.ChangeColor(element.Color, false, ct.Black, false)
		fmt.Println(element.Value)
	}
	ct.ResetColor()
}

func MovementAndCombatTest() {
	var input string
	printFormatedOutput(worldRoomsG[0].getRoomDescription())
	output := make([]FormattedString, 5, 5)
	for {
		read, err := fmt.Scan(&input)
		checkError(err)
		_ = read

		if input == "exit" {
			break
		} else if input == "attack" {
			var target string
			read, err = fmt.Scan(&target)
			checkError(err)

			output = onlinePlayers["Ragnar"].makeAttack(target)
		} else if input == "look" {
			var target string
			read, err = fmt.Scan(&target)
			checkError(err)
			if target == "room" {
				output = worldRoomsG[onlinePlayers["Ragnar"].RoomIN].getRoomDescription()
			} else {
				output = monsterTemplatesG[target].getLookDescription()
			}
		} else if input == "get" {

		} else { //assume movement
			output = onlinePlayers["Ragnar"].moveCharacter(input)
		}

		printFormatedOutput(output)
	}
}

//func LogInWithClientTest() {
//	service := "127.0.0.1:1200"
//	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
//	checkError(err)

//	listener, err := net.ListenTCP("tcp", tcpAddr)
//	checkError(err)

//	conn, err := listener.Accept()
//	checkError(err)
//	fmt.Println("Connection established")
//	encoder := gob.NewEncoder(conn)
//	decoder := gob.NewDecoder(conn)

//	var clientsMessage ClientMessage
//	decoder.Decode(&clientsMessage)

//	fmt.Println("clients message: " + clientsMessage.Value)

//	db, err := sql.Open("mysql",
//		"admin1:admin@tcp(127.0.0.1:3306)/mud-database")
//	checkError(err)
//	defer db.Close()

//	err = db.Ping()
//	checkError(err)

//	rows, err := db.Query("select CharacterNameLI, Password from Login where CharacterNameLI = ?", )
//	defer rows.Close()
//	var name string
//	for rows.Next() {
//		err := rows.Scan(&name)
//		if err != nil {
//			log.Fatal(err)
//		}
//		fmt.Println(name)
//	}

//	checkError(rows.Err())

//	reply := ServerMessage{Value: "This is the servers reply"}
//	encoder.Encode(reply)
//}

//func inventoryAndItemTest() {
//	var i Item
//	i.description = "a cool sword"
//	i.name = "Short Sword"
//	i.itemID = 666

//	var i2 Item
//	i2.description = "a cool stick"
//	i2.name = "Oaken Bo"
//	i2.itemID = 667

//	items := [100]Item{i, i2}

//	inventory := createInventory(items)

//	var item = inventory.getItemByName("Short Sword")

//	fmt.Println(item.name)
//}
