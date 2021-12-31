package main

import (
	"bytes"
	"os"
	"testing"

	"github.com/nealmcc/aoc2021/pkg/alu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// parseUint converts an unsigned int into a slice of 14 digits. It is used
// during testing, to provide input to the ALU and the monad function above.
func parseUint(n uint) []byte {
	bytes := make([]byte, 14)
	var b byte
	for i := 13; n > 0 && i >= 0; i-- {
		n, b = n/10, byte(n%10)
		bytes[i] = b
	}
	return bytes
}

func TestParseUint(t *testing.T) {
	tt := []struct {
		name string
		in   uint
		want []byte
	}{
		{
			"largest value",
			99999999999999,
			[]byte{9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9, 9},
		},
		{
			"12345678901234",
			12345678901234,
			[]byte{1, 2, 3, 4, 5, 6, 7, 8, 9, 0, 1, 2, 3, 4},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := parseUint(tc.in)
			assert.Equal(t, tc.want, got)
		})
	}
}

func TestMonad_vs_ALU(t *testing.T) {
	r := require.New(t)

	src, err := os.ReadFile("input.txt")
	r.NoError(err)

	code, err := alu.Compile(src)
	r.NoError(err)

	calc := alu.New(code)

	tt := []struct {
		name    string
		in      []byte
		want    int
		wantErr bool
	}{
		{
			name: "all 9s",
			in:   parseUint(99999999999999),
		},
		{
			name: "all 5s",
			in:   parseUint(55555555555555),
		},
		{
			name: "all 1s",
			in:   parseUint(11111111111111),
		},
		{
			name: "experiment",
			in:   parseUint(94399898949959),
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r := require.New(t)
			// calc.EnableTrace(zaptest.NewLogger(t).Sugar())

			aluResult, err := calc.Run(bytes.NewReader(tc.in))
			if tc.wantErr {
				r.Error(err)
				return
			}

			comparison := monad(tc.in)
			assert.Equal(t, aluResult, comparison)
			if aluResult == 0 {
				t.Log(tc.in)
			}
		})
	}
}

func TestBackward(t *testing.T) {
	for i := 0; i < 14; i++ {
		for digit := 1; digit <= 9; digit++ {
			for z1 := 0; z1 < 26; z1++ {
				zPrev := backward(z1, digit, i)
				for _, z0 := range zPrev {
					check := forward(z0, digit, i)
					assert.Equal(t, z1, check,
						"got forward(%d, %d, %d) = %d ; want %d",
						z0, digit, i, z1)
				}
			}
		}
	}
}
