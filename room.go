package main

import (
	"log"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
)

type Room struct {
	ID         int                 `json:"id"`
	Upgrader   *websocket.Upgrader `json:"upgrader"`
	Name       string              `json:"name"`
	Players    map[*Player]bool    `json:"players"`
	PlayerData chan Player         `json:"move"`
	Food       []int               `json:"food"`
	Loser1     *Player             `json:"loser1"`
	Loser2     *Player             `json:"loser2"`
}

type Player struct {
	Name      string  `json:"name"`
	Position  [][]int `json:"position"`
	Direction string  `json:"direction"`
	Ws        *websocket.Conn
	Status    int   `json:"status"`
	Food      []int `json:"food"`
}

func (r *Room) startGame() {
	if len(r.Players) < 2 {
		r.PlayerData <- *&Player{Status: 0}
	} else {
		r.generateFood()
		ticker := time.NewTicker(250 * time.Millisecond)
		done := make(chan bool)

		go func() {
			for {
				select {
				case <-done:
					return
				case <-ticker.C:
					lost := 0
					for p := range r.Players {
						p.Status = 1
						if r.Loser1 != p && r.Loser2 != p {
							p.move(r)
						} else {
							lost++
						}
						r.PlayerData <- *p

					}
					if lost == len(r.Players) {
						ticker.Stop()
					}

				}
			}
		}()
	}

}

func (r *Room) generateFood() {
	r.Food = []int{rand.Intn(20), rand.Intn(40)}
}

func (p *Player) changeDirection(direction string) {
	switch direction {
	case "down":
		if p.Direction != "up" {
			p.Direction = direction
		}
	case "up":
		if p.Direction != "down" {
			p.Direction = direction
		}
	case "left":
		if p.Direction != "right" {
			p.Direction = direction
		}
	case "right":
		if p.Direction != "left" {
			p.Direction = direction
		}

	}
}

func (p *Player) move(r *Room) {
	switch p.Direction {
	case "up":
		if p.Position[0][1]-1 < 0 {
			newPos := []int{p.Position[0][0], 39}
			p.Position = append([][]int{newPos}, p.Position...)
		} else {
			newPos := []int{p.Position[0][0], p.Position[0][1] - 1}
			p.Position = append([][]int{newPos}, p.Position...)
		}

	case "down":
		if p.Position[0][1]+1 > 39 {
			newPos := []int{p.Position[0][0], 0}
			p.Position = append([][]int{newPos}, p.Position...)
		} else {
			newPos := []int{p.Position[0][0], p.Position[0][1] + 1}
			p.Position = append([][]int{newPos}, p.Position...)
		}

	case "left":
		if p.Position[0][0]-1 < 0 {
			newPos := []int{19, p.Position[0][1]}
			p.Position = append([][]int{newPos}, p.Position...)
		} else {
			newPos := []int{p.Position[0][0] - 1, p.Position[0][1]}
			p.Position = append([][]int{newPos}, p.Position...)
		}

	case "right":
		if p.Position[0][0]+1 > 19 {
			newPos := []int{0, p.Position[0][1]}
			p.Position = append([][]int{newPos}, p.Position...)

		} else {
			newPos := []int{p.Position[0][0] + 1, p.Position[0][1]}
			p.Position = append([][]int{newPos}, p.Position...)

		}

	}

	if p.Position[0][0] != r.Food[0] || p.Position[0][1] != r.Food[1] {
		p.Position = p.Position[:len(p.Position)-1]
	} else {
		r.generateFood()
	}
	if p.Lost(r) {
		if r.Loser1 != nil {
			r.Loser2 = p
		} else {
			r.Loser1 = p
		}
	}
}

func (p *Player) Lost(r *Room) bool {
	for i := 1; i < len(p.Position); i++ {
		if p.Position[0][0] == p.Position[i][0] && p.Position[0][1] == p.Position[i][1] {
			return true
		}
	}
	enemy := &Player{}

	for e := range r.Players {
		if e != p {
			enemy = e
			break
		}
	}

	for i := 0; i < len(enemy.Position); i++ {
		if p.Position[0][0] == enemy.Position[i][0] && p.Position[0][1] == enemy.Position[i][1] {
			return true
		}
	}
	return false
}

func (r *Room) handleMoves() {
	for {

		playerD := <-r.PlayerData
		if len(r.Players) != 2 {
			playerD.Status = 0
		} else {
			playerD.Status = 1
		}

		if r.Loser1 != nil && r.Loser2 != nil {
			playerD.Status = 2
		}
		playerD.Food = r.Food
		for player := range r.Players {
			if player.Ws == playerD.Ws {
				playerD.Name = "player"
			} else {
				playerD.Name = "enemy"
			}
			err := player.Ws.WriteJSON(&playerD)
			if err != nil {
				log.Printf("error: %s", err)
				player.Ws.Close()
				delete(r.Players, player)
			}
		}
		if r.Loser1 != nil && r.Loser2 != nil {
			for p := range r.Players {
				p.Ws.Close()
			}
			delete(rooms, r.ID)
			break
		}

	}
}

func GenerateNewRoom(id int, name string) *Room {
	players := make(map[*Player]bool)
	positions := make(chan Player)
	return &Room{ID: id, Name: name, Upgrader: &websocket.Upgrader{}, Players: players, PlayerData: positions}
}

type Direction struct {
	NewDirection string `json:"direction"`
}

func (room *Room) handleConnection(p *Player) {
	for {

		var dir Direction
		err := p.Ws.ReadJSON(&dir)
		if err != nil {
			log.Printf("error: %v", err)
			delete(room.Players, p)
			break
		}
		p.changeDirection(dir.NewDirection)
	}
}
