// EventManager
package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"sync"
	"time"
)

type EventManager struct {
	auctn_mutx sync.Mutex
	queue_lock sync.Mutex
	eventQue   []Event
	worldRooms map[int]*Room
	auction    *Auction
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
	cmd := event.getCommand()
	var msgType int
	msgType = GAMEPLAY
	switch {
	case cmd == "auction": //This is to start an auction
		em.auctn_mutx.Lock()
		defer em.auctn_mutx.Unlock()

		if em.auction == nil {
			item, found := cc.character.getItemFromInv(event.Value)
			if found {
				em.auction = newAuction(item)
				go em.runAuction()
				output = newFormattedStringSplice("Your auction was succesfully started.")
			} else {
				output = newFormattedStringSplice("Could not start the auction because you do not have that item.")
			}
		}
	case cmd == "bid":
		fmt.Println("Executing bid command")
		em.auctn_mutx.Lock()
		defer em.auctn_mutx.Unlock()

		if em.auction != nil {
			fmt.Println("Bidding on item")
			output = em.auction.bidOnItem(event.getBid(), cc, time.Now())
			fmt.Println("Done Bidding on item")
		} else {
			output = newFormattedStringSplice("There are not currently any auctions happening.\n")
		}
		fmt.Println("Done Executing bid command")

	case cmd == "unwield":
		output = cc.character.UnWieldWeapon()
	case cmd == "wield":
		output = cc.character.WieldWeaponByName(event.Value)
	case cmd == "unequip":
		output = cc.character.UnEquipArmourByName(event.Value)
	case cmd == "equip":
		output = cc.character.EquipArmorByName(event.Value)
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
		src := em.worldRooms[cc.getCharactersRoomID()]
		dest := src.getConnectedRoom(convertDirectionToInt(event.Value))

		msgType, output = cc.character.moveCharacter(src, dest)
	case cmd == "yell":
		em.sendMessageToWorld(newServerMessageFS(newFormattedStringSplice2(ct.Blue, cc.character.Name+" says \""+event.Value+"\"")))
	case cmd == "say":
		formattedOutput := newFormattedStringSplice2(ct.Blue, cc.character.Name+" says \""+event.Value+"\"")
		em.sendMessageToRoom(cc.character.RoomIN, ServerMessage{Value: formattedOutput})
	}

	if len(output) > 0 {
		cc.sendMsgToClient(newServerMessageTypeFS(msgType, output))
	}
}

func (em *EventManager) sendPeriodicAuctionInfo() {
	for {
		time.Sleep(time.Second * 5)
		if em.auction.timeTillOver() > time.Second*5 {
			em.sendMessageToWorld(em.auction.getAuctionInfo())
		} else {
			break
		}
	}
}

func (em *EventManager) runAuction() {

	go em.sendPeriodicAuctionInfo()

	for {

		time.Sleep(time.Second * 3)

		em.auctn_mutx.Lock()
		if em.auction.timeTillOver() < -time.Second*2 {
			break
		}
		em.auctn_mutx.Unlock()
	}

	winner := em.auction.determineWinner()

	if winner != nil {
		em.auction.awardItemToWinner(winner)
	}

	em.auction = nil
	em.auctn_mutx.Unlock()
}

func (em *EventManager) isAuctionRunning() bool {
	return em.auction != nil
}
