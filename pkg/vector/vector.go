// Package vector implements a 2D coordinate system
package vector

// Coord is a 2-dimensional integer coordinate.
type Coord struct {
	X int
	Y int
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
