package vector

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// Coord is a 2-dimensional integer coordinate.
type Coord struct {
	X int
	Y int
}

// ParseCoords reads one or more input strings in the form of x,y and returns
// the corresponding coordinates.
//
// Example:
//     ParseCoords("8,0", "2,3") => []Coord{{X: 8, Y:0}, {X:2, Y:3}}
func ParseCoords(in ...string) ([]Coord, error) {
	coords := make([]Coord, 0, 2)
	for _, pair := range in {
		pos, err := ParseCoord(pair)
		if err != nil {
			return nil, err
		}
		coords = append(coords, pos)
	}
	return coords, nil
}

// ParseCoord reads one input string in the form of x,y and returns
// the corresponding coordinate.
//
// Example:
//     ParseCoord("8,0") => Coord{X: 8, Y:0}
func ParseCoord(s string) (Coord, error) {
	parts := strings.Split(s, ",")
	if len(parts) != 2 {
		return Coord{}, fmt.Errorf("expected 2 values, got %d", len(parts))
	}
	x, err := strconv.Atoi(parts[0])
	if err != nil {
		return Coord{}, errors.Wrap(err, "failed to parse x coord as integer")
	}
	y, err := strconv.Atoi(parts[1])
	if err != nil {
		return Coord{}, errors.Wrap(err, "failed to parse y coord as integer")
	}
	return Coord{X: x, Y: y}, nil
}

// Add returns the vector sum of a + b.
func Add(a, b Coord) Coord {
	return Coord{
		X: a.X + b.X,
		Y: a.Y + b.Y,
	}
}

// Sub returns the vector difference of a - b.
func Sub(a, b Coord) Coord {
	return Coord{
		X: a.X - b.X,
		Y: a.Y - b.Y,
	}
}

func Neighbours8(v Coord) []Coord {
	return []Coord{
		{X: v.X - 1, Y: v.Y - 1},
		{X: v.X, Y: v.Y - 1},
		{X: v.X + 1, Y: v.Y - 1},

		{X: v.X - 1, Y: v.Y},
		{X: v.X + 1, Y: v.Y},

		{X: v.X - 1, Y: v.Y + 1},
		{X: v.X, Y: v.Y + 1},
		{X: v.X + 1, Y: v.Y + 1},
	}
}

// Reduce returns the shortest vector with the same direction as v,
// that can still be represented with integer values for X and Y.
// Also returns the largest positive integer that evenly divides v.
//
// Example:
//    Reduce({ X: -12, Y: -3 }) => { X: -4, Y: -1}, 3
func Reduce(v Coord) (Coord, int) {
	if (v == Coord{}) {
		return v, 1
	}

	if v.X == 0 {
		if v.Y > 0 {
			return Coord{X: 0, Y: 1}, v.Y
		}
		return Coord{X: 0, Y: -1}, -1 * v.Y
	}

	if v.Y == 0 {
		if v.X > 0 {
			return Coord{X: 1, Y: 0}, v.X
		}
		return Coord{X: -1, Y: 0}, -1 * v.X
	}

	scale := gcd(v.X, v.Y)
	if scale < 0 {
		scale *= -1
	}
	return Coord{
		X: v.X / scale,
		Y: v.Y / scale,
	}, scale
}

// gcd calculates the greatest common divisor of a and b.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
