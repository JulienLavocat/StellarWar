package server

import "encoding/json"

type ClientPacket struct {
	Type    string          `json:"type"`
	Payload json.RawMessage `json:"payload"`

	PlayerId PlayerId
}

type Packet struct {
	Payload interface{} `json:"payload"`
	Type    string      `json:"type"`
}
