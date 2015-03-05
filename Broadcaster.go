// Broadcaster
package main

import "net"

type Broadcaster struct {
	//Some list of who is listening and where they are
	// maybe a pointer to the list of online characters with their locations
}


func (b *Broadcaster) registerNewListener(client net.Conn, char Character) {
	
}

func (b *Broadcaster) broadcastToRoom(message String, room int) {
	
}

func (b *Broadcaster) broadcastToCharacter(message String, characterName string) {
	
}