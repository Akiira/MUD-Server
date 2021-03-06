package main

import (
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"io/ioutil"
	"os"
	"strings"
	"sync"
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

	//TODO change this to an Inventory object
	ItemsInRoom map[string]Item_I
	items_mutex sync.Mutex

	//This is for the monsters native to this room
	monsterTemplateNames []string
}

//This is a constructor that creates a room from xml data
func NewRoom(roomData RoomXML) *Room {
	room := Room{
		Name:        roomData.Name,
		ID:          roomData.ID,
		Description: roomData.Description,
		WorldID:     roomData.WorldID,
		items_mutex: sync.Mutex{},
	}

	for i := 0; i < 10; i++ {
		room.Exits[i] = -1
	}

	for _, roomExit := range roomData.Exits {
		room.Exits[convertDirectionToInt(roomExit.Direction)] = roomExit.ConnectedRoomID
	}

	if room.IsLocal() {
		room.CharactersInRoom = make(map[string]*Character)
		room.MonstersInRoom = make(map[string]*Monster)
		room.ItemsInRoom = make(map[string]Item_I)

		for _, item := range roomData.Items.Items {
			room.AddItem(NewItem(item.(*ItemXML)))
		}

		room.monsterTemplateNames = roomData.Monsters

		room.PopulateMonsters()
		go room.repopulateRoomTick(15)
	}
	return &room
}

//setRoomLink must be called after all rooms are created and is
// responsible for seting the exit pointers to point at the correct rooms
func (room *Room) setRoomLink(roomLink map[int]*Room) {
	for i := 0; i < 10; i++ {
		if room.Exits[i] != -1 {
			room.ExitLinksToRooms[i] = roomLink[room.Exits[i]]
		}
	}
}

func (room *Room) IsAggroed(charName string) bool {
	for _, monster := range room.MonstersInRoom {
		if monster.IsAttackingPlayer(charName) {
			return true
		}
	}
	return false
}

func (room *Room) IsValidDirection(dir int) bool {
	return dir < 10 && dir >= 0 && room.Exits[dir] >= 0
}

func (room *Room) IsLocal() bool {
	return room.WorldID == os.Args[1]
}

func (room *Room) GetConnectedRoom(exit int) *Room {
	if exit != -1 {
		return room.ExitLinksToRooms[exit]
	} else {
		return nil
	}
}

func (room *Room) AddItem(itm Item_I) {
	if itm != nil {
		room.ItemsInRoom[itm.GetName()] = itm
	}
}

func (room *Room) AddPlayer(char *Character) {

	if room.IsLocal() {
		eventManager.SendMessageToRoom(room.ID, newServerMessageFS(newFormattedStringSplice2(ct.Blue, "\n"+char.Name+" has entered this room.")))
		room.CharactersInRoom[strings.ToLower(char.Name)] = char
	}
	char.RoomIN = room.ID
}

func (room *Room) RemovePlayer(charName string) {
	if _, found := room.GetPlayer(charName); found {
		eventManager.SendMessageToRoom(room.ID, newServerMessageFS(newFormattedStringSplice2(ct.Blue, "\n"+charName+" has left the room.")))
		delete(room.CharactersInRoom, strings.ToLower(charName))
	} else {
		fmt.Fprint(os.Stderr, "Failed to find: ", charName, " in Room: ", room.Name, ", in RemovePlayer\n")
	}
}

func (room *Room) UnAggroPlayer(charName string) {
	for _, monster := range room.MonstersInRoom {
		monster.RemoveTarget(charName)
	}
}

func (room *Room) GetPlayer(charName string) (*Character, bool) {
	if room.CharactersInRoom != nil {
		char, found := room.CharactersInRoom[strings.ToLower(charName)]
		return char, found
	} else {
		fmt.Fprint(os.Stderr, "Failed to find: ", charName, " in Room: ", room.Name, ", in GetPlayer\n")
		return nil, false
	}
}

func (room *Room) GetItem(itemName string) (Item_I, bool) {
	for name, item := range room.ItemsInRoom {
		lcName := strings.ToLower(name)

		if strings.Contains(lcName, itemName) {
			return item, true
		}
	}

	fmt.Fprint(os.Stderr, "Failed to find: ", itemName, " in Room: ", room.Name, ", in GetItem\n")
	return nil, false
}

