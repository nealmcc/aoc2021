package vector

import (
	"errors"
)

// I3 is a vector in three-dimensional integer space.
type I3 struct {
	X, Y, Z int
}

// ReflectXY reflects this vector through the x-y plane.
func (v *I3) ReflectXY() {
	v.Z *= -1
}

// RotX90 rotates this vector 90 degrees about the x-axis.
func (v *I3) RotX90() {
	v.Transform([][]int{
		{1, 0, 0},  // x = (1,       0,       0)
		{0, 0, -1}, // y = (0, cos(90), -sin(90))
		{0, 1, 0},  // z = (0, sin(90), cos(90))
	})
}

// I3RotX90 rotates this vector 90 degrees about the x-axis.
func (v *I3) RotY90() {
	v.Transform([][]int{
		{0, 0, 1},  // x = (cos(90),  0,  sin(90))
		{0, 1, 0},  // y = (      0,  1,       0)
		{-1, 0, 0}, // z = (-sin(90), 0, cos(90))
	})
}

// RotZ90 rotates this vector 90-degrees about the z axis.
func (v *I3) RotZ90() {
	v.Transform([][]int{
		{0, -1, 0}, // x = (cos(90), -sin(90), 0)
		{1, 0, 0},  // y = (sin(90),  cos(90), 0)
		{0, 0, 1},  // z = (      0,        0, 1)
	})
}

// I3Transform applies the given linear transformation m to this vector.
// see: https://www.khanacademy.org/math/linear-algebra/matrix-transformations
func (v *I3) Transform(m [][]int) error {
	c, err := CrossProduct(m, [][]int{{v.X}, {v.Y}, {v.Z}})
	if err != nil {
		return err
	}

	v.X = c[0][0]
	v.Y = c[1][0]
	v.Z = c[2][0]

	return nil
}

// Translate this vector by the given amount.
func (v *I3) Translate(t I3) {
	v.X += t.X
	v.Y += t.Y
	v.Z += t.Z
}

// CrossProduct multiples matrix a by matrix b.
// Matrix a must be of size [m][n].
// Matrix b must be of size [n][p].
// The result is a matrix of size[m][p].
// Assumes neither matrix is empty.
func CrossProduct(a, b [][]int) ([][]int, error) {
	m, n, n1, p := len(a), len(a[0]), len(b), len(b[0])
	if n != n1 {
		return nil, errors.New("mismatched matrix sizes")
	}

	out := make([][]int, m)
	for x := 0; x < m; x++ {
		out[x] = make([]int, p)
		for y := 0; y < p; y++ {
			for i := 0; i < n; i++ {
				out[x][y] += a[x][i] * b[i][y]
			}
		}
	}

	return out, nil
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
