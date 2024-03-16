package server

type Packet struct {
	Command string      `json:"command"`
	Data    interface{} `json:"data"`
}
