package main

import (
	"SDR_labo04/src/network"
	"flag"
	"fmt"
)

var (
	id         *int    // id of the server to start
	configPath *string // path to the json configuration file of the network
	debug 	*bool   // debug mode
)

func init() {
	id = flag.Int("id", -1, "id of the server") // -1 is an invalid id
	configPath = flag.String("path", "../data/config.json", "path to the json configuration file of the network")
	debug = flag.Bool("debug", false, "debug mode")
}

func onMessage(message network.Message) {
	if *debug {
		fmt.Printf("| %-10s | %-10s | %-20s | %-10s |\n", "Received", message.Type, message.Data, message.Sender)
	} else {
		fmt.Println("<- ", message.Sender, " : ", message.Type)
	}
}

func onSend(message network.Message) {
	if *debug {
		fmt.Printf("| %-10s | %-10s | %-20s | %-10s |\n", "Sent", message.Type, message.Data, message.Receiver)
	} else {
		fmt.Println("-> ", message.Receiver, " : ", message.Type)
	}
}

func onError(err error) {
	fmt.Println("Usage: main.go -id -path")
	flag.PrintDefaults()
	panic(err)
}

func main() {
	flag.Parse()

	// check if id is valid
	if *id < 0 {
		onError(fmt.Errorf("invalid id %d must be greater than 0", *id))
	}

	network.StartServer(*configPath, *id, onMessage, onSend, onError, *debug)
}
