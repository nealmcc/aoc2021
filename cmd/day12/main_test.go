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

func Test_part2(t *testing.T) {
	r := require.New(t)

	caves, err := read(strings.NewReader(example))
	r.NoError(err)

	got := part2(caves)
	r.Equal(36, got)
}

func Test_pathsFrom(t *testing.T) {
	tt := []struct {
		name       string
		caves      string
		start, end string
		rules      caveRules
		want       []string
	}{
		{
			name:  "two caves, p1",
			caves: "start-end",
			start: "start",
			end:   "end",
			rules: &p1Rules{visited: make(map[string]bool)},
			want:  []string{"start,end"},
		},
		{
			name:  "two caves, p2",
			caves: "start-end",
			start: "start",
			end:   "end",
			rules: &p2Rules{visited: make(map[string]bool)},
			want:  []string{"start,end"},
		},
		{
			name: "three caves with upper, p2",
			caves: `start-end
start-A
end-A
`,
			start: "start",
			end:   "end",
			rules: &p2Rules{visited: make(map[string]bool)},
			want:  []string{"start,end", "start,A,end"},
		},
		{
			name:  "example",
			caves: example,
			start: "start",
			end:   "end",
			rules: &p2Rules{visited: make(map[string]bool)},
			want: []string{
				"start,A,b,A,b,A,c,A,end",
				"start,A,b,A,b,A,end",
				"start,A,b,A,b,end",
				"start,A,b,A,c,A,b,A,end",
				"start,A,b,A,c,A,b,end",
				"start,A,b,A,c,A,c,A,end",
				"start,A,b,A,c,A,end",
				"start,A,b,A,end",
				"start,A,b,d,b,A,c,A,end",
				"start,A,b,d,b,A,end",
				"start,A,b,d,b,end",
				"start,A,b,end",
				"start,A,c,A,b,A,b,A,end",
				"start,A,c,A,b,A,b,end",
				"start,A,c,A,b,A,c,A,end",
				"start,A,c,A,b,A,end",
				"start,A,c,A,b,d,b,A,end",
				"start,A,c,A,b,d,b,end",
				"start,A,c,A,b,end",
				"start,A,c,A,c,A,b,A,end",
				"start,A,c,A,c,A,b,end",
				"start,A,c,A,c,A,end",
				"start,A,c,A,end",
				"start,A,end",
				"start,b,A,b,A,c,A,end",
				"start,b,A,b,A,end",
				"start,b,A,b,end",
				"start,b,A,c,A,b,A,end",
				"start,b,A,c,A,b,end",
				"start,b,A,c,A,c,A,end",
				"start,b,A,c,A,end",
				"start,b,A,end",
				"start,b,d,b,A,c,A,end",
				"start,b,d,b,A,end",
				"start,b,d,b,end",
				"start,b,end",
			},
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
