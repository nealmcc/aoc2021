package ast

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var examples = []struct {
	name      string
	infix     string
	postfix   string
	root      Number
	magnitude int
}{
	{
		name:    "example 1",
		infix:   "[[1,3],[[5,7],9]]",
		postfix: "13+57+9++",
		root: pair(3, 0,
			pair(1, 1, &node{id: 0, value: 1}, &node{id: 2, value: 3}),
			pair(7, 1,
				pair(5, 2, &node{id: 4, value: 5}, &node{id: 6, value: 7}),
				&node{id: 8, depth: 2, value: 9},
			),
		),
		magnitude: 237,
	},
	{
		name: "example 2",
		infix: `[
			[[[6,6],[7,6]],[[7,7],[7,0]]],
			[[[7,7],[7,7]],[[7,8],[9,9]]]
		]`,
		postfix: "66+76++77+70+++77+77++78+99++++",
		root: pair(15, 0,
			pair(7, 1,
				pair(3, 2,
					pair(1, 3, &node{id: 0, value: 6}, &node{id: 2, value: 6}),
					pair(5, 3, &node{id: 4, value: 7}, &node{id: 6, value: 6}),
				),
				pair(11, 2,
					pair(9, 3, &node{id: 8, value: 7}, &node{id: 10, value: 7}),
					pair(13, 3, &node{id: 12, value: 7}, &node{id: 14, value: 0}),
				),
			),
			pair(23, 1,
				pair(19, 2,
					pair(17, 3, &node{id: 16, value: 7}, &node{id: 18, value: 7}),
					pair(21, 3, &node{id: 20, value: 7}, &node{id: 22, value: 7}),
				),
				pair(27, 2,
					pair(25, 3, &node{id: 24, value: 7}, &node{id: 26, value: 8}),
					pair(29, 3, &node{id: 28, value: 9}, &node{id: 30, value: 9}),
				),
			),
		),
		magnitude: 4140,
	},
}

func TestNew(t *testing.T) {
	for _, tc := range examples {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r, a := require.New(t), assert.New(t)

			got, err := New(tc.infix)
			r.NoError(err)

			t.Logf("got tree: %v", got.root)
			a.Equal(tc.postfix, fmt.Sprintf("%+v", got.root))
			a.Equal(got.root, tc.root)
			a.Equal(tc.magnitude, got.Magnitude())
		})
	}
}

func TestErrorHandling(t *testing.T) {
	tt := []struct {
		name  string
		infix string
	}{
		{"Missing closing bracket", "[2,"},
		{"Extra closing bracket", "[2,2]]"},
		{"Missing right-hand child", "[2,[]]"},
		{"Unexpected symbol", "[2,a]"},
		{"Missing operator", "[22]"},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			_, err := New(tc.infix)
			require.Error(t, err)
		})
	}
}

func TestExplode(t *testing.T) {
	tt := []struct {
		name  string
		in    string
		index int
		want  string
	}{
		{
			name:  "left-most node (id = 1)",
			in:    "[[[[[9,8],1],2],3],4]",
			index: 1,
			want:  "[[[[0,9],2],3],4]",
		},
		{
			name:  "a middle node (id = 5)",
			in:    "[4,[[3,[[7,9],6]],8]]",
			index: 5,
			want:  "[4,[[10,[0,15]],8]]",
		},
		{
			name:  "a right-most node (id = 9)",
			in:    "[7,[6,[5,[4,[3,2]]]]]",
			index: 9,
			want:  "[7,[6,[5,[7,0]]]]",
		},
		{
			name:  "right node in a large left subtree (id = 7)",
			in:    "[[3,[2,[1,[7,3]]]],[6,[5,[4,[3,2]]]]]",
			index: 7,
			want:  "[[3,[2,[8,0]]],[9,[5,[4,[3,2]]]]]",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, a := require.New(t), assert.New(t)
			tree, err := New(tc.in)
			r.NoError(err)

			tree.explode(tree.infix[tc.index])

			// in some cases, the 'want' string will not be valid for parsing,
			// because some numbers will be > 9
			// so we compare the string output instead.
			a.Equal(tc.want, fmt.Sprintf("%v", tree.root))
		})
	}
}

func TestReduce(t *testing.T) {
	tt := []struct {
		name string
		in   string
		want string
	}{
		{
			"reduce example 1",
			"[[[[[4,3],4],4],[7,[[8,4],9]]],[1,1]]",
			"[[[[0,7],4],[[7,8],[6,0]]],[8,1]]",
		},
		{
			"example 2",
			"[[[[[1,1],[2,2]],[3,3]],[4,4]],[5,5]]",
			"[[[[3,0],[5,3]],[4,4]],[5,5]]",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, a := require.New(t), assert.New(t)

			tree, err := New(tc.in)
			r.NoError(err)

			tree.reduce()

			want, err := New(tc.want)
			r.NoError(err)

			a.Equal(want.root, tree.root)
		})
	}
}

// pair is a utility function to create a pointer to a pair of nodes.
func pair(id, depth int, left, right *node) *node {
	left.depth = depth + 1
	right.depth = depth + 1
	return &node{id: id, depth: depth, op: opPair, left: left, right: right}
}
