package main

import (
	"testing"

	v "github.com/nealmcc/aoc2021/pkg/vector"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestEnhance(t *testing.T) {
	log := zaptest.NewLogger(t).Sugar()

	alg, img := _testExample.alg, _testExample.img

	img1 := alg.enhance(img, log)
	img2 := alg.enhance(img1, log)

	var want image
	err := want.UnmarshalText([]byte(_testExample.enhance2x))
	require.NoError(t, err)

	assert.Equal(t, want, img2)
}

func TestEnhance_Tricky(t *testing.T) {
	// trickyAlg is tricky because the 0th value is '#'.  This means that when
	// enhancing an infinite image, if the source image is *not* infinitely lit,
	// then the target image will be.  Also, because the 511th value is '.' the
	// converse is true. When enhancing an image that *is* infinitely lit, the
	// resulting image will not be.
	var trickyAlg algorithm
	trickyAlg.UnmarshalText([]byte("#..#.#######.##...##.##.#.#..#..#.....####...####.##.###...##.####......##.###.##...#..##..#.######...###..########.#.##.#.#..#..##.##..####.###.###..#...##.##.###.....###..###....#.####.#..##....#.##...##.#..#.....###.#..#.....##..##.#.#.....#....####.#.#.#....#.#...#.##...#.#.#....#.#.#....##.#.####.##..#####.####.#.####..#...###.###..##...#..###.####...#..#.####.###.##..##....#.####....#.#..##.#..#.##.##..#......###.#...#..#.#.#.##.######.##.##..####.##..#.###.##.....##...#.....#..#....###..####.#.##..#."))

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

	var basicImage image
	basicImage.UnmarshalText(finiteBits)

	var fullyLitImg image
	fullyLitImg.UnmarshalText(finiteBits)
	fullyLitImg.infinitePx = true

	tt := []struct {
		name  string
		img   image
		coord v.Coord
		want  string
	}{
		{
			name:  "fully inside the source image",
			img:   basicImage,
			coord: v.Coord{X: 1, Y: 1},
			want:  "#..\n#..\n##.\n",
		},
		{
			name:  "on the edge of a basic image",
			img:   basicImage,
			coord: v.Coord{},
			want:  "...\n.#.\n.#.\n",
		},
		{
			name:  "fully inside an infinite image",
			img:   fullyLitImg,
			coord: v.Coord{X: 1, Y: 1},
			want:  "#..\n#..\n##.\n",
		},
		{
			name:  "on the edge of an infinite image",
			img:   fullyLitImg,
			coord: v.Coord{X: 4, Y: 4},
			want:  "..#\n###\n###\n",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Log("applying mask to image at", tc.coord)
			t.Logf("\n%v", tc.img)

			log := zaptest.NewLogger(t).Sugar()

			got := getMask(tc.img, tc.coord, log)

			var wantMask image
			wantMask.UnmarshalText([]byte(tc.want))
			assert.Equal(t, "\n"+wantMask.String(), "\n"+got.String(), "masks should match")
		})
	}
}

func TestIsLit(t *testing.T) {
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
