// EventManager
package main

import (
	"github.com/daviddengcn/go-colortext"
)

type EventManager struct {
	//a Q of events
	//a timer
}
type FormattedString struct {
	Color ct.Color
	Value string
}

func executeMove(charName string, direction string) []FormattedString {
	char := onlinePlayers[charName]
	room := worldRoomsG[char.RoomIN]
	dirAsInt := convertDirectionToInt(direction)

	if ( room.Exits[dirAsInt] >= 0 ) {
		room.removePCFromRoom(charName)
		room.ExitLinksToRooms[dirAsInt].addPCToRoom(charName)
		char.RoomIN = room.Exits[dirAsInt]
		return room.ExitLinksToRooms[dirAsInt].getFormattedOutput()
	} else {		
		foo := make([]FormattedString, 1, 1)
		foo[0].Color = ct.Black
		foo[0].Value = "No exit in that direction"
		return foo
	}
}

func executeStandardAttack(charName string, target string) {
	
}


//func executeLook(charName string) string {
	
//}