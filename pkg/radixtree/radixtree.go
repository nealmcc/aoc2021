// Package radixtree implements a radix tree
//
// more info: https://en.wikipedia.org/wiki/Radix_tree
package radixtree

// Node is a node in a radix tree.  The zero value is ready to use.
// This implementation is not safe for concurrent use.
type Node struct {
	hasValue bool
	children map[byte]*Node
}

// Contains determines if the tree contains the exact give string.
func (n *Node) Contains(needle string) bool {
	numFound := 0

	for n != nil && len(n.children) > 0 && numFound < len(needle) {
		n = n.children[needle[numFound]]
		if n != nil {
			numFound++
		}
	}

	return n != nil && n.hasValue && numFound == len(needle)
}

// WithPrefix searches the tree for a subtree of nodes which match as much
// of the prefix as possible.
//
// Ex:
// If the tree contains the following strings:
// - Aloha
// - Alabama
// - Alaska
//
// then searching for 'Alakazam' will return the subtree containing
// - Alabama
// - Alaska
// and a length of 3, indicating that 3 letters matched.
func (n *Node) WithPrefix(prefix string) (*Node, int) {
	numFound := 0

	curr := n
	for numFound < len(prefix) {
		next := curr.children[prefix[numFound]]
		if next == nil {
			break
		}
		numFound++
		curr = next
	}

	return curr, numFound
}

// Insert the given string to the tree.
func (n *Node) Insert(s string) {
	for len(s) > 0 {
		if n.children == nil {
			n.children = make(map[byte]*Node)
		}
		if n.children[s[0]] == nil {
			n.children[s[0]] = &Node{}
		}
		n = n.children[s[0]]
		s = s[1:]
	}
	n.hasValue = true
}

// ToSlice converts this tree into a slice of values. The order is undefined.
func (n *Node) ToSlice(path ...byte) []string {
	// todo: improve efficiency
	words := make([]string, 0)
	if n.hasValue {
		words = append(words, string(path))
	}
	for k, v := range n.children {
		nextPath := append(path, k)
		words = append(words, v.ToSlice(nextPath...)...)
	}
	return words
}

// Size returns the number of values in the tree.
func (n *Node) Size() int {
	sum := 0
	if n.hasValue {
		sum++
	}
	for _, v := range n.children {
		sum += v.Size()
	}
	return sum
}
