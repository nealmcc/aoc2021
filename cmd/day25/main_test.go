package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var examples = []struct {
	name   string
	text   string
	sea    seafloor
	after4 seafloor
}{
	{
		"example 1",
		`v...>>.vv>
.vv>>.vv..
>>.>v>...v
>>v>>.>.v.
v>v.vv.v..
>.>>..v...
.vv..>.>v.
v.v..>>v.v
....v..v.>
`,
		seafloor{
			step:   0,
			width:  10,
			height: 9,
			cukes: [][]byte{
				[]byte("v...>>.vv>"),
				[]byte(".vv>>.vv.."),
				[]byte(">>.>v>...v"),
				[]byte(">>v>>.>.v."),
				[]byte("v>v.vv.v.."),
				[]byte(">.>>..v..."),
				[]byte(".vv..>.>v."),
				[]byte("v.v..>>v.v"),
				[]byte("....v..v.>"),
			},
		},
		seafloor{
			step:   4,
			width:  10,
			height: 9,
			cukes: [][]byte{
				[]byte("v>..v.>>.."),
				[]byte("v.v.>.>.v."),
				[]byte(">vv.>>.v>v"),
				[]byte(">>.>..v>.>"),
				[]byte("..v>v...v."),
				[]byte("..>>.>vv.."),
				[]byte(">.v.vv>v.v"),
				[]byte(".....>>vv."),
				[]byte("vvv>...v.."),
			},
		},
	},
	{
		"wrapping example",
		`...>...
.......
......>
v.....>
......>
.......
..vvv..
`,
		seafloor{
			step:   0,
			width:  7,
			height: 7,
			cukes: [][]byte{
				[]byte("...>..."),
				[]byte("......."),
				[]byte("......>"),
				[]byte("v.....>"),
				[]byte("......>"),
				[]byte("......."),
				[]byte("..vvv.."),
			},
		},
		seafloor{
			step:   4,
			width:  7,
			height: 7,
			cukes: [][]byte{
				[]byte(">......"),
				[]byte("..v...."),
				[]byte("..>.v.."),
				[]byte(".>.v..."),
				[]byte("...>..."),
				[]byte("......."),
				[]byte("v......"),
			},
		},
	},
}

func TestRead(t *testing.T) {
	for _, tc := range examples {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got, err := read(strings.NewReader(tc.text))
			require.NoError(t, err)

			assert.Equal(t, tc.sea, got)
		})
	}
}

func TestMove4(t *testing.T) {
	for _, tc := range examples {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			sea, err := read(strings.NewReader(tc.text))
			require.NoError(t, err)

			for i := 0; i < 4; i++ {
				sea.move()
			}
			assert.Equal(t, tc.after4, sea)
		})
	}
}
