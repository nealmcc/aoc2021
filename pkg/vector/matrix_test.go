package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTimes(t *testing.T) {
	tt := []struct {
		name    string
		a, b    Matrix
		want    Matrix
		wantErr bool
	}{
		{
			name: "the zero matrix times any matrix is the zero matrix",
			a:    Matrix{},
			b: Matrix{
				{1, 2, 3},
			},
		},
		{
			name: "height of a must match width of b",
			a: Matrix{
				{1},
				{2},
			},
			b: Matrix{
				{3},
				{4},
			},
			wantErr: true,
		},
		{
			name: "any matrix times the identity matrix yeilds the original",
			a: Matrix{
				{2, 3, 5},
				{-7, 11, -13},
				{7, 11, 13},
			},
			b: Matrix{
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, 1},
			},
			want: Matrix{
				{2, 3, 5},
				{-7, 11, -13},
				{7, 11, 13},
			},
		},
		{
			name: "a 2x3 matrix times a 3x2 matrix produces a 3x3 matrix",
			a: Matrix{
				{2, 3},
				{5, 7},
				{11, 13},
			},
			b: Matrix{
				{1, 3, 5},
				{2, 4, 6},
			},
			want: Matrix{
				{2 + 6, 6 + 12, 10 + 18},
				{5 + 14, 15 + 28, 25 + 42},
				{11 + 26, 33 + 52, 55 + 78},
			},
		},
		{
			name: "a matrix times a vector transforms that vector",
			a: Matrix{
				{2, 0, 0},
				{0, 2, 0},
				{0, 0, -2},
			},
			b: Matrix{
				{3},
				{5},
				{7},
			},
			want: Matrix{
				{6},
				{10},
				{-14},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := tc.a.Times(tc.b)
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestDeterminant(t *testing.T) {
	tt := []struct {
		name    string
		in      Matrix
		want    int
		wantErr bool
	}{
		{
			name: "the determinant of an empty matrix is 1",
			in:   Matrix{},
			want: 1,
		},
		{
			name:    "a non-square matrix has no determinant",
			in:      Matrix{{1, 2}},
			wantErr: true,
		},
		{
			name: "a 2x2 matrix",
			in: Matrix{
				{1, 2},
				{3, 4},
			},
			want: -2,
		},
		{
			name: "a 3x3 matrix",
			in: Matrix{
				{4, -1, 1},
				{4, 5, 3},
				{-2, 0, 0},
			},
			want: 16,
		},
		{
			name: "a rotation has a determinant of 1",
			in: Matrix{
				// rotation 90 degrees about the z axis
				{0, 1, 0},
				{-1, 0, 0},
				{0, 0, 1},
			},
			want: 1,
		},
		{
			name: "a reflection has a determinant of -1",
			in: Matrix{
				// reflection through the x-y plane
				{1, 0, 0},
				{0, 1, 0},
				{0, 0, -1},
			},
			want: -1,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := tc.in.Determinant()
			if tc.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tc.want, got)
		})
	}
}
