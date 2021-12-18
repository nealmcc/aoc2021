package main

import "github.com/pkg/errors"

// packetType defines the type of a packet
type packetType byte

const (
	_sum packetType = iota
	_product
	_minimum
	_maximum
	_literal
	_greaterThan
	_lessThan
	_equal
)

const (
	intSize = 32 << (^uint(0) >> 63) // 32 or 64

	maxInt = 1<<(intSize-1) - 1
	minInt = -1 << (intSize - 1)
)

// sum returns the sum of all this packet's children.
func (p Packet) sum() (int, error) {
	children, err := p.Children()
	if err != nil {
		return 0, err
	}

	sum := 0
	for _, c := range children {
		val, err := c.Value()
		if err != nil {
			return 0, err
		}
		sum += val
	}
	return sum, nil
}

// product returns the product of all this packet's children.
func (p Packet) product() (int, error) {
	children, err := p.Children()
	if err != nil {
		return 0, err
	}

	prod := 1
	for _, c := range children {
		val, err := c.Value()
		if err != nil {
			return 0, err
		}
		prod *= val
	}

	return prod, nil
}

// min returns the minimum value from this packet's children.
func (p Packet) min() (int, error) {
	children, err := p.Children()
	if err != nil {
		return 0, err
	}

	best := maxInt

	for _, c := range children {
		val, err := c.Value()
		if err != nil {
			return 0, err
		}
		if val < best {
			best = val
		}
	}
	return best, nil
}

// max returns the maximum value from this packet's children.
func (p Packet) max() (int, error) {
	children, err := p.Children()
	if err != nil {
		return 0, err
	}

	best := minInt

	for _, c := range children {
		val, err := c.Value()
		if err != nil {
			return 0, err
		}
		if val > best {
			best = val
		}
	}
	return best, nil
}

// greater returns 1 if this packet's first child is greater than its second.
// Otherwise, greater returns 0.
func (p Packet) greater() (int, error) {
	left, right, err := get2Children(p)
	if err != nil {
		return 0, err
	}

	if left > right {
		return 1, nil
	}
	return 0, nil
}

// less returns 1 if this packet's first child is less than its second.
// Otherwise less returns 0.
func (p Packet) less() (int, error) {
	left, right, err := get2Children(p)
	if err != nil {
		return 0, err
	}

	if left < right {
		return 1, nil
	}
	return 0, nil
}

// equal returns 1 if this packet's first child has the same value as its second.
// Otherwise, equal returns 0.
func (p Packet) equal() (int, error) {
	left, right, err := get2Children(p)
	if err != nil {
		return 0, err
	}

	if left == right {
		return 1, nil
	}
	return 0, nil
}

// get2Children returns the first two children of p.
func get2Children(p Packet) (left, right int, err error) {
	children, err := p.Children()
	if err != nil {
		return 0, 0, err
	}

	if len(children) < 2 {
		return 0, 0, errors.New("invalid operation")
	}

	leftPkt, rightPkt := children[0], children[1]
	left, err = leftPkt.Value()
	if err != nil {
		return 0, 0, err
	}

	right, err = rightPkt.Value()
	if err != nil {
		return 0, 0, err
	}
	return left, right, err
}
