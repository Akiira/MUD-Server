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

