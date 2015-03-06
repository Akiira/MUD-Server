package main

type Room struct {
	Exits [10]int
	ExitLinksToRooms [10]*Room
	Description string
	Location string	
}


func newRoom(exits [10]int, descr string, loc string) Room{
	return Room{ Exits: exits, Description: descr, Location: loc}
}

func (room *Room) setRoomLink(exit int, roomLink *Room){
	room.ExitLinksToRooms[exit] = roomLink
}

func (room *Room) getRoomLink(exit int) *Room{
	return room.ExitLinksToRooms[exit]
}