func (room *Room) GetAndRemoveItem(itemName string) (Item_I, bool) {
	room.items_mutex.Lock()
	defer room.items_mutex.Unlock()

	if item, found := room.GetItem(itemName); found {
		delete(room.ItemsInRoom, item.GetName())
		return item, found
	}

	return nil, false
}

func (room *Room) GiveItemToPlayer(char *Character, itemName string) []FormattedString {

	if item, found := room.GetAndRemoveItem(itemName); found {
		char.AddItem(item)

		return newFormattedStringSplice("You succesfully picked up the item and added it to your invenctory\n")
	} else {
		return newFormattedStringSplice("\nThat item was not found in the room.\n")
	}
}

func (room *Room) GetMonster(monsterName string) *Monster {
	for name, mosnter := range room.MonstersInRoom {

		if strings.Contains(strings.ToLower(name), monsterName) {
			return mosnter
		}
	}

	return nil
}

func (room *Room) GetAgent(name string) (Agenter, bool) {

	if val := room.GetMonster(name); val != nil {
		return val, true
	} else if val, found := room.GetPlayer(name); found {
		return val, true
	} else { // in case, it's already dead, return nil
		return nil, false
	}
}

func (room *Room) KillOffMonster(monsterName string) {
	if monster := room.GetMonster(monsterName); monster != nil {
		drops := room.MonstersInRoom[monsterName].GetLootAndCorpse()
		for _, drop := range drops {
			room.AddItem(drop)
		}
		delete(room.MonstersInRoom, monsterName)
	}
}

func (room *Room) repopulateRoomTick(timeInMinutes time.Duration) {
	for {
		//TODO should also periodically remove corpses and items in room
		//Repopulate the room every x minutes
		time.Sleep(time.Minute * timeInMinutes)
		room.PopulateMonsters()
	}
}

//PopulateMonsters will add monsters to the room based on the preset types
//of monsters allowed in this room. In other words, if the rooms monsters are dead
//or have never been spawned, this will spawn them; if the monsters are all alive
//this will have no effect.
func (room *Room) PopulateMonsters() {

	for _, monsterName := range room.monsterTemplateNames {
		if _, found := room.MonstersInRoom[monsterName]; found == false {
			room.MonstersInRoom[monsterName] = NewMonster(monsterName, room.ID)
		}
	}
}

func (room *Room) GetDescription() []FormattedString {
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
		output += "\n\t" + itemPtr.GetName()
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
	case NORTH:
		return "North"
	case SOUTH:
		return "South"
	case EAST:
		return "East"
	case WEST:
		return "West"
	case NORTH_WEST:
		return "North-West"
	case NORTH_EAST:
		return "North-East"
	case SOUTH_WEST:
		return "South-West"
	case SOUTH_EAST:
		return "South-East"
	case UP:
		return "Up"
	case DOWN:
		return "Down"
	}

	return ""
}

type ExitXML struct {
	XMLName         xml.Name `xml:"Exit"`
	Direction       string   `xml:"Direction"`
	ConnectedRoomID int      `xml:"RoomID"`
}

type RoomXML struct {
	XMLName     xml.Name     `xml:"Room"`
	ID          int          `xml:"ID"`
	Name        string       `xml:"Name"`
	Description string       `xml:"Description"`
	WorldID     string       `xml:"WorldID"`
	Items       InventoryXML `xml:"Inventory"`
	Monsters    []string     `xml:"Monster"`
	Exits       []ExitXML    `xml:"Exit"`
}

type RoomsXML struct {
	XMLName       xml.Name  `xml:"Rooms"`
	Rooms         []RoomXML `xml:"Room"`
	respawnRoomID int       `xml:"RespawnRoomID"`
}

func LoadRooms() map[int]*Room {
	xmlFile, err := os.Open(serverName + ".xml")
	checkErrorWithMessage(err, true, " In load rooms function.")
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var roomsData RoomsXML
	xml.Unmarshal(XMLdata, &roomsData)

	worldRespawnRoomID = roomsData.respawnRoomID
	rooms := make(map[int]*Room)

	for _, roomData := range roomsData.Rooms {
		getRoom := NewRoom(roomData)
		rooms[getRoom.ID] = getRoom
	}

	for _, room := range rooms {
		room.setRoomLink(rooms)
	}

	fmt.Println("Rooms loaded.")

	return rooms
}
