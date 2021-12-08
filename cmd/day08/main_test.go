package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var example string = `be cfbegad cbdgef fgaecd cgeb fdcge agebfd fecdb fabcd edb | fdgacbe cefdb cefbgd gcbe
edbfga begcd cbg gc gcadebf fbgde acbgfd abcde gfcbed gfec | fcgedb cgb dgebacf gc
fgaebd cg bdaec gdafb agbcfd gdcbef bgcad gfac gcb cdgabef | cg cg fdcagb cbg
fbegcd cbd adcefb dageb afcb bc aefdc ecdab fgdeca fcdbega | efabcd cedba gadfec cb
aecbfdg fbg gf bafeg dbefa fcge gcbea fcaegb dgceab fcbdga | gecf egdcabf bgf bfgea
fgeab ca afcebg bdacfeg cfaedg gcfdb baec bfadeg bafgc acf | gebdcfa ecba ca fadegcb
dbcfg fgd bdegcaf fgec aegbdf ecdfab fbedc dacgb gdcebf gf | cefg dcbef fcge gbcadfe
bdfegc cbegaf gecbf dfcage bdacg ed bedf ced adcbefg gebcd | ed bcgafe cdgba cbgef
egadfb cdbfeg cegd fecab cgb gbdefca cg fgcdab egfdb bfceg | gbdfcae bgc cg cgb
gcafb gcf dcaebfg ecagb gf abcdeg gaef cafbge fdbac fegbdc | fgae cfgab fg bagce
`

func Test_read(t *testing.T) {
	r := require.New(t)

	got, err := read(strings.NewReader(example))
	r.NoError(err)

	r.Equal(10, len(got))
	r.Equal(display{
		signals: []signal{
			"be",
			"abcdefg",
			"bcdefg",
			"acdefg",
			"bceg",
			"cdefg",
			"abdefg",
			"bcdef",
			"abcdf",
			"bde",
		},
		digits: map[signal]int{},
		output: []signal{"abcdefg", "bcdef", "bcdefg", "bceg"},
	}, got[0])
}

func Test_part1(t *testing.T) {
	r := require.New(t)

	displays, err := read(strings.NewReader(example))
	r.NoError(err)

	got := part1(displays)

	r.Equal(26, got)
}

func Test_solve(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	displays, err := read(strings.NewReader("acedgfb cdfbe gcdfa fbcad dab cefabd cdfgeb eafb cagedb ab | cdfeb fcadb cdfeb cdbaf"))
	r.NoError(err)

	part1(displays)
	d := displays[0]
	solve(d)

	a.Equal(map[signal]int{
		normal("acedgfb"): 8,
		normal("cdfbe"):   5,
		normal("gcdfa"):   2,
		normal("fbcad"):   3,
		normal("dab"):     7,
		normal("cefabd"):  9,
		normal("cdfgeb"):  6,
		normal("eafb"):    4,
		normal("cagedb"):  0,
		normal("ab"):      1,
	}, d.digits)

	r.Equal(5353, d.value())
}

func Test_part2(t *testing.T) {
	r := require.New(t)

	displays, err := read(strings.NewReader(example))
	r.NoError(err)

	part1(displays)
	p2 := part2(displays)

	r.Equal(61229, p2)
}
