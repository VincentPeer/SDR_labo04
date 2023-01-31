package network

import (
	"SDR_labo04/src/udpserver"
	"fmt"
	"reflect"
)

const (
	debugSleepTime = 1000 // time to sleep in ms when debug is enabled
)

// server represents a UDP server that listens for incoming messages and processes them.
type server struct {
	udp            *udpserver.UDP // udp is a pointer to the UDP connection of the server.
	onMessage      func(Message)  // onMessage is a callback function that is called when a message is received.
	onSend         func(Message)  // onSend is a callback function that is called when a message is sent.
	onError        func(error)    // onError is a callback function that is called when an error occurs.
	config         *networkConfig // config is a pointer to the network configuration.
	configID       int            // configID is the ID of the server in the network configuration.
	debug          bool           // debug is true if debug messages should be printed and messages slow down
	lettersCounted int
	myLetter       string
	nbNeighbors    int
	neighborsChan  map[string]chan map[string]int
	parentId       string
	result         map[string]int
	started        bool
	probeToReceive int
	client *udpserver.UDPAddress
}

// newServer creates a new server with the given UDP connection, onMessage callback, and onError callback.
func newServer(config *networkConfig, configID int, udp *udpserver.UDP, onMessage func(Message), onSend func(Message),
	onError func(error), debug bool, lettersCounted int, myLetter string, nbNeighbors int, parentId string) *server {
	s := &server{
		udp:            udp,
		onMessage:      onMessage,
		onSend:         onSend,
		onError:        onError,
		config:         config,
		configID:       configID,
		debug:          debug,
		lettersCounted: lettersCounted,
		myLetter:       myLetter,
		nbNeighbors:    nbNeighbors,
		parentId:       parentId,
		neighborsChan:  make(map[string]chan map[string]int),
		result:         make(map[string]int),
		started:        false,
		probeToReceive: nbNeighbors,
		client:  	 nil,
	}

	for i := 0; i < len(config.Servers[configID].Neighbors); i++ {
		s.neighborsChan[config.Servers[configID].Neighbors[i]] = make(chan map[string]int)
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

	s := newServer(config, serverID, udp, onMessage, onSend, onError, debug, 0, configServer.Letter, len(configServer.Neighbors), "-1")

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
	fmt.Println(message)
	switch message.Type {
	case typeSendAck: // first server to get the message
		s.client = remoteAddr
		text := message.Data.(string)
		s.lettersCounted = letterCounter(s.myLetter, text)
		s.result[s.myLetter] = s.lettersCounted

		// send Probe to all neighbors
		s.sendToAll(typeProbe, text)
		s.nbNeighbors++
	case typeProbe:
		s.probeToReceive--
		probe := message.Data.(string)
		s.lettersCounted = letterCounter(probe, s.myLetter)
		s.result[s.myLetter] = s.lettersCounted
		if s.probeToReceive == 0 {
			for i := 0; i < s.config.MaxServers; i++ {
				if s.config.Servers[i].ID == message.Sender {
					err := sendToServer(s.udp, s.config, typeEcho, s.result, i)
					if err != nil {
						s.onError(err)
					}
					break
				}
			}
		} else {
		if !s.started {
			s.parentId = message.Sender
			s.started = true
			// send Probe to all neighbors except to the parent
			for i := 0; i < len(s.config.Servers[s.configID].Neighbors); i++ {
				fmt.Println("Neighbor ", s.config.Servers[s.configID].Neighbors[i])
				fmt.Println("parent", s.parentId)
				if s.config.Servers[s.configID].Neighbors[i] != s.parentId {
					fmt.Println("Sending probe to ", s.config.Servers[s.configID].Neighbors[i])
					for j := 0; j < s.config.MaxServers; j++ {
						if s.config.Servers[j].ID == s.config.Servers[s.configID].Neighbors[i] {
							err := sendToServer(s.udp, s.config, typeProbe, probe, j)
							if err != nil {
								println("Error sending probe to ", s.config.Servers[s.configID].Neighbors[i])
							}
							break
						}
					}
				}
			}
		} else {
			for i := 0; i < s.config.MaxServers; i++ {
				if s.config.Servers[i].ID == s.parentId {
					err := sendToServer(s.udp, s.config, typeEcho, s.result, i)
					if err != nil {
						s.onError(err)
					}
					break
				}
			}
		}
	}
	case typeEcho:
		result := message.Data.(map[string]interface{})
		result[s.myLetter] = s.lettersCounted
		for key, value := range result {
			fmt.Println("key:", key, "value:", value)
			if reflect.TypeOf(value).Kind() == reflect.Float64 {
				s.result[key] = int(value.(float64))
			} else {
				s.result[key] = value.(int)
			}
		}
		// send echo to the parent with the result, if we are not the root
		s.nbNeighbors--
		if s.nbNeighbors == 1 {
			if s.parentId != "-1" {
				for i := 0; i < s.config.MaxServers; i++ {
					if s.config.Servers[i].ID == s.parentId {
						err := sendToServer(s.udp, s.config, typeEcho, s.result, i)
						if err != nil {
							s.onError(err)
						}
						break
					}
				}
			} else {
				// if we are the root, we print the result
				fmt.Println("Result: ", s.result)
				msg, err := StringifyMessage(
					Message{
						Type:     typeAck,
						Sender:   s.config.Servers[s.configID].ID,
						Receiver: "client",
						Data:     s.result,
					})
				if err != nil {
					s.onError(err)
				}
				s.udp.Send(s.client, msg)
			}
		}
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

func (s *server) getOutgoingConnection() *udpserver.UDP {
	return s.udp
}

func (s *server) getConfig() *networkConfig {
	return s.config
}

func (s *server) diffusionAlgorithm() {

}

type Probe struct {
	Id      string
	Message string
}

func letterCounter(text string, letter string) int {
	count := 0
	for _, char := range text {
		if string(char) == letter {
			count++
		}
	}
	return count
}
