package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
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
			sensors = append(sensors, sense.normal())

		case line[1] == '-':
			line = bytes.TrimPrefix(line, []byte("--- scanner "))
			line = bytes.TrimSuffix(line, []byte(" ---"))
			id, err := strconv.Atoi(string(line))
			if err != nil {
				return nil, err
			}
			sense = sensor{
				id: id,
				beacons: beaconSet{
					b: make([]v.I3, 0, 16),
				},
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

			sense.beacons.b = append(sense.beacons.b, v.I3{X: x, Y: y, Z: z})
		}
	}
	sensors = append(sensors, sense.normal())

	if err := s.Err(); err != nil {
		return nil, err
	}

	return sensors, nil
}

// sensor is one sensor's perspective of the surrounding ocean.
type sensor struct {
	// id is the sensor's unique identifier
	id int

	// beacons is the set of beacons that this sensor has located.
	beacons beaconSet

	// pos is the position of this sensor relative to the origin of its beaconSet.
	pos v.I3
}

// normal puts this sensor into normal form. A sensor is in normal form when
// the beacons' lower bound is 0,0,0 and all the beacons are sorted
// in ascending order by their X, Y and then Z coordinates.
func (s sensor) normal() sensor {
	move := s.beacons.standardise()
	s.pos.Translate(move)
	return s
}

func (s sensor) String() string {
	var b strings.Builder

	b.WriteString(fmt.Sprintf("--- scanner %d ---\n", s.id))
	b.WriteString(fmt.Sprintf("pos: %v\n", s.pos))
	for _, bn := range s.beacons.b {
		b.WriteString(fmt.Sprintf("%v\n", bn))
	}

	return b.String()
}

// beaconSet defines a cuboid region of the ocean, and the set of beacons
// in that region.
type beaconSet struct {
	// extents defines the size of this region of the ocean.
	extents v.I3

	// b is the list of beacons in this set
	b []v.I3
}

// standardise the given set of beacons by arranging them within a box with
// 0,0,0 at the lower corner, updating the extents of the beaconSet, and sorting
// them in increasing order by their X, Y and Z coordinates.  Returns a vector
// which is how much each beacon was translated by, so that any external frame
// of reference to this beaconSet can be adjusted accordingly.
func (bs *beaconSet) standardise() v.I3 {
	sort.Sort(byXYZ(bs.b))

	bounds := v.Bounds(bs.b)

	move := v.I3{
		X: -1 * bounds.X1,
		Y: -1 * bounds.Y1,
		Z: -1 * bounds.Z1,
	}

	for i := range bs.b {
		bs.b[i].Translate(move)
	}

	bs.extents.X = bounds.X2 - bounds.X1
	bs.extents.Y = bounds.Y2 - bounds.Y1
	bs.extents.Z = bounds.Z2 - bounds.Z1

	return move
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
