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
		root: pair(3,
			pair(1, &node{id: 0, value: 1}, &node{id: 2, value: 3}),
			pair(7,
				pair(5, &node{id: 4, value: 5}, &node{id: 6, value: 7}),
				&node{id: 8, value: 9},
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
		root: pair(15,
			pair(7,
				pair(3,
					pair(1, &node{id: 0, value: 6}, &node{id: 2, value: 6}),
					pair(5, &node{id: 4, value: 7}, &node{id: 6, value: 6}),
				),
				pair(11,
					pair(9, &node{id: 8, value: 7}, &node{id: 10, value: 7}),
					pair(13, &node{id: 12, value: 7}, &node{id: 14, value: 0}),
				),
			),
			pair(23,
				pair(19,
					pair(17, &node{id: 16, value: 7}, &node{id: 18, value: 7}),
					pair(21, &node{id: 20, value: 7}, &node{id: 22, value: 7}),
				),
				pair(27,
					pair(25, &node{id: 24, value: 7}, &node{id: 26, value: 8}),
					pair(29, &node{id: 28, value: 9}, &node{id: 30, value: 9}),
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
	t.SkipNow()

	tt := []struct {
		name, in, want string
	}{
		{
			"left-most node",
			"[[[[[9,8],1],2],3],4]",
			"[[[[0,9],2],3],4]",
		},
		{
			"a middle node",
			"[4,[[3,[[7,9],6]],3]]",
			"[1,0]",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, a := require.New(t), assert.New(t)
			tree, err := New(tc.in)
			r.NoError(err)

			t.Logf("before: %v", tree)

			want, err := New(tc.want)
			r.NoError(err)
			t.Logf("want: %v", want)
			a.Equal(want, tree)
		})
	}
}

func TestReduce(t *testing.T) {
	t.SkipNow()

	tt := []struct {
		name string
		in   string
		want string
	}{
		{
			"ones to fives",
			"[[[[[1,1],[2,2]],[3,3]],[4,4]], [5,5]]",
			"[[[[3,0],[5,3]],[4,4]],[5,5]]",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, a := require.New(t), assert.New(t)

			top, err := New(tc.in)
			r.NoError(err)

			reduce(top)

			want, err := New(tc.want)
			r.NoError(err)

			a.Equal(want, top)
		})
	}
}

// pair is a utility function to create a pointer to a pair of nodes.
func pair(id int, left, right *node) *node {
	return &node{id: id, op: opPair, left: left, right: right}
}
