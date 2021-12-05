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

// Add returns the vector sum of the two coordinates.
func Add(v1, v2 Coord) Coord {
	return Coord{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}

// Scale returns the scalar product of v and n.
func Scale(v Coord, n int) Coord {
	return Coord{
		X: v.X * n,
		Y: v.Y * n,
	}
}
