package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	displays, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(displays)
	fmt.Println("part1", p1)

	p2 := part2(displays)
	fmt.Println("part2", p2)
}

// display is one row of input, with 10 unique signals and a 4-digit output.
// The digits map is used to associate each unique signal with its digit
type display struct {
	signals []signal
	output  []signal
	digits  map[signal]int
}

func read(r io.Reader) ([]display, error) {
	s := bufio.NewScanner(r)

	displays := make([]display, 0, 16)
	for s.Scan() {
		d := display{digits: make(map[signal]int)}
		parts := strings.Split(s.Text(), " | ")
		d.signals = make([]signal, 10)
		for i, v := range strings.Split(parts[0], " ") {
			d.signals[i] = normal(v)
		}
		d.output = make([]signal, 4)
		for i, v := range strings.Split(parts[1], " ") {
			d.output[i] = normal(v)
		}
		displays = append(displays, d)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return displays, nil
}

func part1(displays []display) int {
	count := 0

	for _, d := range displays {
		for _, s := range d.output {
			switch len(s) {
			case 2, 3, 4, 7:
				count++
			default:
				continue
			}
		}
	}

	return count
}

func part2(displays []display) int {
	sum := 0
	for _, d := range displays {
		solve(d)
		sum += d.value()
	}
	return sum
}

func solve(d display) {
	sort.Slice(d.signals, func(i, j int) bool {
		return len(d.signals[i]) < len(d.signals[j])
	})

	// smallest signal is 1
	one := d.signals[0]
	d.digits[one] = 1

	// next smallest is 7
	seven := d.signals[1]
	d.digits[seven] = 7

	// third smallest (length 4) is 4
	four := d.signals[2]
	d.digits[four] = 4

	// largest signal is 8
	eight := d.signals[9]
	d.digits[eight] = 8

	var zero, six, nine signal
	for _, symbol := range d.signals[6:9] {
		switch {
		// 0 + 7 is still 0, and 9 + 7 is still 9, but 6 + 7 is 8
		case union(symbol, seven) != symbol:
			six = symbol
			d.digits[six] = 6
		// 0 - 4 has length 3, but 9 - 4 has length 2
		case len(diff(symbol, four)) == 3:
			zero = symbol
			d.digits[zero] = 0
		default:
			nine = symbol
			d.digits[nine] = 9
		}
	}

	upperRight := diff(one, six)

	var two, three, five signal

	for _, symbol := range d.signals[3:6] {
		switch {
		// 2 + 1 != 2 and 5 + 1 != 5 but 3 + 1 = 3
		case symbol == union(symbol, one):
			three = symbol
			d.digits[three] = 3
		// both 2 and 3 are unchanged by adding 'upper right'
		// but 3 is covered by the case above
		case symbol == union(symbol, upperRight):
			two = symbol
			d.digits[two] = 2
		default:
			five = symbol
			d.digits[five] = 5
		}
	}
}

// signal is a set of segments for a seven-segment display
type signal string

func normal(s string) signal {
	set := make(map[rune]struct{})
	for _, ch := range s {
		set[ch] = struct{}{}
	}

	sig := make([]rune, 0, 7)
	for k := range set {
		sig = append(sig, k)
	}

	sort.Slice(sig, func(i, j int) bool {
		return sig[i] < sig[j]
	})
	return signal(sig)
}

// diff returns the set difference of a - b
func diff(a, b signal) signal {
	// TODO: optimise for a and b both being normalised
	notFound := make([]rune, 0)
	for _, ch := range a {
		if !contains(b, ch) {
			notFound = append(notFound, ch)
		}
	}
	return normal(string(notFound))
}

// union returns the set union of all the given signals
func union(sets ...signal) signal {
	// TODO: optimise for a and b both being normalised
	sum := make(map[rune]struct{})
	for _, bytes := range sets {
		for _, b := range bytes {
			sum[b] = struct{}{}
		}
	}

	keys := make([]rune, 0, 4)
	for k := range sum {
		keys = append(keys, k)
	}

	return normal(string(keys))
}

func contains(sig signal, r rune) bool {
	for _, b := range sig {
		if r == b {
			return true
		}
	}
	return false
}

func (d display) value() int {
	val := d.digits[d.output[0]]
	for _, s := range d.output[1:] {
		val *= 10
		val += d.digits[s]
	}
	return val
}
