package server

import (
	"encoding/json"

	"github.com/StellarWar/world-server/internal/shared"
	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
)

type NetPeer struct {
	s *melody.Session
}

func newNetPeer(s *melody.Session) *NetPeer {
	return &NetPeer{
		s: s,
	}
}

func (np *NetPeer) Send(packet *shared.Packet) {
	data, err := json.Marshal(packet)
	if err != nil {
		log.Error().Err(err).Msg("packet write failed")
		return
	}

	np.s.Write(data)
}
