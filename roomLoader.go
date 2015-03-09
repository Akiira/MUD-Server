// roomLoader
package main

//TODO Should these be moved to the room class?
//		or perhaps some special initialization file that 
//		handles all server start up activities.

import (
	"io/ioutil"
	"encoding/xml"
	"os"
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

func loadRooms() []*Room {
	xmlFile, err := os.Open("roomData.xml")
	checkError(err)
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var roomsData RoomsXML
    xml.Unmarshal(XMLdata, &roomsData)
	
	//var rooms [4]*Room
	rooms := make([]*Room, 4, 4)

	for index, roomData := range roomsData.Rooms {
		rooms[index] = newRoomFromXML(roomData)
	}
	
	for index := range rooms {
		rooms[index].setRoomLink(rooms)
	}

	return rooms
}

