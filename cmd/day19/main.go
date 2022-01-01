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
	"strings"
	"time"

	v "github.com/nealmcc/aoc2021/pkg/vector"
)

// main solves both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	sensors, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(sensors)

	end := time.Now()

	fmt.Println("part1:", p1)
	// fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read the given input, returning all the sensors.
func read(r io.Reader) ([]sensor, error) {
	s := bufio.NewScanner(r)

	sensors := make([]sensor, 0, 8)

	var sense sensor
	for s.Scan() {
		line := s.Bytes()

		switch {
		case len(line) == 0:
			sensors = append(sensors, sense)

		case line[1] == '-':
			line = bytes.TrimPrefix(line, []byte("--- scanner "))
			line = bytes.TrimSuffix(line, []byte(" ---"))
			id, err := strconv.Atoi(string(line))
			if err != nil {
				return nil, err
			}
			sense = sensor{
				id:      id,
				facing:  identity(),
				beacons: make([]v.I3, 0, 16),
			}

		default:
			parts := bytes.Split(line, []byte{','})
			if len(parts) != 3 {
				return nil, errors.New("malformed input")
			}

			x, err := strconv.Atoi(string(parts[0]))
			if err != nil {
				return nil, err
			}
			y, err := strconv.Atoi(string(parts[1]))
			if err != nil {
				return nil, err
			}

			z, err := strconv.Atoi(string(parts[2]))
			if err != nil {
				return nil, err
			}

			sense.beacons = append(sense.beacons, v.I3{X: x, Y: y, Z: z})
		}
	}
	sensors = append(sensors, sense)

	if err := s.Err(); err != nil {
		return nil, err
	}

	return sensors, nil
}

type sensor struct {
	// id is the sensor's unique identifier
	id int

	// facing identifies the current rotation and facing of this sensor.
	// by default, a sensor starts off unrotated or reflected.
	facing transform

	// pos is this sensor's position, relative to an arbitrary origin.
	pos v.I3

	// beacons is the set of beacons that this sensor has located,
	// expressed relative to the sensor's location.
	beacons []v.I3
}

// transform is a 3x3 transformation matrix.
type transform = [][]int

// identity is the identity matrix in I3. Sensors start with this as their
// assumed facing.
func identity() transform {
	return transform{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
}

func (s sensor) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("--- scanner %d ---\n", s.id))
	b.WriteString(fmt.Sprintf("pos: %v\n", s.pos))
	b.WriteString(fmt.Sprintf("transform: %+v\n", s.facing))
	for _, bn := range s.beacons {
		b.WriteString(fmt.Sprintf("%v\n", bn))
	}

	return b.String()
}

func part1(sensors []sensor) int {
	// assume the first sensor is facing the 'normal' way.
	// loop through each of the remaining beacons:
	//   - rotate that other beacon until it 'fits' with this one.
	//   - that is, it must have some rotation, reflection and translation such
	//     that at least 12 of its beacons line up with the cumulative group.
	//   - find that rotation and translation, and
	return 0
}
