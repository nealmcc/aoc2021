package radixtree

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// this example tree contains:
// - test
// - toaster
// - toasting
// - slow
// - slowly
var example Node = Node{
	children: map[byte]*Node{
		't': {children: map[byte]*Node{
			'e': {children: map[byte]*Node{
				's': {children: map[byte]*Node{
					't': {hasValue: true},
				}},
			}},
			'o': {children: map[byte]*Node{
				'a': {children: map[byte]*Node{
					's': {children: map[byte]*Node{
						't': {children: map[byte]*Node{
							'e': {children: map[byte]*Node{
								'r': {hasValue: true},
							}},
							'i': {children: map[byte]*Node{
								'n': {children: map[byte]*Node{
									'g': {hasValue: true},
								}},
							}},
						}},
					}},
				}},
			}},
		}},
		's': {children: map[byte]*Node{
			'l': {children: map[byte]*Node{
				'o': {children: map[byte]*Node{
					'w': {
						hasValue: true,
						children: map[byte]*Node{
							'l': {children: map[byte]*Node{
								'y': {hasValue: true},
							}},
						},
					},
				}},
			}},
		}},
	},
}

func TestContains(t *testing.T) {
	tt := []struct {
		name   string
		needle string
		want   bool
	}{
		{
			name:   "empty string is not found",
			needle: "",
			want:   false,
		},
		{
			name:   "zzzz is not found",
			needle: "zzzz",
			want:   false,
		},
		{
			name:   "slalom is not found",
			needle: "slalom",
			want:   false,
		},
		{
			name:   "test is found",
			needle: "test",
			want:   true,
		},
		{
			name:   "slowly is found",
			needle: "slowly",
			want:   true,
		},
		{
			name:   "slow is found",
			needle: "slow",
			want:   true,
		},
		{
			name:   "toasting is found",
			needle: "toasting",
			want:   true,
		},
		{
			name:   "toast is not found",
			needle: "toast",
			want:   false,
		},
		{
			name:   "toaster is found",
			needle: "toaster",
			want:   true,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := example.Contains(tc.needle)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestWithPrefix(t *testing.T) {
	tt := []struct {
		name      string
		prefix    string
		wantExact bool
		wantLen   int
	}{
		{
			name:    "empty string returns the full tree",
			prefix:  "",
			wantLen: 0,
		},
		{
			name:    "zzzz returns the full tree",
			prefix:  "zzzz",
			wantLen: 0,
		},
		{
			name:    "slalom returns a subtree with 2 letters matched",
			prefix:  "slalom",
			wantLen: 2,
		},
		{
			name:      "test is found",
			prefix:    "test",
			wantExact: true,
			wantLen:   4,
		},
		{
			name:      "slowly is found",
			prefix:    "slowly",
			wantExact: true,
			wantLen:   6,
		},
		{
			name:      "slow is found",
			prefix:    "slow",
			wantExact: true,
			wantLen:   4,
		},
		{
			name:      "toasting is found",
			prefix:    "toasting",
			wantExact: true,
			wantLen:   8,
		},
		{
			name:    "toast returns a subtree, with 5 letters matched",
			prefix:  "toast",
			wantLen: 5,
		},
		{
			name:      "toaster is found",
			prefix:    "toaster",
			wantExact: true,
			wantLen:   7,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			node, length := example.WithPrefix(tc.prefix)
			require.NotNil(t, node)
			assert.Equal(t, tc.wantLen, length)
			assert.Equal(t, tc.wantExact, node.hasValue)
		})
	}
}

func TestInsert(t *testing.T) {
	t.Parallel()
	dict := &Node{}
	dict.Insert("Hello, World")
	world, length := dict.WithPrefix("Hello, ")
	assert.Equal(t, 7, length)
	got := world.Contains("World")
	assert.True(t, got)
}
