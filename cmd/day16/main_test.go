package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var _examples = []struct {
	hex        string
	versionSum int
}{
	{"D2FE28", 6},
	{"8A004A801A8002F478", 16},
	{"620080001611562C8802118E34", 12},
	{"C0015000016115A2E0802F182340", 23},
	{"A0016C880162017C3686B18A3D4780", 31},
}

func Test_part1(t *testing.T) {
	for i, tc := range _examples {
		tc := tc
		t.Run(fmt.Sprintf("example %d", i+1), func(t *testing.T) {
			r, a := require.New(t), assert.New(t)
			packet, err := read(strings.NewReader(tc.hex))
			r.NoError(err)

			a.Equal(tc.versionSum, part1(packet))
		})
	}
}
