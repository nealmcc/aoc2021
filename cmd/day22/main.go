package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	v "github.com/nealmcc/aoc2021/pkg/vector"
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

	r := newReactor()

	for _, op := range boot {
		r.Do(op)
	}
	p1 := r.numLit()

	for _, op := range after {
		r.Do(op)
	}
	p2 := r.numLit()

	end := time.Now()

	fmt.Println("part1:", p1)
	fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// reactor is a collection of activated cells.
type reactor struct {
	cells map[v.Cuboid]struct{}
}

// instruction is one command for the reactor to execute.
type instruction struct {
	isAdd  bool
	x1, x2 int
	y1, y2 int
	z1, z2 int
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
			isAdd: len(opCode) == 2,
			x1:    x1, x2: x2, y1: y1, y2: y2, z1: z1, z2: z2,
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

func newReactor() *reactor {
	r := reactor{
		cells: make(map[v.Cuboid]struct{}, 16),
	}
	return &r
}

// Do performs the given instruction (either add or remove) on this reactor.
func (r *reactor) Do(op instruction) {
	block := v.Cuboid{
		X1: op.x1, X2: op.x2,
		Y1: op.y1, Y2: op.y2,
		Z1: op.z1, Z2: op.z2,
	}.Normal()

	block.X2++
	block.Y2++
	block.Z2++

	r.remove(block)

	if op.isAdd {
		r.cells[block] = struct{}{}
	}
}

// remove the given block from this reactor.
func (r *reactor) remove(block v.Cuboid) {
	keep := make([]v.Cuboid, 0, len(r.cells))
	for curr := range r.cells {
		_, left, _ := v.CuboidOuterJoin(curr, block)
		delete(r.cells, curr)
		keep = append(keep, left...)
	}

	for _, k := range keep {
		r.cells[k] = struct{}{}
	}
}

// numLit returns the number of active cells.
func (r reactor) numLit() int {
	var sum int
	for box := range r.cells {
		sum += box.Volume()
	}
	return sum
}
