package main

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nealmcc/aoc2021/pkg/ast"
)

func TestRead(t *testing.T) {
	r, a := require.New(t), assert.New(t)

	nodes, err := read(strings.NewReader("[[1,2],[[3,4],5]]"))
	r.NoError(err)

	r.Equal(1, len(nodes))
	a.Equal(143, nodes[0].Magnitude())
}

func TestPart1(t *testing.T) {
	tt := []struct {
		name string
		in   []string
		want string
	}{
		{
			"easy example",
			[]string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
			},
			"[[[[1,1],[2,2]],[3,3]],[4,4]]",
		},
		{
			"example with a reduction",
			[]string{
				"[1,1]",
				"[2,2]",
				"[3,3]",
				"[4,4]",
				"[5,5]",
			},
			"[[[[3,0],[5,3]],[4,4]],[5,5]]",
		},
		{
			"example A",
			[]string{
				"[[[0,[4,5]],[0,0]],[[[4,5],[2,6]],[9,5]]]",
				"[7,[[[3,7],[4,3]],[[6,3],[8,8]]]]",
				"[[2,[[0,8],[3,4]]],[[[6,7],1],[7,[1,6]]]]",
				"[[[[2,4],7],[6,[0,5]]],[[[6,8],[2,8]],[[2,1],[4,5]]]]",
				"[7,[5,[[3,8],[1,4]]]]",
				"[[2,[2,2]],[8,[8,1]]]",
				"[2,9]",
				"[1,[[[9,3],9],[[9,0],[0,7]]]]",
				"[[[5,[7,4]],7],1]",
				"[[[[4,2],2],6],[8,7]]",
			},
			"[[[[8,7],[7,7]],[[8,6],[7,7]]],[[[0,7],[6,6]],[8,7]]]",
		},
		{
			"homework",
			[]string{
				"[[[0,[5,8]],[[1,7],[9,6]]],[[4,[1,2]],[[1,4],2]]]",
				"[[[5,[2,8]],4],[5,[[9,9],0]]]",
				"[6,[[[6,2],[5,6]],[[7,6],[4,7]]]]",
				"[[[6,[0,7]],[0,9]],[4,[9,[9,0]]]]",
				"[[[7,[6,4]],[3,[1,3]]],[[[5,5],1],9]]",
				"[[6,[[7,3],[3,2]]],[[[3,8],[5,7]],4]]",
				"[[[[5,4],[7,7]],8],[[8,3],8]]",
				"[[9,3],[[9,9],[6,[4,9]]]]",
				"[[2,[[7,7],7]],[[5,8],[[9,3],[0,2]]]]",
				"[[[[5,2],5],[8,[3,7]]],[[5,[7,5]],[4,4]]]",
			},
			"[[[[6,6],[7,6]],[[7,7],[7,0]]],[[[7,7],[7,7]],[[7,8],[9,9]]]]",
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, a := require.New(t), assert.New(t)

			combined := strings.Join(tc.in, "\n")
			lines, err := read(strings.NewReader(combined))
			r.NoError(err)

			got, err := part1(lines)
			r.NoError(err)

			want, err := ast.New(tc.want)
			r.NoError(err)

			a.Equal(want.Magnitude(), got)
		})
	}
}

func TestPart2(t *testing.T) {
	tt := []struct {
		name string
		in   []string
		want int
	}{
		{
			name: "part2 example",
			in: []string{
				"[[[0,[5,8]],[[1,7],[9,6]]],[[4,[1,2]],[[1,4],2]]]",
				"[[[5,[2,8]],4],[5,[[9,9],0]]]",
				"[6,[[[6,2],[5,6]],[[7,6],[4,7]]]]",
				"[[[6,[0,7]],[0,9]],[4,[9,[9,0]]]]",
				"[[[7,[6,4]],[3,[1,3]]],[[[5,5],1],9]]",
				"[[6,[[7,3],[3,2]]],[[[3,8],[5,7]],4]]",
				"[[[[5,4],[7,7]],8],[[8,3],8]]",
				"[[9,3],[[9,9],[6,[4,9]]]]",
				"[[2,[[7,7],7]],[[5,8],[[9,3],[0,2]]]]",
				"[[[[5,2],5],[8,[3,7]]],[[5,[7,5]],[4,4]]]",
			},
			want: 3993,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, a := require.New(t), assert.New(t)

			combined := strings.Join(tc.in, "\n")
			lines, err := read(strings.NewReader(combined))
			r.NoError(err)

			got, err := part2(lines)
			r.NoError(err)

			a.Equal(tc.want, got)
		})
	}
}
