package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/websocket"
)

type Move struct {
	Locations []int `json:"location"`
}

var idCounter = 0
var roomExists = make(map[int]bool)
var rooms = make(map[int]*websocket.Upgrader)
var chatParticipants = make(map[int][]*websocket.Conn)
var roomMoves = make(map[int][]Move)
var allRooms []int

func Websocket(w http.ResponseWriter, r *http.Request) {
	fmt.Println("here")
	r.ParseForm()
	roomId := r.FormValue("roomId")
	id, _ := strconv.Atoi(roomId)

	if roomId == "" {
		id := idCounter
		idCounter++
		roomExists[id] = true
		allRooms = append(allRooms, id)
		newUpgrader := &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		}

		newUpgrader.CheckOrigin = func(r *http.Request) bool { return true }

		rooms[id] = newUpgrader

		ws, err := newUpgrader.Upgrade(w, r, nil)

		if err != nil {
			fmt.Fprintf(w, "internal error")
		}

		chatParticipants[id] = append(chatParticipants[id], ws)

		fmt.Println("Successfully connected")
		ws.WriteMessage(1001, []byte("Hello"))

		reader(ws, id)
		//TODO create new websocket
	} else if roomExists[id] == true {
		//join to existing websocket
		upgrader := rooms[id]

		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			fmt.Fprintf(w, "internal error")
		}
		go reader(ws, id)
	} else {
		fmt.Fprintf(w, "no room found")
	}
}

func reader(conn *websocket.Conn, id int) {
	for {
		messageType, p, err := conn.ReadMessage()
		print(p)
		if err != nil {
			fmt.Println(err)
			return
		}

		var m Move
		fmt.Println(p)
		err = json.Unmarshal(p, &m)
		fmt.Println(m.Locations)

		roomMoves[id] = append(roomMoves[id], m)

		if err != nil {
			fmt.Println(err)
			return
		}

		marsheled, _ := json.Marshal(m)

		for i := 0; i < len(chatParticipants[id]); i++ {
			if chatParticipants[id][i] == conn {
				continue
			}
			if err = chatParticipants[id][i].WriteMessage(messageType, []byte(marsheled)); err != nil {
				fmt.Println(err)
			}
		}

	}
}
