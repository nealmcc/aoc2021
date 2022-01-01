package vector

import (
	"errors"
)

// I3 is a vector in three-dimensional integer space.
type I3 struct {
	X, Y, Z int
}

// I3RotX90 rotates the given vector 90 degrees about the x-axis.
func I3RotX90(v I3) I3 {
	rot, _ := I3Transform(v, [][]int{
		{1, 0, 0},  // x = (1,       0,       0)
		{0, 0, -1}, // y = (0, cos(90), -sin(90))
		{0, 1, 0},  // z = (0, sin(90), cos(90))
	})
	return rot
}

// I3RotX90 rotates the given vector 90 degrees about the x-axis.
func I3RotY90(v I3) I3 {
	rot, _ := I3Transform(v, [][]int{
		{0, 0, 1},  // x = (cos(90),  0,  sin(90))
		{0, 1, 0},  // y = (      0,  1,       0)
		{-1, 0, 0}, // z = (-sin(90), 0, cos(90))
	})
	return rot
}

// I3RotZ90 rotates the given vector 90-degrees about the z axis.
func I3RotZ90(v I3) I3 {
	rot, _ := I3Transform(v, [][]int{
		{0, -1, 0}, // x = (cos(90), -sin(90), 0)
		{1, 0, 0},  // y = (sin(90),  cos(90), 0)
		{0, 0, 1},  // z = (      0,        0, 1)
	})
	return rot
}

// I3ReflectXY reflects the given vector about the x-y plane.
func I3ReflectXY(v I3) I3 {
	ref, _ := I3Transform(v, [][]int{
		{1, 0, 0},  // x
		{0, 1, 0},  // y
		{0, 0, -1}, // z
	})

	return ref
}

// I3Transform applies the given linear transformation m to the vector v.
// see: https://www.khanacademy.org/math/linear-algebra/matrix-transformations
func I3Transform(v I3, m [][]int) (I3, error) {
	c, err := CrossProduct(m, [][]int{{v.X}, {v.Y}, {v.Z}})
	if err != nil {
		return I3{}, err
	}

	return I3{
		X: c[0][0],
		Y: c[1][0],
		Z: c[2][0],
	}, nil
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
