package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	srv "github.com/StellarWar/world-server/internal/server"
	"github.com/olahol/melody"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const playerId = "id"

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	log.Logger = log.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = "" + time.TimeOnly
	}))

	server := srv.NewServer()
	go server.Run()

	m := melody.New()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	m.HandleConnect(func(s *melody.Session) {
		log.Info().Msgf("Received ws session %s", s.RemoteAddr())

		player := server.CreatePlayer(srv.NewNetPeer(s))
		s.Set(playerId, player.Id)

		server.Register(player)
	})

	m.HandleClose(func(s *melody.Session, code int, reason string) error {
		log.Info().Msgf("Closed ws session %s", s.RemoteAddr())
		server.Unregister(s.MustGet(playerId).(srv.PlayerId))
		return nil
	})

	m.HandleMessage(func(s *melody.Session, msg []byte) {
		var packet srv.ClientPacket
		if err := json.Unmarshal(msg, &packet); err != nil {
			log.Debug().Err(err).Msg("packet decoding failed")
			return
		}

		packet.PlayerId = s.MustGet(playerId).(srv.PlayerId)
		server.PushPacket(&packet)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Info().Msgf("http: listening on port: %s", port)
	log.Fatal().Err(http.ListenAndServe(":"+port, nil))
}
