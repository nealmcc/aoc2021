package main

import (
	"os"
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

func Test_part1(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	molecule, err := read(strings.NewReader(example))
	r.NoError(err)

	for i := 0; i < 10; i++ {
		molecule.replicate()
	}

	a.Equal(1588, molecule.magicNumber())
}

// prevent the compiler from optimising away the magicNumber() call
var n int

func Benchmark_part2(b *testing.B) {
	in, err := os.Open("input.txt")
	if err != nil {
		b.Log(err)
		b.FailNow()
	}
	defer in.Close()
	polymer, err := read(in)
	if err != nil {
		b.Log(err)
		b.FailNow()
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 40; j++ {
			polymer.replicate()
		}
		n = polymer.magicNumber()
	}
}
