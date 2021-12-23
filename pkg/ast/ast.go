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
	infix map[int]*node
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
// the zero node is a value node with a value of 0.
type node struct {
	id    int    // id is this node's index during an in-order traversal of the tree
	op    opcode // op specifies this node's type (either opValue or opPair).
	value int    // if op == opValue, then this holds the node's value
	left  *node  // if op != opValue, this points to the node's left child.
	right *node  // if op != opValue, this points the node's right child.
	depth int    // depth is the number of edges between the root and this node.
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
		infix: make(map[int]*node, len(postfix)),
	}

	s := stack{}
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
	t.root.setDepth(0)
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
	s := stack{}

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

// Sum returns the sum of all the given fish numbers, correctly reducing
// the results of each addition. The input parameters will be unaffected by the
// addition.
func Sum(numbers ...*Tree) (sum *Tree, err error) {
	if len(numbers) == 0 {
		zero := &node{}
		return &Tree{root: zero, infix: map[int]*node{0: zero}}, nil
	}

	sum = numbers[0]
	for i := 1; i < len(numbers); i++ {
		sum, err = Add(*sum, *numbers[i])
		if err != nil {
			return nil, err
		}
	}
	return sum, nil
}

// Add trees a and b, returning a pointer to their sum, which is reduced.
// Trees a and b are unaffected.
func Add(a, b Tree) (*Tree, error) {
	sum := fmt.Sprintf("[%v,%v]", a.root, b.root)
	t, err := New(sum)
	if err != nil {
		return nil, err
	}

	t.reduce()
	return t, nil
}

// reduce adjusts is tree following the rules from Day18:
// To reduce a snailfish number, you must repeatedly do the first action in
// this list that applies to the snailfish number:
//
// - If any pair is nested inside four pairs, the leftmost such pair explodes.
// - If any regular number is 10 or greater, the leftmost such regular number splits.
//
// Once no action in the above list applies, the snailfish number is reduced.
// During reduction, at most one action applies, after which the process returns
// to the top of the list of actions. For example, if split produces a pair that
// meets the explode criteria, that pair explodes before other splits occur.
func (t *Tree) reduce() bool {
	changed := false
	for {
		t.root.setDepth(0)
		if tooDeep := t.firstDepth4(); tooDeep != nil {
			t.explode(tooDeep)
			changed = true
			continue
		}

		tooBig := t.firstOver9()
		if tooBig == nil {
			break
		}
		t.split(tooBig)
		changed = true
	}
	return changed
}

// firstDepth4 finds the left-most pair that is nested inside four pairs.
// Returns nil if there is no such pair.
func (t *Tree) firstDepth4() *node {
	for i := 0; i < len(t.infix); i++ {
		if n, ok := t.infix[i]; ok && n.op == opPair && n.depth == 4 {
			return n
		}
	}
	return nil
}

// firstOver9 finds the if of the left-most leaf that is 10 or greater.
// Returns nil if there is no such leaf.
func (t *Tree) firstOver9() *node {
	for i := 0; i < len(t.infix); i++ {
		if n, ok := t.infix[i]; ok && n.op == opValue && n.value > 9 {
			return n
		}
	}
	return nil
}

// explode the given node, as per the rules on Day18.
// Panics if the given id is not in the tree, or is not a pair.
// Assumes the pair has two value children (which it will, if the rules of
// fish math have been followed.)
func (t *Tree) explode(n *node) {
	if n.op != opPair {
		panic(fmt.Errorf("tring to explode a leaf node: %d", n.id))
	}

	if n.left.op != opValue || n.right.op != opValue {
		panic(fmt.Errorf("trying to explode a pair with a child pair: %d", n.id))
	}

	if t.infix[n.left.id] != n.left {
		panic(fmt.Errorf("left id and left child do not match"))
	}

	if t.infix[n.right.id] != n.right {
		panic(fmt.Errorf("right id and right child do not match"))
	}

	if t.infix[n.id] != n {
		panic(fmt.Errorf("self id and self do not match"))
	}

	// add this node's left child to the leaf that's just to the left
	if leaf := t.findLeafLeftOf(n.id - 1); leaf != nil {
		leaf.value += n.left.value
	}

	// add this node's right child to the leaf that's just to the right
	if leaf := t.findLeafRightOf(n.id + 1); leaf != nil {
		leaf.value += n.right.value
	}

	// set this node to zero, and prevent memory leaks
	n.op, n.value = opValue, 0
	n.left, n.right = nil, nil

	// fix this node's ID
	n.id -= 1
	t.infix[n.id] = n

	// fix IDs of nodes to the right of this one
	lastID := len(t.infix) - 3
	for i := n.id + 1; i <= lastID; i++ {
		n, ok := t.infix[i+2]
		if !ok {
			panic(fmt.Errorf("something is wrong with index %d", i+2))
		}
		n.id = i
		t.infix[n.id] = n
	}

	delete(t.infix, lastID+1)
	delete(t.infix, lastID+2)
}

// split divides the value of a leaf node in half, converting that leaf
// to a pair, with one half of the value in each child.
func (t *Tree) split(n *node) {
	if n.op != opValue {
		panic(fmt.Errorf("trying to split a non-leaf node: %v", n))
	}

	if t.infix[n.id] != n {
		panic(fmt.Errorf("self id and self do not match"))
	}

	// create child nodes (no IDs yet)
	half, rem := n.value/2, n.value%2
	n.op, n.value = opPair, 0
	n.left = &node{value: half, depth: n.depth + 1}
	n.right = &node{value: half + rem, depth: n.depth + 1}

	// shuffle ids for nodes to the right of this pair right by two:
	lastID := len(t.infix) - 1
	for id := lastID; id > n.id; id-- {
		nn := t.infix[id]
		nn.id += 2
		t.infix[nn.id] = nn
	}

	// this node's id increases by one
	n.id += 1
	t.infix[n.id] = n

	// insert children to map (child IDs were set above)
	n.left.id = n.id - 1
	n.right.id = n.id + 1
	t.infix[n.left.id] = n.left
	t.infix[n.right.id] = n.right
}

// findLeafLeftOf finds the next leaf to the left of the given id.
func (t Tree) findLeafLeftOf(id int) *node {
	for id -= 1; id >= 0; id-- {
		if n, ok := t.infix[id]; ok && n.op == opValue {
			return n
		}
	}
	return nil
}

// findLeafRightOf finds the next leaf to the right of the given id.
func (t Tree) findLeafRightOf(id int) *node {
	for id += 1; id < len(t.infix); id++ {
		if n, ok := t.infix[id]; ok && n.op == opValue {
			return n
		}
	}
	return nil
}

// setDepth recursively sets the depth of this node and its children.
func (n *node) setDepth(d int) {
	n.depth = d
	if n.left != nil {
		n.left.setDepth(d + 1)
	}
	if n.right != nil {
		n.right.setDepth(d + 1)
	}
}
