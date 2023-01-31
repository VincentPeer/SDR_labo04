package main

import (
	"SDR_labo04/src/network"
	"flag"
	"fmt"
)

const (
	id                 = "client_0"
	sendCommand        = "send"
	resultCommand	  = "result"
)

var (
	port       *int    // port to listen on
	configPath *string // path to the json configuration file of the network
	command    *string // command to send to server
	serverId   *int    // id of the server to send the command to
	word 	 *string // word to send to server
)

func init() {
	port = flag.Int("port", 8079, "port to listen on")
	configPath = flag.String("path", "../data/config.json", "path to the json configuration file of the network")
	command = flag.String("command", "sendWithAck", "command to send to server (send, receive, etc)") // TODO add real commands
	serverId = flag.Int("server", 1, "id of the server to send the command to")
	word = flag.String("word", "ABCDEF", "word to send to server")
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
		err = client.SendToAll(*word)
		if err != nil {
			fmt.Println(err)
		}
	case resultCommand:
		res, err := client.Result(*serverId)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(res.Data)
	default:
		fmt.Println("Unknown command")
	}

}
