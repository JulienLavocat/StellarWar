package main

import (
	"time"

	srv "github.com/StellarWar/world-server/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	log.Logger = log.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = "" + time.TimeOnly
	}))

	server := srv.NewServer()
	go server.Run()

	srv.NewWebsocketServer(server)
}
