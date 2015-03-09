package main

import (
	"strings"
	"github.com/daviddengcn/go-colortext"
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
	Name string
	ID int
	Description string
	
	// This represents each directions exit and has the room number of the connected
	// room or -1 if no exit in that direction. This can probaly be combined
	// with ExitLinksToRooms.
	Exits [10]int
	ExitLinksToRooms [10]*Room
	
	//This represents the player characters (PCs) in the room
	CharactersInRoom map[string]bool
	
	//This represents the attackable non-playable characters (NPCs) in the room
	MonstersInRoom map[string]*Monster
	
	//May have a third mapping to friendly NPCs like shopkeepers
	//NonCharactersInRoom map[string]*NPC
}

//This is a constructor that creates a room from xml data
func newRoomFromXML( roomData RoomXML) *Room {
	room := Room{
				Name: roomData.Name, 
				ID: roomData.ID, 
				Description: roomData.Description, 
			}
	for i := 0; i < 10; i++ {
		room.Exits[i] = -1
	}
	
	for _, roomExit := range roomData.Exits {
		room.Exits[convertDirectionToInt(roomExit.Direction)] = roomExit.ConnectedRoomID
	}
	
	room.CharactersInRoom = make(map[string]bool)
	room.MonstersInRoom = make(map[string]*Monster)
	return &room
}

//This function must be called after all rooms are created and is
// responsible for seting the exit pointers to point at the correct rooms
func (room *Room) setRoomLink(roomLink []*Room){
	for i := 0; i < 10; i++ {
		if room.Exits[i] != -1 {
			room.ExitLinksToRooms[i] = roomLink[room.Exits[i]]
		}
	}
}

func (room *Room) getRoomLink(exit int) *Room{
	return room.ExitLinksToRooms[exit]
}

func (room *Room) addPCToRoom(charName string) {
	room.CharactersInRoom[charName] = true;
}

func (room *Room) removePCFromRoom(charName string) {
	delete(room.CharactersInRoom, charName)
}

func (room *Room) getMonster(monsterName string) *Monster {
	return room.MonstersInRoom[monsterName]
}

func (room *Room) killOffMonster(monsterName string) {
	delete(room.MonstersInRoom, monsterName)
}

func (room *Room) populateRoomWithMonsters() { //TODO remove hardcoding, maybe load from xml file
	room.MonstersInRoom["Rabbit"] = &Monster{HP: 5, Name: "Rabbit"}
	room.MonstersInRoom["Fox"] = &Monster{HP: 10, Name: "Fox"}
	room.MonstersInRoom["Deer"] = &Monster{HP: 7, Name: "Deer"}
}

func (room *Room) getFormattedOutput() []FormattedString{
	var output string
	formattedString := make([]FormattedString, 4, 4)
	
	formattedString[0].Color = ct.Green
	formattedString[0].Value = room.Name 
	formattedString[1].Color = ct.White;
	formattedString[1].Value = "-------------------------------------------------\n"
	formattedString[1].Value += room.Description
	formattedString[2].Color = ct.Magenta

	output = "Exits: "
	for i:= 0; i < 10; i++ {
		if( room.Exits[i] >= 0 ) {
			output += convertIntToDirection(i) + " "
		}
	}
	formattedString[2].Value = output
	formattedString[3].Color = ct.Red
	output = ""
	for key, _ := range room.MonstersInRoom {
		output += "\n\t" + key 
	}
	formattedString[3].Value = output
	return formattedString
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

func convertIntToDirection(direction int) string {
	
	switch direction {
		case 0 : return "North"
		case 1 : return "South"
		case 2 : return "East"
		case 3 : return "West"
		case 4 : return "North-West"
		case 5 : return "North-East"
		case 6 : return "South-West"
		case 7 : return "South-East"
		case 8 : return "Up"
		case 9 : return "Down"
	}
	
	return ""
}
