// Monster
package main

import (
	"encoding/xml"
	"fmt"
	"github.com/daviddengcn/go-colortext"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
)

type Monster struct {
	Agent
	Defense     int
	description string
	targets     map[string]*target
	em          EventManager
	weapon      Weapon

	lastTarget Agenter
}

//TODO add mutex around targets field

type target struct {
	aggro        int
	attackTarget *ClientConnection
}

var monsterTemplatesG map[string]*Monster

func newMonsterFromXML(monsterData MonsterXML) *Monster {
	m := new(Monster)
	m.Name = monsterData.Name
	m.currentHP = monsterData.HP
	m.Defense = monsterData.Defense
	m.description = monsterData.Description
	m.targets = make(map[string]*target)

	return m
}

func newMonsterFromName(name string) *Monster {
	m := new(Monster)
	m.Name = monsterTemplatesG[name].Name
	m.currentHP = monsterTemplatesG[name].currentHP
	m.Defense = monsterTemplatesG[name].Defense
	m.description = monsterTemplatesG[name].description
	m.targets = make(map[string]*target)

	return m
}

func (m *Monster) fightPlayers() {
	for {
		time.Sleep(2 * time.Second)

		if m.currentHP <= 0 || len(m.targets) <= 0 {
			break
		}

		//Find target with highest aggro
		var attackTarget *ClientConnection
		maxAggro := 0
		for _, targ := range m.targets {
			if targ.aggro > maxAggro {
				attackTarget = targ.attackTarget
			}
		}

		event := newEvent(m, "attack", attackTarget.character, attackTarget)
		m.em.addEvent(event)
	}
}

//TODO implement monsters combat functions

func (m *Monster) addNewTarget(targetCC *ClientConnection, agro int) {

	_, exist := m.targets[targetCC.character.Name]

	if exist {
		m.targets[targetCC.character.Name].aggro += agro
	} else {
		targ := target{aggro: agro, attackTarget: targetCC}
		m.targets[targetCC.character.Name] = &targ

		if len(m.targets) == 1 {
			go m.fightPlayers()
		}
	}
}

func (m *Monster) getAttackRoll() int {
	return (rand.Int() % 20) + m.weapon.attack + m.Strength
}

func (m *Monster) takeDamage(amount int, typeOfDamge int) []FormattedString {
	m.currentHP -= amount
	return nil
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

func (m *Monster) isDead() bool {
	return m.currentHP > 0
}

func (m *Monster) makeAttack(target Agenter) []FormattedString {
	a1 := m.getAttackRoll()
	if a1 >= target.getDefense() {
		target.takeDamage(m.getDamage(), 0)
		output := newFormattedStringCollection()
		output.addMessage(ct.Red, fmt.Sprintf("The %s hit you for %i damage\n", m.Name, m.getDamage()))

		if target.isDead() {
			output.addMessage(ct.Red, "\nYou died!.\n")
		}
		return output.fmtedStrings
	}

	return newFormattedStringSplice2(ct.Red, fmt.Sprintf("The %s's attack missed you.\n", m.Name))
}

func (m *Monster) getClientConnection() *ClientConnection {
	return m.lastTarget.getClientConnection()
}

func (m *Monster) getDamage() int {
	return m.weapon.damage + m.Strength
}

func (m *Monster) getLookDescription() []FormattedString {
	output := make([]FormattedString, 1, 1)

	output[0].Color = ct.Yellow
	output[0].Value = m.description

	return output
}

//------------------Loading Stuff------------------------------
type MonsterXML struct {
	Name        string `xml:"Name"`
	HP          int    `xml:"HP"`
	Defense     int    `xml:"Defense"`
	Description string `xml:"Description"`
}

type MonstersXML struct {
	XMLName  xml.Name     `xml:"Monsters"`
	Monsters []MonsterXML `xml:"Monster"`
}

func loadMonsterData() {
	monsterTemplatesG = make(map[string]*Monster)
	xmlFile, err := os.Open("monsterData.xml")
	checkError(err, true)
	defer xmlFile.Close()

	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var monstersData MonstersXML
	xml.Unmarshal(XMLdata, &monstersData)

	for _, element := range monstersData.Monsters {
		monsterTemplatesG[element.Name] = newMonsterFromXML(element)
	}
}
