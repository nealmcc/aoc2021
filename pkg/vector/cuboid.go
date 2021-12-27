package vector

import (
	"errors"
)

// Cuboid defines a right rectangular prism that begins
// at position (X1, Y1, Y1) in one corner, and ends at
// position (X2, Y2, Z2) in the opposite corner.
type Cuboid struct {
	X1, X2 int
	Y1, Y2 int
	Z1, Z2 int
}

// Normal returns a copy of this cuboid, with X1 <= X2, Y1 <= Y2 and Z1 <= Z2.
func (c Cuboid) Normal() Cuboid {
	if c.X2 < c.X1 {
		c.X1, c.X2 = c.X2, c.X1
	}
	if c.Y2 < c.Y1 {
		c.Y1, c.Y2 = c.Y2, c.Y1
	}
	if c.Z2 < c.Z1 {
		c.Z1, c.Z2 = c.Z2, c.Z1
	}
	return c
}

// Volume returns the volume of this cuboid.
func (c Cuboid) Volume() int {
	c = c.Normal()
	width := c.X2 - c.X1
	height := c.Y2 - c.Y1
	depth := c.Z2 - c.Z1
	return width * height * depth
}

// Contains checks to see if the given point is within the bounds of
// this cuboid.  A point on the surface of a cuboid is considered
// within its bounds.
func (c Cuboid) Contains(x, y, z int) bool {
	c = c.Normal()
	if x <= c.X1 || x >= c.X2 {
		return false
	}
	if y <= c.Y1 || y >= c.Y2 {
		return false
	}
	if z <= c.Z1 || z >= c.Z2 {
		return false
	}
	return true
}

// Explode this cuboid into 8 new cuboids which cover the same region of 3D
// space as the original.  The split will occur at x,y,z.  It is valid to
// explode a cuboid at a point on its surface, in which case the volume of
// some of the resulting pieces will be 0.  It is not valid to explode a
// cuboid at a point it doesn't contain.
func (c Cuboid) Explode(x, y, z int) ([]Cuboid, error) {
	if !c.Contains(x, y, z) {
		return nil, errors.New("can only explode a cuboid from within its bounds")
	}

	two, _ := c.SliceX(x)

	four := make([]Cuboid, 0, 4)
	for _, box := range two {
		boxes, _ := box.SliceY(y)
		four = append(four, boxes...)
	}

	eight := make([]Cuboid, 0, 8)
	for _, box := range four {
		boxes, _ := box.SliceZ(z)
		eight = append(eight, boxes...)
	}

	return eight, nil
}

// SliceX bisects this cuboid into two portions divided by the plane at x.
// Returns false if x is outside the bounds of this cuboid.
func (c Cuboid) SliceX(x int) ([]Cuboid, bool) {
	if x < c.X1 || x > c.X2 {
		return nil, false
	}

	lower := Cuboid{X1: c.X1, X2: x, Y1: c.Y1, Y2: c.Y2, Z1: c.Z1, Z2: c.Z2}
	upper := Cuboid{X1: x, X2: c.X2, Y1: c.Y1, Y2: c.Y2, Z1: c.Z1, Z2: c.Z2}

	return []Cuboid{lower, upper}, true
}

// SliceY bisects this cuboid into two portions divided by the plane at y.
// Returns false if y is outside the bounds of this cuboid.
func (c Cuboid) SliceY(y int) ([]Cuboid, bool) {
	if y < c.Y1 || y > c.Y2 {
		return nil, false
	}

	lower := Cuboid{X1: c.X1, X2: c.X2, Y1: c.Y1, Y2: y, Z1: c.Z1, Z2: c.Z2}
	upper := Cuboid{X1: c.X1, X2: c.X2, Y1: y, Y2: c.Y2, Z1: c.Z1, Z2: c.Z2}

	return []Cuboid{lower, upper}, true
}

// SliceZ bisects this cuboid into two portions divided by the plane at z.
// Returns false if z is outside the bounds of this cuboid.
func (c Cuboid) SliceZ(z int) ([]Cuboid, bool) {
	if z < c.Z1 || z > c.Z2 {
		return nil, false
	}

	lower := Cuboid{X1: c.X1, X2: c.X2, Y1: c.Y1, Y2: c.Y2, Z1: c.Z1, Z2: z}
	upper := Cuboid{X1: c.X1, X2: c.X2, Y1: c.Y1, Y2: c.Y2, Z1: z, Z2: c.Z2}

	return []Cuboid{lower, upper}, true
}

// Intersection finds the intersection of this cuboid with the other, if any.
func (c Cuboid) Intersect(other Cuboid) (Cuboid, bool) {
	if c.IsDisjoint(other) {
		return Cuboid{}, false
	}
	c = c.Normal()
	other = other.Normal()

	// var (
	// 	ourPieces   []Cuboid
	// 	otherPieces []Cuboid
	// )
	return Cuboid{}, true
}

// IsDisjoint returns true if this cuboid and the other cuboid have no
// overlapping space.
func (c Cuboid) IsDisjoint(other Cuboid) bool {
	a, b := c.Normal(), other.Normal()

	if a.X1 > b.X1 {
		a, b = b, a
	}
	if a.X2 <= b.X1 {
		return true
	}

	if a.Y1 > b.Y1 {
		a, b = b, a
	}
	if a.Y2 <= b.Y1 {
		return true
	}

	if a.Z1 > b.Z1 {
		a, b = b, a
	}
	if a.Z2 <= b.Z1 {
		return true
	}

	return false
}
