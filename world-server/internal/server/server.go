package server

import (
	"encoding/json"
	"log"
	"sync"

	"com.stellarwar/world-server/internal/galaxy"
	"com.stellarwar/world-server/internal/packets"
	"github.com/olahol/melody"
)

var (
	mu           sync.Mutex
	nextPlayerId = int32(0)
)

type (
	CommandHandler func(server *Server, client *Player, packet Packet)
	ClientMessage  struct {
		message []byte
		id      int32
	}
)

type Server struct {
	m *melody.Melody

	Galaxy galaxy.Galaxy

	players  map[int32]*Player
	commands map[string]CommandHandler

	inMessages chan ClientMessage
	unregister chan int32
	register   chan *Player
}

func NewServer(galaxy galaxy.Galaxy) Server {
	server := Server{
		Galaxy:     galaxy,
		players:    make(map[int32]*Player),
		inMessages: make(chan ClientMessage),
		unregister: make(chan int32),
		register:   make(chan *Player),
	}

	server.commands = map[string]CommandHandler{
		"set_ready":    SetReadyhandler,
		"claim_system": ClaimStarSystem,
	}

	return server
}

func (s *Server) CreatePlayer(peer *Peer) *Player {
	var playerId int32
	mu.Lock()
	playerId = nextPlayerId
	nextPlayerId++
	mu.Unlock()

	player := NewPlayer(playerId, s, peer)
	s.players[player.Id] = player
	s.register <- player
	return player
}

func (s *Server) HandleDisconnect(playerId int32) {
	s.unregister <- playerId
}

func (s *Server) BroadcastToAll(packet Packet) {
	result, err := packet.Marshall()
	if err != nil {
		log.Fatal(err)
	}
	s.m.Broadcast(result)
}

func (s *Server) BroadcastToOthers(sender int32, packet Packet) {
	for _, p := range s.players {
		if p.Id == sender {
			continue
		}
		p.send <- packet
	}
}

func (s *Server) Run() {
	for {
		select {
		case player := <-s.register:
			player.Logf("connected (online players: %d)", len(s.players))
			s.BroadcastToOthers(player.Id, Packet{
				Command: "player_joined",
				Data: packets.PublicClientPayload{
					Id:             player.Id,
					Color:          player.Color,
					HomeSystem:     player.HomeSystem,
					ClaimedSystems: player.ClaimedSystems,
				},
			})
		case playerId := <-s.unregister:
			if p, ok := s.players[playerId]; ok {
				p.Log("disconnected")

				for _, systemId := range p.ClaimedSystems {
					s.Galaxy.Systems[systemId].Owner = -1
				}
				s.BroadcastToAll(Packet{
					Command: "player_left",
					Data: packets.PlayerLeftPayload{
						Id:           playerId,
						OwnedSystems: p.ClaimedSystems,
					},
				})

				close(s.players[playerId].send)
				delete(s.players, playerId)
			}
		case msg := <-s.inMessages:
			var packet Packet
			err := json.Unmarshal(msg.message, &packet)
			if err != nil {
				continue
			}
			player := s.players[msg.id]
			player.Log("recv:", packet.Command)
			handler, hasHandler := s.commands[packet.Command]
			if !hasHandler {
				player.send <- Packet{
					Command: packet.Command,
					Data: map[string]interface{}{
						"error": "invalid command",
					},
				}
				continue
			}
			handler(s, player, packet)
		}
	}
}

func SetReadyhandler(s *Server, player *Player, packet Packet) {
	players := make(map[int32]packets.PublicClientPayload)
	for i, player := range s.players {
		players[i] = packets.PublicClientPayload{
			Id:             player.Id,
			Color:          player.Color,
			HomeSystem:     player.HomeSystem,
			ClaimedSystems: player.ClaimedSystems,
		}
	}

	player.peer.Send(Packet{
		Command: "sync_client",
		Data: packets.SyncClientPayload{
			Galaxy:  s.Galaxy,
			Players: players,
			Self: packets.PublicClientPayload{
				Id:             player.Id,
				Color:          player.Color,
				HomeSystem:     player.HomeSystem,
				ClaimedSystems: player.ClaimedSystems,
			},
		},
	})
}

func ClaimStarSystem(s *Server, player *Player, packet Packet) {
	systemIdFloat, ok := packet.Data.(float64)
	if !ok {
		return
	}
	systemId := int(systemIdFloat)

	if systemId < 0 || systemId >= int(len(s.Galaxy.Systems)) {
		return
	}

	s.Galaxy.Systems[systemId].Owner = player.Id
	player.ClaimSystem(systemId)

	s.BroadcastToAll(Packet{
		Command: "system_claimed",
		Data: packets.ClaimSystemPayload{
			SystemId: systemId,
			PlayerId: player.Id,
		},
	})
}
