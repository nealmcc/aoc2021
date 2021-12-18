package main

import (
	"fmt"
	"math"

	"github.com/nealmcc/aoc2021/pkg/vector"
)

var target = struct{ left, right, top, bottom int }{
	left: 156, right: 202, top: -69, bottom: -110,
}

// main solves parts 1 and 2
func main() {
	fmt.Println("part 1", part1())
	fmt.Println("part 2", part2())
}

func part1() int {
	// Note that the upward arc and downward arc will have symmetrical y positions.
	// therefore, there will always be a position where y = 0 and dy < 0.
	// When this occurs, the farther up the probe went, the faster down it will
	// be falling.  The fastest down it can be falling is when dy is -110, since
	// that is the lowest point in the target area.
	dyDown := target.bottom // -110

	// When dy is -110, and y is 0, the probe's next movement will land it
	// at the very bottom of the target area.
	// Therefore, the initial y velocity must be 109:
	dyUp := -1*dyDown - 1 // 109

	// now, we just need to find the sum of (1..dyMaxUp)
	sum := dyUp * (dyUp + 1) / 2 // 5995

	return sum
}

func part2() int {
	dxMin, dxMax, dyMin, dyMax := findLimits()

	vectors := make([]vector.Coord, 0, 64)

	for x := dxMin; x <= dxMax; x++ {
		for y := dyMin; y <= dyMax; y++ {
			velocity := vector.Coord{X: x, Y: y}
			if reachesTarget(velocity) {
				vectors = append(vectors, velocity)
			}
		}
	}

	return len(vectors)
}

func findLimits() (int, int, int, int) {
	// largest single step:
	dxMax := target.right
	dyMin := target.bottom

	// largest upward velocity (as established in part 1):
	dyMax := 109

	// let distX be the shortest horizontal distance that the probe must travel:
	distX := target.left

	// let dxMin be the smallest starting x velocity such that the probe
	// barely reaches the left edge of the target.
	// In this scenario, the probe's ending x velocity must be 0, and it decreased
	// by 1 for every step it took.
	// Therefore:
	// sum(dx)[from dx = 0 to dx = dxMin] = distX
	// (dxMin)(dxMin+1)/2 = distX
	// dxMin * dxMin + dxMin -2*distX = 0
	// if w == dxMin, then Aw^2 + Bw + C = 0
	a, b, c := 1, 1, -2*distX
	b24ac := b*b - 4*a*c
	dxMin := int(-1.0*float64(b) + math.Sqrt(float64(b24ac))/2.0)

	if dxMin*(dxMin+1)/2 < distX {
		dxMin++
	}

	return dxMin, dxMax, dyMin, dyMax
}

func reachesTarget(v vector.Coord) bool {
	pos := vector.Coord{X: 0, Y: 0}
	for {
		if pos.X >= target.left && pos.X <= target.right &&
			pos.Y >= target.bottom && pos.Y <= target.top {
			return true
		}

		if pos.X > target.right {
			return false
		}

		if pos.Y < target.bottom {
			return false
		}

		pos = vector.Add(pos, v)

		if v.X > 0 {
			v.X--
		}
		v.Y--
	}
}
