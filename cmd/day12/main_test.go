package main

import (
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

	r.Equal(6, len(caves))
}

func Test_part1(t *testing.T) {
	r := require.New(t)

	caves, err := read(strings.NewReader(example))
	r.NoError(err)

	got := part1(caves)
	r.Equal(10, got)
}

func Test_pathsFrom1(t *testing.T) {
	tt := []struct {
		name       string
		caves      string
		start, end string
		rules      caveRules
		want       []string
	}{
		{
			name:  "two room cave",
			caves: "start-end",
			start: "start",
			end:   "end",
			rules: &p1Rules{visited: make(map[string]bool)},
			want:  []string{"start,end"},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r := require.New(t)
			caves, err := read(strings.NewReader(tc.caves))
			r.NoError(err)
			got := pathsFrom(caves, caves[tc.start], caves[tc.end], tc.rules)
			r.ElementsMatch(tc.want, got)
		})
	}
}
