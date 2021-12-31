package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
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

	seafloor, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(seafloor)

	end := time.Now()

	fmt.Println("part1:", p1)
	// fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read the given input, returning a map of the seafloor.
func read(r io.Reader) (seafloor, error) {
	s := bufio.NewScanner(r)

	sea := seafloor{
		cukes: make([][]byte, 0, 137),
	}

	for s.Scan() {
		in := s.Bytes()
		buf := make([]byte, len(in))
		copy(buf, in)
		sea.cukes = append(sea.cukes, buf)
	}
	if err := s.Err(); err != nil {
		return seafloor{}, err
	}

	sea.height = len(sea.cukes)
	if sea.height > 0 {
		sea.width = len(sea.cukes[0])
	}

	return sea, nil
}

type seafloor struct {
	step   int
	width  int
	height int
	cukes  [][]byte
}

const (
	down  byte = 'v'
	right byte = '>'
	empty byte = '.'
)

func part1(s seafloor) int {
	for {
		numMoved := s.move()
		if numMoved == 0 {
			break
		}
	}
	return s.step
}

// move performs one step of movement, and returns the number of cucumbers that
// moved.
func (s *seafloor) move() int {
	n := s.moveRight()
	n += s.moveDown()
	s.step++
	return n
}

// moveDown shifts the horizontal herd of sea cucumbers right wherever they can,
// and returns the number of cucumbers that moved.
func (s *seafloor) moveRight() int {
	movers := make(map[int][]int, 16)

	for row := 0; row < s.height; row++ {
		for col := 0; col < s.width; col++ {
			if s.cukes[row][col] != right {
				continue
			}
			if s.cukes[row][(col+1)%s.width] == empty {
				movers[row] = append(movers[row], col)
			}
		}
	}

	for row, cols := range movers {
		for _, col := range cols {
			s.cukes[row][col] = empty
			s.cukes[row][(col+1)%s.width] = right
		}
	}

	return len(movers)
}

// moveDown shifts the vertical herd of sea cucumbers down wherever they can,
// and returns the number of cucumbers that moved.
func (s *seafloor) moveDown() int {
	movers := make(map[int][]int, 16)

	for row := 0; row < s.height; row++ {
		for col := 0; col < s.width; col++ {
			if s.cukes[row][col] != down {
				continue
			}
			if s.cukes[(row+1)%s.height][col] == empty {
				movers[row] = append(movers[row], col)
			}
		}
	}

	for row, cols := range movers {
		for _, col := range cols {
			s.cukes[row][col] = empty
			s.cukes[(row+1)%s.height][col] = down
		}
	}

	return len(movers)
}
