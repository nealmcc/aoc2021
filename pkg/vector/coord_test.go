package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReduce(t *testing.T) {
	tt := []struct {
		name      string
		in        Coord
		want      Coord
		wantScale int
	}{
		{
			name:      "when the zero vector reduced, the scale is 1",
			in:        Coord{},
			want:      Coord{},
			wantScale: 1,
		},
		{
			name:      "when both x and y are positive, the scale is positive",
			in:        Coord{12, 3},
			want:      Coord{4, 1},
			wantScale: 3,
		},
		{
			name:      "when both x and y are negative, the scale is positive",
			in:        Coord{-12, -3},
			want:      Coord{-4, -1},
			wantScale: 3,
		},
		{
			name:      "when x is negative and y is positive, the scale is positive",
			in:        Coord{-12, 3},
			want:      Coord{-4, 1},
			wantScale: 3,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)

			got, scale := Reduce(tc.in)
			a.Equal(tc.want, got)
			a.Equal(tc.wantScale, scale)
		})
	}
}
