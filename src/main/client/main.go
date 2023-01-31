package main

import (
	"SDR_labo04/src/network"
	"flag"
	"fmt"
)

const (
	id          = "client_0"
	sendCommand = "send"
)

var (
	port       *int    // port to listen on
	configPath *string // path to the json configuration file of the network
	command    *string // command to send to server
	serverId   *int    // id of the server to send the command to
)

func init() {
	port = flag.Int("port", 8079, "port to listen on")
	configPath = flag.String("path", "../data/config.json", "path to the json configuration file of the network")
	command = flag.String("command", "sendWithAck", "command to send to server (send, receive, etc)") // TODO add real commands
	serverId = flag.Int("server", 1, "id of the server to send the command to")
}

func usage() {
	fmt.Println("Usage: main.go -port -path -command -server")
	flag.PrintDefaults()
}

func main() {
	flag.Parse()

	client, err := network.NewClient(8078, *port, id, *configPath)

	if err != nil {
		fmt.Println(err)
		usage()
		return
	}

	switch *command {
	case sendCommand:
		err = client.Send("Hello World!", *serverId)
		if err != nil {
			fmt.Println(err)
		}
	default:
		fmt.Println("Unknown command")
	}

}
