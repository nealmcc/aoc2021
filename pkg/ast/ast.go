// Package ast describes an abstract syntax tree for really wonky
// arithmetic, as performed by the snailfish on Day18:
// https://adventofcode.com/2021/day/18
package ast

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

// Number is a snailfish number.
type Number interface {
	// Magnitude calculates a derived value for this number, which varies
	// according to its underlying tree structure.
	Magnitude() int
}

// Tree holds the full Tree for a Number.
type Tree struct {
	// root is a reference to the top node in the tree.
	root *node

	// infix holds all the nodes in their infix order.
	// A node's ID is its index in this slice.
	infix []*node
}

var _ Number = Tree{}

// Magnitude implements Number.
func (t Tree) Magnitude() int {
	if t.root == nil {
		return 0
	}
	return t.root.Magnitude()
}

// node is a single element in the tree. Each node is also a Number.
// a node may be a pair (in which case it will have an operator and two children)
// or else it will be a leaf node, with just a value.
// The node's ID is the sequence order this node is encountered during an
// in-order traversal of the tree.
type node struct {
	id    int    // id is this node's index during an in-order traversal of the tree
	op    opcode // op specifies what type of node this is.
	value int    // if op == opValue, then this holds the node's value
	left  *node  // if op != opValue, this points to the node's left child.
	right *node  // if op != opVAlue, this points the node's right child.
}

var _ Number = node{}

// Magnitude implements Number.
func (n node) Magnitude() int {
	if n.op == opValue {
		return n.value
	}
	return 3*n.left.Magnitude() + 2*n.right.Magnitude()
}

// format implements fmt.Formatter to simplify debugging.
// Using the '+' flag will result in postfix notation, while
// using no flag will return infix notation, identical to the original
// input string.
func (n node) Format(s fmt.State, verb rune) {
	switch n.op {
	case opValue:
		s.Write([]byte(strconv.Itoa(n.value)))

	default:
		if s.Flag('+') {
			// output using postfix notation:
			s.Write([]byte(fmt.Sprintf("%+v%+v+", n.left, n.right)))
		} else {
			// by default, use infix notation:
			s.Write([]byte(fmt.Sprintf("[%v,%v]", n.left, n.right)))
		}
	}
}

// opcode is an enum that specifies the type of node. The zero value is a
// value node with no children.
type opcode rune

const (
	opValue      opcode = 0   // a value literal
	opPair       opcode = ',' // a pair.
	opGroupStart opcode = '[' // open bracket -  only used while parsing
	opGroupEnd   opcode = ']' // close bracket - only used while parsing
)

// New parses a number and returns it.
//
// Example:
//     tree := New("[[1,2],[[3,4],5]]")
//     Magnitude(tree) // output: 143
func New(text string) (*Tree, error) {
	postfix, err := shuntingYard(text)
	if err != nil {
		return nil, err
	}

	t := Tree{
		infix: make([]*node, len(postfix)),
	}

	s := Stack{}
	for _, n := range postfix {
		t.infix[n.id] = n

		if n.op == opValue {
			s.Push(n)
			continue
		}

		if s.Length() < 2 {
			return nil, errors.New("missing operand")
		}
		right := s.Pop().(*node)

		left := s.Pop().(*node)

		n.left, n.right = left, right
		s.Push(n)
	}

	if s.Length() != 1 {
		return nil, errors.New("missing operator")
	}

	top := s.Pop().(*node)
	t.root = top

	return &t, nil
}

// shuntingYard reads tokens from the infix notation and converts
// them to postfix notation so that they can more easily be evaluated.
// Discards whitespace, but otherwise requires all input to be valid.
// Assumes all numeric inputs are a single digit.
// The returned nodes will have their id, value, and operator defined, but
// pairs won't have their left and right children pointers set yet. Therefore,
// this function should really only be called from within New() which adds the
// child pointers, making the tree useful.
func shuntingYard(infix string) (postfix []*node, err error) {
	// reverse polish notation
	rpn := make([]*node, 0, 64)
	s := Stack{}

	var nextID int

	for i, token := range infix {
		switch token {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			rpn = append(rpn, &node{id: nextID, value: int(token - '0')})
			nextID++

		case ',':
			s.Push(&node{id: nextID, op: opPair})
			nextID++

		case '[':
			// brackets don't get any ID - they will be discarded
			s.Push(&node{op: opGroupStart})

		case ']':
			done := false
			for !done && s.Length() > 0 {
				top := s.Pop().(*node)
				if top.op == opGroupStart {
					done = true
					break
				}
				rpn = append(rpn, top)
			}
			if !done {
				return nil, fmt.Errorf("index %d: mismatched closing bracket", i)
			}

		default:
			if unicode.IsSpace(token) {
				continue
			}
			return nil, fmt.Errorf("index %d: unexpected input %q", i, token)
		}
	}

	for s.Length() > 0 {
		operator := s.Pop()
		rpn = append(rpn, operator.(*node))
	}

	return rpn, nil
}

func reduce(n *Tree) {
}
