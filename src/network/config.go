package network

import (
	"encoding/json"
	"os"
)

// config is the representation of the server configuration
type config struct {
	// The id of the server
	ID string `json:"id"`
	// The port to listen on
	Port int `json:"port"`
	// The address to listen on
	Address string `json:"address"`
}

// networkConfig is the representation of the network config file
type networkConfig struct {
	// The servers in the network
	Servers []config `json:"servers"`
	// The maximum number of servers in the network
	MaxServers int `json:"maxServers"`
	// Timeout in milliseconds before a server is considered dead
	Timeout int `json:"timeout"`
}

// fromJSON reads the network config from a JSON file
func fromJSON(path string) (*networkConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	config := &networkConfig{}
	err = decoder.Decode(config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
