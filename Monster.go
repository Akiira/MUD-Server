// Monster
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"io/ioutil"
	"math/rand"
	"os"
	"sync"
	"time"
)

type Monster struct {
	Agent
	Defense     int
	description string
	targets     map[string]*target
	em          EventManager
	weapon      *Weapon
	lootDrops   []loot

	fightingPlayersMutex sync.Mutex
	fightingPlayers      bool
	lastTarget           Agenter
}

//TODO add mutex around targets?

type target struct {
	aggro        int
	attackTarget *Character
}

type loot struct {
	item     Item_I
	dropRate int
}

var monsterTemplatesG map[string]*Monster

//------------------MONSTER CONSTRUCTORS------------------------------

func newMonsterFromXML(monsterData MonsterXML) *Monster {
	m := new(Monster)
	m.Name = monsterData.Name
	m.currentHP = monsterData.HP
	m.Defense = monsterData.Defense
	m.description = monsterData.Description
	m.targets = make(map[string]*target)
	m.weapon = monsterData.EquipedWeapon.toItem().(*Weapon)

	for index, itm := range monsterData.Loot.LootItem.Items {
		drop := loot{item: (itm).(ItemXML_I).toItem(), dropRate: monsterData.Loot.DropRates[index]}
		m.lootDrops = append(m.lootDrops, drop)
	}

	return m
}

func newMonsterFromName(name string, roomID int) *Monster {
	m := new(Monster)
	*m = *monsterTemplatesG[name]
	m.targets = make(map[string]*target)
	m.RoomIN = roomID

	for _, lootItem := range monsterTemplatesG[name].lootDrops {
		m.lootDrops = append(m.lootDrops, lootItem)
	}

	return m
}

//------------------MONSTER COMBAT FUNCTIONS------------------------------

func (m *Monster) fightPlayers() {
	for {
		time.Sleep(2 * time.Second)

		m.fightingPlayersMutex.Lock()
		if m.currentHP <= 0 || len(m.targets) <= 0 {
			m.fightingPlayers = false
			m.fightingPlayersMutex.Unlock()
			break
		}

		//Find target with highest aggro
		var attackTarget *Character
		maxAggro := 0
		for _, targ := range m.targets { //TODO handle targets that moved rooms or logged off.

			if targ.attackTarget.isDead() {
				delete(m.targets, targ.attackTarget.getName())
				continue
			}

			if targ.aggro > maxAggro {
				attackTarget = targ.attackTarget
			}
		}

		if attackTarget != nil {
			event := newEvent(m, "attack", attackTarget.Name)
			eventManager.addEvent(event)
		}
		m.fightingPlayersMutex.Unlock()
	}
}

func (m *Monster) addTarget(targetChar Agenter) {

	_, exist := m.targets[targetChar.getName()]

	if exist {
		m.targets[targetChar.getName()].aggro += 5
	} else {
		fmt.Println("\tAdding player to target list.")
		targ := target{aggro: 5, attackTarget: targetChar.(*Character)}
		m.targets[targetChar.getName()] = &targ

		m.fightingPlayersMutex.Lock()
		if len(m.targets) == 1 && m.fightingPlayers == false {
			m.fightingPlayers = true
			go m.fightPlayers()
		}
		m.fightingPlayersMutex.Unlock()
	}
}

func (m *Monster) makeAttack(target Agenter) []FormattedString {
	fmt.Println("\t\tMonster making attack against player.")
	output := newFormattedStringCollection()
	a1 := m.getAttackRoll()
	if a1 >= target.getDefense() {

		output.addMessage(ct.Red, fmt.Sprintf("The %s hit you for %d damage\n", m.Name, m.getDamage()))
		target.takeDamage(m.getDamage(), 0)

		if target.isDead() {
			output.addMessages(target.respawn())
			delete(m.targets, target.getName())

		}
		target.sendMessage(newServerMessageFS(output.fmtedStrings))
		return output.fmtedStrings
	}
	output.addMessage(ct.Red, fmt.Sprintf("The %s's attack missed you.\n", m.Name))
	target.sendMessage(newServerMessageFS(output.fmtedStrings))
	return output.fmtedStrings
}

func (m *Monster) takeDamage(amount int, typeOfDamge int) {
	m.currentHP -= amount
}

func (m *Monster) respawn() *FmtStrCollection {
	return new(FmtStrCollection)
}

//------------------MONSTER GETTERS------------------------------

func (m *Monster) getAttackRoll() int {
	return (rand.Int() % 20) + m.weapon.attack + m.Strength
}

func (c *Monster) getRoomID() int {
	return c.RoomIN
}
func (m *Monster) getDefense() int {
	return m.Defense
}

func (m *Monster) getName() string {
	return m.Name
}

func (m *Monster) getCorpse() *Item {
	return &Item{name: m.Name + " corpse", description: "A freshly kill " + m.Name + " corpse."}
}

func (m *Monster) getLootAndCorpse() []Item_I {
	return append(m.getLoot(), m.getCorpse())
}

func (m *Monster) getLoot() []Item_I {
	lootItems := make([]Item_I, 0)
	if len(m.lootDrops) > 0 {
		roll := rand.Intn(1000)

		for _, itm := range m.lootDrops {
			if roll <= itm.dropRate {
				lootItems = append(lootItems, itm.item.getCopy())
			}
		}
	}
	return lootItems
}

func (m *Monster) getDamage() int {
	return m.weapon.getDamage() + m.Strength
}

func (m *Monster) getLookDescription() []FormattedString {
	return newFormattedStringSplice2(ct.Yellow, m.description)
}

func (m *Monster) sendMessage(msg ServerMessage) {
	//Do nothing.
}

func (m *Monster) isDead() bool {
	return m.currentHP < 0
}

func (m *Monster) isPlayerAggroed(name string) bool {
	for _, targets := range m.targets {
		if targets.attackTarget.Name == name {
			return true
		}
	}
	return false
}

//------------------Loading Stuff------------------------------
type MonsterXML struct {
	Name          string    `xml:"Name"`
	HP            int       `xml:"HP"`
	Defense       int       `xml:"Defense"`
	Description   string    `xml:"Description"`
	EquipedWeapon WeaponXML `xml:"Weapon"`
	Loot          LootXML   `xml:"Loot"`
}

type MonstersXML struct {
	XMLName  xml.Name     `xml:"Monsters"`
	Monsters []MonsterXML `xml:"Monster"`
}

type LootXML struct {
	XMLName   xml.Name     `xml:"Loot"`
	LootItem  InventoryXML `xml:"Inventory"`
	DropRates []int        `xml:"DropRate"`
}

func loadMonsterData() {
	monsterTemplatesG = make(map[string]*Monster)
	xmlFile, err := os.Open("monsterData.xml")
	checkErrorWithMessage(err, true, " In Load monster data.")
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var monstersData MonstersXML
	err = xml.Unmarshal(XMLdata, &monstersData)
	checkErrorWithMessage(err, true, " In load Monster Data function.")

	for _, element := range monstersData.Monsters {
		monsterTemplatesG[element.Name] = newMonsterFromXML(element)
	}
}
