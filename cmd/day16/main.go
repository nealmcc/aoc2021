package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	packet, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	versionSum := part1(packet)

	end := time.Now()

	fmt.Println("part1:", versionSum)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read a packet from the given reader
func read(r io.Reader) (Packet, error) {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	p := new(rawPacket)
	err = p.UnmarshalText(buf.Bytes())
	return p, err
}

func part1(p Packet) int {
	sum := p.Version()
	for _, child := range p.Children() {
		sum += child.Version()
	}
	return sum
}

type Packet interface {
	Version() int
	Value() int
	Children() []Packet
}
