package main

import (
	"fmt"
	"time"
)

// main solves both part 1 and part 2.
func main() {
	start := time.Now()
	p1 := solve(true /* find largest */)
	p2 := solve(false /* find smallest */)
	end := time.Now()

	fmt.Println("part1:", p1)
	fmt.Println("part2:", p2)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// solve finds the largest (or smallest) 14-digit number that contains no 0's
// and is accepted by an ALU running the monad.
func solve(findLargest bool) int {
	// zs is the set of z values that we want to arrive at
	zs := make(intSet, 1)
	zs.add(0)

	// zDigits stores the 'best' digits required to reach a given value of z
	zDigits := make(map[int][]int, 9*14*26)

	// the order we try the digits determines which result we keep for z=0
	digits := []int{9, 8, 7, 6, 5, 4, 3, 2, 1}
	if findLargest {
		digits = []int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	}

	for i := 13; i >= 0; i-- {
		targetZs := make(intSet, 26)
		for _, digit := range digits {
			for z1 := range zs {
				prevZs := backward(z1, digit, i)
				for _, z0 := range prevZs {
					// Store the sequence of digits required to arrive at this
					// value of z.  Results from 'better' digits will
					// overwrite results from previous digits, based on the
					// order we try them.
					zDigits[z0] = append([]int{digit}, zDigits[z1]...)
					targetZs.add(z0)
				}
			}
		}
		zs = targetZs
	}

	n := 0
	for _, digit := range zDigits[0] {
		n = 10*n + digit
	}

	return n
}

// intSet is used as a set of integers.
type intSet map[int]struct{}

func (s *intSet) add(n int) {
	(*s)[n] = struct{}{}
}

// These constants differentiate each 'loop' in the given assembly.
// Note that half the time, m1 is zero or negative, and those are the same
// times when div is 26.  Therefore, on half of the loops, we will be multiplying
// z by 26, (plus w plus m2). On the other half of the loops, we will be
// dividing z by 26.
var (
	div    = []int{1, 1, 1, 26, 26, 1, 26, 26, 1, 26, 1, 1, 26, 26}
	magic1 = []int{12, 14, 11, -9, -7, 11, -1, -16, 11, -15, 10, 12, -4, 0}
	magic2 = []int{15, 12, 15, 12, 15, 2, 11, 15, 10, 2, 0, 0, 15, 15}
)

// monad is the equivalent of the input program, but written more concisely.
// It's not actually part of the solution for day24, but is useful to ensure
// the 'backward' function is correct.
func monad(input []byte) int {
	var z int
	for i := 0; i < 14; i++ {
		w := int(input[i])
		z = forward(z, w, i)
	}
	return z
}

// forward performs the calculation for a single digit.
// Only the value of z is used in subsequent steps.
func forward(z0, w, i int) (z1 int) {
	var (
		d  = div[i]
		m1 = magic1[i]
		m2 = magic2[i]
	)

	// This can only true whenever d is 26, because whenever d is 1, w-m1 is <0.
	// Note that when d is 26, w-m1 is always the range [0,26)
	// We will use this fact in the backward() function.
	if z0%26 == w-m1 {
		return z0 / d
	}

	// Since w ∈ [1,9] and m2 ∈ [0,15] then w + m2 ∈ [1,24]
	// therefore, z1 ≡ w+m2 (mod 26)
	// We will use this fact in the backward() function.
	z1 = (z0/d)*26 + w + m2
	return z1
}

// backward is the inverse of forward.  This is the crux of the solution.
// It finds up to two possible previous values for zPrev which will produce the
// given value of zNext when passed in to forward().  The two previous values
// come from the two branches of execution in forward().
// It is possible (common) that there will be no possible z0 which can produce
// the given z1, in which case, the return list will be empty.  This speeds up
// our search a lot.
func backward(z1, w, i int) (zPrev []int) {
	var (
		d  = div[i]
		m1 = magic1[i]
		m2 = magic2[i]
	)

	// this is the inverse of the first return statement from forward()
	// This can only be a possible answer if w-m1 ∈ [0,26)
	if 0 <= w-m1 && w-m1 < 26 {
		z0 := z1*d + w - m1
		zPrev = append(zPrev, z0)
	}

	// this is the inverse of the second return statement from forward()
	// This can only be a possible answer if (z1 - w - m2) is some multiple of 26.
	x := z1 - w - m2
	if x%26 == 0 {
		z0 := d * x / 26
		zPrev = append(zPrev, z0)
	}
	return zPrev
}
