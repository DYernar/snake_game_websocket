package main

import (
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/rs/cors"
)

var rooms = make(map[int]*Room)
var roomID = 0

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "7070"
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		w.WriteHeader(200)
		w.Write([]byte("hello world"))
	})

	mux.HandleFunc("/ws", Websocket)

	handler := cors.Default().Handler(mux)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal("Listen and serve err: ", err)
	}
}

func Websocket(w http.ResponseWriter, r *http.Request) {
	newRoom := &Room{}
	free := false
	for _, r := range rooms {
		if len(r.Players) == 1 {
			newRoom = r
			free = true
			break
		}
	}
	position := [][]int{{3, 5}, {3, 6}}
	if !free {
		strID := strconv.Itoa(roomID)
		newRoom = GenerateNewRoom(roomID, "room "+strID)
		rooms[roomID] = newRoom
		position = [][]int{{6, 5}, {6, 6}}
		go newRoom.handleMoves()
	}

	ws, _ := newRoom.Upgrader.Upgrade(w, r, nil)
	p := &Player{Ws: ws, Direction: "down", Position: position}
	f := []int{9, 9}
	newRoom.Players[p] = true
	newRoom.Food = f
	go newRoom.handleConnection(p)
	newRoom.startGame()

	// go newRoom.handleConnection(ws)
}
