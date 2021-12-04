// Package radixtrie implements a radix tree or trie
//
// more info: https://en.wikipedia.org/wiki/Radix_tree
package radixtrie

import (
	"strings"
)

type Edge struct {
	target Node
	label  string
}

type Node struct {
	edges    []*Edge
	hasValue bool
}

// Contains searches the tree for the given prefix, and returns a subtree
// of nodes which match as much of the prefix as possible, and the length
// of the prefix that was matched.
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
func (n *Node) Contains(needle string) bool {
	numFound := 0

	for n != nil && len(n.edges) > 0 && numFound < len(needle) {
		var next *Edge
		for _, e := range n.edges {
			if strings.HasPrefix(needle[numFound:], e.label) {
				next = e
				break
			}
		}

		if next == nil {
			n = nil
			break
		}

		n = &next.target
		numFound += len(next.label)
	}

	return n != nil && n.hasValue && numFound == len(needle)
}
