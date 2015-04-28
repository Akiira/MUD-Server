package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

var (
	servers            map[string]string = make(map[string]string)
	eventManager       *EventManager
	serverName         string
	worldRespawnRoomID int
)

func main() {
	gob.Register(WeaponXML{})
	gob.Register(ArmourXML{})
	gob.Register(ArmourSetXML{})
	gob.Register(ItemXML{})

	if len(os.Args) < 2 {
		fmt.Println(os.Args[0] + " requires 1 arguments, worldname")
		os.Exit(1)
	}

	ReadServerAddresses()
	go RunServer()
	GetInputFromUser()
}

func GetInputFromUser() {
	var input string
	for {

		fmt.Scan(&input)
		input = strings.TrimSpace(input)

		if input == "exit" {
			eventManager.SaveAllCharacters()
			os.Exit(1)
		} else if input == "exitNoSave" {
			os.Exit(1)
		} else if input == "refreshserver" {
			ReadServerAddresses()
		} else {
			fmt.Println("Bad input.")
		}
	}
}

//ReadServerAddresses reads the list of server names and their corresponding addresses
//in from a txt file. The names and addresses are stored in the global servers variable.
func ReadServerAddresses() error {
	file, err := os.Open("serverConfig/serverList.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		readData := strings.Fields(scanner.Text())
		fmt.Println(readData)
		servers[readData[0]] = readData[1]
	}

	return scanner.Err()
}

//NewServerListener takes an internet address - could just be a port number - and
//starts a server on that address. The protocol used is tcp and the object returned
//is a listener, which will listen to that address waiting for connections.
func NewServerListener(addr string) *net.TCPListener {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", addr)
	checkError(err, true)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err, true)
	return listener
}

func RunServer() {
	LoadMonsterData()
	serverName = os.Args[1]
	eventManager = NewEventManager()

	listener := NewServerListener(servers[serverName])

	for {
		conn, err := listener.Accept()
		checkError(err, false)

		if err == nil {
			err = HandleConnection(conn)
			checkError(err, false)
		}
	}
}

func HandleConnection(conn net.Conn) (err error) {
	var clientResponse ClientMessage
	decoder := gob.NewDecoder(conn)

	err = decoder.Decode(&clientResponse)
	checkError(err, false)

	if clientResponse.Command == "heartbeat" {
		return HandleHeartBeatConnection(conn)
	} else if clientResponse.Command == "refreshserver" {
		return HandleServerRefresh(conn, clientResponse.Value)
	} else {
		return HandlePlayerConnection(conn, decoder, clientResponse.getUsername())
	}
}

func HandleServerRefresh(conn net.Conn, updatedAddress string) (err error) {
	defer conn.Close()

	fmt.Println("refresh server Connection established")
	err = ioutil.WriteFile("./serverConfig/serverList.txt", []byte(updatedAddress), 0666)

	//Since failing to update these addresses means the server cant do its job
	//any more, its best to just shut it down.
	checkError(err, true)

	return ReadServerAddresses()
}

func HandleHeartBeatConnection(conn net.Conn) (err error) {
	defer conn.Close()

	fmt.Println("heartbeat Connection established")

	//Simply send a beat back to let the central server know
	//this server is still alive
	err = gob.NewEncoder(conn).Encode(newServerMessageTypeS(REPLYPING, "beat"))

	return err
}

func HandlePlayerConnection(conn net.Conn, decoder *gob.Decoder, charsName string) (err error) {
	fmt.Println("Player Connection established")
	var clientsCharacter *Character

	if clientsCharacter, err = GetCharacterFromStorage(charsName); err == nil {
		NewClientConnection(conn, clientsCharacter, decoder, gob.NewEncoder(conn))
		eventManager.AddPlayerToRoom(clientsCharacter)
	}

	return err
}

//TODO is this more appropiate here or in the character class?
func SendCharactersXML(charData *CharacterXML) {

	conn, err := net.Dial("tcp", servers["characterStorage"])
	checkError(err, true)
	defer conn.Close()

	encdr := gob.NewEncoder(conn)
	err = encdr.Encode(&ServerMessage{MsgType: SAVEFILE, Value: newFormattedStringSplice(charData.Name)})
	checkError(err, true)

	err = encdr.Encode(*charData)
	checkError(err, false)
}

func checkError(err error, exitIfError bool) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		if exitIfError {
			os.Exit(1)
		}
	}
}

func checkErrorWithMessage(err error, exitIfError bool, messageIfError string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		fmt.Fprintf(os.Stderr, "Additional Message: %s", messageIfError)
		if exitIfError {
			os.Exit(1)
		}
	}
}
