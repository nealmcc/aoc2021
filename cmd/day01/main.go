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

	count1, err := part1(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("part 1 count", count1)

	file.Seek(0, io.SeekStart)

	count2, err := part2(file)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("part 2 count", count2)
}

func part1(r io.Reader) (int, error) {
	s := bufio.NewScanner(r)

	curr, err := nextInt(s)
	if err != nil {
		return 0, err
	}

	count := 0
	for s.Scan() {
		next, err := strconv.Atoi(s.Text())
		if err != nil {
			return 0, err
		}
		if next > curr {
			count++
		}
		curr = next
	}

	if err := s.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func part2(r io.Reader) (int, error) {
	s := bufio.NewScanner(r)

	a, err := nextInt(s)
	if err != nil {
		return 0, err
	}
	b, err := nextInt(s)
	if err != nil {
		return 0, err
	}
	c, err := nextInt(s)
	if err != nil {
		return 0, err
	}

	count := 0
	for s.Scan() {
		next, err := strconv.Atoi(s.Text())
		if err != nil {
			return 0, err
		}
		// the current window will be larger than the previous iff
		// b + c + next > a + b + c
		if next > a {
			count++
		}
		a, b, c = b, c, next
	}

	if err := s.Err(); err != nil {
		return 0, err
	}

	return count, nil
}

func nextInt(s *bufio.Scanner) (int, error) {
	if !s.Scan() {
		return 0, errors.New("input too short")
	}
	return strconv.Atoi(s.Text())
}
