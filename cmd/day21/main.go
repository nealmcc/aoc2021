package main

import (
	"fmt"
	"time"
)

// main solves parts 1 and 2 using the above input
func main() {
	start := time.Now()

	p1 := part1(game{
		P1:    player{Name: "player 1", Pos: 8},
		P2:    player{Name: "player 2", Pos: 2},
		ToWin: 1000,
	})

	p2 := part2(game{
		P1:    player{Name: "player 1", Pos: 8},
		P2:    player{Name: "player 2", Pos: 2},
		ToWin: 21,
	})

	end := time.Now()

	fmt.Println("part1:", p1)
	fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// part1 solves part 1 of the problem.
func part1(g game) int {
	d := new(d100)
	var done bool
	for !done {
		sum := d.roll() + d.roll() + d.roll()
		done = g.playTurn(sum)
	}
	loser := g.loser()
	return loser.Score * g.Turn * 3
}

// part2 solves part 2 of the problem.
func part2(g game) int {
	m := &multiverse{
		done:    make(map[game]int, 1024),
		playing: map[game]int{g: 1},
	}

	var done bool
	for !done {
		done = m.playTurn()
	}

	gamesWon := make(map[string]int, 2)
	for game, count := range m.done {
		gamesWon[game.winner().Name] += count
	}

	var best int
	for _, total := range gamesWon {
		if total > best {
			best = total
		}
	}

	return best
}

// d100 is a deterministic die that rolls from 1-100, in that order
type d100 int

// roll returns a number between 1-100. The first roll will be 1,
// the second 2, and so on.
func (d *d100) roll() int {
	*d = ((*d) % 100) + 1
	return int(*d)
}

// player stores the state of a player in a game.
type player struct {
	Name  string
	Score int
	Pos   int
}

// advance moves this player forward the given number of squares.
func (p *player) advance(n int) {
	p.Pos += n % 10
	if p.Pos > 10 {
		p.Pos -= 10
	}
}

// clone creates a deep copy of this player
func (p player) clone() player {
	return player{
		Name:  p.Name,
		Score: p.Score,
		Pos:   p.Pos,
	}
}

// game is the state of a game.  It is suitable to use as a key in a map.
type game struct {
	P1    player
	P2    player
	Turn  int
	ToWin int
}

// playTurn plays the next player's turn using the given dice roll, and returns
// true if the game is over.
func (g *game) playTurn(sum int) bool {
	g.Turn++

	p := &g.P1
	if g.Turn%2 == 0 {
		p = &g.P2
	}
	p.advance(sum)

	p.Score += p.Pos
	return p.Score >= g.ToWin
}

// winner returns this game's winning player.
func (g game) winner() player {
	if g.P1.Score > g.P2.Score {
		return g.P1
	}
	return g.P2
}

// loser returns this game's losing player.
func (g game) loser() player {
	if g.P1.Score < g.P2.Score {
		return g.P1
	}
	return g.P2
}

// clone creates a deep copy of this game.
func (g game) clone() game {
	return game{
		P1:    g.P1.clone(),
		P2:    g.P2.clone(),
		Turn:  g.Turn,
		ToWin: g.ToWin,
	}
}

// multiverse tracks the number of universes that contain each game state.
type multiverse struct {
	done    map[game]int // done keeps track completed games
	playing map[game]int // playing keeps track of in-progress games
}

// playTurn causes each incomplete game in the multiverse to play 1 turn.
// Returns true if all games are complete.
func (m *multiverse) playTurn() bool {
	next := make(map[game]int, 16)

	for g, count := range m.playing {
		for sum, chance := range _3d3Chance {
			g := g.clone()
			gameOver := g.playTurn(sum)
			if gameOver {
				m.done[g] += count * int(chance)
			} else {
				next[g] += count * int(chance)
			}
		}
	}

	m.playing = next
	return len(next) == 0
}

// _3d3Chance is a map from the sum of three 3-sided dice to the
// probability (out of 27) of rolling that sum.
var _3d3Chance = map[int]int{
	3: 1, // 1,1,1
	4: 3, // 1,1,2 (x3)
	5: 6, // 1,1,3 (x3) and 1,2,2 (x3)
	6: 7, // 1,2,3 (x6) and 2,2,2
	7: 6, // 2,2,3 (x3) and (1,3,3) (x3)
	8: 3, // 3,3,2 (x3)
	9: 1, // 3,3,3
}
