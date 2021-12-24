package main

import (
	"encoding"
	"fmt"

	v "github.com/nealmcc/aoc2021/pkg/vector"
	"go.uber.org/zap"
)

// image is a square bitmap of pixels, surrounded by an infinite area of
// pixels which are either all lit, or all unlit.
type image struct {
	size   int
	numLit int
	pixels [][]bool
	// infinitePx determines if the pixels outside of the map are lit.
	infinitePx bool
}

// compile-time interface checks
var (
	_ encoding.TextUnmarshaler = (*image)(nil)
	_ fmt.Formatter            = image{}
)

func (img image) isLit(p v.Coord, logs ...*zap.SugaredLogger) bool {
	inBounds := 0 <= p.X && p.X < img.size &&
		0 <= p.Y && p.Y < img.size

	if inBounds {
		return img.pixels[p.Y][p.X]
	}
	return img.infinitePx
}

func (img *image) UnmarshalText(data []byte) error {
	img.pixels = make([][]bool, 0, 8)
	img.numLit = 0

	line := make([]bool, 0, 8)
	for i, ch := range data {
		switch ch {
		case '\n':
			img.pixels = append(img.pixels, line)
			next := make([]bool, 0, len(line))
			line = next

		case '#':
			line = append(line, true)
			img.numLit++

		case '.':
			line = append(line, false)

		default:
			return fmt.Errorf("invalid input at column %d", i)
		}
	}
	img.size = len(img.pixels)
	return nil
}

// Format implements fmt.Formatter
func (img image) Format(s fmt.State, verb rune) {
	width, ok := s.Width()
	if !ok || width < img.size {
		width = img.size
	}

	buf := make([]byte, width)
	pad := (width - img.size) / 2
	for row := -1 * pad; row < img.size+pad; row++ {
		for col := -1 * pad; col < img.size+pad; col++ {
			if img.isLit(v.Coord{X: col, Y: row}) {
				buf[col+pad] = '#'
			} else {
				buf[col+pad] = '.'
			}
		}
		s.Write(append(buf, '\n'))
	}
}

func (img image) String() string {
	return fmt.Sprintf("%v", img)
}
