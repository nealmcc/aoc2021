package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/nealmcc/aoc2021/pkg/fishtree"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	tokens, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	_, err = fishtree.ShuntingYard(tokens)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()

	fmt.Println("part1:", 0)

	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read the given input, returning a byte reader.  The byte reader can then
// be assembled into an abstract fishtree using fishtree.ShuntingYard()
func read(r io.Reader) (*bytes.Reader, error) {
	buf, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buf), nil
}
