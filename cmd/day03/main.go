package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
)

func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	diag, err := part1(in)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("part 1 product", diag.gamma()*diag.epsilon())
}

func part1(r io.Reader) (*meter, error) {
	s := bufio.NewScanner(r)
	if !s.Scan() {
		return nil, errors.New("input too short")
	}

	first := s.Text()
	m := newMeter(len(first))
	err := m.sample(first)
	if err != nil {
		return nil, err
	}

	for s.Scan() {
		if err := m.sample(s.Text()); err != nil {
			return nil, err
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return m, nil
}
