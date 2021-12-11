package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var example string = `[({(<(())[]>[[{[]{<()<>>
[(()[<>])]({[<{<<[]>>(
{([(<{}[<>[]}>{[]{[(<()>
(((({<>}<{<{<>}{[]{[]{}
[[<[([]))<([[{}[[()]]]
[{[{({}]{}}([{[{{{}}([]
{<[[]]>}<{[{[{[]{()[[[]
[<(<(<(<{}))><([]([]()
<{([([[(<>()){}]>(<<{{
<{([{{}}[<[[[<>{}]]]>[]]
`

func Test_part1(t *testing.T) {
	r := require.New(t)

	nav, err := read(strings.NewReader(example))
	r.NoError(err)

	got, _ := part1(nav)

	r.Equal(26397, got)
}

func Test_part2(t *testing.T) {
	tt := []struct {
		name      string
		in        string
		wantRest  string
		wantScore int
	}{
		{
			name:      "example 1",
			in:        "[({(<(())[]>[[{[]{<()<>>",
			wantRest:  "}}]])})]",
			wantScore: 288957,
		},
		{
			name:      "example 2",
			in:        "[(()[<>])]({[<{<<[]>>(",
			wantRest:  ")}>]})",
			wantScore: 5566,
		},
		{
			name:      "example 3",
			in:        "(((({<>}<{<{<>}{[]{[]{}",
			wantRest:  "}}>}>))))",
			wantScore: 1480781,
		},
		{
			name:      "example 4",
			in:        "{<[[]]>}<{[{[{[]{()[[[]",
			wantRest:  "]]}}]}]}>",
			wantScore: 995444,
		},
		{
			name:      "example 5",
			in:        "<{([{{}}[<[[[<>{}]]]>[]]",
			wantRest:  "])}>",
			wantScore: 294,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			l := line{raw: []bracket(tc.in)}
			err := l.parse()

			if _, ok := err.(ErrIncomplete); !ok {
				t.Logf("want ErrIncomplete, but got %s\n", err)
				t.FailNow()
			}

			remainder := l.suggest()
			assert.Equal(t, tc.wantRest, string(remainder))
			p2 := part2([]line{{raw: []byte(tc.in)}})
			assert.Equal(t, tc.wantScore, p2)
		})
	}
}
