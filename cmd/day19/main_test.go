package main

import (
	"fmt"
	"sort"
	"strings"
	"testing"

	v "github.com/nealmcc/aoc2021/pkg/vector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRead(t *testing.T) {
	example := []struct {
		name  string
		in    string
		boxes []block
	}{
		{
			"first example",
			`--- scanner 0 ---
0,2,0
4,1,0
3,3,0

--- scanner 1 ---
-1,-1,0
-5,0,0
-2,1,0
`,
			[]block{
				{
					Sensors: map[int]v.I3{0: {}},
					Beacons: []v.I3{
						{X: 0, Y: 2},
						{X: 3, Y: 3},
						{X: 4, Y: 1},
					},
				},
				{
					Sensors: map[int]v.I3{1: {}},
					Beacons: []v.I3{
						{X: -5, Y: 0},
						{X: -2, Y: 1},
						{X: -1, Y: -1},
					},
				},
			},
		},
	}
	for _, tc := range example {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := read(strings.NewReader(tc.in))
			require.NoError(t, err)
			for i, box := range got {
				sort.Sort(byXYZ(box.Beacons))
				assert.Equal(t, tc.boxes[i], box)
			}
		})
	}
}

func TestBeaconsMatch(t *testing.T) {
	var (
		box1 = block{
			Sensors: map[int]v.I3{0: {}},
			Beacons: []v.I3{
				{X: 0, Y: 2},
				{X: 3, Y: 3},
				{X: 4, Y: 1},
			},
		}

		box2 = block{
			Sensors: map[int]v.I3{1: {}},
			Beacons: []v.I3{
				{X: -5, Y: 0},
				{X: -2, Y: 1},
				{X: -1, Y: -1},
			},
		}

		extraBeacon = block{
			Sensors: map[int]v.I3{2: {}},
			Beacons: append(box1.Beacons, v.I3{X: -1, Y: -1, Z: -1}),
		}
	)
	sort.Sort(byXYZ(extraBeacon.Beacons))

	tt := []struct {
		name       string
		a, b       block
		minQty     int
		wantOffset v.I3
		wantMatch  bool
	}{

		{
			name:       "example 0",
			a:          box1,
			b:          box2,
			minQty:     3,
			wantOffset: v.I3{X: 5, Y: 2},
			wantMatch:  true,
		},
		{
			name:      "insufficient matches => false",
			a:         box1,
			b:         box2,
			minQty:    4,
			wantMatch: false,
		},
		{
			name:       "an extra beacon in a => true",
			a:          extraBeacon,
			b:          box2,
			minQty:     3,
			wantOffset: v.I3{X: 5, Y: 2},
			wantMatch:  true,
		},
		{
			name:       "an extra beacon in b => true",
			a:          box2,
			b:          extraBeacon,
			minQty:     3,
			wantOffset: v.I3{X: -5, Y: -2},
			wantMatch:  true,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			offset, isMatch := beaconsMatch(tc.a.Beacons, tc.b.Beacons, tc.minQty)
			assert.Equal(t, tc.wantOffset, offset)
			assert.Equal(t, tc.wantMatch, isMatch)
		})
	}
}

func TestRotations(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	unique := make(map[v.I3]int, 24)
	vect := v.I3{X: 2, Y: 3, Z: 5}

	for i, rot := range getRotations() {
		det, err := rot.Determinant()
		r.NoError(err)
		r.Equal(1, det, "i=%d: rotations must preserve scale and signedness", i)

		vect2, err := vect.Transform(rot)
		r.NoError(err)

		if prev, exists := unique[vect2]; exists {
			t.Logf("\nalready found rotation %d:\n\t%v\n\t%v\n\t%v",
				i, rot[0], rot[1], rot[2])
			t.Logf("previously found at i: %d", prev)
			t.Fail()
		}

		unique[vect2] = i
	}

	a.Equal(24, len(unique))
}

func TestRotateToMatch(t *testing.T) {
	boxes, err := read(strings.NewReader(`--- scanner 0 ---
-1,-1,1
-2,-2,2
-3,-3,3
-2,-3,1
5,6,-4
8,0,7

--- scanner 0 ---
1,-1,1
2,-2,2
3,-3,3
2,-1,3
-5,4,-6
-8,-7,0

--- scanner 0 ---
-1,-1,-1
-2,-2,-2
-3,-3,-3
-1,-3,-2
4,6,5
-7,0,8

--- scanner 0 ---
1,1,-1
2,2,-2
3,3,-3
1,3,-2
-4,-6,5
7,0,8

--- scanner 0 ---
1,1,1
2,2,2
3,3,3
3,1,2
-6,-4,-5
0,7,-8
`))

	require.NoError(t, err)

	for i := 0; i < len(boxes)-1; i++ {
		for j := i + 1; j < len(boxes); j++ {
			t.Run(fmt.Sprintf("%d_%d", i, j), func(t *testing.T) {
				a := assert.New(t)
				box2, ok := rotateToMatch(boxes[i], boxes[j], 6)
				a.Truef(ok, "boxes %d and %d should match", i, j)
				if ok {
					t.Log("box1", boxes[i])
					t.Log("box2", box2)
				}
			})
		}
	}
}
