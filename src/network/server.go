package network

import (
	"SDR_labo04/src/udpserver"
	"fmt"
)

const (
	debugSleepTime = 1000 // time to sleep in ms when debug is enabled
)

// server represents a UDP server that listens for incoming messages and processes them.
type server struct {
	udp       *udpserver.UDP // udp is a pointer to the UDP connection of the server.
	onMessage func(Message)  // onMessage is a callback function that is called when a message is received.
	onSend    func(Message)  // onSend is a callback function that is called when a message is sent.
	onError   func(error)    // onError is a callback function that is called when an error occurs.
	config    *networkConfig // config is a pointer to the network configuration.
	configID  int            // configID is the ID of the server in the network configuration.
	debug     bool           // debug is true if debug messages should be printed and messages slow down
	lettersCounted int 	  // lettersCounted is the number of letters counted by the server
	letter string 		  // letter is the letter that the server is counting
	activesNeighbours map[string]bool // activesNeighbours is a map of the active neighbours of the server
	neighborsChan map[string]chan map[string]interface{}
	result map[string]int
}

// newServer creates a new server with the given UDP connection, onMessage callback, and onError callback.
func newServer(config *networkConfig, configID int, udp *udpserver.UDP, onMessage func(Message), onSend func(Message), onError func(error), debug bool, letter string) *server {
	s := &server{
		udp:       udp,
		onMessage: onMessage,
		onSend:    onSend,
		onError:   onError,
		config:    config,
		configID:  configID,
		debug:     debug,
		lettersCounted: 0,
		letter: letter,
		activesNeighbours: make(map[string]bool),
		neighborsChan: make(map[string]chan map[string]interface{}),
		result: make(map[string]int),
	}

	for i := 0; i < len(config.Servers[configID].Neighbors); i++ {
		s.neighborsChan[config.Servers[configID].Neighbors[i]] = make(chan map[string]interface{})
		s.activesNeighbours[config.Servers[configID].Neighbors[i]] = true
	}
	return s
}

// StartServer starts a UDP server using the configuration parameters specified in the given JSON file.
// If a message is received, the onMessage function is called with the message and the sender's address.
// If an error occurs, the onError function is called with the error.
func StartServer(configFile string, serverID int, onMessage func(Message), onSend func(Message), onError func(error), debug bool) {
	// Load configuration from JSON file
	config, err := fromJSON(configFile)
	if err != nil {
		onError(err)
		return
	}

	// Check validity of server ID
	if serverID < 0 || serverID > config.MaxServers {
		onError(fmt.Errorf("invalid server id, must be between 0 and %d specified in config", config.MaxServers))
		return
	}

	// Get configuration for specified server
	configServer := config.Servers[serverID]

	// Create new UDP server with server configuration
	udp := udpserver.NewUDP(configServer.Address, configServer.Port, configServer.ID)

	s := newServer(config, serverID, udp, onMessage, onSend, onError, debug, configServer.Letter)

	// Sends wazzup messages to all servers to inform them that this server is alive
	s.sendToAll(typeWazzup, "")

	// Start listening for UDP messages
	udp.Listen(
		// Parses the message and calls the onMessage function
		func(message string, remoteAddr *udpserver.UDPAddress) {
			parsedMessage, err := ParseMessage(message)
			if err != nil {
				onError(err)
				return
			}
			s.handleMessage(parsedMessage, remoteAddr)

			s.onMessage(parsedMessage)
		},
		onError)
}

// Stop stops the server.
func (s *server) Stop() {
	s.udp.Started = false
}

// handleMessage processes an incoming message and sends an acknowledgement message if appropriate.
func (s *server) handleMessage(message Message, remoteAddr *udpserver.UDPAddress) {
	// Faire des dingueries

	switch message.Type {
		case typeSend:
			word := message.Data.(string)
			s.lettersCounted = letterCounter(s.letter, word)
			s.result[s.letter] = s.lettersCounted
			go s.waveAlgorithm()
		case typeWave:
			waveMessage := message.Data.(map[string]interface{})
			fmt.Println("Received wave message ", waveMessage)
			s.neighborsChan[message.Sender] <- waveMessage
		case typeResult:
			resultMessage, err := StringifyMessage(
				Message{
					Type:     typeAck,
					Sender:   s.config.Servers[s.configID].ID,
					Receiver: message.Sender,
					Data:     s.result,
				})
				if err != nil {
					s.onError(err)
					return
				}
				s.udp.Send(remoteAddr, resultMessage)
		}
}

func (s *server) sendToAll(msgType string, msgContent interface{}) {

	for i := 0; i < len(s.config.Servers); i++ {
		if s.config.Servers[i].ID != s.udp.ID {
			err := sendToServer(s.udp, s.config, msgType, msgContent, i)
			if err == nil {
				s.onSend(Message{
					Type:     msgType,
					Sender:   s.config.Servers[s.configID].ID,
					Receiver: s.config.Servers[i].ID,
					Data:     msgContent,
				})
			}	
		}
	}
}

type WaveMessage struct {
	Result map[string]int
	Server_id string
	Active bool
}

func (s *server) waveAlgorithm() {
	for {
		s.sendToAll(typeWave, WaveMessage{
			Result: s.result,
			Server_id: s.config.Servers[s.configID].ID,
			Active: true,
		})
		for k, a := range s.activesNeighbours {
			if a {
				fmt.Println("Waiting for message from ", k)
				msg := <- s.neighborsChan[k]
				s.activesNeighbours[k] = msg["Active"].(bool)
				
				for key, value := range msg["Result"].(map[string]interface{}) {
					if _, ok := s.result[key]; !ok {
						s.result[key] = (int)(value.(float64))
					} else if s.result[key] < (int)(value.(float64)) {
						s.result[key] = (int)(value.(float64))
					}
				}
			}
		}
		if len(s.result) == len(s.config.Servers) {
			break
		}
	}
	fmt.Println("I think I have the result")
	s.sendToAll(typeWave, WaveMessage{
		Result: s.result,
		Server_id: s.config.Servers[s.configID].ID,
		Active: false,
	})
	for k, a := range s.activesNeighbours {
		if a {
			msg := <- s.neighborsChan[k]
			s.activesNeighbours[k] = msg["Active"].(bool)
		}
	}
	fmt.Println("Finished wave algorithm")
}

func (s *server) getOutgoingConnection() *udpserver.UDP {
	return s.udp
}

func (s *server) getConfig() *networkConfig {
	return s.config
}

func letterCounter(letter string, text string) int {
	count := 0
	for _, char := range text {
		if string(char) == letter {
			count++
		}
	}
	return count
}
