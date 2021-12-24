package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
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

	p1 := part1(alg, img)
	if err != nil {
		log.Fatal(err)
	}

	p2 := part2(alg, img)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()

	fmt.Println("part1:", p1)
	fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read the given input, returning an enhancement algorithm and an image.
func read(r io.Reader) (algorithm, image, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return algorithm{}, image{}, err
	}

	alg := algorithm{}
	data, err := buf.ReadBytes('\n')
	if err != nil {
		return algorithm{}, image{}, err
	}
	if err := alg.UnmarshalText(data); err != nil {
		return algorithm{}, image{}, err
	}

	// discard a newline
	_, err = buf.ReadByte()
	if err != nil {
		return algorithm{}, image{}, err
	}

	var img image
	err = img.UnmarshalText(buf.Bytes())
	return alg, img, err
}

func part1(alg algorithm, src image) int {
	format := "after %d\n%" + strconv.Itoa(src.size+6) + "v\n"

	fmt.Printf(format, 0, src)

	img1 := alg.enhance(src)
	fmt.Printf(format, 1, img1)

	img2 := alg.enhance(img1)
	fmt.Printf(format, 2, img2)

	return img2.numLit
}

func part2(alg algorithm, img image) int {
	for i := 0; i < 50; i++ {
		img = alg.enhance(img)
	}
	return img.numLit
}
