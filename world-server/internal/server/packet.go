package server

import "encoding/json"

type Packet struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data"`
}

func (p *Packet) Marshall() ([]byte, error) {
	return json.Marshal(p)
}
