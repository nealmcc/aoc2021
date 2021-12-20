// Package fish describes an abstract syntax tree for really wonky
// arithmetic, as performed by the snailfish on Day18:
// https://adventofcode.com/2021/day/18
package fish

import (
	"errors"
	"fmt"
	"strconv"
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

// Pair is the interface for operators with two children.
type Pair interface {
	Node
	Left() Node
	Right() Node
}

// Leaf is the interface for value nodes with no children.
type Leaf interface {
	Node
	Set(x int)
}

// New parses one fish number and returns its top node.
//
// Example:
//     tree := New("[[1,2],[[3,4],5]]")
//     Magnitude(tree) // output: 143
func New(s string) (Node, error) {
	rpn, err := shuntingYard(s)
	if err != nil {
		return nil, err
	}

	return parsePostfix(rpn)
}

func Reduce(p Pair) {
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
func parsePostfix(postfix []rune) (Pair, error) {
	if len(postfix) < 2 {
		return nil, errors.New("input too short")
	}

	s := Stack{}
	for i, token := range postfix {
		switch token {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			n := Num(token - '0')
			s.Push(&n)

		case '+':
			if s.Length() < 2 {
				return nil, fmt.Errorf("index %d: want operand ; got %q", i, token)
			}
			right := s.Pop()
			left := s.Pop()
			s.Push(&Add{L: left, R: right})

		default:
			return nil, fmt.Errorf("index %d: unexpected input %q", i, token)
		}
	}

	if length := s.Length(); length != 1 {
		return nil, fmt.Errorf("expected 1 root ; got %d", length)
	}

	n := s.Pop().(*Add)
	return n, nil
}

// Add is a Pair with a value that is the sum of its two children.
// An Add node has a magnitude of 3x its left child plus 2x its right.
type Add struct {
	L, R Node
}

// compile-time interface checks.
var (
	_ Pair          = new(Add)
	_ fmt.Formatter = new(Add)
)

// Left implements Pair (along with Right)
func (a *Add) Left() Node { return a.L }

// Right implements Pair (along with Left)
func (a *Add) Right() Node { return a.R }

// Value implements Node. The value of Add is the sum of its children.
func (a *Add) Value() int {
	return a.L.Value() + a.R.Value()
}

// Magnitude implements Node. The magnitude of Add is 3x the magnitude
// of its left child plus 2x the magnitude of its right child.
func (a *Add) Magnitude() int {
	return 3*(a.L.Magnitude()) + 2*(a.R.Magnitude())
}

// Format implements fmt.Formatter
// If the + flag has been set, then Format prints using postfix notation.
// Otherwise, Format prints using infix notation.
func (a *Add) Format(s fmt.State, verb rune) {
	if s.Flag('+') {
		s.Write([]byte(fmt.Sprintf("%+v%+v+", a.L, a.R)))
		return
	}
	s.Write([]byte(fmt.Sprintf("[%v+%v]", a.L, a.R)))
}

// Num is a leaf node.
type Num int

// N is a convenience function that creates a pointer to a new Num
func N(n int) *Num {
	leaf := Num(n)
	return &leaf
}

// compile-time interface checks.
var (
	_ Leaf         = new(Num)
	_ fmt.Stringer = new(Num)
)

// Set implements Leaf
func (n *Num) Set(x int) { *n = Num(x) }

// Value implements Node.Value().
func (n *Num) Value() int { return int(*n) }

// Magnitude implements Node.Magnitude. The magnitude of a Num is its value.
func (n *Num) Magnitude() int { return int(*n) }

// String implements fmt.Stringer so that it's easier to read debug output.
func (n *Num) String() string {
	return strconv.Itoa(int(*n))
}

// // tree holds the full tree as a slice of nodes. The root is at index 0,
// // and the left and right children of an element at index i are at indices
// // 2i+1 and 2i+2, respectively.
// type tree struct {
// 	el []element
// }

// // element is a wrapper for any node in the tree.
// type element struct {
// 	n  Node
// 	id int
// }

// // Init returns a tree
// func Init(root Pair) *tree {
// 	s := Stack{}
// 	s.Push(root.(Node))
// 	elements := make([]element, 0, 16)

// 	for s.Length() > 0 {
// 		n := s.Pop()
// 		if _, isLeaf := n.(Leaf); isLeaf {
// 			elements = append(elements, &n)
// 		}

// 		if p, isPair := n.(Pair); isPair {
// 			s.Push(p.Left())
// 			s.Push(p.Right())
// 		}
// 	}

// 	return ix
// }
