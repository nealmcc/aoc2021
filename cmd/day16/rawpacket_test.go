package main

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRawPacket_Unmarshal(t *testing.T) {
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

func TestRawPacket_Version(t *testing.T) {
	p := new(rawPacket)
	err := p.UnmarshalText([]byte("D2FE28"))
	require.NoError(t, err)

	assert.Equal(t, 6, p.Version())
}

func TestRawPacket_Value(t *testing.T) {
	p := new(rawPacket)
	err := p.UnmarshalText([]byte("D2FE28"))
	require.NoError(t, err)

	assert.Equal(t, _literalValue, p.packetType())
	assert.Equal(t, 2021, p.Value())
}

func TestRawPacket_bits(t *testing.T) {
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

			p := new(rawPacket)
			err := p.UnmarshalText([]byte(tc.in))
			r.NoError(err)

			got, err := p.bits(tc.bitIndex)
			a.Equal(tc.want, got)
			a.Equal(tc.wantErr, err)
		})
	}
}
