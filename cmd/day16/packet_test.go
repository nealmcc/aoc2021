package main

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPacket_Unmarshal(t *testing.T) {
	for i, tc := range _examples {
		tc := tc
		t.Run(fmt.Sprintf("example %d", i), func(t *testing.T) {
			r, a := require.New(t), assert.New(t)
			packet, err := read(strings.NewReader(tc.hex))
			r.NoError(err)

			a.Equal(strings.ToLower(tc.hex), fmt.Sprintf("%s", packet))
		})
	}
}

func TestPacket_Version(t *testing.T) {
	p := new(Packet)
	err := p.UnmarshalText([]byte("D2FE28"))
	require.NoError(t, err)

	assert.Equal(t, 6, p.Version())
}

func TestPacket_Value(t *testing.T) {
	p := new(Packet)
	err := p.UnmarshalText([]byte("D2FE28"))
	require.NoError(t, err)

	assert.Equal(t, _literalValue, p.packetType())
	assert.Equal(t, 2021, p.Value())
}

func TestPacket_BitLength(t *testing.T) {
	tt := []struct {
		name       string
		in         string
		wantValue  int
		wantLength int
	}{
		{
			name:       "literal value 1 has 11 bits",
			in:         "5020",
			wantValue:  1,
			wantLength: 11,
		},
		{
			name:       "literal value 2 has 11 bits",
			in:         "9040",
			wantValue:  2,
			wantLength: 11,
		},
		{
			name:       "literal value 3 has 11 bits",
			in:         "9060",
			wantValue:  3,
			wantLength: 11,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			p := new(Packet)
			err := p.UnmarshalText([]byte(tc.in))
			require.NoError(t, err)

			assert.Equal(t, tc.wantValue, p.Value())
			assert.Equal(t, tc.wantLength, p.bitLength())
		})
	}
}

func TestPacket_eightBits(t *testing.T) {
	tt := []struct {
		name     string
		in       string
		bitIndex int
		want     byte
		wantErr  error
	}{
		{
			name:     "read one byte",
			in:       "FF",
			bitIndex: 0,
			want:     0xFF,
		},
		{
			name:     "get EOF from an empty packet",
			in:       "",
			bitIndex: 0,
			wantErr:  io.EOF,
		},
		{
			name:     "read from the start of a packet",
			in:       "F0F0",
			bitIndex: 0,
			want:     0xF0,
		},
		{
			name:     "read an easy overlapping byte",
			in:       "F0F0",
			bitIndex: 4,
			want:     0x0F,
		},
		{
			name:     "read an awkward overlapping byte",
			in:       "F0F0", // 1111 0000 1111 0000
			bitIndex: 2,
			want:     0xC3, // 1100 0011
		},
		{
			name:     "read from the end of a packet",
			in:       "F0F0",
			bitIndex: 8,
			want:     0xF0,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			r, a := require.New(t), assert.New(t)

			p := new(Packet)
			err := p.UnmarshalText([]byte(tc.in))
			r.NoError(err)

			got, err := p.eightBits(tc.bitIndex)
			a.Equal(tc.want, got)
			a.Equal(tc.wantErr, err)
		})
	}
}
