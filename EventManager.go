// EventManager
package main

import (
	"github.com/daviddengcn/go-colortext"
)

//event manager should only receive event from either monster / player and echo to all that monster / player in the room
// then those player / monster will decide by themselve to get hit or not
// with this concept of oop it should let us handle both eventmanager and play easily
type EventManager struct {
	//a Q of events
	//a timer
}
type FormattedString struct {
	Color ct.Color
	Value string
}

//this method should be moved to the character class (remote stub)
func executeMove(charName string, direction string) []FormattedString {
	char := onlinePlayers[charName]
	room := worldRoomsG[char.RoomIN]
	dirAsInt := convertDirectionToInt(direction)

	if room.Exits[dirAsInt] >= 0 {
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

//this method should be moved to the character class (remote stub)
func executeStandardAttack(charName string, targetName string) []FormattedString {

	char := onlinePlayers[charName]
	target := worldRoomsG[char.RoomIN].getMonster(targetName)
	output := make([]FormattedString, 2, 2)

	a1 := char.getAttack()
	if a1 >= target.Defense {
		target.HP -= 2
		output[0].Value = "\nYou hit the " + targetName + "!"
	} else {
		output[0].Value = "\nYou missed the " + targetName + "!"
	}

	a2 := target.getAttackRoll()

	if a2 >= char.Defense {
		char.HitPoints -= 1
		output[1].Value = "\nThe " + targetName + " hit you!"
	} else {
		output[1].Value = "\nThe " + targetName + " narrowly misses you!"
	}

	return output
}

//func executeLook(charName string) string {

//}
