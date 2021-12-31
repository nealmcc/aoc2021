package alu

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestRun(t *testing.T) {
	tt := []struct {
		name      string
		program   string
		input     []byte
		output    int
		wantErr   bool
		registers [4]int
	}{
		{
			name: "store 4 bits of input in registers w,x,y,z",
			program: `inp w
add z w
mod z 2
div w 2
add y w
mod y 2
div w 2
add x w
mod x 2
div w 2
mod w 2
`,
			input:     []byte{15},
			output:    1,
			registers: [4]int{1, 1, 1, 1},
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, a := require.New(t), assert.New(t)

			code, err := Compile([]byte(tc.program))
			r.NoError(err)
			alu := New(code)

			alu.EnableTrace(zaptest.NewLogger(t).Sugar())

			got, err := alu.Run(bytes.NewReader(tc.input))
			if tc.wantErr {
				r.Error(err)
			} else {
				r.NoError(err)
				a.Equal(tc.output, got)
			}

			for i := 0; i < 4; i++ {
				a.Equal(tc.registers[i], alu.reg[i], "register %c", byte(i)+byte(regW))
			}
		})
	}
}
