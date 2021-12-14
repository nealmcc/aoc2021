package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var example string = `NNCB

CH -> B
HH -> N
CB -> H
NH -> C
HB -> C
HC -> B
HN -> C
NN -> C
BH -> H
NC -> B
NB -> B
BN -> B
BB -> N
BC -> B
CC -> N
CN -> C
`

func Test_read(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	polymer, rules, err := read(strings.NewReader(example))
	r.NoError(err)

	want := compound{
		pairs: map[pair]int{
			{'N', 'N'}: 1,
			{'N', 'C'}: 1,
			{'C', 'B'}: 1,
		},
		elements: map[byte]int{
			'N': 2,
			'C': 1,
			'B': 1,
		},
	}

	a.Equal(want, *polymer)
	a.Equal(16, len(rules))
}

func Test_part1(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	template, rules, err := read(strings.NewReader(example))
	r.NoError(err)

	a.Equal(1588, part1(template, rules))
}
