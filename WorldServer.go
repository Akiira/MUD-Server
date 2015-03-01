/* DaytimeServer
 */
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
	"log"
	"net"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

type ClientMessage struct {
	Command int
	Value   string
}

type ServerMessage struct {
	Value string
}

func main() {

	//databaseTest()

	service := "127.0.0.1:1200"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)

	conn, err := listener.Accept()
	checkError(err)
	fmt.Println("Connection established")
	encoder := gob.NewEncoder(conn)
	decoder := gob.NewDecoder(conn)

	var clientsMessage ClientMessage
	decoder.Decode(&clientsMessage)

	fmt.Println("clients message: " + clientsMessage.Value)
	reply := ServerMessage{Value: "This is the servers reply"}
	encoder.Encode(reply)

	conn.Close()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func databaseTest() {
	db, err := sql.Open("mysql",
		"admin1:admin@tcp(127.0.0.1:3306)/sakila")

	if err != nil {
		panic(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}
	fmt.Println("After if stmt")
	defer db.Close()

	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("After 2nd if stmt")
	var name string
	rows, err := db.Query("select title from film")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	for rows.Next() {
		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(name)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
}

func inventoryAndItemTest() {
	var i Item
	i.description = "a cool sword"
	i.name = "Short Sword"
	i.itemID = 666

	var i2 Item
	i2.description = "a cool stick"
	i2.name = "Oaken Bo"
	i2.itemID = 667

	items := [100]Item{i, i2}

	inventory := createInventory(items)

	var item = inventory.getItemByName("Short Sword")

	fmt.Println(item.name)
}
