// Monster
package main

import (
	"math/rand"
	"io/ioutil"
	"encoding/xml"
	"os"
)
type Monster struct {
	Name string
	HP int	
	Defense int
}

type MonsterXML struct {
	Name string `xml:"Name"`
	HP int `xml:"HP"`
	Defense int `xml:"Defense"`
}
type MonstersXML struct {
	XMLName xml.Name `xml:"Monsters"`
	Monsters []MonsterXML `xml:"Monster"`
}

func newMonsterFromXML(monsterData MonsterXML) *Monster {
	m := new(Monster)
	m.Name = monsterData.Name
	m.HP = monsterData.HP
	m.Defense = monsterData.Defense
	
	return m
}

func (m *Monster) getAttackRoll() int {
	return rand.Int() % 6
}

func loadMonsterData() {
	xmlFile, err := os.Open("monsterData.xml")
	checkError(err)
	defer xmlFile.Close()
	
	XMLdata, _ := ioutil.ReadAll(xmlFile)

	var monstersData MonstersXML
    xml.Unmarshal(XMLdata, &monstersData)
	
	for _, element := range monstersData.Monsters {
		monsterTemplatesG[element.Name] = newMonsterFromXML(element)
	}
}