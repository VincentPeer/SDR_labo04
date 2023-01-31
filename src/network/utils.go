package network

import (
	"SDR_labo04/src/udpserver"
	"fmt"
	"time"
)

const (
	timeoutErrorMessage = "the connection to the server timed out"
)

type ackable interface {
	getAckChannel(connection *udpserver.UDPAddress) chan Message
	createAckChannel(connection *udpserver.UDPAddress) chan Message
	closeAckChannel(connection *udpserver.UDPAddress)
	getOutgoingConnection() *udpserver.UDP
	getConfig() *networkConfig
}

// sendToServer sends a message of the given type with the given message string to the
// specified server ID. It returns an error if any occurred.
func sendToServer(conn *udpserver.UDP, config *networkConfig, msgType string, message interface{}, serverID int) error {
	// Check if server id is valid
	if serverID < 0 || serverID > config.MaxServers {
		return fmt.Errorf("invalid server id, must be between 0 and %d specified in config", config.MaxServers)
	}
	// Get configuration for specified server
	configServer := config.Servers[serverID]

	// Check if the servers are neighbours
	for i, neighbour := range configServer.Neighbors {
		if neighbour == conn.ID {
			break
		} else if i == len(configServer.Neighbors)-1 {
			return fmt.Errorf("server %s is not a neighbour of server %d", conn.ID, serverID)
		}
	}

	remoteCon := udpserver.NewUDP(configServer.Address, configServer.Port, configServer.ID)

	jsonMessage, err := StringifyMessage(Message{
		Type:     msgType,
		Sender:   conn.ID,
		Receiver: remoteCon.ID,
		Data:     message,
	})

	if err != nil {
		return err
	}

	// Send message to server
	return conn.Send(&remoteCon.UDPAddress, jsonMessage)
}

func sendWithAck(ackHandler ackable, msgType string, message interface{}, serverID int, onAck func(Message), onErr func(error)) {

	config := ackHandler.getConfig()
	serverConn := udpserver.NewUDPConn(config.Servers[serverID].Address, config.Servers[serverID].Port)
	callbackChan := ackHandler.createAckChannel(&serverConn)

	// Send message to server
	err := sendToServer(ackHandler.getOutgoingConnection(), config, msgType, message, serverID)
	if err != nil {
		onErr(err)
		return
	}

	// Wait for acknowledgement
	go func() {
		select {
		case ack := <-callbackChan:
			onAck(ack)
		case <-time.After(time.Duration(config.Timeout) * time.Millisecond):
			onErr(fmt.Errorf("%s for message %s : %s", timeoutErrorMessage, msgType, message))
		}
		ackHandler.closeAckChannel(&serverConn)
	}()
}