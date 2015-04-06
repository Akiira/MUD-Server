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
	worldRooms []*Room
}

func newEventManager(worldName string) *EventManager {
	em := new(EventManager)
	em.listeners = make(map[string]Listener)
	em.eventQue = make([]Event, 0, 10)
	em.worldRooms = loadRooms(worldName)

	go em.waitForTick()

	return em
}

func (em *EventManager) sendMessageToRoom(roomID int, msg ServerMessage) {
	room := em.worldRooms[roomID]

	for _, client := range room.CharactersInRoom {
		client.myClientConn.sendMsgToClient(msg)
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
			output = event.agent.makeAttack(event.target)
		}

		event.agent.getClientConnection().sendMsgToClient(newServerMessage(output))

		if event.target.isDead() {

		}
	}

	em.eventQue = em.eventQue[0:0]
}

func (em *EventManager) executeNonCombatEvent(cc *ClientConnection, event *ClientMessage) {
	var output []FormattedString
	eventRoom := em.worldRooms[cc.character.RoomIN]
	cmd := event.Command
	var msgType int
	msgType = GAMEPLAY
	switch {
	case cmd == "inv":
		output = cc.character.PersonalInvetory.getInventoryDescription()
	case cmd == "save" || cmd == "exit":
		saveCharacterToFile(cc.getCharacter())
		sendCharactersFile(cc.getCharactersName())
		output = newFormattedStringSplice("Character succesfully saved.\n")
	case cmd == "stats":
		output = cc.getCharacter().getStatsPage()
	case cmd == "look":
		output = eventRoom.getRoomDescription()
	case cmd == "get":
		output = eventRoom.getItem(cc.character, event.Value)
	case cmd == "move":
		msgType, output = cc.character.moveCharacter(event.Value)
	case cmd == "say":
		formattedOutput := newFormattedStringSplice2(ct.Blue, cc.character.Name+" says \""+event.Value+"\"")
		em.sendMessageToRoom(cc.character.RoomIN, ServerMessage{Value: formattedOutput})
	}

	if len(output) > 0 {
		cc.sendMsgToClient(newServerMessageWithType(msgType, output))
	}
}
