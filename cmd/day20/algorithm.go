package main

import (
	"encoding"
	"errors"

	v "github.com/nealmcc/aoc2021/pkg/vector"
	"go.uber.org/zap"
)

// algorithm defines the rules for enhancing an image
type algorithm [8]uint64

// compile-time interface check.
var _ encoding.TextUnmarshaler = (*algorithm)(nil)

// UnmarshalText implements encoding.TextUnmarshaler
func (alg *algorithm) UnmarshalText(data []byte) error {
	if len(data) < 512 {
		return errors.New("algorithms require 512 bytes of data")
	}

	for i, j := 0, 64; j <= 512; i, j = i+64, j+64 {
		var n uint64
		for b := i; b < j; b++ {
			n <<= 1
			if data[b] == '#' {
				n += 1
			}
		}
		alg[i/64] = n
	}
	return nil
}

func (alg algorithm) enhance(src image, logs ...*zap.SugaredLogger) image {
	buf := make([][]bool, src.size+2)
	for i := 0; i < src.size+2; i++ {
		buf[i] = make([]bool, src.size+2)
	}

	dest := image{
		size:   src.size + 2,
		pixels: buf,
	}

	// translate a point from the coordinate system of the destination
	// to the coordinate system of the source image.
	translate := func(destCoord v.Coord) (srcCoord v.Coord) {
		return v.Coord{
			X: destCoord.X - 1,
			Y: destCoord.Y - 1,
		}
	}

	// setPixel assigns the correct value on the destination image, based on
	// applying the algorithm to the source image.
	setPixel := func(target v.Coord) {
		mask := getMask(src, translate(target))
		isLit := alg.isLit(indexFor(mask))
		dest.pixels[target.Y][target.X] = isLit
		if isLit {
			dest.numLit++
		}
	}

	for x := 0; x < dest.size; x++ {
		for y := 0; y < dest.size; y++ {
			setPixel(v.Coord{X: x, Y: y})
		}
	}

	if src.infinitePx {
		dest.infinitePx = alg.isLit(511)
	} else {
		dest.infinitePx = alg.isLit(0)
	}

	return dest
}

// genMask creates a 3x3 section of the given image centered at the given point.
func getMask(img image, center v.Coord, logs ...*zap.SugaredLogger) image {
	buf := make([][]bool, 3)
	for i := 0; i < 3; i++ {
		buf[i] = make([]bool, 3)
	}

	out := image{
		size:   3,
		pixels: buf,
	}

	topLeft := v.Add(center, v.Coord{X: -1, Y: -1})

	points := append(v.Neighbours8(center), center)

	for _, p := range points {
		dest := v.Sub(p, topLeft)
		if img.isLit(p, logs...) {
			out.pixels[dest.Y][dest.X] = true
		}
	}

	return out
}

// indexFor calculates an algorithm index [0..511] based on the given 3x3 image.
func indexFor(mask image) int {
	var key int

	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			key <<= 1
			if mask.pixels[row][col] {
				key += 1
			}
		}
	}

	return key
}

// isLit looks up the pixel value for the given key.
func (alg algorithm) isLit(key int) bool {
	quo, rem := key/64, key%64
	word := alg[quo]
	var mask uint64 = 1 << (63 - rem)
	isLit := word & mask
	return isLit > 0
}
