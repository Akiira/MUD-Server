// Monster
package main

import (
	"math/rand"
	"io/ioutil"
	"encoding/xml"
	"os"
	"github.com/daviddengcn/go-colortext"
)
type Monster struct {
	Name string
	HP int	
	Defense int
	description string
}

type MonsterXML struct {
	Name string `xml:"Name"`
	HP int `xml:"HP"`
	Defense int `xml:"Defense"`
	Description string `xml:"Description"`
}
type MonstersXML struct {
	XMLName xml.Name `xml:"Monsters"`
	Monsters []MonsterXML `xml:"Monster"`
}

var monsterTemplatesG map[string]*Monster

func newMonsterFromXML(monsterData MonsterXML) *Monster {
	m := new(Monster)
	m.Name = monsterData.Name
	m.HP = monsterData.HP
	m.Defense = monsterData.Defense
	m.description = monsterData.Description

	return m
}

func newMonsterFromName(name string) *Monster {
	m := new(Monster)
	m.Name = monsterTemplatesG[name].Name
	m.HP = monsterTemplatesG[name].HP
	m.Defense = monsterTemplatesG[name].Defense
	m.description = monsterTemplatesG[name].description

	return m
}

func (m *Monster) getAttackRoll() int {
	return rand.Int() % 6
}

func (m *Monster) getLookDescription() []FormattedString {
	output := make([]FormattedString, 1, 1)
	
	output[0].Color = ct.Yellow
	output[0].Value = m.description
	
	return output
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