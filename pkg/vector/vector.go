// Package vector implements a 2D coordinate system
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

// Parse reads one or more input strings in the form of x,y and returns
// the corresponding coordinates.
//
// Example:
//     ParseCoords("8,0", "2,3") => []Coord{{X: 8, Y:0}, {X:2, Y:3}}
func ParseCoords(in ...string) ([]Coord, error) {
	coords := make([]Coord, 0, 2)
	for _, pair := range in {
		parts := strings.Split(pair, ",")
		if len(parts) != 2 {
			return nil, fmt.Errorf("expected 2 values, got %d", len(parts))
		}
		x, err := strconv.Atoi(parts[0])
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse x coord as integer")
		}
		y, err := strconv.Atoi(parts[1])
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse y coord as integer")
		}
		coords = append(coords, Coord{X: x, Y: y})
	}
	return coords, nil
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

func IsZero(v Coord) bool {
	return v == Coord{}
}

// Reduce returns the shortest vector with the same direction as v,
// that can still be represented with integer values for X and Y.
// Also returns the largest positive integer that evenly divides v.
//
// Example:
//    Reduce{{X: 12, Y: 3}} => {X: 4, Y: 1}, 3
//
// Note that Reduce(v1) == (v2, n) if and only if Scale(v2, n) == v1.
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

// Scale returns the scalar product of v and n.
func Scale(v Coord, n int) Coord {
	return Coord{
		X: v.X * n,
		Y: v.Y * n,
	}
}

// gcd calculates the greatest common divisor of a and b.
func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}
