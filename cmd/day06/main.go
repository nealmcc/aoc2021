package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

var input string = "1,1,1,3,3,2,1,1,1,1,1,4,4,1,4,1,4,1,1,4,1,1,1,3,3,2,3,1,2,1,1,1,1,1,1,1,3,4,1,1,4,3,1,2,3,1,1,1,5,2,1,1,1,1,2,1,2,5,2,2,1,1,1,3,1,1,1,4,1,1,1,1,1,3,3,2,1,1,3,1,4,1,2,1,5,1,4,2,1,1,5,1,1,1,1,4,3,1,3,2,1,4,1,1,2,1,4,4,5,1,3,1,1,1,1,2,1,4,4,1,1,1,3,1,5,1,1,1,1,1,3,2,5,1,5,4,1,4,1,3,5,1,2,5,4,3,3,2,4,1,5,1,1,2,4,1,1,1,1,2,4,1,2,5,1,4,1,4,2,5,4,1,1,2,2,4,1,5,1,4,3,3,2,3,1,2,3,1,4,1,1,1,3,5,1,1,1,3,5,1,1,4,1,4,4,1,3,1,1,1,2,3,3,2,5,1,2,1,1,2,2,1,3,4,1,3,5,1,3,4,3,5,1,1,5,1,3,3,2,1,5,1,1,3,1,1,3,1,2,1,3,2,5,1,3,1,1,3,5,1,1,1,1,2,1,2,4,4,4,2,2,3,1,5,1,2,1,3,3,3,4,1,1,5,1,3,2,4,1,5,5,1,4,4,1,4,4,1,1,2"

// main solves the both part 1 and part 2, reading from the above string
func main() {
	fish, err := read(input)
	if err != nil {
		log.Fatal(err)
	}

	fish.cycle(80)
	fmt.Printf("part1: %d\n", fish.count())

	fish.cycle(256 - 80)
	fmt.Printf("part2: %d\n", fish.count())
}

// read a school of fish from the given input
func read(in string) (school, error) {
	parts := strings.Split(in, ",")
	fish := school{}

	for _, a := range parts {
		i, err := strconv.Atoi(a)
		if err != nil {
			return school{}, err
		}
		fish.days[i]++
	}

	return fish, nil
}

// school is a group of fish
type school struct {
	// each element in the array indicates how many days are at the given
	// stage in their reproductive cycle.
	days [9]int
}

// cycle executes the given number of daily cycles of reproduction
func (s *school) cycle(days int) {
	for n := 0; n < days; n++ {
		next := s.days[0]

		for i := 0; i < len(s.days)-1; i++ {
			s.days[i] = s.days[i+1]
		}
		s.days[6] += next
		s.days[8] = next
	}
}

// count the total number of fish in the school
func (s *school) count() int64 {
	var sum int64
	for _, n := range s.days {
		sum += int64(n)
	}
	return sum
}

// Format the school of fish for printing
func (s *school) Format(state fmt.State, v rune) {
	for i, n := range s.days {
		fmt.Fprintf(state, "%d:%2d   ", i, n)
	}
	fmt.Fprintf(state, "total: %d", s.count())
}
