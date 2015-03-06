package main

import (
	"fmt"
	"strings"
	"strconv"
)

// Enumeration for movement/exit directions
const (
	NORTH = 0
	SOUTH = 1
	EAST = 2
	WEST = 3
	NORTH_EAST = 4
	NORTH_WEST = 5
	SOUTH_EAST = 6
	SOUTH_WEST = 7
	UP = 8
	DOWN = 9
)

type Room struct {
	ID int
	Description string
	Exits [10]int
	ExitLinksToRooms [10]*Room
	//Location string	
}

func newRoom(id string, descr string, exitWithLinks string) *Room{
	
	temp := strings.Split(id, "\r\n")

	id = temp[0]
	rID, err := strconv.Atoi(id)
	checkError(err)
	room := Room{ ID: rID, 
				Description: descr, 
			}

	for i := 0; i < 10; i++ {
		room.Exits[i] = -1
	}


	for index, element := range strings.Split(exitWithLinks, "\r\n"){
		_ = index
		dirAndRoom := strings.Split(element, " ")
		if len(dirAndRoom) != 2 {
			break
		}
		fmt.Println("Element: |", element, "| ")
		fmt.Println("DirAndRoom: ", dirAndRoom)
		dir := convertDirectionToInt(dirAndRoom[0])
		temp = strings.Split(dirAndRoom[1], "\r\n")
		roomID, err := strconv.Atoi(temp[0])
		checkError(err)
		
		room.Exits[dir] = roomID
	}

	return &room
}

func convertDirectionToInt(direction string) int {
	
	switch strings.ToLower(direction) {
		case "n" , "n\r\n" , "n\n" : return NORTH
		case "s" , "s\r\n" , "s\n" : return SOUTH
		case "e" , "e\r\n" , "e\n" : return EAST
		case "w" , "w\r\n" , "w\n" : return WEST
		case "nw", "nw\r\n", "nw\n": return NORTH_WEST
		case "ne", "ne\r\n", "ne\n": return NORTH_EAST
		case "sw", "sw\r\n", "sw\n": return SOUTH_WEST
		case "se", "se\r\n", "se\n": return SOUTH_EAST
		case "u" , "u\r\n" , "u\n" : return UP
		case "d" , "d\r\n" , "d\n" : return DOWN
	}
	
	return -1
}

func (room *Room) setRoomLink(roomLink [4]*Room){
	for i := 0; i < 10; i++ {
		if room.Exits[i] != -1 {
			//fmt.Println("Add: ", room.Exits[i], ", for room: ", room.ID)
			//fmt.Println("\tAdd: ", roomLink[room.Exits[i]].ID)
			room.ExitLinksToRooms[i] = roomLink[room.Exits[i]]
		}
	}
}

func (room *Room) getRoomLink(exit int) *Room{
	return room.ExitLinksToRooms[exit]
}