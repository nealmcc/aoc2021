package main

import (
	"fmt"
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
	alg: "..#.#..#####.#.#.#.###.##.....###.##.#..###.####..#####..#....#..#..##..###..######.###...####..#..#####..##..#.#####...##.#.#..#.##..#.#......#.###.######.###.####...#.##.##..#..#..#####.....#.#....###..#.##......#.....#..#..#..##..#...##.######.####.####.#.#...#.......#..#.#.#...####.##.#......#..#...##.#.##..#...##.#.##..###.#......#.#.......#.#.#.####.###.##...#.....####.#..#..#.##.#....##..#.####....##...##..#...#......#.#.......#.......##..####..#...#.#.#...##..#.#..###..#####........#..####......#..#",
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

	p1 := part1(alg, img)

	a.Equal(35, p1.numLit)
}

func TestFormat(t *testing.T) {
	got := fmt.Sprintf("%v", _testExample.img)

	assert.Equal(t, `#..#.
#....
##..#
..#..
..###
`, got)
}

func TestEnhance(t *testing.T) {
	alg, img := _testExample.alg, _testExample.img

	img1 := alg.enhance(img)
	img2 := alg.enhance(img1)

	want, err := newImage([]byte(_testExample.enhance2x))
	require.NoError(t, err)

	assert.Equal(t, want, img2)
}

func TestEnhance_Tricky(t *testing.T) {
	// trickyAlg is tricky because the 0th value is '#'.  This means that when
	// enhancing an infinite image, if the source image is *not* infinitely lit,
	// then the target image will be.  Also, because the 511th value is '.' the
	// converse is true. When enhancing an image that *is* infinitely lit, the
	// resulting image will not be.
	trickyAlg := algorithm("#..#.#######.##...##.##.#.#..#..#.....####...####.##.###...##.####......##.###.##...#..##..#.######...###..########.#.##.#.#..#..##.##..####.###.###..#...##.##.###.....###..###....#.####.#..##....#.##...##.#..#.....###.#..#.....##..##.#.#.....#....####.#.#.#....#.#...#.##...#.#.#....#.#.#....##.#.####.##..#####.####.#.####..#...###.###..##...#..###.####...#..#.####.###.##..##....#.####....#.#..##.#..#.##.##..#......###.#...#..#.#.#.##.######.##.##..####.##..#.###.##.....##...#.....#..#....###..####.#.##..#.")

	img := _testExample.img

	img1 := trickyAlg.enhance(img)
	assert.True(t, img1.infinitePx, "image 1 is infinitely lit")

	img2 := trickyAlg.enhance(img1)
	assert.False(t, img2.infinitePx, "image 2 is not infinitely lit")
}

func TestGetMask(t *testing.T) {
	finiteBits := []byte(`#..#.
#....
##..#
..#..
..###
`)

	basicImage, err := newImage(finiteBits)
	require.NoError(t, err)

	fullyLitImg, err := newImage(finiteBits)
	require.NoError(t, err)
	fullyLitImg.infinitePx = true

	tt := []struct {
		name     string
		img      image
		row, col int
		want     string
	}{
		{
			name: "fully inside the source image",
			img:  basicImage,
			row:  1,
			col:  1,
			want: "#..\n#..\n##.\n",
		},
		{
			name: "on the edge of a basic image",
			img:  basicImage,
			want: "...\n.#.\n.#.\n",
		},
		{
			name: "fully inside an infinite image",
			img:  fullyLitImg,
			row:  1,
			col:  1,
			want: "#..\n#..\n##.\n",
		},
		{
			name: "on the edge of an infinite image",
			img:  fullyLitImg,
			row:  4,
			col:  4,
			want: "..#\n###\n###\n",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Log("applying mask to image at", tc.row, tc.col)
			t.Logf("\n%v", tc.img)

			got := getMask(tc.img, tc.row, tc.col)

			wantMask, err := newImage([]byte(tc.want))
			require.NoError(t, err)
			assert.Equal(t, wantMask, got)
		})
	}
}

func TestAlgIsLit(t *testing.T) {
	tt := []struct {
		name string
		key  int
		want bool
	}{
		{"first bit", 0, false},
		{"0th word, second-last bit, ", 62, true},
		{"last bit", 511, true},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			got := _testExample.alg.isLit(tc.key)
			assert.Equal(t, tc.want, got)
		})
	}
}
