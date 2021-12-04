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
	// samples is a radix tree containing all the data that this meter collected
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

// oxygen returns this meter's oxygen reading.
func (m *meter) oxygen() (int64, error) {
	chooseOnes := m.ones[0] >= m.count/2
	var prefix []byte
	if chooseOnes {
		prefix = []byte{'1'}
	} else {
		prefix = []byte{'0'}
	}

	tree, _ := m.samples.WithPrefix(string(prefix))
	for tree.Size() > 1 {
		oneTree, _ := tree.WithPrefix("1")
		zeroTree, _ := tree.WithPrefix("0")

		if zeroTree.Size() > oneTree.Size() {
			tree = zeroTree
			prefix = append(prefix, '0')
		} else {
			tree = oneTree
			prefix = append(prefix, '1')
		}
	}

	o2 := tree.ToSlice(prefix...)
	return strconv.ParseInt(o2[0], 2, 16)
}

// carbonDioxide returns this meter's carbon dioxide reading.
func (m *meter) carbonDioxide() (int64, error) {
	// todo: remove duplication between co2() and o2()
	chooseOnes := m.ones[0] < m.count/2
	var prefix []byte
	if chooseOnes {
		prefix = []byte{'1'}
	} else {
		prefix = []byte{'0'}
	}

	tree, _ := m.samples.WithPrefix(string(prefix))
	for tree.Size() > 1 {
		oneTree, _ := tree.WithPrefix("1")
		zeroTree, _ := tree.WithPrefix("0")

		if zeroTree.Size() <= oneTree.Size() {
			tree = zeroTree
			prefix = append(prefix, '0')
		} else {
			tree = oneTree
			prefix = append(prefix, '1')
		}
	}

	co2 := tree.ToSlice(prefix...)
	return strconv.ParseInt(co2[0], 2, 16)
}
