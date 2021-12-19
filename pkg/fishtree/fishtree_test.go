package fishtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var examples = []struct {
	name      string
	infix     string
	postfix   string
	root      Node
	value     int
	magnitude int
}{
	{
		name:    "example 1",
		infix:   "[[1,2],[[3,4],5]]",
		postfix: "12+34+5++",
		root: Add{
			L: Add{L: Num(1), R: Num(2)},
			R: Add{
				L: Add{L: Num(3), R: Num(4)},
				R: Num(5),
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
		root: Add{
			L: Add{
				L: Add{
					L: Add{L: Num(6), R: Num(6)},
					R: Add{L: Num(7), R: Num(6)},
				},
				R: Add{
					L: Add{L: Num(7), R: Num(7)},
					R: Add{L: Num(7), R: Num(0)},
				},
			},
			R: Add{
				L: Add{
					L: Add{L: Num(7), R: Num(7)},
					R: Add{L: Num(7), R: Num(7)},
				},
				R: Add{
					L: Add{L: Num(7), R: Num(8)},
					R: Add{L: Num(9), R: Num(9)},
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

func Test_inputValidation(t *testing.T) {
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
