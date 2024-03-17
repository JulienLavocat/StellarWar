package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"com.stellarwar/world-server/internal/galaxy"
	srv "com.stellarwar/world-server/internal/server"
	"github.com/olahol/melody"
)

var server = srv.NewServer(GenerateMap())

func main() {
	go server.Run()

	m := melody.New()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		m.HandleRequest(w, r)
	})

	m.HandleConnect(func(s *melody.Session) {
		player := server.CreatePlayer(srv.NewPeer(s))

		s.Set("id", player.Id)
	})

	m.HandleDisconnect(func(s *melody.Session) {
		if value, exists := s.Get("id"); exists {
			server.HandleDisconnect(value.(int32))
		}
	})

	m.HandleMessage(func(s *melody.Session, b []byte) {
		// TODO: Handle
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
	}
	log.Println("http: listening on port " + port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func Upgrade(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}

	player := server.CreatePlayer(c)

	go player.ReadPump()
	go player.WritePump()
}

func GenerateMap() galaxy.Galaxy {
	start := time.Now()
	generator := galaxy.GalaxyGenerator{
		Width:          200,
		Height:         200,
		SystemsDensity: 50,
		Seed:           "test",
	}
	log.Printf("Generated galaxy in %s", time.Since(start))

	return generator.Generate()
}
