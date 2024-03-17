package server

import (
	"github.com/StellarWar/world-server/internal/galaxy"
	"github.com/StellarWar/world-server/internal/packets"
	"github.com/StellarWar/world-server/internal/shared"
	"github.com/rs/zerolog/log"
)

type ClientPacketHandler func(s *Server, p *Player, packet *shared.ClientPacket)

type Server struct {
	galaxy galaxy.Galaxy

	players  map[shared.PlayerId]*Player
	handlers map[string]ClientPacketHandler

	unregister chan shared.PlayerId
	messages   chan *shared.ClientPacket
}

func NewServer(g galaxy.Galaxy) *Server {
	return &Server{
		galaxy:     g,
		unregister: make(chan shared.PlayerId),
		messages:   make(chan *shared.ClientPacket),
		players:    make(map[shared.PlayerId]*Player),
		handlers: map[string]ClientPacketHandler{
			"register":   registerHandler,
			"unregister": unregisterHandler,
			"test_packet": func(s *Server, p *Player, packet *shared.ClientPacket) {
				log.Info().Any("packet", packet).Msg("got test packet")
			},
		},
	}
}

func (s *Server) CreatePlayer(peer *NetPeer) *Player {
	player := newPlayer(peer)
	s.players[player.Id] = player

	s.messages <- packets.NewRegisterPacket(player.Id)

	return player
}

func (s *Server) Run() {
	for {
		msg := <-s.messages
		log.Info().Any("packet", msg).Msg("received")

		player, exists := s.players[msg.PlayerId]
		if !exists {
			continue
		}

		if handler, exists := s.handlers[msg.Type]; exists {
			handler(s, player, msg)
		}
	}
}

func (s *Server) PushPacket(packet *shared.ClientPacket) {
	s.messages <- packet
}

func registerHandler(s *Server, p *Player, packet *shared.ClientPacket) {
	log.Info().Msgf("Player %d registered", packet.PlayerId)
}

func unregisterHandler(s *Server, p *Player, packet *shared.ClientPacket) {
	log.Info().Msgf("Player %d unregistered", packet.PlayerId)
	delete(s.players, p.Id)
}
