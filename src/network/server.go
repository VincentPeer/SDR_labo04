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
}

// newServer creates a new server with the given UDP connection, onMessage callback, and onError callback.
func newServer(config *networkConfig, configID int, udp *udpserver.UDP, onMessage func(Message), onSend func(Message), onError func(error), debug bool) *server {
	s := &server{
		udp:       udp,
		onMessage: onMessage,
		onSend:    onSend,
		onError:   onError,
		config:    config,
		configID:  configID,
		debug:     debug,
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

	s := newServer(config, serverID, udp, onMessage, onSend, onError, debug)

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
