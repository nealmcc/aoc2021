package alu

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCompile(t *testing.T) {
	tt := []struct {
		name      string
		in        string
		wantCount int
		wantErr   bool
	}{
		{
			name:      "can parse inp instruction",
			in:        "inp w\n",
			wantCount: 1,
		},
		{
			name:      "can parse add instruction with register",
			in:        "add y z\n",
			wantCount: 1,
		},
		{
			name:      "can parse add instruction with value",
			in:        "add y -42\n",
			wantCount: 1,
		},
		{
			name:      "can parse mul instruction with register",
			in:        "mul y z\n",
			wantCount: 1,
		},
		{
			name:      "can parse mul instruction with value",
			in:        "mul y -42\n",
			wantCount: 1,
		},
		{
			name:      "can parse div instruction with register",
			in:        "div y z\n",
			wantCount: 1,
		},
		{
			name:      "can parse div instruction with value",
			in:        "div y -42\n",
			wantCount: 1,
		},
		{
			name:    "prevents division by 0 when ALU is initialised",
			in:      "div y 0\n",
			wantErr: true,
		},
		{
			name:      "can parse mod instruction with register",
			in:        "mod y z\n",
			wantCount: 1,
		},
		{
			name:      "can parse mod instruction with value",
			in:        "mod y -42\n",
			wantCount: 1,
		},
		{
			name:    "prevents modulo 0 when ALU is initialised",
			in:      "mod y 0\n",
			wantErr: true,
		},
		{
			name:      "can parse eql instruction with register",
			in:        "eql y z\n",
			wantCount: 1,
		},
		{
			name:      "can parse eql instruction with value",
			in:        "eql y -42\n",
			wantCount: 1,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			r, a := require.New(t), assert.New(t)
			got, err := Compile([]byte(tc.in))
			if tc.wantErr {
				r.Error(err)
				return
			}

			r.NoError(err)
			a.Equal(tc.wantCount, len(got))
		})
	}
}
