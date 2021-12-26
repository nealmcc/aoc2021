package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var examples = []struct {
	name    string
	in      string
	boot    []instruction
	after   []instruction
	cubesOn int
}{
	{
		name: "reboot steps",
		in: `on x=10..12,y=10..12,z=10..12
on x=11..13,y=11..13,z=11..13
off x=9..11,y=9..11,z=9..11
on x=10..10,y=10..10,z=10..10
`,
		boot: []instruction{
			{on: true, x1: 10, x2: 12, y1: 10, y2: 12, z1: 10, z2: 12},
			{on: true, x1: 11, x2: 13, y1: 11, y2: 13, z1: 11, z2: 13},
			{on: false, x1: 9, x2: 11, y1: 9, y2: 11, z1: 9, z2: 11},
			{on: true, x1: 10, x2: 10, y1: 10, y2: 10, z1: 10, z2: 10},
		},
		after:   []instruction{},
		cubesOn: 39,
	},
	{
		name: "a larger example",
		in: `on x=-20..26,y=-36..17,z=-47..7
on x=-20..33,y=-21..23,z=-26..28
on x=-22..28,y=-29..23,z=-38..16
on x=-46..7,y=-6..46,z=-50..-1
on x=-49..1,y=-3..46,z=-24..28
on x=2..47,y=-22..22,z=-23..27
on x=-27..23,y=-28..26,z=-21..29
on x=-39..5,y=-6..47,z=-3..44
on x=-30..21,y=-8..43,z=-13..34
on x=-22..26,y=-27..20,z=-29..19
off x=-48..-32,y=26..41,z=-47..-37
on x=-12..35,y=6..50,z=-50..-2
off x=-48..-32,y=-32..-16,z=-15..-5
on x=-18..26,y=-33..15,z=-7..46
off x=-40..-22,y=-38..-28,z=23..41
on x=-16..35,y=-41..10,z=-47..6
off x=-32..-23,y=11..30,z=-14..3
on x=-49..-5,y=-3..45,z=-29..18
off x=18..30,y=-20..-8,z=-3..13
on x=-41..9,y=-7..43,z=-33..15
on x=-54112..-39298,y=-85059..-49293,z=-27449..7877
on x=967..23432,y=45373..81175,z=27513..53682
`,
		boot: []instruction{
			{on: true, x1: -20, x2: 26, y1: -36, y2: 17, z1: -47, z2: 7},
			{on: true, x1: -20, x2: 33, y1: -21, y2: 23, z1: -26, z2: 28},
			{on: true, x1: -22, x2: 28, y1: -29, y2: 23, z1: -38, z2: 16},
			{on: true, x1: -46, x2: 7, y1: -6, y2: 46, z1: -50, z2: -1},
			{on: true, x1: -49, x2: 1, y1: -3, y2: 46, z1: -24, z2: 28},
			{on: true, x1: 2, x2: 47, y1: -22, y2: 22, z1: -23, z2: 27},
			{on: true, x1: -27, x2: 23, y1: -28, y2: 26, z1: -21, z2: 29},
			{on: true, x1: -39, x2: 5, y1: -6, y2: 47, z1: -3, z2: 44},
			{on: true, x1: -30, x2: 21, y1: -8, y2: 43, z1: -13, z2: 34},
			{on: true, x1: -22, x2: 26, y1: -27, y2: 20, z1: -29, z2: 19},
			{on: false, x1: -48, x2: -32, y1: 26, y2: 41, z1: -47, z2: -37},
			{on: true, x1: -12, x2: 35, y1: 6, y2: 50, z1: -50, z2: -2},
			{on: false, x1: -48, x2: -32, y1: -32, y2: -16, z1: -15, z2: -5},
			{on: true, x1: -18, x2: 26, y1: -33, y2: 15, z1: -7, z2: 46},
			{on: false, x1: -40, x2: -22, y1: -38, y2: -28, z1: 23, z2: 41},
			{on: true, x1: -16, x2: 35, y1: -41, y2: 10, z1: -47, z2: 6},
			{on: false, x1: -32, x2: -23, y1: 11, y2: 30, z1: -14, z2: 3},
			{on: true, x1: -49, x2: -5, y1: -3, y2: 45, z1: -29, z2: 18},
			{on: false, x1: 18, x2: 30, y1: -20, y2: -8, z1: -3, z2: 13},
			{on: true, x1: -41, x2: 9, y1: -7, y2: 43, z1: -33, z2: 15},
		},
		after: []instruction{
			{on: true, x1: -54112, x2: -39298, y1: -85059, y2: -49293, z1: -27449, z2: 7877},
			{on: true, x1: 967, x2: 23432, y1: 45373, y2: 81175, z1: 27513, z2: 53682},
		},
		cubesOn: 590784,
	},
}

func TestRead(t *testing.T) {
	for _, tc := range examples {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, a := require.New(t), assert.New(t)
			boot, after, err := read(strings.NewReader(tc.in))
			r.NoError(err)

			a.Equal(tc.boot, boot)
			a.Equal(tc.after, after)
		})
	}
}

func TestPart1(t *testing.T) {
	for _, tc := range examples {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := part1(tc.boot)
			assert.Equal(t, tc.cubesOn, got)
		})
	}
}
