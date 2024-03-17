package server

import (
	"log"

	"github.com/olahol/melody"
)

type Peer struct {
	session *melody.Session
}

func NewPeer(s *melody.Session) *Peer {
	return &Peer{
		session: s,
	}
}

func (p *Peer) Send(packet Packet) {
	result, err := packet.Marshall()
	if err != nil {
		log.Fatal(err)
	}

	p.session.Write(result)
}
