package main

const (
	PLAYER  = true
	MONSTER = false
)

type Event struct { //TODO use an Agent interface to avoid two pointers
	playerORMonster bool
	nameOfAgent     string
	action          string
	valueOrTarget   string
	character       *Character
	monster         *Monster
}


//TODO fix these functions to use only *Agent rahter then charName, *Character or *Monster
func newEventFromMessage(msg ClientMessage, charName string) Event {
	return newEvent(PLAYER, charName, msg.Command, msg.Value)
}
func newEventFromMessage2(msg ClientMessage, char *Character) Event {
	return newEvent(PLAYER, charName, msg.Command, msg.Value, character: char, monster: nil)
}

func newEvent(agent bool, nameOfAgent string, action string, value string) Event {
	event := new(Event)
	event.playerORMonster = agent
	event.nameOfAgent = nameOfAgent
	event.action = action
	event.valueOrTarget = value
	return *event
}
func newEvent2(agent bool, agent *Character, action string, value string) Event {
	event := new(Event)
	event.playerORMonster = agent
	event.nameOfAgent = agent.Name
	event.action = action
	event.valueOrTarget = value
	event.character = agent
	event.monster = nil
	return *event
}
func newEvent3(agent bool, agent *Monster, action string, value string) Event {
	event := new(Event)
	event.playerORMonster = agent
	event.nameOfAgent = agent.Name
	event.action = action
	event.valueOrTarget = value
	event.character = nil
	event.monster = agent
	return *event
}
