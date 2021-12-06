package main

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"log"
	"os"
	"strings"

	v "github.com/nealmcc/aoc2021/pkg/vector"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	segments, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	d1 := render(segments, false)
	fmt.Printf("part1: %d\n", count(d1))

	png1, err := os.Create("./part1.png")
	if err != nil {
		log.Fatal(err)
	}
	defer png1.Close()
	writePNG(d1, png1)

	d2 := render(segments, true)
	fmt.Printf("part2: %d\n", count(d2))

	png2, err := os.Create("./part2.png")
	if err != nil {
		log.Fatal(err)
	}
	defer png2.Close()
	writePNG(d2, png2)
}

// segment is a line segment from a to b.
type segment struct {
	a, b v.Coord
}

// bitmap stores how many fissures are at each coordinate
type bitmap struct {
	points map[v.Coord]int
	width  int
	height int
}

// read a list of line segments from the given input
func read(r io.Reader) ([]segment, error) {
	s := bufio.NewScanner(r)

	segments := make([]segment, 0, 16)
	for s.Scan() {
		coords, err := v.ParseCoords(strings.Split(s.Text(), " -> ")...)
		if err != nil {
			return nil, err
		}
		segments = append(segments, segment{a: coords[0], b: coords[1]})
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return segments, nil
}

// render plots the list of segments, creating a bitmap showing how many
// fissures are present at each position.  The includeDiag parameteter
// determines whether to include diagonal fissures (part 2) or not (part 1)
func render(segments []segment, includeDiag bool) bitmap {
	b := bitmap{
		points: make(map[v.Coord]int),
	}

	for _, seg := range segments {
		delta := v.Sub(seg.b, seg.a)
		if !includeDiag && delta.X != 0 && delta.Y != 0 {
			continue
		}

		if seg.a.X > b.width {
			b.width = seg.a.X
		}
		if seg.b.X > b.width {
			b.width = seg.b.X
		}
		if seg.a.Y > b.height {
			b.height = seg.a.Y
		}
		if seg.b.Y > b.height {
			b.height = seg.b.Y
		}

		curr := seg.a
		b.points[curr]++

		unit, _ := v.Reduce(delta)
		for curr != seg.b {
			curr = v.Add(curr, unit)
			b.points[curr]++
		}
	}

	if len(b.points) > 0 {
		b.width++
		b.height++
	}
	return b
}

// count the number of points in the given bitmap that have 2 or more fissures
func count(b bitmap) int {
	sum := 0
	for _, n := range b.points {
		if n >= 2 {
			sum++
		}
	}
	return sum
}

// compile-time interface check
var _ fmt.Formatter = bitmap{}

// Format implements the fmt.Formatter interface, so we can easily view
// the bitmap, just by using fmt.Println() or fmt.Printf() etc.
// the formatting options and verb are ignored.
func (b bitmap) Format(state fmt.State, verb rune) {
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			count := b.points[v.Coord{X: x, Y: y}]
			if count > 0 {
				state.Write([]byte{'0' + byte(count)})
			} else {
				state.Write([]byte{'.'})
			}
		}
		if b.height-y > 1 {
			state.Write([]byte{'\n'})
		}
	}
}

func writePNG(b bitmap, w io.Writer) error {
	img := image.NewNRGBA(image.Rect(0, 0, b.width, b.height))
	bg := image.Uniform{color.White}
	draw.Draw(img, img.Bounds(), &bg, image.Point{}, draw.Src)
	for y := 0; y < b.height; y++ {
		for x := 0; x < b.width; x++ {
			count := b.points[v.Coord{X: x, Y: y}]
			if count == 0 {
				continue
			}
			img.Set(x, y, color.NRGBA{
				B: uint8(255),
				A: uint8(255) / uint8(count),
			})
		}
	}
	return png.Encode(w, img)
}
