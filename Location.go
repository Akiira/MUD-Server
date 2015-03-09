package main

//TODO So a room can be uniquely identified by its roomID
//		which is an int so this class may not be needed.
//		Its still not clear how we are going to distinguish
//		rooms on different servers though. It could be we use
//		a struct like this with the fields Server and roomID.

type Location struct {
	x int
	y int
}
