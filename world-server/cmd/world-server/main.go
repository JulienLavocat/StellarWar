package main

import (
	"time"

	"github.com/StellarWar/world-server/internal/galaxy"
	srv "github.com/StellarWar/world-server/internal/server"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMicro
	log.Logger = log.Output(zerolog.NewConsoleWriter(func(w *zerolog.ConsoleWriter) {
		w.TimeFormat = "" + time.TimeOnly
	}))

	server := srv.NewServer(GenerateMap())
	go server.Run()

	srv.NewWebsocketServer(server)
}

func GenerateMap() galaxy.Galaxy {
	start := time.Now()
	generator := galaxy.GalaxyGenerator{
		Width:          200,
		Height:         200,
		SystemsDensity: 50,
		Seed:           "test",
	}
	log.Debug().Msgf("Generated galaxy in %s", time.Since(start))

	return generator.Generate()
}
