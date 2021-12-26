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
)

// main solves both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	boot, _, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(boot)

	end := time.Now()

	fmt.Println("part1:", p1)
	// fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// instruction is one command for the reactor to execute.
type instruction struct {
	on     bool
	x1, x2 int
	y1, y2 int
	z1, z2 int
}

// reactor is a set of 3d coordinates
type reactor map[int]map[int]map[int]struct{}

func newReactor() *reactor {
	r := make(reactor, 101)
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
func part1(boot []instruction) int {
	r := newReactor()

	for _, act := range boot {
		act.Do(r)
	}
	return r.numLit()
}

func (act instruction) Do(r *reactor) {
	for x := act.x1; x <= act.x2; x++ {
		for y := act.y1; y <= act.y2; y++ {
			for z := act.z1; z <= act.z2; z++ {
				if act.on {
					r.set(x, y, z)
				} else {
					r.unset(x, y, z)
				}
			}
		}
	}
}

func (r reactor) numLit() int {
	var sum int
	for _, ys := range r {
		for _, zs := range ys {
			sum += len(zs)
		}
	}
	return sum
}

func (r *reactor) set(x, y, z int) {
	if (*r)[x] == nil {
		(*r)[x] = make(map[int]map[int]struct{}, 101)
	}

	if (*r)[x][y] == nil {
		(*r)[x][y] = make(map[int]struct{}, 101)
	}

	(*r)[x][y][z] = struct{}{}
}

func (r *reactor) unset(x, y, z int) {
	ys, ok := (*r)[x]
	if !ok {
		return
	}

	zs, ok := ys[y]
	if !ok {
		return
	}

	delete(zs, z)
}
