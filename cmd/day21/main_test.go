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
	p1 := player{Name: "p1", Pos: 4}
	p2 := player{Name: "p2", Pos: 8}
	g := game{
		players: []*player{&p1, &p2},
		d:       new(d100),
	}

	got := part1(g)
	r.Equal(739785, got)
}
