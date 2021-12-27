package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCuboidVolume(t *testing.T) {
	tt := []struct {
		name string
		box  Cuboid
		want int
	}{
		{
			"a cube from the origin",
			Cuboid{X2: 10, Y2: 10, Z2: 10},
			1000,
		},
		{
			"a cube with negative width has positive volume",
			Cuboid{X2: -10, Y2: 10, Z2: 10},
			1000,
		},
		{
			"a cube with non-zero postions",
			Cuboid{X1: -10, X2: 10, Y1: -10, Y2: 10, Z1: -10, Z2: 10},
			8000,
		},
		{
			"a rectangle has no volume",
			Cuboid{X1: 10, X2: 10, Y1: -10, Y2: 10, Z1: -10, Z2: 10},
			0,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := tc.box.Volume()
			if got != tc.want {
				t.Logf("%#v Volume() = %d ; want %d", tc.box, got, tc.want)
				t.Fail()
			}
		})
	}
}

func TestCuboidExplode(t *testing.T) {
	tt := []struct {
		name    string
		box     Cuboid
		x, y, z int
		want    []Cuboid
		wantErr bool
	}{
		{
			name:    "outisde a box",
			box:     Cuboid{X2: 2, Y2: 2, Z2: 2},
			x:       -1,
			y:       1,
			z:       1,
			wantErr: true,
		},
		{
			name:    "on the surface of a box",
			box:     Cuboid{X2: 2, Y2: 2, Z2: 2},
			x:       1,
			y:       0,
			z:       1,
			wantErr: true,
		},
		{
			name:    "on the edge of a box",
			box:     Cuboid{X2: 2, Y2: 2, Z2: 2},
			x:       1,
			y:       0,
			z:       0,
			wantErr: true,
		},
		{
			name:    "at the corner of a box",
			box:     Cuboid{X2: 2, Y2: 2, Z2: 2},
			wantErr: true,
		},
		{
			name: "inside a box",
			box:  Cuboid{X2: 4, Y2: 4, Z2: 4},
			x:    2,
			y:    2,
			z:    2,
			want: []Cuboid{
				{
					X1: 0, X2: 2,
					Y1: 0, Y2: 2,
					Z1: 0, Z2: 2,
				},
				{
					X1: 2, X2: 4,
					Y1: 0, Y2: 2,
					Z1: 0, Z2: 2,
				},
				{
					X1: 0, X2: 2,
					Y1: 2, Y2: 4,
					Z1: 0, Z2: 2,
				},
				{
					X1: 0, X2: 2,
					Y1: 0, Y2: 2,
					Z1: 2, Z2: 4,
				},
				{
					X1: 0, X2: 2,
					Y1: 2, Y2: 4,
					Z1: 2, Z2: 4,
				},
				{
					X1: 2, X2: 4,
					Y1: 0, Y2: 2,
					Z1: 2, Z2: 4,
				},
				{
					X1: 2, X2: 4,
					Y1: 2, Y2: 4,
					Z1: 0, Z2: 2,
				},
				{
					X1: 2, X2: 4,
					Y1: 2, Y2: 4,
					Z1: 2, Z2: 4,
				},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := tc.box.Explode(tc.x, tc.y, tc.z)
			if tc.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.ElementsMatch(t, tc.want, got)
		})
	}
}

func TestCuboidIsDisjoint(t *testing.T) {
	var (
		cube10    = Cuboid{X2: 10, Y2: 10, Z2: 10}
		cubeNeg10 = Cuboid{X2: -10, Y2: -10, Z2: -10}
	)

	tt := []struct {
		name string
		a, b Cuboid
		want bool
	}{
		{
			name: "the same cube",
			a:    cube10,
			b:    cube10,
			want: false,
		},
		{
			name: "cubes that share a vertex, but nothing else",
			a:    cube10,
			b:    cubeNeg10,
			want: true,
		},
		{
			name: "cubes that share y and z planes, but are separate on the x-axis",
			a:    cube10,
			b:    Cuboid{X1: 15, X2: 25, Z1: 5, Z2: 15, Y1: 5, Y2: 15},
			want: true,
		},
		{
			name: "cubes that share x and y planes, but are separate on the z-axis",
			a:    cube10,
			b:    Cuboid{X1: 5, X2: 15, Y1: 5, Y2: 15, Z1: -15, Z2: -25},
			want: true,
		},
		{
			name: "cubes that share x and z planes, but are separate on the y-axis",
			a:    cube10,
			b:    Cuboid{X1: 5, X2: 15, Y1: -15, Y2: -25, Z1: 5, Z2: 15},
			want: true,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := tc.a.IsDisjoint(tc.b)
			assert.Equal(t, tc.want, got)
		})
	}
}
