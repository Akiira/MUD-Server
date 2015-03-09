// Broadcaster
package main

import "net"

type Broadcaster struct {
	//Some list of who is listening and where they are
	// maybe a pointer to the list of online characters with their locations

	//Maybe we can just use an array of size 100
	//removing something in the middle might only just swap it with the last one
	//the order shouldn't cause trouble to our model of broadcasting msg
	//since we are going to forward to all of them anyway
	//-Hand
}

func (b *Broadcaster) registerNewListener(client net.Conn, char Character) {

}

func (b *Broadcaster) broadcastToRoom(message string, room int) {

}

func (b *Broadcaster) broadcastToCharacter(message string, characterName string) {

}
