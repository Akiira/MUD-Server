package main

const (
	PLAYER  = true
	MONSTER = false
)

type Event struct {
	playerORMonster bool
	action          string
	valueOrTarget   Agenter //TODO consider changing this to an Ageneter as well
	agent           Agenter
	client          *ClientConnection
}

func newEventFromMessage(msg ClientMessage, agent Agenter, cc *ClientConnection) Event {
	target := eventManager.worldRooms[cc.getCharactersRoomID()].getMonster(msg.Value)
	return newEvent(PLAYER, agent, msg.Command, target, cc)
}

func newEvent(playerORMonster bool, agent Agenter, action string, target Agenter, cc *ClientConnection) Event {
	event := new(Event)
	event.playerORMonster = playerORMonster
	event.agent = agent
	event.action = action
	event.valueOrTarget = target
	event.client = cc
	return *event
}
