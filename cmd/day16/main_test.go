package main

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func Test_part1(t *testing.T) {
	tt := []struct {
		hex        string
		versionSum int
	}{
		{"D2FE28", 6},
		{"38006F45291200", 9},
		{"8A004A801A8002F478", 16},
		{"620080001611562C8802118E34", 12},
		{"C0015000016115A2E0802F182340", 23},
		{"A0016C880162017C3686B18A3D4780", 31},
	}
	for _, tc := range tt {
		tc := tc
		t.Run(fmt.Sprintf("example %s", tc.hex), func(t *testing.T) {
			t.Parallel()

			r, a := require.New(t), assert.New(t)
			packet, err := read(strings.NewReader(tc.hex))
			r.NoError(err)

			log := zaptest.NewLogger(t).Sugar()

			p1, err := part1(*packet, log)
			r.NoError(err)

			a.Equal(tc.versionSum, p1)
		})
	}
}

func Test_part2(t *testing.T) {
	tt := []struct {
		name  string
		hex   string
		value int
	}{
		{
			"sum of 1 and 2",
			"C200B40A82", 3,
		},
		{
			"product of 6 and 9",
			"04005AC33890", 54,
		},
		{
			"min of 7,8,9",
			"880086C3E88112", 7,
		},
		{
			"max of 7,8,9",
			"CE00C43D881120", 9,
		},
		{
			"5 less than 15",
			"D8005AC2A8F0", 1,
		},
		{
			"5 greater than 15",
			"F600BC2D8F", 0,
		},
		{
			"5 equal to 15",
			"9C005AC2F8F0", 0,
		},
		{
			"1+3 == 2*2",
			"9C0141080250320F1802104A08", 1,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			r, a := require.New(t), assert.New(t)
			packet, err := read(strings.NewReader(tc.hex))
			r.NoError(err)

			log := zaptest.NewLogger(t).Sugar()
			packet.SetLogger(log)

			got, err := packet.Value()
			r.NoError(err)

			a.Equal(tc.value, got)
		})
	}
}
