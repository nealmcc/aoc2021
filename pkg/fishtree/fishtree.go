// Package fishtree describes an abstract syntax tree for really wonky
// arithmetic, as performed by the snailfish on Day18.
package fishtree

import (
	"fmt"
	"unicode"

	"github.com/nealmcc/aoc2021/pkg/stack"
)

// Node is the interface for all nodes in the tree.
type Node interface {
	// Value calculates the arithmetic result for the tree.
	Value() int
	// Magnitude calculates the derived (asymmetric) value for the tree.
	Magnitude() int
}

// Num is an atomic node with no children.
type Num int

var _ Node = Num(0)

// Value implements Node.Value.
func (n Num) Value() int { return int(n) }

// Magnitude implements Node.Magnitude.
func (n Num) Magnitude() int { return int(n) }

// Add is a Node with two children that evaluates by returning their sum,
// and has a magnitude of 3x its left node and 2x its right.
type Add struct {
	L, R Node
}

var _ Node = Add{}

// Value implements Node.Value.
func (a Add) Value() int {
	return a.L.Value() + a.R.Value()
}

// Magnitude implements Node.Magnitude.
func (a Add) Magnitude() int {
	return 3*a.L.Magnitude() + 2*a.R.Magnitude()
}

// New parses one fish number and returns it as the root of a fishtree.
//
// Example:
//     tree := New("[[1,2],[[3,4],5]]")
//     tree.Magnitude() // output: 143
func New(s string) (Node, error) {
	rpn, err := shuntingYard(s)
	if err != nil {
		return nil, err
	}
	return parsePostfix(rpn)
}

// shuntingYard reads tokens from the infix notation and converts
// them to postfix notation so that they can be evaluated as a tree.
// Converts commas to '+' and discards whitespace, but otherwise requires
// all input to be valid.  Expects all numeric inputs to be a single digit.
func shuntingYard(infix string) ([]rune, error) {
	// reverse polish notation
	rpn := make([]rune, 0, 64)
	s := stack.ByteStack{}

	for i, token := range infix {
		switch token {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			rpn = append(rpn, token)

		case ',', '+':
			s.Push('+')

		case '[':
			s.Push('[')

		case ']':
			top, ok := s.Peek()
			for ok && top != '[' {
				rpn = append(rpn, rune(s.Pop()))
				top, ok = s.Peek()
			}
			if !ok || top != '[' {
				return nil, fmt.Errorf("index %d: mismatched closing bracket", i)
			}
			// discard the opening bracket from the stack
			s.Pop()

		default:
			if unicode.IsSpace(token) {
				continue
			}
			return nil, fmt.Errorf("index %d: unexpected input %q", i, token)
		}
	}

	for s.Length() > 0 {
		operator := s.Pop()
		rpn = append(rpn, rune(operator))
	}

	return rpn, nil
}

// parsePostfix creates a new fish tree based on the given postfix notation.
// Expects all numeric inputs to be a single digit.
func parsePostfix(postfix []rune) (Node, error) {
	s := NodeStack{}

	for i, token := range postfix {
		switch token {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			s.Push(Num(token - '0'))

		case '+':
			if s.Length() < 2 {
				return nil, fmt.Errorf("index %d: want operand ; got %q", i, token)
			}
			n := Add{}
			n.R = s.Pop()
			n.L = s.Pop()
			s.Push(n)

		default:
			return nil, fmt.Errorf("index %d: unexpected input %q", i, token)
		}
	}

	if length := s.Length(); length != 1 {
		return nil, fmt.Errorf("expected 1 root ; got %d", length)
	}

	return s.Pop(), nil
}
