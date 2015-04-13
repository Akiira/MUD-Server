package main

import (
	"encoding/xml"
	"github.com/daviddengcn/go-colortext"
	"io/ioutil"
	"os"
	"strings"
	"time"
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

	LocalWorld = "local"
)

type Room struct {
	Name        string
	ID          int
	Description string
	WorldID     string

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

	//This is for the monsters native to this room
	monsterTemplateNames []string
}

//This is a constructor that creates a room from xml data
func newRoomFromXML(roomData RoomXML) *Room {
	room := Room{
		Name:        roomData.Name,
		ID:          roomData.ID,
		Description: roomData.Description,
		WorldID:     roomData.WorldID,
	}

	for i := 0; i < 10; i++ {
		room.Exits[i] = -1
	}

	for _, roomExit := range roomData.Exits {
		room.Exits[convertDirectionToInt(roomExit.Direction)] = roomExit.ConnectedRoomID
	}

	if room.isLocal() {
		room.CharactersInRoom = make(map[string]*Character)
		room.MonstersInRoom = make(map[string]*Monster)
		room.ItemsInRoom = make(map[string]*Item)
		room.monsterTemplateNames = roomData.Monsters

		room.populateRoomWithMonsters()
		go room.repopulateRoomTick(15)
	}
	return &room
}

//This function must be called after all rooms are created and is
// responsible for seting the exit pointers to point at the correct rooms
func (room *Room) setRoomLink(roomLink map[int]*Room) {
	for i := 0; i < 10; i++ {
		if room.Exits[i] != -1 {
			room.ExitLinksToRooms[i] = roomLink[room.Exits[i]]
		}
	}
}

func (room *Room) isValidDirection(dir int) bool {
	return dir < 10 && dir >= 0 && room.Exits[dir] >= 0
}

func (room *Room) isLocal() bool {
	return room.WorldID == LocalWorld
}

func (room *Room) getConnectedRoom(exit int) *Room {
	return room.ExitLinksToRooms[exit]
}

func (room *Room) addPCToRoom(char *Character) {

	if room.isLocal() {
		room.CharactersInRoom[char.Name] = char
	}
	char.RoomIN = room.ID
}

func (room *Room) removePCFromRoom(charName string) {
	if char, found := room.getPC(charName); found {
		char.RoomIN = -1
		delete(room.CharactersInRoom, charName)
	}
}

func (room *Room) getPC(charName string) (*Character, bool) {
	if room.CharactersInRoom != nil {
		char, found := room.CharactersInRoom[charName]

		return char, found
	} else {
		return nil, false
	}
}

func (room *Room) getItem(char *Character, itemName string) []FormattedString {

	item, found := room.ItemsInRoom[itemName]

	if found {
		char.addItemToInventory(item)
		delete(room.ItemsInRoom, itemName)

		return newFormattedStringSplice("You succesfully picked up the item and added it to your invenctory")
	} else {
		return newFormattedStringSplice("That item was not found in the room.")
	}

}

func (room *Room) getMonster(monsterName string) *Monster {

	if val, ok := room.MonstersInRoom[monsterName]; ok {
		//fmt.Println(room.MonstersInRoom[monsterName])
		return val
	} else { // in case, it's already dead, return nil
		return nil
	}
}

func (room *Room) getAgentInRoom(name string) Agenter {

	if val := room.getMonster(name); val != nil {
		return val
	} else if val, found := room.getPC(name); found {
		return val
	} else { // in case, it's already dead, return nil
		return nil
	}
}

func (room *Room) killOffMonster(monsterName string) {
	delete(room.MonstersInRoom, monsterName)
	room.ItemsInRoom[monsterName] = &Item{name: monsterName + " corpse", description: "A freshly kill " + monsterName + " corpse."}
}

func (room *Room) repopulateRoomTick(timeInMinutes time.Duration) {
	for {
		//TODO should also periodically remove corpses and items in room
		//Repopulate the room every x minutes
		time.Sleep(time.Minute * timeInMinutes)
		room.populateRoomWithMonsters()
	}
}

func (room *Room) populateRoomWithMonsters() {

	for _, monsterName := range room.monsterTemplateNames {
		if _, found := room.MonstersInRoom[monsterName]; found == false {
			room.MonstersInRoom[monsterName] = newMonsterFromName(monsterName, room.ID)
		}
	}
}

func (room *Room) getRoomDescription() []FormattedString {
	var output string
	fs := newFormattedStringCollection()

	fs.addMessage(ct.Green, room.Name)
	fs.addMessage(ct.White, "\n-------------------------------------------------\n")
	fs.addMessage2(room.Description)

	output = "\nExits: "
	for i := 0; i < 10; i++ {
		if room.Exits[i] >= 0 {
			output += convertIntToDirection(i) + " "
		}
	}

	fs.addMessage(ct.Magenta, output)
	output = ""

	for _, itemPtr := range room.ItemsInRoom {
		output += "\n\t" + itemPtr.description
	}
	fs.addMessage(ct.Yellow, output)
	output = ""
	for key, _ := range room.MonstersInRoom {
		output += "\n\t" + key
	}
	fs.addMessage(ct.Red, output)
	output = ""
	for key, _ := range room.CharactersInRoom {
		output += "\n\t" + key
	}
	fs.addMessage(ct.Blue, output+"\n")
	return fs.fmtedStrings
}

//===================="STATIC" FUNCTIONS======================//

func convertDirectionToInt(direction string) int {

	switch strings.ToLower(direction) {
	case "n", "n\r\n", "n\n", "north":
		return NORTH
	case "s", "s\r\n", "s\n", "south":
		return SOUTH
	case "e", "e\r\n", "e\n", "east":
		return EAST
	case "w", "w\r\n", "w\n", "west":
		return WEST
	case "nw", "nw\r\n", "nw\n", "northwest":
		return NORTH_WEST
	case "ne", "ne\r\n", "ne\n", "northeast":
		return NORTH_EAST
	case "sw", "sw\r\n", "sw\n", "southwest":
		return SOUTH_WEST
	case "se", "se\r\n", "se\n", "southeast":
		return SOUTH_EAST
	case "u", "u\r\n", "u\n", "up":
		return UP
	case "d", "d\r\n", "d\n", "down":
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

type ExitXML struct {
	XMLName         xml.Name `xml:"Exit"`
	Direction       string   `xml:"Direction"`
	ConnectedRoomID int      `xml:"RoomID"`
}

//TODO add localWorld field
type RoomXML struct {
	XMLName     xml.Name  `xml:"Room"`
	ID          int       `xml:"ID"`
	Name        string    `xml:"Name"`
	Description string    `xml:"Description"`
	WorldID     string    `xml:"WorldID"`
	Monsters    []string  `xml:"Monster"`
	Exits       []ExitXML `xml:"Exit"`
}

type RoomsXML struct {
	XMLName xml.Name  `xml:"Rooms"`
	Rooms   []RoomXML `xml:"Room"`
}

func loadRooms(worldName string) map[int]*Room {
	xmlFile, err := os.Open(worldName + ".xml")
	checkErrorWithMessage(err, true, " In load rooms function.")
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var roomsData RoomsXML
	xml.Unmarshal(XMLdata, &roomsData)

	rooms := make(map[int]*Room)

	for _, roomData := range roomsData.Rooms {
		getRoom := newRoomFromXML(roomData)
		rooms[getRoom.ID] = getRoom
	}

	for _, room := range rooms {
		room.setRoomLink(rooms)
	}
	return rooms
}
