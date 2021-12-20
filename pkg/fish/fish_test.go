package fish

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var examples = []struct {
	name      string
	infix     string
	postfix   string
	root      Pair
	value     int
	magnitude int
}{
	{
		name:    "example 1",
		infix:   "[[1,2],[[3,4],5]]",
		postfix: "12+34+5++",
		root: &Add{
			L: &Add{L: N(1), R: N(2)},
			R: &Add{
				L: &Add{L: N(3), R: N(4)},
				R: N(5),
			},
		},
		value:     15,
		magnitude: 143,
	},
	{
		name: "example 2",
		infix: `[
			[[[6,6],[7,6]],[[7,7],[7,0]]],
			[[[7,7],[7,7]],[[7,8],[9,9]]]
		]`,
		postfix: "66+76++77+70+++77+77++78+99++++",
		root: &Add{
			L: &Add{
				L: &Add{
					L: &Add{L: N(6), R: N(6)},
					R: &Add{L: N(7), R: N(6)},
				},
				R: &Add{
					L: &Add{L: N(7), R: N(7)},
					R: &Add{L: N(7), R: N(0)},
				},
			},
			R: &Add{
				L: &Add{
					L: &Add{L: N(7), R: N(7)},
					R: &Add{L: N(7), R: N(7)},
				},
				R: &Add{
					L: &Add{L: N(7), R: N(8)},
					R: &Add{L: N(9), R: N(9)},
				},
			},
		},
		value:     107,
		magnitude: 4140,
	},
}

func TestAdd(t *testing.T) {
	for _, tc := range examples {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			a := assert.New(t)

			a.Equal(tc.magnitude, tc.root.Magnitude())
			a.Equal(tc.value, tc.root.Value())
		})
	}
}

func TestShuntingYard(t *testing.T) {
	for _, tc := range examples {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r, a := require.New(t), assert.New(t)

			got, err := shuntingYard(tc.infix)
			r.NoError(err)

			a.Equal(tc.postfix, string(got))
		})
	}
}

func TestParsePostfix(t *testing.T) {
	for _, tc := range examples {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r, a := require.New(t), assert.New(t)

			got, err := parsePostfix([]rune(tc.postfix))
			r.NoError(err)

			a.Equal(tc.root, got)
			t.Log(got)
		})
	}
}

func TestNew(t *testing.T) {
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

			p, isPair := top.(Pair)
			r.True(isPair)
			Reduce(p)

			want, err := New(tc.want)
			r.NoError(err)

			a.Equal(want, top)
		})
	}
}
