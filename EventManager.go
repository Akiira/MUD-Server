// EventManager
package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"sync"
)

type Listener interface {
	//set its current eventmanager
	setCurrentEventManager(em *EventManager)
	sendMsgToClient(msg ServerMessage)
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
	tmp[0].Color = ct.Blue
	tmp[0].Value = msg

	newMsg.Value = tmp

	fmt.Println("Number of listeners: ", em.numListener)

	for i := 0; i < em.numListener; i++ {
		em.myListener[i].sendMsgToClient(newMsg)
	}
}

func (em *EventManager) subscribeListener(newListener Listener) {

	em.queue_lock.Lock()
	em.myListener[em.numListener] = newListener
	em.numListener++
	fmt.Println("Num: ", em.numListener)
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
}

// The client connection class what should receive the clients message;
//	it can then parse it and determine what event to add here.
//	Then the event manager will call the appropiate room or character functions
func (em *EventManager) addEvent() {

}

func (em *EventManager) executeNonCombatEvent(cc *ClientConnection, event *ClientMessage) {
	var output []FormattedString
	_ = output
	cmd := event.Command
	roomID := cc.character.RoomIN

	switch {
	case cmd == "look":
		output = worldRoomsG[roomID].getRoomDescription()
	case cmd == "get":
		output = worldRoomsG[roomID].getItem(cc.character, event.Value)
	case cmd == "move":
		output = cc.character.moveCharacter(event.Value)
	case cmd == "say":
		output = make([]FormattedString, 1, 1)
		output[0].Color = ct.Black
		output[0].Value = ""
		em.dummySentMsg("Hello everyone in this one room!")
		//TODO error message in exe  non combat move
		//	case true :
		//		output = errorMessage
	}

	fmt.Println("Sending message: ", output)

	if output[0].Value != "" {
		cc.sendMsgToClient(ServerMessage{Value: output})
	}
}
