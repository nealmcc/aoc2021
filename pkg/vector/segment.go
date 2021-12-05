package vector

import "math"

// Segment is a line segment from a to b.
type Segment struct {
	A, B Coord
}

// Slope returns the slope of the line segment from A to B.
// if A == B then the slope is is not defined.
func (s Segment) Slope() float64 {
	delta := Sub(s.B, s.A)
	if (delta == Coord{}) {
		return math.NaN()
	}

	if delta.X == 0 {
		return math.Inf(delta.Y)
	}

	return float64(delta.Y) / float64(delta.X)
}
