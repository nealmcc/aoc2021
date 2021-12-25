package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

// algorithm defines the rules for enhancing an image
type algorithm string

// image is a square bitmap of pixels, surrounded by an infinite area of
// pixels which are either all lit, or all unlit.
type image struct {
	// size is the width and height of this image's finite section.
	size int
	// numLit is the count of lit pixels within this image's finite section.
	numLit int
	// pixels holds the status of pixels within this image's finite section.
	pixels [][]bool
	// infinitePx defines the status of pixels outside the finite section.
	infinitePx bool
}

// compile-time interface checks
var (
	_ fmt.Formatter = image{}
	_ fmt.Stringer  = image{}
)

// main solves both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	alg, img, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	img1 := part1(alg, img)
	if err != nil {
		log.Fatal(err)
	}

	img2 := part2(alg, img1)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()

	fmt.Println("part1:", img1.numLit)
	fmt.Println("part2:", img2.numLit)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read the given input, returning an enhancement algorithm and an image.
func read(r io.Reader) (algorithm, image, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return "", image{}, err
	}

	alg, err := buf.ReadString('\n')
	if err != nil {
		return "", image{}, err
	}

	// trim trailing newline from algorithm
	alg = alg[:len(alg)-1]

	// discard a newline from input
	_, err = buf.ReadByte()
	if err != nil {
		return "", image{}, err
	}

	img, err := newImage(buf.Bytes())
	return algorithm(alg), img, err
}

func part1(alg algorithm, src image) image {
	format := "after %d\n%" + strconv.Itoa(src.size+8) + "v\n"

	fmt.Printf(format, 0, src)

	img1 := alg.enhance(src)
	fmt.Printf(format, 1, img1)

	img2 := alg.enhance(img1)
	fmt.Printf(format, 2, img2)

	return img2
}

func part2(alg algorithm, img image) image {
	for i := 0; i < 48; i++ {
		img = alg.enhance(img)
	}
	return img
}

// enhance returns a new image based on the source image.
func (alg algorithm) enhance(src image) image {
	buf := make([][]bool, src.size+2)
	for i := 0; i < src.size+2; i++ {
		buf[i] = make([]bool, src.size+2)
	}

	dest := image{
		size:   src.size + 2,
		pixels: buf,
	}

	for row := 0; row < dest.size; row++ {
		for col := 0; col < dest.size; col++ {
			mask := getMask(src, row-1, col-1)
			dest.set(row, col, alg.isLit(indexFor(mask)))
		}
	}

	if src.infinitePx {
		dest.infinitePx = alg.isLit(511)
	} else {
		dest.infinitePx = alg.isLit(0)
	}

	return dest
}

// getMask copies a 3x3 section of the given image centered at the given point.
func getMask(img image, row, col int) image {
	buf := make([][]bool, 3)
	for i := 0; i < 3; i++ {
		buf[i] = make([]bool, 3)
	}

	out := image{
		size:   3,
		pixels: buf,
	}

	top, left := row-1, col-1
	for r := top; r < top+3; r++ {
		for c := left; c < left+3; c++ {
			out.set(r-top, c-left, img.isLit(r, c))
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
	return alg[key] == '#'
}

// newImage parses the given text and constructs a new image from it.
// The infinite pixels of the resulting image will be unlit.
func newImage(data []byte) (image, error) {
	img := image{
		pixels: make([][]bool, 0, 8),
	}
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
			return image{}, fmt.Errorf("invalid input at column %d", i)
		}
	}
	img.size = len(img.pixels)
	return img, nil
}

// isLit checks to see if the given pixel is lit.  When the position is outside
// the finite bounds in the matrix, then this image's 'infinite pixel status'
// is used.
func (img image) isLit(row, col int) bool {
	if img.inBounds(row, col) {
		return img.pixels[row][col]
	}

	return img.infinitePx
}

// inBounds checks to see if the given coordinate is within the bounds of
// the finite portion of the image.
func (img image) inBounds(row, col int) bool {
	return 0 <= row && row < img.size && 0 <= col && col < img.size
}

// set the pixel value of the given coordinate.
func (img *image) set(row, col int, newVal bool) error {
	if !img.inBounds(row, col) {
		return errors.New("cannot set an out of bounds pixel value")
	}

	prev := img.pixels[row][col]
	if prev == newVal {
		return nil
	}

	img.pixels[row][col] = newVal
	if newVal {
		img.numLit++
	} else {
		img.numLit--
	}
	return nil
}

// Format implements fmt.Formatter.
// The width is used to increase padding around the image if desired.
func (img image) Format(s fmt.State, verb rune) {
	width, ok := s.Width()
	if !ok || width < img.size {
		width = img.size
	}

	buf := make([]byte, width)
	pad := (width - img.size) / 2
	for row := -1 * pad; row < img.size+pad; row++ {
		for col := -1 * pad; col < img.size+pad; col++ {
			if img.isLit(row, col) {
				buf[col+pad] = '#'
			} else {
				buf[col+pad] = '.'
			}
		}
		s.Write(append(buf, '\n'))
	}
}

// String implements fmt.Stringer.
func (img image) String() string {
	return fmt.Sprintf("%v", img)
}
