package network

import (
	"SDR_labo04/src/udpserver"
	"fmt"
	"time"
)

const (
	address = "127.0.0.1"
)

type Client struct {
	// Udp connection of the client to send messages
	udpConnSend *udpserver.UDP
	// Udp connection of the client to receive messages
	udpConnReceive *udpserver.UDP
	// Configuration of the network the client is connected to
	config *networkConfig
	// Map of channels to receive acknowledgements from the server by server id
	ackChannels map[udpserver.UDPAddress]chan Message
}

// NewClient creates a new client with the given address, port and id. It also takes
// a file path to a JSON configuration file that specifies the network configuration and
// an onError callback function that is called when an error occurs. It returns a pointer to
// the created client and an error if any occurred.
func NewClient(portSend int, portReceive int, id string, networkConfigPath string) (*Client, error) {
	// Check if port is valid
	if portSend < 0 || portSend > 65535 || portReceive < 0 || portReceive > 65535 {
		return nil, fmt.Errorf("invalid port, must be between 0 and 65535")
	}

	// Check if id is valid
	if id == "" {
		return nil, fmt.Errorf("invalid id, must not be empty")
	}

	// Load configuration from JSON file
	config, err := fromJSON(networkConfigPath)
	if err != nil {
		return nil, fmt.Errorf("can't read config file: %v", err)
	}

	client := Client{
		udpConnSend:    udpserver.NewUDP(address, portSend, id),
		udpConnReceive: udpserver.NewUDP(address, portReceive, id),
		config:         config,
		ackChannels:    make(map[udpserver.UDPAddress]chan Message),
	}

	return &client, nil

}

func (c *Client) SendToAll(message string) error {
	for i, _ := range c.config.Servers {
		err := c.Send(message, i)
		if err != nil {
			return err
		}
	}
	return nil
}

// Send sends a message to a server. It returns an error if any occurred.
func (c *Client) Send(message string, serverID int) error {
	return sendToServer(c.udpConnSend, c.config, typeSend, message, serverID)
}

func (c *Client) Result(serverID int) (Message, error) {
	return c.sendWithAckSync(typeResult, "", serverID)
}

// SendWithAckSync sends a message to a server and waits for an acknowledgement from the server, with a specified timeout.
// It returns the received acknowledgement message and an error if any occurred.
func (c *Client) SendWithAckSync(message string, serverID int) (Message, error) {
	return c.sendWithAckSync(typeSendAck, message, serverID)
}

// Stop sends a message to a server requesting to stop. It waits for an acknowledgement from the server,
// with a specified timeout. It returns the received acknowledgement message and an error if any occurred.
func (c *Client) Stop(serverID int) (Message, error) {
	return c.sendWithAckSync(typeStop, "", serverID)
}

// sendWithAckSync sends a message to a server and waits for an acknowledgement from the server,
// with a specified timeout. It returns the received acknowledgement message and an error if any occurred.
func (c *Client) sendWithAckSync(msgType string, message string, serverID int) (Message, error) {
	ackChannel := make(chan Message)
	errChannel := make(chan error)

	relayToErrorChannel := func(err error) {
		errChannel <- err
	}

	sendWithAck(c, msgType, message, serverID, func(ack Message) {
		ackChannel <- ack
	}, relayToErrorChannel)

	go c.listen(relayToErrorChannel)

	select {
	case ack := <-ackChannel:
		return ack, nil
	case err := <-errChannel:
		return Message{}, err
	case <-time.After(time.Duration(c.config.Timeout) * time.Millisecond):
		return Message{}, fmt.Errorf(timeoutErrorMessage)
	}
}

// listen listens for incoming messages from servers and processes them.
func (c *Client) listen(onErr func(error)) {
	// Start listening for UDP messages
	c.udpConnSend.Listen(
		// Checks if the message is an acknowledgement and calls the corresponding channel
		func(message string, remoteAddr *udpserver.UDPAddress) {
			parsedMessage, err := ParseMessage(message)
			if err != nil {
				onErr(err)
				return
			}
			if parsedMessage.Type == typeAck {
				if channel, ok := c.ackChannels[*remoteAddr]; ok {
					channel <- parsedMessage
				}
			}
		}, onErr)
}

func (c *Client) getAckChannel(remoteAddr *udpserver.UDPAddress) chan Message {
	if channel, ok := c.ackChannels[*remoteAddr]; ok {
		return channel
	}
	return nil
}

func (c *Client) createAckChannel(remoteAddr *udpserver.UDPAddress) chan Message {
	c.ackChannels[*remoteAddr] = make(chan Message)
	return c.ackChannels[*remoteAddr]
}

func (c *Client) closeAckChannel(remoteAddr *udpserver.UDPAddress) {
	close(c.ackChannels[*remoteAddr])
	delete(c.ackChannels, *remoteAddr)
}

func (c *Client) getOutgoingConnection() *udpserver.UDP {
	return c.udpConnSend
}

func (c *Client) getConfig() *networkConfig {
	return c.config
}
