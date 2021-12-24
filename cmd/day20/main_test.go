package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _testExample = struct {
	in        string
	alg       algorithm
	img       image
	enhance2x string
}{
	in: `..#.#..#####.#.#.#.###.##.....###.##.#..###.####..#####..#....#..#..##..###..######.###...####..#..#####..##..#.#####...##.#.#..#.##..#.#......#.###.######.###.####...#.##.##..#..#..#####.....#.#....###..#.##......#.....#..#..#..##..#...##.######.####.####.#.#...#.......#..#.#.#...####.##.#......#..#...##.#.##..#...##.#.##..###.#......#.#.......#.#.#.####.###.##...#.....####.#..#..#.##.#....##..#.####....##...##..#...#......#.#.......#.......##..####..#...#.#.#...##..#.#..###..#####........#..####......#..#

#..#.
#....
##..#
..#..
..###
`,
	alg: algorithm{
		0b_0010100111110101010111011000001110110100111011110011111001000010,
		0b_0100110011100111111011100011110010011111001100101111100011010100,
		0b_1011001010000001011101111110111011110001011011001001001111100000,
		0b_1010000111001011000000100000100100100110010001101111110111101111,
		0b_0101000100000001001010100011110110100000010010001101011001000110,
		0b_1011001110100000010100000001010101111011101100010000011110100100,
		0b_1011010000110010111100001100011001000100000010100000001000000011,
		0b_0011110010001010100011001010011100111110000000010011110000001001,
	},
	img: image{
		size:       5,
		numLit:     10,
		infinitePx: false,
		pixels: [][]bool{
			{true, false, false, true, false},
			{true, false, false, false, false},
			{true, true, false, false, true},
			{false, false, true, false, false},
			{false, false, true, true, true},
		},
	},
	enhance2x: `.......#.
.#..#.#..
#.#...###
#...##.#.
#.....#.#
.#.#####.
..#.#####
...##.##.
....###..
`,
}

func TestRead(t *testing.T) {
	r, a := require.New(t), assert.New(t)
	alg, img, err := read(strings.NewReader(_testExample.in))
	r.NoError(err)

	a.Equal(_testExample.alg, alg)
	a.Equal(_testExample.img, img)
}

func TestPart1(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	alg, img, err := read(strings.NewReader(_testExample.in))
	r.NoError(err)

	a.Equal(35, part1(alg, img))
}
