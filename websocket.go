package main

import (
	"log"
	"net/http"
	"net/url"
	"snake_game/game"

	"github.com/gorilla/websocket"
)

func Websocket(w http.ResponseWriter, r *http.Request) {
	ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)

	if _, ok := err.(*websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		return
	}

	playerName := "Player"
	params, _ := url.ParseQuery(r.URL.RawQuery)
	if len(params["name"]) > 0 {
		playerName = params["name"][0]
	}

	var room *room

	if len(freeRooms) > 0 {
		for _, r := range freeRooms {
			room = r
			break
		}
	} else {
		room = NewRoom("")
	}

	player := game.NewPlayer(playerName)
	pConn := NewPlayerConn(ws, player, room)
	room.join <- pConn

	log.Printf("Player: %s has joined to room %s", pConn.Name, room.name)

}
