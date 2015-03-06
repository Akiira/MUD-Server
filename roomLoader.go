// roomLoader
package main

import (
	"strings"
	"io/ioutil"
	"fmt"
)
import (
	_ "fmt"
	_ "os"
	_ "strings"
)

const (
	ROOM_ID = 0
	ROOM_DESCRIPTION = 1
	ROOM_EXITS = 2
)


func LoadRooms() [4]*Room {
	bytes, err := ioutil.ReadFile("testRooms1.txt")
	checkError(err)
	var rooms [4]*Room
	
	foo := string(bytes)
	
	//fmt.Println(foo)
	
	roomsString := strings.Split(foo, "#R\r\n")

	for index,element := range roomsString {
		roomStr := strings.Split(element, "#\r\n")
		fmt.Println(roomStr[0])
		//fmt.Println(roomStr[1])
		//fmt.Println(roomStr[2])
		
		rooms[index] = newRoom(roomStr[ROOM_ID], roomStr[ROOM_DESCRIPTION], roomStr[ROOM_EXITS])
	}
	
	for index := range rooms {
		rooms[index].setRoomLink(rooms)
	}
	
	return rooms
}
