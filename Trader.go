package main

type Trader struct {
	isSelected  bool
	isConfirmed bool
	itemMap     map[string]int
	dealerCC    *ClientConnection
}

func Trading(p1 *Trader, p2 *Trader) bool {
	//p1 clientconnection is store in p2 trader
	p1Inv := p2.dealerCC.character.PersonalInvetory
	p1Ready := p1Inv.checkAvailableItemMap(p1.itemMap)
	//p2 clientconnection is store in p1 trader
	p2Inv := p1.dealerCC.character.PersonalInvetory
	p2Ready := p2Inv.checkAvailableItemMap(p2.itemMap)

	if p1Ready && p2Ready {
		deductItemP1 := p1Inv.removeItemMapFromInventory(p1.itemMap)
		deductItemP2 := p2Inv.removeItemMapFromInventory(p2.itemMap)
		p1Inv.addItemListToInventory(deductItemP2)
		p2Inv.addItemListToInventory(deductItemP1)
		return true
	} else {
		return false
	}

}
