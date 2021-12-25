package main

import (
	"fmt"
	"time"

	"go.uber.org/zap"
)

// main solves parts 1 and 2 using the above input
func main() {
	start := time.Now()

	player1 := player{
		Name: "player 1",
		Pos:  8,
	}

	player2 := player{
		Name: "player 2",
		Pos:  2,
	}

	p1 := part1(game{
		players: []*player{&player1, &player2},
		d:       new(d100),
	})

	end := time.Now()

	fmt.Println("part1:", p1)
	// fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

type player struct {
	Name  string
	Score int
	Pos   int
}

type die interface {
	roll() int
}

// d100 is a deterministic die that rolls from 1-100, in that order
type d100 int

var _ die = (*d100)(nil)

// roll returns a number between 1-100. The first roll will be 1,
// the second 2, and so on.
func (d *d100) roll() int {
	*d = ((*d) % 100) + 1
	return int(*d)
}

type game struct {
	players []*player
	d       die
	turn    int
}

func part1(g game, logs ...*zap.SugaredLogger) int {
	var done bool
	for !done {
		done = g.playTurn(logs...)
	}
	loser := g.loser()
	return loser.Score * g.turn * 3
}

func (g *game) playTurn(logs ...*zap.SugaredLogger) bool {
	g.turn++
	i := (g.turn - 1) % 2
	p := g.players[i]
	for _, log := range logs {
		log.Infow("start turn", "turn", g.turn, "player", p.Name)
	}
	r1, r2, r3 := g.d.roll(), g.d.roll(), g.d.roll()
	for _, log := range logs {
		log.Infow("dice rolled", "r1", r1, "r2", r2, "r3", r3)
	}
	for _, log := range logs {
		log.Infow("before move", "player", p)
	}
	p.advance(r1 + r2 + r3)
	p.Score += p.Pos
	for _, log := range logs {
		log.Infow("after move", "player", p)
	}
	return p.Score >= 1000
}

func (g game) loser() player {
	var lost *player

	for _, p := range g.players {
		if lost == nil || p.Score < lost.Score {
			lost = p
		}
	}

	return *lost
}

func (p *player) advance(n int) {
	p.Pos += n
	for p.Pos > 10 {
		p.Pos -= 10
	}
}
