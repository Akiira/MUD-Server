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

func NewEventManager() *EventManager {
	em := new(EventManager)
	em.eventQue = make([]Event, 0, 10)
	em.worldRooms = loadRooms()
	em.traders = make(map[string]bool)
	go em.StartCombatRounds()

	return em
}

func (em *EventManager) SendMessageToWorld(msg ServerMessage) {
	for _, room := range em.worldRooms {
		for _, char := range room.CharactersInRoom {
			char.SendMessage(msg)
		}
	}
}

func (em *EventManager) SendMessageToRoom(roomID int, msg ServerMessage) {
	room := em.GetRoom(roomID)

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
		go em.ExecuteCombatRound()
	}
}

func (em *EventManager) ExecuteCombatRound() {
	em.queue_lock.Lock()
	var output []FormattedString
	alreadyActed := make(map[string]bool)

	for _, event := range em.eventQue {
		action, agent := event.action, event.agent

		target, _ := em.worldRooms[agent.GetRoomID()].getAgentInRoom(event.target)

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

func (em *EventManager) ExecuteNonCombatEvent(cc *ClientConnection, event *ClientMessage) {
	var output []FormattedString
	var msgType int = GAMEPLAY

	switch event.getCommand() {
	case "auction": //This is to start an auction
		output = em.StartAuction(cc.getCharacter(), event.Value)
	case "bid":
		output = em.BidOnAuction(cc, event.getBid())
	case "unwield":
		output = cc.character.UnWieldWeapon()
	case "wield":
		output = cc.character.WieldWeapon(event.Value)
	case "unequip", "remove":
		output = cc.character.UnEquipArmourByName(event.Value)
	case "equip", "wear":
		output = cc.character.EquipArmorByName(event.Value)
	case "equipment", "eq":
		output = cc.character.GetEquipment()
	case "inventory", "inv":
		output = cc.character.PersonalInvetory.getInventoryDescription()
	case "save", "exit":
		SendCharactersXML(cc.getCharacter().toXML())
		output = newFormattedStringSplice("Character succesfully saved.\n")
	case "stats":
		output = cc.getCharacter().GetStats()
	case "look":
		room := em.GetRoom(cc)
		if event.Value == "room" || event.Value == "" {
			output = room.GetDescription()
		} else if _, found := room.GetItem(event.Value); found { //TODO change return type of getDescr..
			//output = item.getDescription()
		} else if _, found := cc.getCharacter().GetItem(event.Value); found {
			//output = item.getDescription()
		} else if _, found := room.getAgentInRoom(event.Value); found {
			//output = agent.getDescription()
		} else {
			output = newFormattedStringSplice("That item or creature could not be found anywhere.\n")
		}

	case "get":
		output = em.GetRoom(cc.getCharactersRoomID()).GiveItemToPlayer(cc.character, event.Value)
	case "drop":
		if em.IsTrading(cc.getCharactersName()) {
			output = newFormattedStringSplice2(ct.Red, "\nYou can not drop items while you are trading.\n")
		} else if item, found := cc.getCharacter().GetAndRemoveItem(event.Value); found {
			em.GetRoom(cc.getCharactersRoomID()).AddItem(item)
			output = newFormattedStringSplice("You dropped the " + item.getName() + " on the ground.\n")
		} else {
			output = newFormattedStringSplice("You dropped the " + item.getName() + " on the ground.\n")
		}
	case "move":
		if em.IsTrading(cc.getCharactersName()) {
			output = newFormattedStringSplice2(ct.Red, "\nYou can not move rooms while you are trading.\n")
		} else if em.IsInCombat(cc.getCharactersName()) {
			output = newFormattedStringSplice2(ct.Red, "\nYou can not move rooms while you are in combat. If your need to get away try 'flee'.\n")
		} else {
			src := em.worldRooms[cc.getCharactersRoomID()]
			dest := src.getConnectedRoom(convertDirectionToInt(event.Value))

			msgType, output = cc.character.moveCharacter(src, dest)
		}
	case "flee":
		if em.IsTrading(cc.getCharactersName()) {
			output = newFormattedStringSplice2(ct.Red, "\nYou can not flee from a room while you are trading.\n")
		} else {
			src := em.worldRooms[cc.getCharactersRoomID()]
			dest := src.getConnectedRoom(convertDirectionToInt(event.Value))

			src.UnAggroPlayer(cc.getCharactersName())

			msgType, output = cc.character.moveCharacter(src, dest)
			output = append(cc.getCharacter().ApplyFleePenalty(), output...)
		}
	case "yell":
		em.SendMessageToWorld(newServerMessageFS(newFormattedStringSplice2(ct.Blue, cc.character.Name+" says \""+event.Value+"\"")))
	case "say":
		formattedOutput := newFormattedStringSplice2(ct.Blue, cc.character.Name+" says \""+event.Value+"\"")
		em.SendMessageToRoom(cc.character.RoomIN, ServerMessage{Value: formattedOutput})
	case "trade":
		go em.ExecuteTradeEvent(cc.getCharacter(), event)
	case "help":
		output = newFormattedStringSplice("\nYou can use the following commands\nattack\ninv\nlook\nyell\nsay\ntrade\nbid\nwield\nunwield\nequip\nget\nmove\nauction\n")
	}

	if len(output) > 0 {
		cc.Write(newServerMessageTypeFS(msgType, output))
	}
}

func (em *EventManager) SaveAllCharacters() {
	for _, room := range em.worldRooms {
		if room.IsLocal() {
			for _, player := range room.CharactersInRoom {
				SendCharactersXML(player.toXML())
			}
		}
	}
}

func (em *EventManager) GetRoom(input interface{}) *Room {
	var roomID int

	switch input := input.(type) {
	default:
		fmt.Printf("Unexpected type %T in EventManager.GetRoom", input)
		fmt.Println(" with value: ", input)
		return nil
	case int:
		roomID = input
	case *ClientConnection:
		roomID = input.getCharactersRoomID()
	case *Character:
		roomID = input.GetRoomID()
	}

	if room, found := em.worldRooms[roomID]; found {
		return room
	} else {
		return nil
	}
}

func (em *EventManager) AddPlayerToRoom(char *Character) {
	if room := em.GetRoom(char); room != nil {
		fmt.Println("\t\tAdding Player to room.")
		room.AddPlayer(char)
		if room.IsLocal() {
			fmt.Println("\t\tSending message.")
			char.SendMessage(room.GetDescription())
		}
	}
}

func (em *EventManager) RemovePlayerFromRoom(charName string) {
	if room := em.GetRoom(charName); room != nil {
		room.RemovePlayer(charName)
	}
}

func (em *EventManager) IsInCombat(charName string) bool {
	if room := em.GetRoom(charName); room != nil {
		return room.IsAggroed(charName)
	}

	return false
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

func (em *EventManager) ExecuteTradeEvent(trader *Character, event *ClientMessage) {
	em.SetPlayerToTrading(trader.GetName())
	defer em.SetPlayerToNotTrading(trader.GetName())

	//If other player rejects trade we return
	var tradee *Character
	if tradee = em.AskOtherPlayerToTrade(trader, event.Value); tradee == nil {
		return
	}

	em.SetPlayerToTrading(tradee.GetName())
	defer em.SetPlayerToNotTrading(tradee.GetName())

	trader.SendMessage("The trade is opened, what items would you like to trade? Type 'add' [item name] to add an item to the trade or 'add' [quantity] [item name].\n")
	tradee.SendMessage("The trade is opened, what items would you like to trade? Type 'add' [item name] to add an item to the trade or 'add' [quantity] [item name].\n")

	tradersItems, tradeesItems := newInventory(), newInventory()
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

	tradee, found := em.GetRoom(trader).GetPlayer(tradeeName)

	if !found {
		fmt.Println("\tDid not find other player.")
		trader.SendMessage("That player is not in this room.\n")
		return nil
	}

	tradee.SendMessage("Would you like to trade with " + trader.GetName() + "? Type accept or reject.\n")

	if tradee.GetTradeResponse() != "accept" {
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

func (em *EventManager) SendFinalTradeTerms(trader, tradee *Character, trderInv, trdeeInv *Inventory) {

	tradersMsg := newFormattedStringCollection()
	tradersMsg.addMessage2("Here are the final terms of the trade, you will receive:\n")
	tradersMsg.addMessages2(trdeeInv.getInventoryDescription())
	tradersMsg.addMessage2("And are trading away:\n")
	tradersMsg.addMessages2(trderInv.getInventoryDescription())
	tradersMsg.addMessage2("If you accept the terms of this trade then type 'accept' else type 'reject'.\n")

	tradeesMsg := newFormattedStringCollection()
	tradeesMsg.addMessage2("Here are the final terms of the trade, you will receive:\n")
	tradeesMsg.addMessages2(trderInv.getInventoryDescription())
	tradeesMsg.addMessage2("And are trading away:\n")
	tradeesMsg.addMessages2(trdeeInv.getInventoryDescription())
	tradeesMsg.addMessage2("If you accept the terms of this trade then type 'accept' else type 'reject'.\n")

	tradee.SendMessage(tradeesMsg.fmtedStrings)
	trader.SendMessage(tradersMsg.fmtedStrings)
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

func (em *EventManager) StartAuction(char *Character, itemName string) (output []FormattedString) {
	em.auctn_mutx.Lock()

	if !em.IsAuctionRunning() {
		if item, found := char.GetAndRemoveItem(itemName); found {
			em.auction = newAuction(item)
			go em.RunAuction()
			output = newFormattedStringSplice("Your auction was succesfully started.")
		} else {
			output = newFormattedStringSplice("Could not start the auction because you do not have that item.")
		}
	} else {
		output = newFormattedStringSplice("Could not start the auction because there is already an auction running.")
	}

	em.auctn_mutx.Unlock() //Do not defer me
	return output
}

func (em *EventManager) BidOnAuction(cc *ClientConnection, bid int) []FormattedString {
	em.auctn_mutx.Lock()
	defer em.auctn_mutx.Unlock()

	if em.IsAuctionRunning() {
		return em.auction.bidOnItem(bid, cc, time.Now())
	} else {
		return newFormattedStringSplice("There are not currently any auctions happening.\n")
	}
}

func (em *EventManager) SendPeriodicAuctionInfo() {
	for {
		time.Sleep(time.Second * 5)
		if em.auction.timeTillOver() > time.Second*5 {
			em.SendMessageToWorld(em.auction.getAuctionInfo())
		} else {
			break
		}
	}
}

func (em *EventManager) RunAuction() {

	go em.SendPeriodicAuctionInfo()

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

func (em *EventManager) IsAuctionRunning() bool {
	return em.auction != nil
}
