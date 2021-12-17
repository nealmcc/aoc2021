package main

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

// rawPacket is a compact slice of bytes, with a predefined schema.
// All packets have a header, of varying length
// (either 22 bits or 18 bits, depending on the length type)
//
// The first three bits [0:3] are the version number of the rawPacket.
// The next three bits [3:6] are the type of the rawPacket.
// The 7th bit (index = 6) defines the length type of the Packet as follows:
//
// - If the length type is 0, then the following 15 bits define the total length
// in bits of the sub-packets contained by this Packet.  In this case, the
// header of this Packet will be 22 bits long.
//
// - If the 7th bit is a 1, then the following 11 bits define the number
// of sub-packets immediately contained by this Packet.  In this case, the
// header of this Packet will be 18 bits long.
type rawPacket []byte

// compile-time interface checks
var (
	_ Packet                   = rawPacket(nil)
	_ encoding.TextUnmarshaler = &rawPacket{}
	_ encoding.TextMarshaler   = rawPacket(nil)
	_ fmt.Formatter            = rawPacket(nil)
)

// UnmarshalText implements encoding.TextUnmarshaler
func (p *rawPacket) UnmarshalText(hex []byte) error {
	*p = make([]byte, 0, 64)
	for i := 0; i+1 < len(hex); i += 2 {
		// each pair of 'hexits' is one byte of data
		val, err := strconv.ParseUint(string(hex[i:i+2]), 16, 8)
		if err != nil {
			return err
		}
		*p = append(*p, byte(val))
	}
	return nil
}

// MarshalText implements encoding.TextMarshaler
func (p rawPacket) MarshalText() ([]byte, error) {
	// each byte of data requires two 'hexits'
	text := make([]byte, 0, 2*len(p))

	for _, b := range p {
		hexits := strconv.FormatUint(uint64(b), 16)
		if len(hexits) == 1 {
			text = append(text, '0')
		}
		text = append(text, []byte(hexits)...)
	}
	return text, nil
}

// Format implements fmt.Formatter
func (p rawPacket) Format(f fmt.State, verb rune) {
	text, _ := p.MarshalText()
	f.Write([]byte(text))
}

// Version implements Packet.Version()
func (p rawPacket) Version() int {
	const (
		bitIndex int  = 0
		numBits  byte = 3
		shift    byte = 8 - numBits
		mask     byte = 2<<numBits - 1
	)

	bits, err := p.bits(bitIndex)
	if err != nil {
		log.Fatal(err)
	}

	return int(bits >> shift & mask)
}

// Value implements Packet.Value()
func (p rawPacket) Value() int {
	switch p.packetType() {
	case _literalValue:
		const (
			first_index   = 6    // the starting bit index for data
			shift         = 3    // shift right after reading 5 bits into one byte
			mask_val      = 0x0F // read the lower 4 bits of the value
			mask_continue = 0x10 // check the 5th least significant bit to see if we'll continue
		)

		var sum int
		i := first_index
		for {
			data, err := p.bits(i)
			if err != nil {
				log.Fatal(err)
			}

			val := int(data >> shift & mask_val)
			sum = sum<<4 + val

			if data&mask_continue == 0 {
				break
			}
			i += 5
		}
		return sum

	default:
		panic(errors.New("not implemented"))
	}
}

// Children implements Packet.Children()
func (p rawPacket) Children() []Packet {
	if p.packetType() == _literalValue {
		return nil
	}

	lenType := p.lengthType()
	b1, err := p.bits(7)
	if err != nil {
		log.Fatal(err)
	}
	b2, err := p.bits(15)
	if err != nil {
		log.Fatal(err)
	}

	switch lenType {
	case _lengthInBits:
		// bits[7:22] hold the size in bits
		size := int(b1)<<8 + int(b2)
		fmt.Println("fixed internal size", size)
		return nil

	case _lengthInPackets:
		// bits[7:18] hold the number of packets
		b2 := b2 >> 5
		numPackets := int(b1)<<8 + int(b2)
		fmt.Println(numPackets, " packet(s) inside this one")
		return nil

	default:
		log.Fatal(errors.New("invalid length type"))
	}

	return nil
}

// packetType determines the type of a packet
type packetType byte

const (
	_literalValue packetType = 4
)

// packetType returns the type of packet that this is (value or an operator)
func (p rawPacket) packetType() packetType {
	const (
		bitIndex int  = 3
		numBits  byte = 3
		shift    byte = 8 - numBits
		mask     byte = 2<<numBits - 1
	)

	bits, err := p.bits(bitIndex)
	if err != nil {
		log.Fatal(err)
	}

	return packetType(bits >> shift & mask)
}

// packetType determines which method this packet uses to indicate how long
// its contents are.
func (p rawPacket) lengthType() lengthType {
	const (
		bitIndex int  = 6
		numBits  byte = 1
		shift    byte = 8 - numBits
		mask     byte = 2<<numBits - 1
	)

	bits, err := p.bits(bitIndex)
	if err != nil {
		log.Fatal(err)
	}
	return lengthType(bits >> shift & mask)
}

// these bitmasks are used when evaluating the header of a raw packet
const (
	MASK_LENGTH_BITS  uint32 = 0x_1F_FF_C0_00 // 15 bits [7:22]
	SHIFT_LENGTH_BITS int    = 10             // shift the bit length this many bits after masking 4 bytes

	MASK_LENGTH_PACKETS  uint32 = 0x_1F_FC_00_00 // 11 bits [7:18]
	SHIFT_LENGTH_PACKETS int    = 14             // shift the packet length this many bits after masking 4 bytes
)

// lengthType is an enum which determines how to find the size of a packet
type lengthType byte

const (
	_lengthInBits    lengthType = iota // operator packets, which define the total size of their child packets in bits [8:24]
	_lengthInPackets                   // operator packets which define the number of child packets in bits [8:19]
	_lengthDynamic                     // value packets have a dynamic length
)

// bits returns the eight bits of data that start at bitIndex
// (if the raw packet were indexed by bit, rather than byte)
func (p rawPacket) bits(bitIndex int) (byte, error) {
	i, rem := bitIndex/8, bitIndex%8
	if i >= len(p) || i+1 == len(p) && rem != 0 {
		return 0, io.EOF
	}

	bits := p[i]

	if rem != 0 {
		next := p[i+1]
		bits = bits << rem
		next = next >> (8 - rem)
		bits = bits | next
	}

	return bits, nil
}
