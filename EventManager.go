// EventManager
package main

import (
//"github.com/daviddengcn/go-colortext"
)

type Listener interface {
	//set its current eventmanager
	setCurrentEventManager(em *EventManager)
	getEventMessage(msg ClientMessage)
}

//event manager should only receive event from either monster / player and echo to all that monster / player in the room
// then those player / monster will decide by themselve to get hit or not
// with this concept of oop it should let us handle both eventmanager and play easily
type EventManager struct {
	//a Q of events
	//a timer
	//messageQueue  [100]ClientMessage
	myBroadcaster Broadcaster
	myListener    [100]Listener
	numListener   int
}

type FormattedString struct {
	//Color ct.Color
	Value string
}

func (em *EventManager) dummySentMsg(msg string) {

	//for num := 0; ; num++ {
	var newMsg ClientMessage
	newMsg.Value = msg

	for i := 0; i < em.numListener; i++ {
		go em.myListener[i].getEventMessage(newMsg)
	}

	//	time.Sleep(1 * time.Second)
	//}

}

func (em *EventManager) subscribeListener(newListener Listener) {
	em.myListener[em.numListener] = newListener
	em.numListener++
}

func (em *EventManager) receiveMessage(msg ClientMessage) {

	(*em).dummySentMsg(msg.Value)
}
