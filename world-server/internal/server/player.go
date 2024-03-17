package server

import "sync/atomic"

type PlayerId uint32

var playerIdCounter = atomic.Uint32{}

type Player struct {
	peer *NetPeer
	Id   PlayerId
}

func newPlayer(peer *NetPeer) *Player {
	return &Player{
		peer: peer,
		Id:   PlayerId(playerIdCounter.Add(1)),
	}
}

func (p *Player) Send(packet *Packet) {
	p.peer.Send(packet)
}
