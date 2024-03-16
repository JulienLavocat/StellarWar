package server

import (
	"fmt"
	"log"
	"strings"
	"time"

	"com.stellarwar/world-server/utils"
	"github.com/gorilla/websocket"
)

const (
	writeWait  = 10 * time.Second
	pongWait   = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10
)

type Player struct {
	Id             int32
	Color          [3]int
	HomeSystem     int
	ClaimedSystems []int

	server *Server
	socket *websocket.Conn

	send chan Packet
}

func NewPlayer(playerId int32, server *Server, socket *websocket.Conn) *Player {
	homeSystem := server.Galaxy.AssignNextHomeSystem(playerId)
	return &Player{
		Id:             playerId,
		Color:          utils.RandomHSV(),
		server:         server,
		socket:         socket,
		send:           make(chan Packet),
		HomeSystem:     homeSystem,
		ClaimedSystems: []int{homeSystem},
	}
}

func (p *Player) Send(packet Packet) {
	p.socket.WriteJSON(packet)
}

func (p *Player) ClaimSystem(systemId int) {
	p.ClaimedSystems = append(p.ClaimedSystems, systemId)
}

func (p *Player) Logf(format string, v ...interface{}) {
	log.Printf("[%d] %s", p.Id, fmt.Sprintf(format, v...))
}
func (p *Player) Log(line ...string) {
	log.Printf("[%d] %s", p.Id, strings.Join(line, " "))
}

func (p *Player) ReadPump() {
	defer func() {
		p.server.unregister <- p.Id
		p.socket.Close()
	}()

	p.socket.SetReadDeadline(time.Now().Add(pongWait))
	p.socket.SetPongHandler(func(string) error { p.socket.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := p.socket.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure, websocket.CloseNormalClosure) {
				p.Logf("error: %v", err)
			}
			break
		}
		p.server.inMessages <- ClientMessage{
			id:      p.Id,
			message: message,
		}
	}
}

func (p *Player) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		p.socket.Close()
	}()

	for {
		select {
		case packet, ok := <-p.send:
			p.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				p.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			p.socket.WriteJSON(packet)
		case <-ticker.C:
			p.socket.SetWriteDeadline(time.Now().Add(writeWait))
			if err := p.socket.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
