package main

import (
	"bufio"
	"container/heap"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/nealmcc/aoc2021/pkg/queue"
	"github.com/nealmcc/aoc2021/pkg/vector"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	cave, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(cave)
	p2 := part2(cave)

	end := time.Now()

	fmt.Println("part1:", p1)
	fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read a cave from the given input
func read(r io.Reader) (cavern, error) {
	s := bufio.NewScanner(r)

	if !s.Scan() {
		return cavern{}, errors.New("empty input")
	}
	if err := s.Err(); err != nil {
		return cavern{}, err
	}

	row := s.Bytes()
	cave := newCavern(len(row))

	setRisk := func(x, y int, b byte) {
		cave.risk[vector.Coord{X: x, Y: y}] = int(b - '0')
	}

	for x := 0; x < len(row); x++ {
		setRisk(x, 0 /* y */, row[x])
	}

	for y := 1; s.Scan(); y++ {
		row := s.Bytes()
		for x := 0; x < len(row); x++ {
			setRisk(x, y, row[x])
		}
	}
	if err := s.Err(); err != nil {
		return cavern{}, err
	}

	return cave, nil
}

// cavern is a room with coordinates ranging from (0, 0) at the top left
// to (width-1, width-1) at the bottom right.
// Each position in the cave has a variable risk level from 0 to 9.
type cavern struct {
	width int
	risk  map[vector.Coord]int
}

func newCavern(width int) cavern {
	return cavern{
		width: width,
		risk:  make(map[vector.Coord]int, width*width),
	}
}

func part1(c cavern) int {
	var (
		from = vector.Coord{X: 0, Y: 0}
		to   = vector.Coord{X: c.width - 1, Y: c.width - 1}
	)
	cost, _ := c.shortestPath(from)
	return cost[to]
}

func part2(c cavern) int {
	c.scale(5)
	return part1(c)
}

// scale increases the size of the cavern by n in both the x and y directions.
// Adjusts the risk factor of each position according to the rules in part 2.
func (c *cavern) scale(factor int) {
	template := make(map[vector.Coord]int, len(c.risk))
	for k, v := range c.risk {
		template[k] = v
	}

	limit := c.width * factor
	for x := 0; x < limit; x++ {
		for y := 0; y < limit; y++ {
			k1 := vector.Coord{X: x % c.width, Y: y % c.width}
			k2 := vector.Coord{X: x, Y: y}
			risk := template[k1] + tileScore(k2, c.width)
			for risk >= 10 {
				risk -= 9
			}
			c.risk[k2] = risk
		}
	}

	c.width *= factor
}

// shortestPath computes the shortest distance and route from start to all
// other squares in the cavern. The cavern is assumed to be square.
// Uses Dijkstra's Algorithm.
//
// dist holds the minimal distance from the start to the given coordinate.
//
// prev holds the predecessor of each coordinate, when visiting it from the start.
func (c cavern) shortestPath(start vector.Coord) (cost map[vector.Coord]int, prev map[vector.Coord]vector.Coord) {
	cost = make(map[vector.Coord]int, len(c.risk))
	prev = make(map[vector.Coord]vector.Coord, len(c.risk))
	cost[start] = 0

	// push each node on the graph into the queue, with an initial distance
	q := new(queue.Coord)
	// save a pointer to each node, so we can update it in the queue
	pointers := make([]*queue.CoordNode, c.width*c.width)
	const infinity = 1<<63 - 1
	for pos := range c.risk {
		if pos != start {
			cost[pos] = infinity
		}
		node := &queue.CoordNode{Value: pos, Priority: -1 * cost[pos]}
		pointers[pos.X*c.width+pos.Y] = node
		heap.Push(q, node)
	}

	// the first node that we definitely know the cost of is the start. It has
	// the highest priority (0 - all others are negative) so we pop it off the
	// queue, and find the cost of all adjacent nodes.  We keep working out from
	// the next nearest node until all nodes have a defined cost.
	for q.Len() > 0 {
		node := heap.Pop(q).(*queue.CoordNode)
		curr, cumulativeRisk := node.Value, node.Priority*-1

		for _, next := range c.neighbours(curr) {
			nextRisk, ok := c.risk[next]
			if !ok {
				continue
			}
			// the distance from start to next, if we arrive via curr:
			alt := cumulativeRisk + nextRisk
			if alt < cost[next] {
				cost[next] = alt
				p := pointers[next.X*c.width+next.Y]
				q.Update(p, p.Value, -1*alt)
				// also save the route to get here:
				prev[next] = curr
			}
		}
	}

	return cost, prev
}

// neighbours returns a slice of all positions within the cavern that
// are adjacent to p.
// Squares are only adjacent vertically and horizontally - not diagonally.
func (c cavern) neighbours(p vector.Coord) []vector.Coord {
	return []vector.Coord{
		{X: p.X, Y: p.Y - 1},
		{X: p.X, Y: p.Y + 1},
		{X: p.X - 1, Y: p.Y},
		{X: p.X + 1, Y: p.Y},
	}
}

func (c cavern) Format(f fmt.State, verb rune) {
	// create a buffer to write values into
	display := make([][]byte, c.width)
	for i := 0; i < c.width; i++ {
		display[i] = make([]byte, c.width)
	}

	// write risk values into the buffer
	for pos, risk := range c.risk {
		display[pos.Y][pos.X] = '0' + byte(risk)
	}

	// print the buffer
	for _, row := range display {
		f.Write(row)
		f.Write([]byte{'\n'})
	}
}

// tileScore finds the incremental amount that a given position's risk
// needs to be adjusted by for part 2
func tileScore(pos vector.Coord, width int) int {
	return pos.X/width + pos.Y/width
}
