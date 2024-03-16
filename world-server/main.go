package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"com.stellarwar/world-server/galaxy"
	srv "com.stellarwar/world-server/server"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}
var server = srv.NewServer(GenerateMap())

func main() {
	http.HandleFunc("/", Upgrade)
	go server.Run()

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
