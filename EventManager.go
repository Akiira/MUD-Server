// EventManager
package main

import (
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"sync"
	"time"
)

type EventManager struct {
	trade_lock sync.Mutex
	auctn_mutx sync.Mutex
	queue_lock sync.Mutex
	eventQue   []Event
	worldRooms map[int]*Room

	auction *Auction
	traders map[string]bool
}

func newEventManager() *EventManager {
	em := new(EventManager)
	em.eventQue = make([]Event, 0, 10)
	em.worldRooms = loadRooms()
	em.traders = make(map[string]bool)
	go em.StartCombatRounds()

	return em
}

func (em *EventManager) sendMessageToWorld(msg ServerMessage) {

	for _, room := range em.worldRooms {
		for _, char := range room.CharactersInRoom {
			char.SendMessage(msg)
		}
	}
}

func (em *EventManager) sendMessageToRoom(roomID int, msg ServerMessage) {
	room := em.worldRooms[roomID]

	for _, char := range room.CharactersInRoom {
		char.SendMessage(msg)
	}
}

func (em *EventManager) AddEvent(event Event) {
	em.queue_lock.Lock()
	em.eventQue = append(em.eventQue, event)
	em.queue_lock.Unlock()
}

func (em *EventManager) StartCombatRounds() {
	for {
		time.Sleep(time.Second * 6)
		go em.executeCombatRound()
	}
}

func (em *EventManager) executeCombatRound() {
	em.queue_lock.Lock()
	var output []FormattedString
	alreadyActed := make(map[string]bool)

	for _, event := range em.eventQue {
		//TODO sort events by initiative stat before executing them
		action := event.action
		agent := event.agent

		target := em.worldRooms[agent.GetRoomID()].getAgentInRoom(event.target)

		if _, found := alreadyActed[agent.GetName()]; !found {
			alreadyActed[agent.GetName()] = true

			switch {
			case action == "attack":
				output = agent.makeAttack(target)
			}

			agent.SendMessage(newServerMessageFS(output))
		}
	}

	em.eventQue = em.eventQue[:0]
	em.queue_lock.Unlock()
}

func (em *EventManager) executeNonCombatEvent(cc *ClientConnection, event *ClientMessage) {
	var output []FormattedString
	var msgType int = GAMEPLAY

	switch event.getCommand() {
	case "auction": //This is to start an auction
		em.auctn_mutx.Lock()
		defer em.auctn_mutx.Unlock()

		if em.auction == nil {
			item, found := cc.character.GetItem(event.Value)
			if found {
				em.auction = newAuction(item)
				go em.runAuction()
				output = newFormattedStringSplice("Your auction was succesfully started.")
			} else {
				output = newFormattedStringSplice("Could not start the auction because you do not have that item.")
			}
		}
	case "bid":
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

	case "unwield":
		output = cc.character.UnWieldWeapon()
	case "wield":
		output = cc.character.WieldWeapon(event.Value)
	case "unequip":
		output = cc.character.UnEquipArmourByName(event.Value)
	case "equip":
		output = cc.character.EquipArmorByName(event.Value)
	case "inv":
		output = cc.character.PersonalInvetory.getInventoryDescription()
	case "save", "exit":
		sendCharactersXML(cc.getCharacter().toXML())
		output = newFormattedStringSplice("Character succesfully saved.\n")
	case "stats":
		output = cc.getCharacter().GetStats()
	case "look":
		output = em.GetRoom(cc.getCharactersRoomID()).GetDescription()
	case "get":
		output = em.GetRoom(cc.getCharactersRoomID()).GiveItemToPlayer(cc.character, event.Value)
	case "move":
		src := em.worldRooms[cc.getCharactersRoomID()]
		dest := src.getConnectedRoom(convertDirectionToInt(event.Value))

		msgType, output = cc.character.moveCharacter(src, dest)
	case "yell":
		em.sendMessageToWorld(newServerMessageFS(newFormattedStringSplice2(ct.Blue, cc.character.Name+" says \""+event.Value+"\"")))
	case "say":
		formattedOutput := newFormattedStringSplice2(ct.Blue, cc.character.Name+" says \""+event.Value+"\"")
		em.sendMessageToRoom(cc.character.RoomIN, ServerMessage{Value: formattedOutput})
	case "trade":
		go em.ExecuteTradeEvent(cc.getCharacter(), event)
	}

	if len(output) > 0 {
		cc.Write(newServerMessageTypeFS(msgType, output))
	}
}

func (em *EventManager) GetRoom(roomID int) *Room {
	room, found := em.worldRooms[roomID]
	if found {
		return room
	} else {
		return nil
	}
}

func (em *EventManager) AddPlayerToRoom(char *Character, roomID int) {
	if room := em.GetRoom(roomID); room != nil {
		room.AddPlayer(char)
		if room.isLocal() {
			char.SendMessage(room.GetDescription())
		}
	}
}

func (em *EventManager) RemovePlayerFromRoom(charName string, roomID int) {
	if room := em.GetRoom(roomID); room != nil {
		room.RemovePlayer(charName)
	}
}

//-----------------------------TRADING EVENT FUNCTIONS------------------------//
func (em *EventManager) IsTrading(charName string) bool {
	val, found := em.traders[charName]

	return found && val
}

