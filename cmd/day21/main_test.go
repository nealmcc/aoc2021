package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestD100_roll(t *testing.T) {
	r := require.New(t)
	d := new(d100)

	for want := 1; want <= 100; want++ {
		got := d.roll()
		r.Equal(want, got)
	}

	for want := 1; want <= 100; want++ {
		got := d.roll()
		r.Equal(want, got)
	}
}

func TestPart1(t *testing.T) {
	r := require.New(t)

	g := game{
		P1:    player{Name: "p1", Pos: 4},
		P2:    player{Name: "p2", Pos: 8},
		ToWin: 1000,
	}

	got := part1(g)
	r.Equal(739785, got)
}

func TestPart2(t *testing.T) {
	r := require.New(t)

	g := game{
		P1:    player{Name: "p1", Pos: 4},
		P2:    player{Name: "p2", Pos: 8},
		ToWin: 21,
	}

	got := part2(g)
	r.Equal(444356092776315, got)
}
