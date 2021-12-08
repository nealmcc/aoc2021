package main

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	_, err = read(in)
	if err != nil {
		log.Fatal(err)
	}
}

type digit int

const (
	zero digit = iota
	one
	two
	three
	four
	five
	six
	seven
	eight
	nine
)

// display is one row of input, with 10 unique signals and a 4-digit output.
// The digits map is used to associate each unique signal with its digit
type display struct {
	signals []string
	digits  map[string]digit
	output  []string
}

func read(r io.Reader) ([]display, error) {
	s := bufio.NewScanner(r)

	displays := make([]display, 0, 16)
	for s.Scan() {
		d := display{digits: make(map[string]digit)}
		parts := strings.Split(s.Text(), " | ")
		d.signals = strings.Split(parts[0], " ")
		d.output = strings.Split(parts[1], " ")
		displays = append(displays, d)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return displays, nil
}