func (em *EventManager) SetPlayerToTrading(name string) {
	em.trade_lock.Lock()
	em.traders[name] = true
	em.trade_lock.Unlock()
}

func (em *EventManager) SetPlayerToNotTrading(name string) {
	em.trade_lock.Lock()
	em.traders[name] = false
	em.trade_lock.Unlock()
}

//TODO check players are not in combat before starting trade
func (em *EventManager) ExecuteTradeEvent(trader *Character, event *ClientMessage) {

	em.SetPlayerToTrading(trader.GetName())
	defer em.SetPlayerToNotTrading(trader.GetName())

	//If other player rejects trade we return
	var tradee *Character
	if tradee = em.AskOtherPlayerToTrade(trader, event.Value); tradee == nil {
		return
	}

	em.SetPlayerToTrading(tradee.GetName()) //TODO should we double check they are still in the room?
	defer em.SetPlayerToNotTrading(tradee.GetName())

	trader.SendMessage("The trade is opened, what items would you like to trade? Type 'add' [item name] to add an item to the trade or 'add' [quantity] [item name].\n")
	tradee.SendMessage("The trade is opened, what items would you like to trade? Type 'add' [item name] to add an item to the trade or 'add' [quantity] [item name].\n")

	tradersItems := newInventory()
	tradeesItems := newInventory()
	accepted := false

	defer func() { //This ensures that no matter how the function returns, the correct person gets their items
		if accepted {
			trader.AddInventoryToInventory(tradeesItems)
			tradee.AddInventoryToInventory(tradersItems)
			trader.SendMessage("Both parties accepted the trade, you receieved your items.\n")
			tradee.SendMessage("Both parties accepted the trade, you receieved your items.\n")
		} else {
			tradee.AddInventoryToInventory(tradeesItems)
			trader.AddInventoryToInventory(tradersItems)
			trader.SendMessage("The trade was ended, any entered items have been returned to your inventory.\n")
			tradee.SendMessage("The trade was ended, any entered items have been returned to your inventory.\n")
		}
	}()

	em.GetTradeItemsFromPlayers(trader, tradee, tradersItems, tradeesItems)

	em.SendFinalTradeTerms(trader, tradee, tradersItems, tradeesItems)

	accepted = em.AskFinalTradePrompt(trader, tradee)
}

func (em *EventManager) AskOtherPlayerToTrade(trader *Character, tradeeName string) *Character {
	// ask other player if they want to trade
	room := em.GetRoom(trader.GetRoomID())
	tradee, found := room.GetPlayer(tradeeName)

	if !found {
		fmt.Println("\tDid not find other player.")
		trader.SendMessage("That player is not in this room.\n")
		return nil
	}

	tradee.SendMessage("Would you like to trade with " + trader.GetName() + "? Type accept or reject.\n")
	response := tradee.GetTradeResponse()

	if response != "accept" {
		fmt.Println("\tOther player did not responde with accept.")
		trader.SendMessage("\nThe other player did not accept your trade.\n")
		return nil
	}

	return tradee
}

func (em *EventManager) GetTradeItemsFromPlayers(trader, tradee *Character, trderInv, trdeeInv *Inventory) {
	var wg sync.WaitGroup
	wg.Add(2)

	// get msg of items from players
	go trader.GetItemsToTrade(trderInv, &wg)
	go tradee.GetItemsToTrade(trdeeInv, &wg)

	wg.Wait()
}

//TODO combine msg to same person into one server msg to make it appear cleaner on clients side
func (em *EventManager) SendFinalTradeTerms(trader, tradee *Character, trderInv, trdeeInv *Inventory) {
	//send final terms out to players
	trader.SendMessage("Here are the final terms of the trade, you will receive:\n")
	trader.SendMessage(trdeeInv.getInventoryDescription())
	trader.SendMessage("And are trading away:\n")
	trader.SendMessage(trderInv.getInventoryDescription())

	tradee.SendMessage("Here are the final terms of the trade, you will receive:\n")
	tradee.SendMessage(trderInv.getInventoryDescription())
	tradee.SendMessage("And are trading away:\n")
	tradee.SendMessage(trdeeInv.getInventoryDescription())

	trader.SendMessage("If you accept the terms of this trade then type 'accept' else type 'reject'.\n")
	tradee.SendMessage("If you accept the terms of this trade then type 'accept' else type 'reject'.\n")
}

func (em *EventManager) AskFinalTradePrompt(trader, tradee *Character) bool {
	// players accept or reject trade.
	response := trader.GetTradeResponse()
	if response != "accept" {
		fmt.Println("\tTrader did not accept.")
		trader.SendMessage("You rejected the trade.\n")
		tradee.SendMessage("The other player did not accept the final terms of the trade.\n")
		return false
	}

	response = tradee.GetTradeResponse()
	if response != "accept" {
		fmt.Println("\tTradee did not accept.")
		trader.SendMessage("You rejected the trade.\n")
		tradee.SendMessage("The other player did not accept the final terms of the trade.\n")
		return false
	}

	fmt.Println("\tTrade accepted by both parties.")
	return true
}

//-----------------------------AUCTION EVENT FUNCTIONS------------------------//

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
