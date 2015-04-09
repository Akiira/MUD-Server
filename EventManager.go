// EventManager
package main

import (
	"github.com/daviddengcn/go-colortext"
	"sync"
	"time"
)

type EventManager struct {
	queue_lock sync.Mutex
	eventQue   []Event
	worldRooms map[int]*Room
}

func newEventManager(worldName string) *EventManager {
	em := new(EventManager)
	em.eventQue = make([]Event, 0, 10)
	em.worldRooms = loadRooms(worldName)

	go em.waitForTick()

	return em
}

func (em *EventManager) sendMessageToWorld(msg ServerMessage) {

	for _, room := range em.worldRooms {
		for _, char := range room.CharactersInRoom {
			char.sendMessage(msg)
		}
	}
}

func (em *EventManager) sendMessageToRoom(roomID int, msg ServerMessage) {
	room := em.worldRooms[roomID]

	for _, char := range room.CharactersInRoom {
		char.sendMessage(msg)
	}
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
	alreadyActed := make(map[string]bool)
	for _, event := range em.eventQue {
		//TODO sort events by initiative stat before executing them
		action := event.action
		agent := event.agent

		target := em.worldRooms[agent.getRoomID()].getAgentInRoom(event.target)

		if _, found := alreadyActed[agent.getName()]; !found {
			alreadyActed[agent.getName()] = true

			switch {
			case action == "attack":
				output = agent.makeAttack(target)
			}

			agent.sendMessage(newServerMessageFS(output))
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
	case cmd == "bid":
		//TODO
		// check if there is a running auction
		// if so then: auction.bidOnItem(event.getBid(), cc, time.Now()) //Or get time from Timestamp?
	case cmd == "inv":
		output = cc.character.PersonalInvetory.getInventoryDescription()
	case cmd == "save" || cmd == "exit":
		sendCharactersXML(cc.character.toXML())
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
		cc.sendMsgToClient(newServerMessageTypeFS(msgType, output))
	}
}
