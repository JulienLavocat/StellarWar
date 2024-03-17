package server

import (
	"github.com/StellarWar/world-server/internal/galaxy"
	"github.com/rs/zerolog/log"
)

type Server struct {
	galaxy galaxy.Galaxy

	players map[PlayerId]*Player

	register   chan *Player
	unregister chan PlayerId
	messages   chan *ClientPacket
}

func NewServer() *Server {
	return &Server{
		galaxy:     galaxy.Galaxy{},
		register:   make(chan *Player),
		unregister: make(chan PlayerId),
		messages:   make(chan *ClientPacket),
		players:    make(map[PlayerId]*Player),
	}
}

func (s *Server) CreatePlayer(peer *NetPeer) *Player {
	player := newPlayer(peer)
	s.players[player.Id] = player

	return player
}

func (s *Server) Run() {
	for {
		select {
		case player := <-s.register:
			log.Info().Msgf("Player %d registered", player.Id)

		case playerId := <-s.unregister:
			log.Info().Msgf("Player %d unregistered", playerId)

		case msg := <-s.messages:
			log.Info().Any("packet", msg).Msg("received")

			player, exists := s.players[msg.PlayerId]
			if !exists {
				return
			}

			player.Send(&Packet{
				Type:    "echo",
				Payload: msg,
			})
		}
	}
}

func (s *Server) Register(player *Player) {
	s.register <- player
}

func (s *Server) Unregister(playerId PlayerId) {
	s.unregister <- playerId
}

func (s *Server) PushPacket(packet *ClientPacket) {
	s.messages <- packet
}
