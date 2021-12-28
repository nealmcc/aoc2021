package vector

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

// Intersection finds the intersection of this cuboid with the other, if any.
// If there is no intersection, then the zero Cuboid is returned.
func (c Cuboid) Intersect(other Cuboid) Cuboid {
	intersection, _, _ := CuboidOuterJoin(c, other)
	return intersection
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

// Slice this cuboid into new cuboids which cover the same region of 3D
// space as the original.  The split will occur at the planes x y, and z.
// It is valid to s
// explode a cuboid at a point on its surface, in which case the volume of
// some of the resulting pieces will be 0.  It is not valid to explode a
// cuboid at a point it doesn't contain.
func (c Cuboid) Slice(x, y, z int) []Cuboid {
	two := c.SliceX(x)

	four := make([]Cuboid, 0, 4)
	for _, box := range two {
		boxes := box.SliceY(y)
		four = append(four, boxes...)
	}

	eight := make([]Cuboid, 0, 8)
	for _, box := range four {
		boxes := box.SliceZ(z)
		eight = append(eight, boxes...)
	}

	return eight
}

// SliceX bisects this cuboid into two portions divided by the plane at x.
// If the plane at x does not intersect this cuboid, then this cuboid is returned.
func (c Cuboid) SliceX(x int) []Cuboid {
	c = c.Normal()
	if x <= c.X1 || x >= c.X2 {
		return []Cuboid{c}
	}

	lower := Cuboid{X1: c.X1, X2: x, Y1: c.Y1, Y2: c.Y2, Z1: c.Z1, Z2: c.Z2}
	upper := Cuboid{X1: x, X2: c.X2, Y1: c.Y1, Y2: c.Y2, Z1: c.Z1, Z2: c.Z2}

	return []Cuboid{lower, upper}
}

// SliceY bisects this cuboid into two portions divided by the plane at y.
// If the plane at y does not intersect this cuboid, then this cuboid is returned.
func (c Cuboid) SliceY(y int) []Cuboid {
	c = c.Normal()
	if y <= c.Y1 || y >= c.Y2 {
		return []Cuboid{c}
	}

	lower := Cuboid{X1: c.X1, X2: c.X2, Y1: c.Y1, Y2: y, Z1: c.Z1, Z2: c.Z2}
	upper := Cuboid{X1: c.X1, X2: c.X2, Y1: y, Y2: c.Y2, Z1: c.Z1, Z2: c.Z2}

	return []Cuboid{lower, upper}
}

// SliceZ bisects this cuboid into two portions divided by the plane at z.
// If the plane at z does not intersect this cuboid, then this cuboid is returned.
func (c Cuboid) SliceZ(z int) []Cuboid {
	c = c.Normal()
	if z <= c.Z1 || z >= c.Z2 {
		return []Cuboid{c}
	}

	lower := Cuboid{X1: c.X1, X2: c.X2, Y1: c.Y1, Y2: c.Y2, Z1: c.Z1, Z2: z}
	upper := Cuboid{X1: c.X1, X2: c.X2, Y1: c.Y1, Y2: c.Y2, Z1: z, Z2: c.Z2}

	return []Cuboid{lower, upper}
}

// CuboidOuterJoin returns the intersection of cuboids a and b
// as well as the difference of a-b and b-a.
// The sum of all the return values is a + b.
func CuboidOuterJoin(a, b Cuboid) (intersect Cuboid, left, right []Cuboid) {
	if a.IsDisjoint(b) {
		return Cuboid{}, []Cuboid{a}, []Cuboid{b}
	}
	a, b = a.Normal(), b.Normal()

	var (
		aPieces = make(map[Cuboid]struct{}, 16)
		bPieces = make(map[Cuboid]struct{}, 16)
	)

	eight := a.Slice(b.X1, b.Y1, b.Z1)
	for _, box := range eight {
		slices := box.Slice(b.X2, b.Y2, b.Z2)
		for _, box := range slices {
			if box.Volume() > 0 {
				aPieces[box.Normal()] = struct{}{}
			}
		}
	}

	eight = b.Slice(a.X1, a.Y1, a.Z1)
	for _, box := range eight {
		slices := box.Slice(a.X2, a.Y2, a.Z2)
		for _, box := range slices {
			if box.Volume() > 0 {
				bPieces[box.Normal()] = struct{}{}
			}
		}
	}

	for boxA := range aPieces {
		if _, match := bPieces[boxA]; match {
			intersect = boxA
			delete(aPieces, boxA)
			delete(bPieces, boxA)
			break
		}
	}

	left = make([]Cuboid, 0, len(aPieces))
	for key := range aPieces {
		left = append(left, key)
	}

	right = make([]Cuboid, 0, len(bPieces))
	for key := range bPieces {
		right = append(right, key)
	}

	return intersect, left, right
}
