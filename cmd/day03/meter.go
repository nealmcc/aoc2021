package main

import (
	"errors"
	"strconv"

	"github.com/nealmcc/aoc2021/pkg/radixtree"
)

// meter represents a diagnostic meter which samples data.
type meter struct {
	// ones tracks the count of ones that have been set for each position
	ones []int
	// count is the total number of samples this meter has taken
	count int
	// samples is a dictionary of all the data that this meter collected
	samples *radixtree.Node
}

// newMeter creates a new diagnostic meter, rated for
// reading diagnostic samples with n bits of data.
func newMeter(n int) *meter {
	return &meter{
		ones:    make([]int, n),
		samples: &radixtree.Node{},
	}
}

// sample a single value, storing the results in the meter.
func (m *meter) sample(s string) error {
	n, err := strconv.ParseInt(s, 2, 16)
	if err != nil {
		return err
	}
	if n > m.maxSample() {
		return errors.New("input too large")
	}

	m.samples.Insert(s)

	var rem int64
	for i := len(m.ones) - 1; i >= 0; i-- {
		n, rem = n/2, n%2
		m.ones[i] += int(rem)
	}
	m.count++

	return nil
}

// maxSample is the largest possible value this meter is capable of reading.
func (m *meter) maxSample() int64 {
	return 1<<len(m.ones) - 1
}

// gamma returns this meter's current value of gamma.
func (m *meter) gamma() int64 {
	var g int64

	mid := m.count / 2
	for _, n := range m.ones {
		g = g << 1
		if n > mid {
			g += 1
		}
	}

	return g
}

// epsilon returns this meter's current value of epsilon.
func (m *meter) epsilon() int64 {
	return m.maxSample() ^ m.gamma()
}
