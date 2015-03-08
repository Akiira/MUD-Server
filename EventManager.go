// EventManager
package main

import (

)

type EventManager struct {
	//a Q of events
	//a timer
}


func executeMove(charName string, direction string) string {
	char := onlinePlayers[charName]
	room := worldRoomsG[char.RoomIN]
	dirAsInt := convertDirectionToInt(direction)

	if ( room.Exits[dirAsInt] >= 0 ) {
		room.removePCFromRoom(charName)
		room.ExitLinksToRooms[dirAsInt].addPCToRoom(charName)
		char.RoomIN = room.Exits[dirAsInt]
		return room.ExitLinksToRooms[dirAsInt].getFormattedOutput()
	} else {
		return "No exit in that direction"
	}
}