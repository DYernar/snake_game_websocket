package main

import (
	"log"
	"snake_game/game"

	"github.com/alehano/wsgame/utils"
)

var allRooms = make(map[string]*room)
var freeRooms = make(map[string]*room)
var roomsCount int

type room struct {
	name        string
	playerConns map[*playerConn]bool
	updateAll   chan bool
	join        chan *playerConn
	leave       chan *playerConn
}

func (r *room) run() {
	for {
		select {
		case c := <-r.join:
			r.playerConns[c] = true
			r.updateAllPlayers()
			if len(r.playerConns) == 2 {
				delete(freeRooms, r.name)
				var p []*game.Player
				for k, _ := range r.playerConns {
					p = append(p, k.Player)
				}
				game.PairPlayers(p[0], p[1])
			}
		case c := <-r.leave:
			c.Player.GiveUp()
			r.updateAllPlayers()
			delete(r.playerConns, c)
			if len(r.playerConns) == 0 {
				goto Exit
			}
		case <-r.updateAll:
			r.updateAllPlayers()
		}
	}

Exit:
	delete(allRooms, r.name)
	delete(freeRooms, r.name)
	roomsCount -= 1
	log.Print("room closed: ", r.name)
}

func (r *room) updateAllPlayers() {
	for c := range r.playerConns {
		c.sendState()
	}
}

func NewRoom(name string) *room {
	if name == "" {
		name = utils.RandString(16)
	}

	room := &room{
		name:        name,
		playerConns: make(map[*playerConn]bool),
		updateAll:   make(chan bool),
		join:        make(chan *playerConn),
		leave:       make(chan *playerConn),
	}

	allRooms[name] = room
	freeRooms[name] = room
	roomsCount++
	go room.run()
	return room
}
