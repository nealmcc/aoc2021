package main

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
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
	r := require.New(t)
	p := new(Packet)
	err := p.UnmarshalText([]byte("D2FE28"))
	r.NoError(err)

	got, err := p.Version()
	r.NoError(err)

	assert.Equal(t, 6, got)
}

func TestPacket_Value(t *testing.T) {
	r := require.New(t)
	p := new(Packet)
	err := p.UnmarshalText([]byte("D2FE28"))
	r.NoError(err)

	pt, err := p.packetType()
	r.NoError(err)

	assert.Equal(t, _literal, pt)

	val, err := p.Value()
	r.NoError(err)
	assert.Equal(t, 2021, val)
}

func TestPacket_bitLength(t *testing.T) {
	tt := []struct {
		name       string
		in         Packet
		wantLength int
	}{
		{
			"literal value 1 has 11 bits",
			testPacket("5020", 0),
			11,
		},
		{
			"literal value 2 has 11 bits",
			testPacket("9040", 0),
			11,
		},
		{
			"literal value 3 has 11 bits",
			testPacket("9060", 0),
			11,
		},
		{
			"literal value 3 has 11 bits",
			testPacket("9060", 0),
			11,
		},
		{
			"literal value 13 has 11 bits, even with a bit offset",
			testPacket("8e34", 3),
			11,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			log := zaptest.NewLogger(t).Sugar()
			r, a := require.New(t), assert.New(t)

			tc.in.SetLogger(log)

			length, err := tc.in.bitLength()
			r.NoError(err)
			a.Equal(tc.wantLength, length)
		})
	}
}

func TestPacket_Children(t *testing.T) {
	r := require.New(t)
	var p Packet
	err := p.UnmarshalText([]byte("EE00D40C823060"))
	r.NoError(err)
	children, err := p.Children()
	r.NoError(err)
	r.Equal(3, len(children))

	tt := []struct {
		child   int
		version int
		value   int
	}{
		{child: 0, version: 2, value: 1},
		{child: 1, version: 4, value: 2},
		{child: 2, version: 1, value: 3},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(fmt.Sprintf("%dth child", tc.child), func(t *testing.T) {
			r, a := require.New(t), assert.New(t)
			child := children[tc.child]

			ver, err := child.Version()
			r.NoError(err)
			a.Equal(tc.version, ver, "version mismatch")

			val, err := child.Value()
			r.NoError(err)
			a.Equal(tc.value, val, "value mismatch")
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

			got, err := p.nBits(tc.bitIndex, 8)
			a.Equal(tc.want, got)
			cause := errors.Cause(err)
			a.ErrorIs(tc.wantErr, cause)
		})
	}
}

func TestPacket_nBits(t *testing.T) {
	tt := []struct {
		name     string
		p        Packet
		bitIndex int
		numBits  fieldWidth
		want     byte
	}{
		{
			name:     "outer packet version",
			p:        testPacket("8A004A801A8002F478", 0),
			bitIndex: 0,
			numBits:  3,
			want:     4,
		},
		{
			name:     "outer packet type",
			p:        testPacket("8A004A801A8002F478", 0),
			bitIndex: 3,
			numBits:  3,
			want:     2, // operator
		},
		{
			name:     "outer packet size type ID",
			p:        testPacket("8A004A801A8002F478", 0),
			bitIndex: 6,
			numBits:  1,
			want:     1, // immediate
		},
		{
			name:     "outer packet size upper",
			p:        testPacket("8A004A801A8002F478", 0),
			bitIndex: 7,
			numBits:  3,
			want:     0,
		},
		{
			name:     "outer packet size lower",
			p:        testPacket("8A004A801A8002F478", 0),
			bitIndex: 10,
			numBits:  8,
			want:     1,
		},
		{
			name:     "child c0 version",
			p:        testPacket("8A004A801A8002F478", 18),
			bitIndex: 0,
			numBits:  3,
			want:     1,
		},
		{
			name:     "child c0 type",
			p:        testPacket("8A004A801A8002F478", 18),
			bitIndex: 3,
			numBits:  3,
			want:     2, // operator
		},
		{
			name:     "child c0 sizeID",
			p:        testPacket("8A004A801A8002F478", 18),
			bitIndex: 6,
			numBits:  1,
			want:     1, // num packets
		},
		{
			name:     "child c0 size upper",
			p:        testPacket("8A004A801A8002F478", 18),
			bitIndex: 7,
			numBits:  3,
			want:     0,
		},
		{
			name:     "child c0 size lower",
			p:        testPacket("8A004A801A8002F478", 18),
			bitIndex: 10,
			numBits:  8,
			want:     1,
		},
		{
			name:     "grandchild gc0 version",
			p:        testPacket("8A004A801A8002F478", 36),
			bitIndex: 0,
			numBits:  3,
			want:     5,
		},
		{
			name:     "grandchild gc0 type",
			p:        testPacket("8A004A801A8002F478", 36),
			bitIndex: 3,
			numBits:  3,
			want:     2, // operator
		},
		{
			name:     "grandchild gc0 size type ID",
			p:        testPacket("8A004A801A8002F478", 36),
			bitIndex: 6,
			numBits:  1,
			want:     0, // immediate
		},
		{
			name:     "grandchild gc0 size upper byte",
			p:        testPacket("8A004A801A8002F478", 36),
			bitIndex: 7,
			numBits:  7,
			want:     0,
		},
		{
			name:     "grandchild gc0 size lower byte",
			p:        testPacket("8A004A801A8002F478", 36),
			bitIndex: 14,
			numBits:  8,
			want:     0b0000_1011, // 11
		},
		{
			name:     "great-grandchild ggc0 version",
			p:        testPacket("8A004A801A8002F478", 58),
			bitIndex: 0,
			numBits:  3,
			want:     6,
		},
		{
			name:     "great-grandchild ggc0 packet type",
			p:        testPacket("8A004A801A8002F478", 58),
			bitIndex: 3,
			numBits:  3,
			want:     4, // literal value
		},
		{
			name:     "great-grandchild ggc0 value chunk 0",
			p:        testPacket("8A004A801A8002F478", 58),
			bitIndex: 6,
			numBits:  5,
			want:     0b0000_1111, // 0 = stop, val = 15, overall = 15
		},
		{
			name:     "read a literal chunk from end of a byte slice, even with an offset",
			p:        testPacket("8e34", 3),
			bitIndex: 6,
			numBits:  5,
			want:     0b0000_1101, // 0 = stop, val = 13
		},
		{
			name:     "read from the start of a one-byte slice",
			p:        testPacket("7c", 0),
			bitIndex: 0,
			numBits:  3,
			want:     0b_0000_0011,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			log := zaptest.NewLogger(t).Sugar()
			tc.p.SetLogger(log)

			b, err := tc.p.nBits(tc.bitIndex, tc.numBits)
			require.NoError(t, err)
			assert.Equal(t, tc.want, b)
		})
	}
}

func testPacket(hex string, bitIndex int) Packet {
	var p Packet
	p.UnmarshalText([]byte(hex))
	byteIndex := bitIndex / 8
	first := bitIndex % 8
	p.buf = p.buf[byteIndex:]
	p.firstBit = first
	return p
}
