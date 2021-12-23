package main

import (
	"bufio"
	"encoding"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/nealmcc/aoc2021/pkg/vector"
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

	p1, err := part1(*alg, *img)
	if err != nil {
		log.Fatal(err)
	}

	// p2, err := part2(alg, img)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	end := time.Now()

	fmt.Println("part1:", p1)
	// fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read the given input, returning an enhancement algorithm and an image.
func read(r io.Reader) (*algorithm, *image, error) {
	s := bufio.NewScanner(r)

	s.Scan()
	if err := s.Err(); err != nil {
		return nil, nil, err
	}

	alg := algorithm{}
	if err := alg.UnmarshalBinary(s.Bytes()); err != nil {
		return nil, nil, err
	}

	// discard empty line
	s.Scan()

	img := image{
		pixels: make(map[vector.Coord]px, 10000),
	}
	var row int
	for s.Scan() {
		data := s.Bytes()
		for col, b := range data {
			if b == '#' {
				img.pixels[vector.Coord{X: col, Y: row}] = px(true)
			}
		}
		row++
	}
	img.size = row

	if err := s.Err(); err != nil {
		return nil, nil, err
	}

	return &alg, &img, nil
}

func part1(alg algorithm, img image) (int, error) {
	return 0, errors.New("not implemented")
}

func part2(alg algorithm, img image) (int, error) {
	return 0, errors.New("not implemented")
}

// algorithm defines the rules for enhancing an image
type algorithm [8]uint64

var _ encoding.BinaryUnmarshaler = (*algorithm)(nil)

func (alg *algorithm) UnmarshalBinary(data []byte) error {
	if len(data) != 512 {
		return errors.New("algorithms are always exactly 512 bytes long")
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

// image is a square bitmap of pixels
type image struct {
	size   int
	pixels map[vector.Coord]px
}

var _ fmt.Formatter = image{}

func (img image) Format(s fmt.State, verb rune) {
	buf := make([][]byte, img.size)
	for i := 0; i < img.size; i++ {
		buf[i] = make([]byte, img.size)
		for j := 0; j < img.size; j++ {
			if img.pixels[vector.Coord{X: i, Y: j}] {
				buf[i][j] = '#'
			} else {
				buf[i][j] = '.'
			}
		}
	}

	for _, row := range buf {
		s.Write(append(row, '\n'))
	}
}

// px is a monochrome px - true means lit, false means unlit
type px bool
