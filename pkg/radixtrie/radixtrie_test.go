package radixtrie

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// this example tree contains:
// - test
// - toaster
// - toasting
// - slow
// - slowly
var example Node = Node{
	edges: []*Edge{
		{
			label: "t",
			target: Node{
				edges: []*Edge{
					{label: "est", target: Node{hasValue: true}},
					{label: "oast", target: Node{
						edges: []*Edge{
							{label: "er", target: Node{hasValue: true}},
							{label: "ing", target: Node{hasValue: true}},
						},
					}},
				},
			},
		},
		{
			label: "slow",
			target: Node{
				hasValue: true,
				edges: []*Edge{
					{label: "ly", target: Node{hasValue: true}},
				},
			},
		},
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
