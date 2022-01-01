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
				id:     0,
				facing: identity(),
				beacons: []v.I3{
					{X: 0, Y: 2, Z: 0},
					{X: 4, Y: 1, Z: 0},
					{X: 3, Y: 3, Z: 0},
				},
			},
			{
				id:     1,
				facing: identity(),
				beacons: []v.I3{
					{X: -1, Y: -1, Z: 0},
					{X: -5, Y: 0, Z: 0},
					{X: -2, Y: 1, Z: 0},
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

			for _, s := range got {
				t.Log(s.String())
			}
		})
	}
}
