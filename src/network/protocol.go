package network

import (
	"encoding/json"
	"fmt"
)

const (
	// typeSend is the type of a send message
	typeSend = "send"
	// typeSendAck is the type of a send ack message
	typeSendAck = "sendAck"
	// typeAck is the type of an ack message
	typeAck = "ack"
	// typeCharge is the type of a charge message
	typeCharge = "charge"
	// typeStop is the type of a stop message
	typeStop = "stop"
	// typeElect is the type of an elect message
	typeElect = "elect"
	// typeAnnounce is the type of an announce message
	typeAnnounce = "announce"
	// typeResult is the type of a result message
	typeResult = "result"
	// typeFinish is the type of a finish message
	typeFinish = "finish"
	// typeAskLeader is the type of an ask leader message
	typeAskLeader = "askLeader"
	// typeWazzup is the type of a wazzup message
	typeWazzup = "wazzup"
)

// Message est la structure de données utilisée pour représenter
// les messages envoyés et reçus par l'algorithme Chang et Roberts.
type Message struct {
	// Type indique le type de message (request ou reply).
	Type string `json:"type"`
	// Sender est l'identifiant de l'expéditeur du message.
	Sender string `json:"sender"`
	// Receiver est l'identifiant du destinataire du message.
	Receiver string `json:"receiver"`
	// Data est le contenu du message.
	Data interface{} `json:"data"`
}

// ParseMessage prend en entrée une chaîne de caractères JSON et renvoie
// un objet Message correspondant. Si le parsing échoue, une erreur est renvoyée.
func ParseMessage(jsonStr string) (Message, error) {
	var msg Message
	err := json.Unmarshal([]byte(jsonStr), &msg)
	if err != nil {
		return Message{}, fmt.Errorf("failed to parse message: %v", err)
	}
	return msg, nil
}

// StringifyMessage prend en entrée un objet Message et renvoie une chaîne de caractères
// JSON correspondante. Si la conversion échoue, une erreur est renvoyée.
func StringifyMessage(msg Message) (string, error) {
	jsonBytes, err := json.Marshal(msg)
	if err != nil {
		return "", fmt.Errorf("failed to stringify message: %v", err)
	}
	return string(jsonBytes), nil
}
