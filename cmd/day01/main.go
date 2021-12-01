package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

func main() {
	file, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	count, err := part1(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("part 1 count", count)
}

func part1(r io.Reader) (int, error) {
	s := bufio.NewScanner(r)
	if !s.Scan() {
		return 0, errors.New("empty input")
	}

	prev, err := strconv.Atoi(s.Text())
	if err != nil {
		return 0, err
	}

	count := 0
	for s.Scan() {
		if err := s.Err(); err != nil {
			return 0, err
		}
		curr, err := strconv.Atoi(s.Text())
		if err != nil {
			return 0, err
		}
		if curr > prev {
			count++
		}
		prev = curr
	}

	if err := s.Err(); err != nil {
		return 0, err
	}

	return count, nil
}
