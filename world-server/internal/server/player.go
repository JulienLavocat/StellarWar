package server

import (
	"sync/atomic"

	"github.com/StellarWar/world-server/internal/shared"
)

var playerIdCounter = atomic.Uint32{}

type Player struct {
	peer *NetPeer
	Id   shared.PlayerId
}

func newPlayer(peer *NetPeer) *Player {
	return &Player{
		peer: peer,
		Id:   shared.PlayerId(playerIdCounter.Add(1)),
	}
}

func (p *Player) Send(packet *shared.Packet) {
	p.peer.Send(packet)
}
