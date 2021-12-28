package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/nealmcc/aoc2021/pkg/vector"
	"go.uber.org/zap"
)

// main solves both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	boot, after, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	r := part1(boot)
	p1 := r.numLit()

	r = part2(r, after)
	p2 := r.numLit()

	end := time.Now()

	fmt.Println("part1:", p1)
	fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// instruction is one command for the reactor to execute.
type instruction struct {
	on     bool
	x1, x2 int
	y1, y2 int
	z1, z2 int
}

// reactor is a collection of activated cells.
type reactor struct {
	cells map[vector.Cuboid]struct{}
}

func newReactor() *reactor {
	r := reactor{
		cells: make(map[vector.Cuboid]struct{}, 16),
	}
	return &r
}

// read the given input, returning an enhancement algorithm and an image.
func read(r io.Reader) (boot, after []instruction, err error) {
	boot = make([]instruction, 0, 20)
	after = make([]instruction, 0, 400)

	isBoot := func(i instruction) bool {
		in50 := func(n int) bool {
			return -50 <= n && n <= 50
		}
		return in50(i.x1) && in50(i.x2) &&
			in50(i.y1) && in50(i.y2) &&
			in50(i.z1) && in50(i.z2)
	}

	splitRange := func(data []byte) (int, int, error) {
		rangeSep := []byte{'.', '.'}
		parts := bytes.Split(data, rangeSep)
		if len(parts) != 2 {
			return 0, 0, errors.New("malformed input: expected 2 range values")
		}
		from, err := strconv.Atoi(string(parts[0]))
		if err != nil {
			return 0, 0, err
		}
		to, err := strconv.Atoi(string(parts[1]))
		if err != nil {
			return 0, 0, err
		}

		return from, to, nil
	}

	s := bufio.NewScanner(r)
	for s.Scan() {
		parts := bytes.Split(s.Bytes(), []byte{' '})
		if len(parts) != 2 {
			return nil, nil, errors.New("malformed input: expected opCode and params")
		}

		opCode, params := parts[0], parts[1]

		parts = bytes.Split(params, []byte{','})
		if len(parts) != 3 {
			return nil, nil, errors.New("malformed input: expected 3 ranges")
		}

		xs, ys, zs := parts[0], parts[1], parts[2]
		x1, x2, err := splitRange(xs[2:])
		if err != nil {
			return nil, nil, err
		}

		y1, y2, err := splitRange(ys[2:])
		if err != nil {
			return nil, nil, err
		}

		z1, z2, err := splitRange(zs[2:])
		if err != nil {
			return nil, nil, err
		}

		step := instruction{
			on: len(opCode) == 2,
			x1: x1, x2: x2, y1: y1, y2: y2, z1: z1, z2: z2,
		}

		if isBoot(step) {
			boot = append(boot, step)
		} else {
			after = append(after, step)
		}
	}

	if err := s.Err(); err != nil {
		return nil, nil, err
	}

	return boot, after, nil
}

// part1 solves part 1.
func part1(boot []instruction, logs ...*zap.SugaredLogger) *reactor {
	r := newReactor()

	for _, act := range boot {
		act.Do(r, logs...)
	}
	return r
}

// part2 solves part 2.
func part2(r *reactor, after []instruction, logs ...*zap.SugaredLogger) *reactor {
	sort.Slice(after, func(i, j int) bool {
		return after[i].on && !after[j].on
	})

	for _, act := range after {
		act.Do(r, logs...)
	}
	return r
}

// Do performs this instruction on the reactor, and uses the optional logger(s)
// to print debugging information.
func (act instruction) Do(r *reactor, logs ...*zap.SugaredLogger) {
	infow := func(msg string, keysAndValues ...interface{}) {
		for _, l := range logs {
			l.Infow(msg, keysAndValues...)
		}
	}

	infow("beginning action", "numCubes",
		(act.x2-act.x1)*(act.y2-act.y1)*(act.z2-act.z1))

	start := time.Now()

	box := vector.Cuboid{
		X1: act.x1, X2: act.x2,
		Y1: act.y1, Y2: act.y2,
		Z1: act.z1, Z2: act.z2,
	}
	if act.on {
		r.set(box)
	} else {
		r.unset(box)
	}

	end := time.Now()
	infow("finished action", "time spent", end.Sub(start))
}

// numLit returns the number of lit cells within the reactor.
func (r reactor) numLit() int {
	var sum int
	for box := range r.cells {
		sum += box.Volume()
	}
	return sum
}

func (r *reactor) set(newCell vector.Cuboid) {
	keep := make([]vector.Cuboid, 16)

	for cell := range r.cells {
		overlap, ok := cell.Intersect(newCell)
		if !ok {
			continue
		}
		delete(r.cells, cell)
		boxes := cell.Remove(overlap)
		for _, b := range boxes {
			keep = append(keep, b)
		}
	}
	r.cells[newCell] = struct{}{}
	for _, cell := range keep {
		r.cells[cell] = struct{}{}
	}
}

func (r *reactor) unset(box vector.Cuboid) {
}
