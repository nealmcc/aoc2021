package vector

import (
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestCuboidSlice(t *testing.T) {
	tt := []struct {
		name    string
		box     Cuboid
		x, y, z int
		want    []Cuboid
	}{
		{
			name: "no intersections with the box => just return the box",
			box:  Cuboid{X2: 2, Y2: 2, Z2: 2},
			x:    -1,
			y:    -1,
			z:    -1,
			want: []Cuboid{{X2: 2, Y2: 2, Z2: 2}},
		},
		{
			name: "at the corner of a box => just return the box",
			box:  Cuboid{X2: 2, Y2: 2, Z2: 2},
			want: []Cuboid{{X2: 2, Y2: 2, Z2: 2}},
		},
		{
			name: "on the edge of a box => 2 pieces",
			box:  Cuboid{X2: 2, Y2: 2, Z2: 2},
			x:    1,
			want: []Cuboid{
				{X1: 0, Y1: 0, Z1: 0, X2: 1, Y2: 2, Z2: 2},
				{X1: 1, Y1: 0, Z1: 0, X2: 2, Y2: 2, Z2: 2},
			},
		},
		{
			name: "in the middle of one face of the box => 4 pieces",
			box:  Cuboid{X2: 2, Y2: 2, Z2: 2},
			x:    1,
			y:    0,
			z:    1,
			want: []Cuboid{
				{X1: 0, Y1: 0, Z1: 0, X2: 1, Y2: 2, Z2: 1},
				{X1: 1, Y1: 0, Z1: 0, X2: 2, Y2: 2, Z2: 1},
				{X1: 0, Y1: 0, Z1: 1, X2: 1, Y2: 2, Z2: 2},
				{X1: 1, Y1: 0, Z1: 1, X2: 2, Y2: 2, Z2: 2},
			},
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

			got := tc.box.Slice(tc.x, tc.y, tc.z)
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
			t.Parallel()
			got := tc.a.IsDisjoint(tc.b)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestCuboidOuterJoin(t *testing.T) {
	var (
		zero      = Cuboid{}
		zeroOne   = Cuboid{X2: 1, Y2: 1, Z2: 1}
		zeroTwo   = Cuboid{X2: 2, Y2: 2, Z2: 2}
		zeroThree = Cuboid{X2: 3, Y2: 3, Z2: 3}
		oneTwo    = Cuboid{X1: 1, Y1: 1, Z1: 1, X2: 2, Y2: 2, Z2: 2}
		oneFour   = Cuboid{X1: 1, Y1: 1, Z1: 1, X2: 4, Y2: 4, Z2: 4}
	)

	tt := []struct {
		name        string
		a, b        Cuboid
		intersect   Cuboid
		left, right []Cuboid
		volume      int
	}{
		{
			name:      "disjoint cubes",
			a:         zeroOne,
			b:         oneTwo,
			volume:    2,
			intersect: zero,
			left:      []Cuboid{zeroOne},
			right:     []Cuboid{oneTwo},
		},
		{
			name:      "one cube extends the other",
			a:         zeroThree,
			b:         zeroOne,
			volume:    27,
			intersect: zeroOne,
			left: []Cuboid{
				{X1: 0, Y1: 0, Z1: 1, X2: 1, Y2: 1, Z2: 3},
				{X1: 0, X2: 1, Y1: 1, Y2: 3, Z1: 0, Z2: 1},
				{X1: 0, X2: 1, Y1: 1, Y2: 3, Z1: 1, Z2: 3},
				{X1: 1, X2: 3, Y1: 0, Y2: 1, Z1: 0, Z2: 1},
				{X1: 1, X2: 3, Y1: 0, Y2: 1, Z1: 1, Z2: 3},
				{X1: 1, X2: 3, Y1: 1, Y2: 3, Z1: 0, Z2: 1},
				{X1: 1, X2: 3, Y1: 1, Y2: 3, Z1: 1, Z2: 3},
			},
			right: []Cuboid{},
		},
		{
			name:      "one cube surrounds the other",
			a:         zeroThree,
			b:         oneTwo,
			volume:    27,
			intersect: oneTwo,
			left: []Cuboid{
				{X1: 0, X2: 1, Y1: 0, Y2: 1, Z1: 0, Z2: 1},
				{X1: 0, X2: 1, Y1: 0, Y2: 1, Z1: 1, Z2: 2},
				{X1: 0, X2: 1, Y1: 0, Y2: 1, Z1: 2, Z2: 3},
				{X1: 0, X2: 1, Y1: 1, Y2: 2, Z1: 0, Z2: 1},
				{X1: 0, X2: 1, Y1: 2, Y2: 3, Z1: 0, Z2: 1},
				{X1: 0, X2: 1, Y1: 1, Y2: 2, Z1: 1, Z2: 2},
				{X1: 0, X2: 1, Y1: 1, Y2: 2, Z1: 2, Z2: 3},
				{X1: 0, X2: 1, Y1: 2, Y2: 3, Z1: 1, Z2: 2},
				{X1: 0, X2: 1, Y1: 2, Y2: 3, Z1: 2, Z2: 3},
				{X1: 1, X2: 2, Y1: 0, Y2: 1, Z1: 0, Z2: 1},
				{X1: 2, X2: 3, Y1: 0, Y2: 1, Z1: 0, Z2: 1},
				{X1: 1, X2: 2, Y1: 0, Y2: 1, Z1: 1, Z2: 2},
				{X1: 1, X2: 2, Y1: 0, Y2: 1, Z1: 2, Z2: 3},
				{X1: 2, X2: 3, Y1: 0, Y2: 1, Z1: 1, Z2: 2},
				{X1: 2, X2: 3, Y1: 0, Y2: 1, Z1: 2, Z2: 3},
				{X1: 1, X2: 2, Y1: 1, Y2: 2, Z1: 0, Z2: 1},
				{X1: 1, X2: 2, Y1: 2, Y2: 3, Z1: 0, Z2: 1},
				{X1: 2, X2: 3, Y1: 1, Y2: 2, Z1: 0, Z2: 1},
				{X1: 2, X2: 3, Y1: 2, Y2: 3, Z1: 0, Z2: 1},
				{X1: 1, X2: 2, Y1: 1, Y2: 2, Z1: 2, Z2: 3},
				{X1: 1, X2: 2, Y1: 2, Y2: 3, Z1: 1, Z2: 2},
				{X1: 1, X2: 2, Y1: 2, Y2: 3, Z1: 2, Z2: 3},
				{X1: 2, X2: 3, Y1: 1, Y2: 2, Z1: 1, Z2: 2},
				{X1: 2, X2: 3, Y1: 1, Y2: 2, Z1: 2, Z2: 3},
				{X1: 2, X2: 3, Y1: 2, Y2: 3, Z1: 1, Z2: 2},
				{X1: 2, X2: 3, Y1: 2, Y2: 3, Z1: 2, Z2: 3},
			},
			right: []Cuboid{},
		},
		{
			name:      "the cuboids partially overlap",
			a:         zeroTwo,
			b:         oneFour,
			volume:    2*2*2 + 3*3*3 - 1,
			intersect: oneTwo,
			left: []Cuboid{
				{X1: 0, X2: 1, Y1: 0, Y2: 1, Z1: 0, Z2: 1},
				{X1: 0, X2: 1, Y1: 0, Y2: 1, Z1: 1, Z2: 2},
				{X1: 0, X2: 1, Y1: 1, Y2: 2, Z1: 0, Z2: 1},
				{X1: 0, X2: 1, Y1: 1, Y2: 2, Z1: 1, Z2: 2},
				{X1: 1, X2: 2, Y1: 0, Y2: 1, Z1: 0, Z2: 1},
				{X1: 1, X2: 2, Y1: 0, Y2: 1, Z1: 1, Z2: 2},
				{X1: 1, X2: 2, Y1: 1, Y2: 2, Z1: 0, Z2: 1},
			},
			right: []Cuboid{
				{X1: 1, X2: 2, Y1: 1, Y2: 2, Z1: 2, Z2: 4},
				{X1: 1, X2: 2, Y1: 2, Y2: 4, Z1: 1, Z2: 2},
				{X1: 1, X2: 2, Y1: 2, Y2: 4, Z1: 2, Z2: 4},
				{X1: 2, X2: 4, Y1: 1, Y2: 2, Z1: 1, Z2: 2},
				{X1: 2, X2: 4, Y1: 1, Y2: 2, Z1: 2, Z2: 4},
				{X1: 2, X2: 4, Y1: 2, Y2: 4, Z1: 1, Z2: 2},
				{X1: 2, X2: 4, Y1: 2, Y2: 4, Z1: 2, Z2: 4},
			},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)

			intersect, left, right := CuboidOuterJoin(tc.a, tc.b)

			volume := intersect.Volume()
			for _, box := range left {
				volume += box.Volume()
			}
			for _, box := range right {
				volume += box.Volume()
			}

			a.Equal(tc.volume, volume)
			a.Equal(tc.intersect, intersect)
			a.Equal(tc.intersect, tc.a.Intersect(tc.b))

			a.ElementsMatch(tc.left, left, "left elements (a-b)")
			a.ElementsMatch(tc.right, right, "right elements (b-a)")
		})
	}
}
