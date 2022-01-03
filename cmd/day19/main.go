package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"time"

	v "github.com/nealmcc/aoc2021/pkg/vector"
)

// block is a cuboid region of the ocean which contains a collection of beacons,
// and one or more sensors.
// The coordinate system of the block is internally consistent.
type block struct {
	// Sensors is the set of sensors which have combined to observe this block.
	Sensors map[int]v.I3

	// Beacons is the collection of beacons in this block.
	Beacons []v.I3
}

// main solves both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	blocks, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(blocks)

	end := time.Now()

	fmt.Println("part1:", p1)
	// fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read the given input, returning the blocks of ocean that were scanned.
func read(r io.Reader) ([]block, error) {
	s := bufio.NewScanner(r)

	blocks := make([]block, 0, 8)

	var (
		id int
		b  block
	)
	for s.Scan() {
		line := s.Bytes()

		switch {
		case len(line) == 0:
			blocks = append(blocks, b)

		case line[1] == '-':
			b = block{
				Sensors: map[int]v.I3{id: {}},
				Beacons: make([]v.I3, 0, 16),
			}
			id++

		default:
			pos, err := v.ParseI3(line)
			if err != nil {
				return nil, err
			}
			b.Beacons = append(b.Beacons, pos)
		}
	}
	blocks = append(blocks, b)

	if err := s.Err(); err != nil {
		return nil, err
	}

	return blocks, nil
}

func part1(blocks []block) int {
	// assume the first block is facing the 'normal' way.
	// loop through each of the remaining blocks:
	//   - rotate that other block until it 'fits' with this one.
	//   - that is, it must have some rotation and reflection such
	//     that at least 12 of its beacons line up with the cumulative group.
	//   - find that rotation and translation, apply it to the other block,
	//     and merge them together.
	return 0
}

// transform creates a copy of this block, with all of its vectors transformed
// using the given linear transformation
func (box block) transform(m v.Matrix) block {
	b2 := block{
		Sensors: make(map[int]v.I3, len(box.Sensors)),
		Beacons: make([]v.I3, len(box.Beacons)),
	}
	for i, bn := range box.Beacons {
		bn2, _ := bn.Transform(m)
		b2.Beacons[i] = bn2
	}

	for id, s := range box.Sensors {
		s2, _ := s.Transform(m)
		b2.Sensors[id] = s2
	}

	return b2
}

// translate makes a copy of this box that has been repositioned by delta.
func (box block) translate(delta v.I3) block {
	b2 := block{
		Sensors: make(map[int]v.I3, len(box.Sensors)),
		Beacons: make([]v.I3, len(box.Beacons)),
	}

	for i, bn := range box.Beacons {
		b2.Beacons[i] = bn.Add(delta)
	}

	for id, s := range box.Sensors {
		b2.Sensors[id] = s.Add(delta)
	}

	return b2
}

// rotateToMatch repeatedly attempts to rotate box2 until it can match box1,
// or else there are no more rotations to attempt.  Returns the first
// transformation of box2 that matches box1.
// Returns false if no match can be found.
func rotateToMatch(box1, box2 block, minQty int) (block, bool) {
	sort.Sort(byXYZ(box1.Beacons))
	sort.Sort(byXYZ(box2.Beacons))

	for _, rot := range getRotations() {
		box2t := box2.transform(rot)
		sort.Sort(byXYZ(box2t.Beacons))

		b1, b2 := box1.Beacons, box2t.Beacons

		if offset, isMatch := beaconsMatch(b1, b2, minQty); isMatch {
			box2t = box2t.translate(offset)
			return box2t, true
		}
	}

	return box2, false
}

// beaconsMatch finds the offset between the two slices of beacons, if they have
// at least minQty elements in common when that offset added to all elements in
// b2. Assumes that both slices are sorted, and that minQty is at least 1.
func beaconsMatch(beacons1, beacons2 []v.I3, minQty int) (offset v.I3, isMatch bool) {
	if len(beacons1) < minQty || len(beacons2) < minQty {
		return v.I3{}, false
	}

	// choose a subset of the beacons in each slice so we can get an accurate
	// offset even if the first element in either list doesn't match
	for j := 0; j < len(beacons1)-minQty+1; j++ {
		b1 := beacons1[j:]

		for k := 0; k < len(beacons2)-minQty+1; k++ {
			b2 := beacons2[k:]

			// find the vector from b2 to b1
			offset = b1[0].Subtract(b2[0])

			// walk both lists at the same time, counting identical beacons
			var i1, i2, n int
			for i1 < len(b1) && i2 < len(b2) {
				if b1[i1] == b2[i2].Add(offset) {
					i1, i2, n = i1+1, i2+1, n+1
					if n >= minQty {
						return offset, true
					}
					continue
				}

				if isLess(b1[i1], b2[i2]) {
					i1++
					continue
				}

				i2++
			}
		}
	}

	return v.I3{}, false
}

// getRotations returns the 24 linear transformations that define each distinct
// rotation in I3.
func getRotations() []v.Matrix {
	rot := []v.Matrix{
		{
			{1, 0, 0},
			{0, 1, 0},
			{0, 0, 1},
		},
		{
			{1, 0, 0},
			{0, 0, -1},
			{0, 1, 0},
		},
		{
			{1, 0, 0},
			{0, -1, 0},
			{0, 0, -1},
		},
		{
			{1, 0, 0},
			{0, 0, 1},
			{0, -1, 0},
		},
		{
			{0, -1, 0},
			{1, 0, 0},
			{0, 0, 1},
		},
		{
			{0, 0, 1},
			{1, 0, 0},
			{0, 1, 0},
		},
		{
			{0, 1, 0},
			{1, 0, 0},
			{0, 0, -1},
		},
		{
			{0, 0, -1},
			{1, 0, 0},
			{0, -1, 0},
		},
		{
			{-1, 0, 0},
			{0, -1, 0},
			{0, 0, 1},
		},
		{
			{-1, 0, 0},
			{0, 0, -1},
			{0, -1, 0},
		},
		{
			{-1, 0, 0},
			{0, 1, 0},
			{0, 0, -1},
		},
		{
			{-1, 0, 0},
			{0, 0, 1},
			{0, 1, 0},
		},
		{
			{0, 1, 0},
			{-1, 0, 0},
			{0, 0, 1},
		},
		{
			{0, 0, 1},
			{-1, 0, 0},
			{0, -1, 0},
		},
		{
			{0, -1, 0},
			{-1, 0, 0},
			{0, 0, -1},
		},
		{
			{0, 0, -1},
			{-1, 0, 0},
			{0, 1, 0},
		},
		{
			{0, 0, -1},
			{0, 1, 0},
			{1, 0, 0},
		},
		{
			{0, 1, 0},
			{0, 0, 1},
			{1, 0, 0},
		},
		{
			{0, 0, 1},
			{0, -1, 0},
			{1, 0, 0},
		},
		{
			{0, -1, 0},
			{0, 0, -1},
			{1, 0, 0},
		},
		{
			{0, 0, -1},
			{0, -1, 0},
			{-1, 0, 0},
		},
		{
			{0, -1, 0},
			{0, 0, 1},
			{-1, 0, 0},
		},
		{
			{0, 0, 1},
			{0, 1, 0},
			{-1, 0, 0},
		},
		{
			{0, 1, 0},
			{0, 0, -1},
			{-1, 0, 0},
		},
	}

	return rot
}
