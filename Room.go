package main

import (
	"github.com/daviddengcn/go-colortext"
	"strings"
)

// Enumeration for movement/exit directions
const (
	NORTH      = 0
	SOUTH      = 1
	EAST       = 2
	WEST       = 3
	NORTH_EAST = 4
	NORTH_WEST = 5
	SOUTH_EAST = 6
	SOUTH_WEST = 7
	UP         = 8
	DOWN       = 9
)

type Room struct {
	Name        string
	ID          int
	Description string

	// This represents each directions exit and has the room number of the connected
	// room or -1 if no exit in that direction. This can probaly be combined
	// with ExitLinksToRooms.
	Exits            [10]int
	ExitLinksToRooms [10]*Room

	//This represents the player characters (PCs) in the room
	CharactersInRoom map[string]*Character

	//This represents the attackable non-playable characters (NPCs) in the room
	MonstersInRoom map[string]*Monster

	//May have a third mapping to friendly NPCs like shopkeepers
	//NonCharactersInRoom map[string]*NPC

	ItemsInRoom map[string]*Item
}

//This is a constructor that creates a room from xml data
func newRoomFromXML(roomData RoomXML) *Room {
	room := Room{
		Name:        roomData.Name,
		ID:          roomData.ID,
		Description: roomData.Description,
	}
	for i := 0; i < 10; i++ {
		room.Exits[i] = -1
	}

	for _, roomExit := range roomData.Exits {
		room.Exits[convertDirectionToInt(roomExit.Direction)] = roomExit.ConnectedRoomID
	}

	room.CharactersInRoom = make(map[string]*Character)
	room.MonstersInRoom = make(map[string]*Monster)
	room.ItemsInRoom = make(map[string]*Item)
	return &room
}

//This function must be called after all rooms are created and is
// responsible for seting the exit pointers to point at the correct rooms
func (room *Room) setRoomLink(roomLink []*Room) {
	for i := 0; i < 10; i++ {
		if room.Exits[i] != -1 {
			room.ExitLinksToRooms[i] = roomLink[room.Exits[i]]
		}
	}
}

func (room *Room) getRoomLink(exit int) *Room {
	return room.ExitLinksToRooms[exit]
}

func (room *Room) addPCToRoom(char *Character) {
	room.CharactersInRoom[char.Name] = char
}

func (room *Room) removePCFromRoom(charName string) {
	delete(room.CharactersInRoom, charName)
}

func (room *Room) getItem(char *Character, itemName string) []FormattedString {

	item := room.ItemsInRoom[itemName]
	char.addItemToInventory(*item)

	delete(room.ItemsInRoom, itemName)
	output := make([]FormattedString, 1, 1)

	output[0].Color = ct.White
	output[0].Value = "You succesfully picked up the item and added it to your invenctory"

	return output
}

func (room *Room) getMonster(monsterName string) *Monster {
	return room.MonstersInRoom[monsterName]
}

func (room *Room) killOffMonster(monsterName string) {
	delete(room.MonstersInRoom, monsterName)
	room.ItemsInRoom[monsterName] = &Item{name: monsterName + " corpse", description: "A freshly kill " + monsterName + " corpse."}
}

func (room *Room) populateRoomWithMonsters() { //TODO remove hardcoding, maybe load from xml file
	room.MonstersInRoom["Rabbit"] = newMonsterFromName("Rabbit")
	room.MonstersInRoom["Fox"] = newMonsterFromName("Deer")
	room.MonstersInRoom["Deer"] = newMonsterFromName("Fox")
}

func (room *Room) getRoomDescription() []FormattedString {
	var output string
	formattedString := make([]FormattedString, 5, 5)

	formattedString[0].Color = ct.Green
	formattedString[0].Value = room.Name
	formattedString[1].Color = ct.White
	formattedString[1].Value = "-------------------------------------------------\n"
	formattedString[1].Value += room.Description
	formattedString[2].Color = ct.Magenta

	output = "Exits: "
	for i := 0; i < 10; i++ {
		if room.Exits[i] >= 0 {
			output += convertIntToDirection(i) + " "
		}
	}
	formattedString[2].Value = output

	output = ""
	formattedString[3].Color = ct.Yellow
	for _, itemPtr := range room.ItemsInRoom {
		output += "\n\t" + itemPtr.description
	}
	formattedString[3].Value = output
	formattedString[4].Color = ct.Red
	output = ""
	for key, _ := range room.MonstersInRoom {
		output += "\n\t" + key
	}
	formattedString[4].Value = output
	return formattedString
}

func convertDirectionToInt(direction string) int {

	switch strings.ToLower(direction) {
	case "n", "n\r\n", "n\n":
		return NORTH
	case "s", "s\r\n", "s\n":
		return SOUTH
	case "e", "e\r\n", "e\n":
		return EAST
	case "w", "w\r\n", "w\n":
		return WEST
	case "nw", "nw\r\n", "nw\n":
		return NORTH_WEST
	case "ne", "ne\r\n", "ne\n":
		return NORTH_EAST
	case "sw", "sw\r\n", "sw\n":
		return SOUTH_WEST
	case "se", "se\r\n", "se\n":
		return SOUTH_EAST
	case "u", "u\r\n", "u\n":
		return UP
	case "d", "d\r\n", "d\n":
		return DOWN
	}

	return -1
}

func convertIntToDirection(direction int) string {

	switch direction {
	case 0:
		return "North"
	case 1:
		return "South"
	case 2:
		return "East"
	case 3:
		return "West"
	case 4:
		return "North-West"
	case 5:
		return "North-East"
	case 6:
		return "South-West"
	case 7:
		return "South-East"
	case 8:
		return "Up"
	case 9:
		return "Down"
	}

	return ""
}
