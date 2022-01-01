package main

import (
	"strings"
	"testing"

	v "github.com/nealmcc/aoc2021/pkg/vector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var examples = []struct {
	name    string
	in      string
	sensors []sensor
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
		[]sensor{
			{
				id:  0,
				pos: v.I3{X: 0, Y: -1},
				beacons: beaconSet{
					extents: v.I3{X: 4, Y: 2},
					b: []v.I3{
						{X: 0, Y: 1},
						{X: 3, Y: 2},
						{X: 4, Y: 0},
					},
				},
			},
			{
				id:  1,
				pos: v.I3{X: 5, Y: 1},
				beacons: beaconSet{
					extents: v.I3{X: 4, Y: 2},
					b: []v.I3{
						{X: 0, Y: 1},
						{X: 3, Y: 2},
						{X: 4, Y: 0},
					},
				},
			},
		},
	},
}

func TestRead(t *testing.T) {
	for _, tc := range examples {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := read(strings.NewReader(tc.in))
			require.NoError(t, err)
			assert.Equal(t, tc.sensors, got)

			for _, sensor := range got {
				t.Log(sensor.String())
			}
		})
	}
}
