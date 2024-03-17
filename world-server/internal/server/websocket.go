package server

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/StellarWar/world-server/internal/packets"
	"github.com/StellarWar/world-server/internal/shared"
	"github.com/olahol/melody"
	"github.com/rs/zerolog/log"
)

const playerId = "id"

type WebsocketServer struct {
	server *Server
	m      *melody.Melody
}

func NewWebsocketServer(server *Server) {
	wsServer := &WebsocketServer{
		m:      melody.New(),
		server: server,
	}

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		wsServer.m.HandleRequest(w, r)
	})

	wsServer.m.HandleConnect(wsServer.handleConnect)
	wsServer.m.HandleClose(wsServer.handleClose)
	wsServer.m.HandleMessage(
		wsServer.handleMessage)

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Info().Msgf("http: listening on port: %s", port)
	log.Fatal().Err(http.ListenAndServe(":"+port, nil))
}

func (ws *WebsocketServer) handleConnect(s *melody.Session) {
	log.Info().Msgf("Received ws session %s", s.RemoteAddr())

	player := ws.server.CreatePlayer(newNetPeer(s))
	s.Set(playerId, player.Id)
}

func (ws *WebsocketServer) handleClose(s *melody.Session, code int, reason string) error {
	log.Info().Msgf("Closed ws session %s", s.RemoteAddr())
	ws.server.PushPacket(packets.NewUnregisterPacket(s.MustGet(playerId).(shared.PlayerId)))
	return nil
}

func (ws *WebsocketServer) handleMessage(s *melody.Session, msg []byte) {
	var packet shared.ClientPacket
	if err := json.Unmarshal(msg, &packet); err != nil {
		log.Debug().Err(err).Msg("packet decoding failed")
		return
	}

	packet.PlayerId = s.MustGet(playerId).(shared.PlayerId)
	ws.server.PushPacket(&packet)
}
