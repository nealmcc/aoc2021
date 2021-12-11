package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_score(t *testing.T) {
	tt := []struct {
		name string
		in   errCorrupted
		want int
	}{
		{"round", errCorrupted{got: ')'}, 3},
		{"square", errCorrupted{got: ']'}, 57},
		{"curly", errCorrupted{got: '}'}, 1197},
		{"angle", errCorrupted{got: '>'}, 25137},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			got := tc.in.Score()
			require.Equal(t, tc.want, got)
		})
	}
}
