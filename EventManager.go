// EventManager
package main

import (
	"github.com/daviddengcn/go-colortext"
	"sync"
	"time"
)

type Listener interface {
	setCurrentEventManager(em *EventManager)
	sendMsgToClient(msg ServerMessage)
	getCharactersName() string
}

type EventManager struct {
	listeners  map[string]Listener
	queue_lock sync.Mutex
	eventQue   []Event
	room       *Room
	worldRooms []*Room
}

func newEventManager() *EventManager {
	em := new(EventManager)
	em.listeners = make(map[string]Listener)
	em.eventQue = make([]Event, 0, 10)
	em.worldRooms = loadRooms()

	go em.waitForTick()

	return em
}

func newEventManagerForRoom(room *Room) *EventManager {
	em := new(EventManager)
	em.listeners = make(map[string]Listener)
	em.eventQue = make([]Event, 0, 10)
	em.room = room
	em.worldRooms = loadRooms()
	return em
}

func (em *EventManager) sendMessageToRoom(msg string) {
	var newMsg ServerMessage
	tmp := make([]FormattedString, 1, 1)

	tmp[0] = FormattedString{Color: ct.Blue, Value: msg}
	newMsg.Value = tmp

	for _, listener := range em.listeners {
		listener.sendMsgToClient(newMsg)
	}
}

func (em *EventManager) subscribeListener(newListener Listener) {

	em.queue_lock.Lock()
	em.listeners[newListener.getCharactersName()] = newListener
	em.queue_lock.Unlock()
}

func (em *EventManager) unsubscribeListener(prevListener Listener) {

	em.queue_lock.Lock()
	delete(em.listeners, prevListener.getCharactersName())
	em.queue_lock.Unlock()

}

func (em *EventManager) receiveMessage(msg ClientMessage) {
	em.sendMessageToRoom(msg.Value)
}

func (em *EventManager) addEvent(event Event) {
	em.queue_lock.Lock()
	em.eventQue = append(em.eventQue, event)
	em.queue_lock.Unlock()
}

func (em *EventManager) waitForTick() {
	for {
		time.Sleep(time.Second * 2)
		go em.executeCombatRound()
	}
}

func (em *EventManager) executeCombatRound() {
	var output []FormattedString

	for _, event := range em.eventQue {
		//TODO sort events by initiative stat before executing them
		action := event.action
		switch {
		case action == "attack":
			output = event.agent.makeAttack(event.valueOrTarget)
		}

		if event.client != nil {
			var servMsg ServerMessage
			servMsg.Value = output
			event.client.sendMsgToClient(servMsg)
		}
	}

	em.eventQue = em.eventQue[0:0]
}

func (em *EventManager) executeNonCombatEvent(cc *ClientConnection, event *ClientMessage) {
	var output []FormattedString
	eventRoom := em.worldRooms[cc.character.RoomIN]
	cmd := event.Command
	switch {
	case cmd == "look":
		output = eventRoom.getRoomDescription()
	case cmd == "get":
		output = eventRoom.getItem(cc.character, event.Value)
	case cmd == "move":
		output = cc.character.moveCharacter(event.Value)
	case cmd == "say":
		str := cc.character.Name + " says \"" + event.Value + "\""
		em.sendMessageToRoom(str)
	}

	if len(output) > 0 {
		cc.sendMsgToClient(ServerMessage{Value: output})
	}
}
