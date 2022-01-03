package vector

import "errors"

// Matrix is used either to representing a linear transformation or a sequence
// of vector variables.
type Matrix [][]int

// Equals returns true iff matrices a and b contain the same elements in the
// same positions. Therefore a nil Matrix and an empty Matrix are considered equal.
func (a Matrix) Equals(b Matrix) bool {
	height := len(a)
	if len(b) != height {
		return false
	}
	if height == 0 {
		return true
	}

	for i := 0; i < height; i++ {
		aRow, bRow := a[i], b[i]
		width := len(aRow)
		if len(bRow) != width {
			return false
		}
		if width == 0 {
			continue
		}
		for j, el := range aRow {
			if bRow[j] != el {
				return false
			}
		}
	}
	return true
}

// Times calculates the matrix multiplication of a x b.
// Matrix a must be of size [m][n].
// Matrix b must be of size [n][p].
// The result is a matrix of size[m][p].
func (a Matrix) Times(b Matrix) (Matrix, error) {
	if len(a) == 0 {
		return nil, nil
	}

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

// Determinant calculates the determinant of the square matrix a.
// The determinant of a matrix gives some useful information about the effect
// of transforming a vector by a. In geometric terms, it measures the scaling
// and reflectivity of applying the transformation.
//
// see more:
// https://www.khanacademy.org/math/multivariable-calculus/thinking-about-multivariable-function/x786f2022:vectors-and-matrices/a/determinants-mvc
func (a Matrix) Determinant() (int, error) {
	size := len(a)
	if size == 0 {
		// by definition, the determinant of an empty matrix is 1.
		return 1, nil
	}

	if len(a[0]) != size {
		return 0, errors.New("matrix must be square")
	}

	subMatrix := func(i int) Matrix {
		rows := make(Matrix, 0, size-1)
		for n := 1; n < size; n++ {
			row := make([]int, size-1)
			copy(row[:i], a[n][:i])
			copy(row[i:], a[n][i+1:])
			rows = append(rows, row)
		}
		return rows
	}

	var sum int
	for i := 0; i < size; i++ {
		n := a[0][i]
		if i%2 == 1 {
			n *= -1
		}
		d2, _ := subMatrix(i).Determinant()
		sum += n * d2
	}

	return sum, nil
}
