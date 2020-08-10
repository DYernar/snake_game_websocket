package game

import "log"

type Player struct {
	Name     string
	Enemy    *Player
	Position []int
}

func NewPlayer(name string) *Player {
	return &Player{Name: name}
}

func PairPlayers(p1 *Player, p2 *Player) {
	p1.Enemy, p2.Enemy = p2, p1
}

func (p *Player) Command(command string) {
	log.Print("Command: '", command, "'received by player ", p.Name)
}

func (p *Player) GetState() []int {
	return p.Position
}

func (p *Player) GiveUp() {
	log.Print("Player gave up: ", p.Name)
}
