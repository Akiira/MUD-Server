// roomLoader
package main

import (
	//"strings"
	"io/ioutil"
	//"fmt"
	"encoding/xml"
	"os"
)

const (
	ROOM_ID = 0
	ROOM_DESCRIPTION = 1
	ROOM_EXITS = 2
)

type ExitXML struct {
	XMLName xml.Name `xml:"Exit"`
	Direction string `xml:"Direction"`
	ConnectedRoomID int `xml:"RoomID"`
}

type RoomXML struct {
	XMLName xml.Name `xml:"Room"`
	ID int `xml:"ID"`
	Name string `xml:"Name"`
	Description string `xml:"Description"`
	Exits []ExitXML `xml:"Exit"`
}

type RoomsXML struct {
	XMLName xml.Name `xml:"Rooms"`
	Rooms []RoomXML `xml:"Room"`
}

func loadRooms() [4]*Room {
	xmlFile, err := os.Open("roomData.xml")
	checkError(err)
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var roomsData RoomsXML
    xml.Unmarshal(XMLdata, &roomsData)
	
	var rooms [4]*Room

	for index, roomData := range roomsData.Rooms {
		rooms[index] = newRoomFromXML(roomData)
	}
	
	for index := range rooms {
		rooms[index].setRoomLink(rooms)
	}

	return rooms
}

//func LoadRooms() [4]*Room {
//	bytes, err := ioutil.ReadFile("testRooms1.txt")
//	checkError(err)
//	var rooms [4]*Room
	
//	foo := string(bytes)
	
//	//fmt.Println(foo)
	
//	roomsString := strings.Split(foo, "#R\r\n")

//	for index,element := range roomsString {
//		roomStr := strings.Split(element, "#\r\n")
//		fmt.Println(roomStr[0])
//		//fmt.Println(roomStr[1])
//		//fmt.Println(roomStr[2])
		
//		rooms[index] = newRoom(roomStr[ROOM_ID], roomStr[ROOM_DESCRIPTION], roomStr[ROOM_EXITS])
//	}
	
//	for index := range rooms {
//		rooms[index].setRoomLink(rooms)
//	}
	
//	return rooms
//}
