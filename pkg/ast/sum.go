package ast

import "errors"

// Sum returns the sum of all the given fish numbers, correctly reducing
// the results of each addition.
func Sum(numbers ...Number) (Number, error) {
	return nil, errors.New("not implemented")
	// if len(numbers) == 0 {
	// 	return N(0), nil
	// }

	// if len(numbers) == 1 {
	// 	return numbers[0], nil
	// }

	// s := Stack{}

	// // infix sequence
	// infix := -1
	// mn := func(n Number) node {
	// 	infix++
	// 	return node{n: n, infixID: infix}
	// }

	// s.Push(mn(numbers[0]))

	// for i := 1; i < len(numbers); i++ {
	// 	top := s.Pop()
	// 	pair := &add{
	// 		l: top.n,
	// 		r: numbers[i],
	// 	}
	// 	s.Push(pair.reduce())
	// }
	// return s.Pop(), nil
}
