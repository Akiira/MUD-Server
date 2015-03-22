// EventManager
package main

import (
	"github.com/daviddengcn/go-colortext"
	"sync"
)

type Listener interface {
	//set its current eventmanager
	setCurrentEventManager(em *EventManager)
	getEventMessage(msg ServerMessage)
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
	queue_lock    sync.Mutex
}

func (em *EventManager) dummySentMsg(msg string) {

	var newMsg ServerMessage
	tmp := make([]FormattedString, 1, 1)
	tmp[0].Color = ct.Black
	tmp[0].Value = msg

	newMsg.Value = tmp

	for i := 0; i < em.numListener; i++ {
		go em.myListener[i].getEventMessage(newMsg)
	}
}

func (em *EventManager) subscribeListener(newListener Listener) {

	em.queue_lock.Lock()
	em.myListener[em.numListener] = newListener
	em.numListener++
	em.queue_lock.Unlock()
}

func (em *EventManager) unsubscribeListener(prevListener Listener) {

	em.queue_lock.Lock()
	for i := 0; i < em.numListener; i++ {
		if em.myListener[i] == prevListener {
			em.myListener[i] = em.myListener[em.numListener-1]
			em.myListener[em.numListener-1] = nil
			em.numListener--
			break
		}
	}
	em.queue_lock.Unlock()

}

func (em *EventManager) receiveMessage(msg ClientMessage) {
	em.dummySentMsg(msg.Value)
	//	command := msg.Command

	//	if command == "exit" {
	//		break
	//	} else if command == "attack" {
	//		target := msg.Value
	//		output = onlinePlayers["Ragnar"].makeAttack(target) //TODO remove hard coded character names
	//	} else if command == "look" {

	//		target := msg.Value
	//		if target == "room" {
	//			output = worldRoomsG[onlinePlayers["Ragnar"].RoomIN].getFormattedOutput()
	//		} else {
	//			output = monsterTemplatesG[target].getLookDescription()
	//		}
	//	} else if command == "get" {

	//	} else { //assume movement
	//		output = onlinePlayers["Ragnar"].moveCharacter(command)
	//	}

	//	printFormatedOutput(output)
}
