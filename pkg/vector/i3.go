package vector

import (
	"bytes"
	"errors"
	"strconv"
)

// I3 is a vector in three-dimensional integer space.
type I3 struct {
	X, Y, Z int
}

// ParseI3 interprets the given slice of bytes as a triplet of integers
// and returns the corresponding x,y,z vector.
func ParseI3(text []byte) (I3, error) {
	parts := bytes.Split(text, []byte{','})
	if len(parts) != 3 {
		return I3{}, errors.New("malformed input")
	}

	x, err := strconv.Atoi(string(parts[0]))
	if err != nil {
		return I3{}, err
	}
	y, err := strconv.Atoi(string(parts[1]))
	if err != nil {
		return I3{}, err
	}

	z, err := strconv.Atoi(string(parts[2]))
	if err != nil {
		return I3{}, err
	}

	return I3{X: x, Y: y, Z: z}, nil
}

// Add returns a copy of the vector sum of (v + v2).
func (v I3) Add(v2 I3) I3 {
	v.X += v2.X
	v.Y += v2.Y
	v.Z += v2.Z
	return v
}

// Subtract returns the vector difference (v - v2).
func (v I3) Subtract(v2 I3) I3 {
	v.X -= v2.X
	v.Y -= v2.Y
	v.Z -= v2.Z
	return v
}

// ToMatrix converts this vector to a matrix in column vector form, suitable
// for use in a cross product with a linear transformation.
func (v I3) ToMatrix() Matrix {
	return Matrix{{v.X}, {v.Y}, {v.Z}}
}

// RotX90 returns a copy of this vector rotated 90 degrees about the x-axis.
func (v I3) RotX90() I3 {
	v, _ = v.Transform(Matrix{
		{1, 0, 0},  // x = (1,       0,       0)
		{0, 0, -1}, // y = (0, cos(90), -sin(90))
		{0, 1, 0},  // z = (0, sin(90), cos(90))
	})
	return v
}

// I3RotY90 returns a copy of this vector rotated 90 degrees about the y-axis.
func (v I3) RotY90() I3 {
	v, _ = v.Transform(Matrix{
		{0, 0, 1},  // x = (cos(90),  0,  sin(90))
		{0, 1, 0},  // y = (      0,  1,       0)
		{-1, 0, 0}, // z = (-sin(90), 0, cos(90))
	})
	return v
}

// RotZ90 returns a copy of this vector rotated 90-degrees about the z axis.
func (v I3) RotZ90() I3 {
	v, _ = v.Transform(Matrix{
		{0, -1, 0}, // x = (cos(90), -sin(90), 0) x becomes negative y
		{1, 0, 0},  // y = (sin(90),  cos(90), 0) y becomes x
		{0, 0, 1},  // z = (      0,        0, 1) z is unaffected
	})
	return v
}

// I3Transform returns a copy of this vector transformed by m.
// see: https://www.khanacademy.org/math/linear-algebra/matrix-transformations
func (v I3) Transform(m Matrix) (I3, error) {
	c, err := m.Times(v.ToMatrix())
	if err != nil {
		return I3{}, err
	}

	v.X = c[0][0]
	v.Y = c[1][0]
	v.Z = c[2][0]

	return v, nil
}

// IdentityI3 returns the identity matrix in I3
func IdentityI3() Matrix {
	return Matrix{
		{1, 0, 0},
		{0, 1, 0},
		{0, 0, 1},
	}
}

// Bounds finds a bounding box that encloses all of the given coordinates.
func Bounds(coords []I3) Cuboid {
	const high int = 1<<63 - 1
	const low int = -1 * high

	x1, y1, z1 := high, high, high
	x2, y2, z2 := low, low, low

	for _, pos := range coords {
		if pos.X < x1 {
			x1 = pos.X
		}
		if pos.X > x2 {
			x2 = pos.X
		}

		if pos.Y < y1 {
			y1 = pos.Y
		}
		if pos.Y > y2 {
			y2 = pos.Y
		}

		if pos.Z < z1 {
			z1 = pos.Z
		}
		if pos.Z > z2 {
			z2 = pos.Z
		}
	}

	return Cuboid{X1: x1, Y1: y1, Z1: z1, X2: x2, Y2: y2, Z2: z2}
}
