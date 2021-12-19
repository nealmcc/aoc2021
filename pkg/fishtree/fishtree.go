// Package fishtree describes an abstract syntax tree for really wonky
// arithmetic, as performed by the snailfish on Day18.
//
// The easiest way to create a tree is to write it as a string using infix
// notation (as in the Day18 examples) and then pass that input to
// ShuntingYard(), followed by New().
package fishtree

import (
	"bytes"
	"errors"
	"fmt"
	"unicode"

	"github.com/nealmcc/aoc2021/pkg/stack"
)

// Node is the interface for all nodes in the tree.
type Node interface {
	// Value calculates the arithmetic result for the tree.
	Value() int
	// Magnitude calculates the derived (assymmetric) value for the tree.
	Magnitude() int
}

// Num is an atomic node with no children
type Num int

var _ Node = Num(0)

// Value implements Valuer
func (n Num) Value() int { return int(n) }

// Magnitude implements Valuer.
func (n Num) Magnitude() int { return int(n) }

// Add is a binary Node that evaluates by returning the sum of its operands,
// and has a magnitude of 3x its left node and 2x its right.
type Add struct {
	L, R Node
}

var _ Node = Add{}

// Value implements Valuer
func (a Add) Value() int {
	return a.L.Value() + a.R.Value()
}

// Magnitude implements Valuer
func (a Add) Magnitude() int {
	return 3*a.L.Magnitude() + 2*a.R.Magnitude()
}

// ShuntingYard reads tokens from the reader in infix notation, and converts
// them to postfix notation, so that they can be evaluated as a tree.
// Converts commas to '+' and discards whitespace, but otherwise requires
// all input to be valid.
func ShuntingYard(r *bytes.Reader) ([]rune, error) {
	// reverse polish notation
	rpn := make([]rune, 0, 64)
	s := stack.ByteStack{}

	for r.Len() > 0 {
		token, _, err := r.ReadRune()
		if err != nil {
			return nil, err
		}

		switch token {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			rpn = append(rpn, token)

		case ',', '+':
			s.Push('+')

		case '[':
			s.Push('[')

		case ']':
			operator := s.Pop()
			for operator != '[' {
				rpn = append(rpn, rune(operator))
				operator = s.Pop()
			}

		default:
			if unicode.IsSpace(token) {
				continue
			}
			return nil, fmt.Errorf("unexpected input: %c", token)
		}
	}

	for s.Length() > 0 {
		operator := s.Pop()
		rpn = append(rpn, rune(operator))
	}

	return rpn, nil
}

// New creates a new fish tree based on the given postfix notation.
func New(postfix []rune) (Node, error) {
	s := NodeStack{}

	for _, r := range postfix {
		switch r {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s.Push(Num(r - '0'))

		case '+':
			if s.Length() < 2 {
				return nil, errors.New("received operator without two operands")
			}
			n := Add{}
			n.R = s.Pop()
			n.L = s.Pop()
			s.Push(n)

		default:
			return nil, fmt.Errorf("unexpected input: %c", r)
		}
	}

	if length := s.Length(); length != 1 {
		return nil, fmt.Errorf("expected 1 root ; got %d", length)
	}

	return s.Pop(), nil
}
