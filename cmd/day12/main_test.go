package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

var example string = `start-A
start-b
A-c
A-b
b-d
A-end
b-end
`

func Test_read(t *testing.T) {
	r := require.New(t)

	caves, err := read(strings.NewReader(example))
	r.NoError(err)

	r.Equal(6, len(caves.names))
	r.Equal(6, len(caves.nodes))
}

func Test_part1(t *testing.T) {
	r := require.New(t)

	caves, err := read(strings.NewReader(example))
	r.NoError(err)

	got := part1(caves)
	r.Equal(10, got)
}

func Test_part2(t *testing.T) {
	r := require.New(t)

	caves, err := read(strings.NewReader(example))
	r.NoError(err)

	got := part2(caves)
	r.Equal(36, got)
}

func Test_countPaths(t *testing.T) {
	tt := []struct {
		name       string
		caves      string
		start, end string
		useP1      bool
		want       int
	}{
		{
			name:  "two caves, p1",
			caves: "start-end",
			start: "start",
			end:   "end",
			useP1: true,
			want:  1,
		},
		{
			name:  "two caves, p2",
			caves: "start-end",
			start: "start",
			end:   "end",
			want:  1,
		},
		{
			name: "three caves with upper, p2",
			caves: `start-end
start-A
end-A
`,
			start: "start",
			end:   "end",
			want:  2,
		},
		{
			name:  "example",
			caves: example,
			start: "start",
			end:   "end",
			want:  36,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			// t.Parallel()
			r := require.New(t)
			caves, err := read(strings.NewReader(tc.caves))
			r.NoError(err)
			var rules caveRules
			if tc.useP1 {
				rules = newP1(len(caves.names))
			} else {
				rules = newP2(len(caves.names))
			}

			start := caves.getOrCreate(tc.start)
			end := caves.getOrCreate(tc.end)

			got := countPaths(caves, start, end, rules)
			r.Equal(tc.want, got)
		})
	}
}

func Benchmark_countPaths(b *testing.B) {
	in, _ := os.Open("input.txt")
	defer in.Close()

	caves, _ := read(in)

	for n := 0; n < b.N; n++ {
		rules := newP2(len(caves.names))
		start := caves.getOrCreate("start")
		end := caves.getOrCreate("end")
		countPaths(caves, start, end, rules)
	}
}
