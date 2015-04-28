// Monster
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"sync"
	"time"
)

type Monster struct {
	Agent
	Defense     int
	description string

	em        EventManager
	weapon    *Weapon
	lootDrops []loot

	targets       map[string]*target
	targets_mutex sync.Mutex

	fightingPlayers  bool
	fightPlyrs_Mutex sync.Mutex
}

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

func NewMonsterFromXML(monsterData MonsterXML) *Monster {
	m := new(Monster)
	m.Name = monsterData.Name
	m.currentHP = monsterData.HP
	m.Defense = monsterData.Defense
	m.description = monsterData.Description
	m.targets = make(map[string]*target)
	m.weapon = monsterData.EquipedWeapon.ToItem().(*Weapon)

	for index, itm := range monsterData.Loot.LootItem.Items {
		drop := loot{item: (itm).(ItemXML_I).ToItem(), dropRate: monsterData.Loot.DropRates[index]}
		m.lootDrops = append(m.lootDrops, drop)
	}

	return m
}

func NewMonster(name string, roomID int) *Monster {
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

		m.fightPlyrs_Mutex.Lock()
		if m.currentHP <= 0 || len(m.targets) <= 0 {
			m.fightingPlayers = false
			m.fightPlyrs_Mutex.Unlock()
			break
		}

		//Find target with highest aggro
		var attackTarget *Character
		maxAggro := 0
		for _, targ := range m.targets {

			if targ.attackTarget.IsDead() {
				delete(m.targets, targ.attackTarget.GetName())
				continue
			}

			if targ.aggro > maxAggro {
				attackTarget = targ.attackTarget
			}
		}

		if attackTarget != nil {
			event := NewEvent(m, "attack", attackTarget.Name)
			eventManager.AddEvent(event)
		}
		m.fightPlyrs_Mutex.Unlock()
	}
}

func (m *Monster) AddTarget(targetChar Agenter) {
	m.targets_mutex.Lock()
	defer m.targets_mutex.Unlock()

	if _, exist := m.targets[targetChar.GetName()]; exist {
		m.targets[targetChar.GetName()].aggro += 5
	} else {
		targ := target{aggro: 5, attackTarget: targetChar.(*Character)}
		m.targets[targetChar.GetName()] = &targ

		m.fightPlyrs_Mutex.Lock()
		if len(m.targets) == 1 && m.fightingPlayers == false {
			m.fightingPlayers = true
			go m.fightPlayers()
		}
		m.fightPlyrs_Mutex.Unlock()
	}
}

func (m *Monster) RemoveTarget(name string) {
	m.targets_mutex.Lock()
	defer m.targets_mutex.Unlock()

	if _, found := m.targets[name]; found {
		delete(m.targets, name)
	}
}

func (m *Monster) Attack(target Agenter) []FormattedString {
	m.targets_mutex.Lock()
	defer m.targets_mutex.Unlock()

	//If the player died or fled then the monster cant attack
	if !m.IsAttackingPlayer(target.GetName()) || target.GetRoomID() != m.RoomIN {
		return nil
	}

	output := newFormattedStringCollection()

	if m.GetAttackRoll() >= target.GetDefense() {
		dmg := m.GetDamage()
		output.addMessage(ct.Red, fmt.Sprintf("The %s hit you for %d damage\n", m.Name, dmg))
		target.TakeDamage(dmg, 0)

		if target.IsDead() {
			output.addMessages2(target.Respawn())
			delete(m.targets, target.GetName())
		}

		target.SendMessage(newServerMessageFS(output.fmtedStrings))
		return output.fmtedStrings
	}

	output.addMessage(ct.Red, fmt.Sprintf("The %s's attack missed you.\n", m.Name))
	target.SendMessage(newServerMessageFS(output.fmtedStrings))

	return output.fmtedStrings
}

func (m *Monster) TakeDamage(amount int, typeOfDamge int) {
	m.currentHP -= amount
}

func (m *Monster) Respawn() []FormattedString {
	return nil
}

//------------------MONSTER GETTERS------------------------------

func (m *Monster) GetAttackRoll() int {
	return (rand.Int() % 20) + m.weapon.attack + m.Strength
}

func (c *Monster) GetRoomID() int {
	return c.RoomIN
}

func (m *Monster) GetDefense() int {
	return m.Defense
}

func (m *Monster) GetName() string {
	return m.Name
}

func (m *Monster) GetCorpse() *Item {
	return &Item{name: m.Name + " corpse", description: "A freshly kill " + m.Name + " corpse."}
}

func (m *Monster) GetLootAndCorpse() []Item_I {
	return append(m.GetLoot(), m.GetCorpse())
}

func (m *Monster) GetLoot() []Item_I {
	lootItems := make([]Item_I, 0)
	if len(m.lootDrops) > 0 {
		roll := rand.Intn(1000)

		for _, itm := range m.lootDrops {
			if roll <= itm.dropRate {
				lootItems = append(lootItems, itm.item.GetCopy())
			}
		}
	}
	return lootItems
}

func (m *Monster) GetDamage() int {
	return m.weapon.GetDamage() + m.Strength
}

func (m *Monster) GetDescription() string {
	return m.description
}

func (m *Monster) SendMessage(msg interface{}) {
	//Do nothing, required for agenter interface
}

func (m *Monster) IsDead() bool {
	return m.currentHP < 0
}

func (m *Monster) IsAttackingPlayer(name string) bool {
	for _, targets := range m.targets {
		if strings.EqualFold(targets.attackTarget.Name, name) {
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

func LoadMonsterData() {
	monsterTemplatesG = make(map[string]*Monster)
	xmlFile, err := os.Open("monsterData.xml")
	checkErrorWithMessage(err, true, " In Load monster data.")
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var monstersData MonstersXML
	err = xml.Unmarshal(XMLdata, &monstersData)
	checkErrorWithMessage(err, true, " In load Monster Data function.")

	for _, element := range monstersData.Monsters {
		monsterTemplatesG[element.Name] = NewMonsterFromXML(element)
	}

	fmt.Println("Monster Data Loaded.")
}
