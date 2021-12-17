package main

import (
	"encoding"
	"errors"
	"fmt"
	"io"
	"log"
	"strconv"
)

// Packet is a compact slice of bytes, with a predefined schema.
// All packets have a header, of varying length
// (either 22 bits or 18 bits, depending on the length type)
//
// The first three bits [0:3] are the version number of the Packet.
// The next three bits [3:6] are the type of the Packet.
// The 7th bit (index = 6) defines the length type of the Packet as follows:
//
// - If the length type is 0, then the following 15 bits define the total length
// in bits of the sub-packets contained by this Packet.  In this case, the
// header of this Packet will be 22 bits long.
//
// - If the 7th bit is a 1, then the following 11 bits define the number
// of sub-packets immediately contained by this Packet.  In this case, the
// header of this Packet will be 18 bits long.
type Packet struct {
	buf []byte
	// firstbit is an offset from 0, defining where this packet begins
	// within the bufffer.
	firstBit int
}

// compile-time interface checks
var (
	_ encoding.TextUnmarshaler = &Packet{}
	_ encoding.TextMarshaler   = Packet{}
	_ fmt.Formatter            = Packet{}
)

// Version number for this packet.
func (p Packet) Version() int {
	const (
		bitIndex int  = 0
		numBits  byte = 3
	)
	bits, err := p.nBits(bitIndex, numBits)
	if err != nil {
		log.Fatal(err)
	}

	return int(bits)
}

// Value returns the cumulative value of this packet including its children.
func (p Packet) Value() int {
	switch p.packetType() {
	case _literalValue:
		const (
			first_index   = 6    // the starting bit index for data
			mask_val      = 0x0F // read the lower 4 bits of the value
			mask_continue = 0x10 // check the 5th least significant bit to see if we'll continue
		)

		var sum int
		i := first_index
		for {
			data, err := p.nBits(i, 5)
			if err != nil {
				log.Fatal(err)
			}

			val := int(data & mask_val)
			sum = sum<<4 + val

			if data&mask_continue == 0 {
				break
			}
			i += 5
		}
		return sum

	default:
		log.Fatal(errors.New("operator values not implemented"))
	}

	return 0
}

// Children returns a slice of the immediate children of this packet.
func (p Packet) Children() []Packet {
	if p.packetType() == _literalValue {
		return nil
	}

	sizeUpper, err := p.eightBits(7)
	if err != nil {
		log.Fatal(err)
	}
	sizeLower, err := p.eightBits(15)
	if err != nil {
		log.Fatal(err)
	}

	if p.lengthType() == _lengthInBits {
		// bits[7:23) hold the subPacketTotal in bits
		subPacketTotal := int(sizeUpper)<<8 + int(sizeLower)
		const bitIndex = 23
		var (
			byteIndex = (p.firstBit + bitIndex) / 8
			startBit  = (p.firstBit + bitIndex) % 8
			children  = make([]Packet, 0, 2)
		)
		var sum int
		for sum < subPacketTotal {
			child := Packet{
				buf:      p.buf[byteIndex:],
				firstBit: startBit,
			}
			children = append(children, child)
			sum += child.bitLength()
		}
		return children
	}

	// bits[7:18] hold the number of child packets
	numChildren := int(sizeUpper)<<8 + int(sizeLower>>5)
	children := make([]Packet, 0, numChildren)
	for n := 0; n < numChildren; n++ {
		child := Packet{
			buf:      p.buf,
			firstBit: p.firstBit + 18,
		}
		children = append(children, child)
	}
	return children
}

// Format implements fmt.Formatter. Ignores verbs and state.
func (p Packet) Format(f fmt.State, verb rune) {
	text, _ := p.MarshalText()
	f.Write([]byte(text))
}

// MarshalText implements encoding.TextMarshaler.  Always returns a nil error.
func (p Packet) MarshalText() ([]byte, error) {
	// each byte of data requires two 'hexits'
	text := make([]byte, 0, 2*len(p.buf))

	for _, b := range p.buf {
		hexits := strconv.FormatUint(uint64(b), 16)
		if len(hexits) == 1 {
			text = append(text, '0')
		}
		text = append(text, []byte(hexits)...)
	}
	return text, nil
}

// UnmarshalText implements encoding.TextUnmarshaler, and should be used to
// initialize a packet from its text representation.
func (p *Packet) UnmarshalText(hex []byte) error {
	(*p).buf = make([]byte, 0, 64)
	for i := 0; i+1 < len(hex); i += 2 {
		// each pair of 'hexits' is one byte of data
		val, err := strconv.ParseUint(string(hex[i:i+2]), 16, 8)
		if err != nil {
			return err
		}
		(*p).buf = append(p.buf, byte(val))
	}
	return nil
}

// packetType determines the type of a packet
type packetType byte

const (
	_literalValue packetType = 4
)

// packetType returns the type of packet that this is (value or an operator)
func (p Packet) packetType() packetType {
	const (
		bitIndex int  = 3
		numBits  byte = 3
	)
	bits, err := p.nBits(bitIndex, numBits)
	if err != nil {
		log.Fatal(err)
	}

	return packetType(bits)
}

// bitLength is the size in bits of this raw packet including its children
func (p Packet) bitLength() int {
	var length int
	if p.packetType() == _literalValue {
		length = 6
		const chunk_size = 5
		for {
			data, err := p.nBits(length, 1)
			if err != nil {
				log.Fatal(err)
			}
			length += chunk_size

			const mask_continue = 0x01
			if data&mask_continue == 0 {
				return length
			}
		}
	}

	length = 7

	sizeUpper, err := p.eightBits(7)
	if err != nil {
		log.Fatal(err)
	}
	sizeLower, err := p.eightBits(15)
	if err != nil {
		log.Fatal(err)
	}

	if p.lengthType() == _lengthInBits {
		// bits[7:23) hold the subPacketTotal in bits
		return length + int(sizeUpper)<<8 + int(sizeLower)
	}

	// bits[7:18] hold the number of child packets
	numChildren := int(sizeUpper)<<8 + int(sizeLower>>5)
	for n := 0; n < numChildren; n++ {
		child := &Packet{
			buf:      p.buf,
			firstBit: p.firstBit + 18,
		}
		length += child.bitLength()
	}
	return length
}

type lengthType byte

const (
	_lengthInBits    lengthType = 0 // operator, 0
	_lengthInPackets lengthType = 1 // operator, 1
	_lengthDynamic   lengthType = 2 // value, 2
)

// packetType determines which method this packet uses to indicate how long
// its contents are.
func (p Packet) lengthType() lengthType {
	if p.packetType() == _literalValue {
		return _lengthDynamic
	}

	const (
		bitIndex int  = 6
		numBits  byte = 1
	)
	b, err := p.nBits(bitIndex, numBits)
	if err != nil {
		log.Fatal(err)
	}
	return lengthType(b)
}

// eightBits returns the eight bits of data that start at the given logical
// bitIndex (if the raw packet were indexed by bit, rather than byte).
func (p Packet) eightBits(bitIndex int) (byte, error) {
	bitIndex += p.firstBit
	i, rem := bitIndex/8, bitIndex%8
	if i >= len(p.buf) || i+1 == len(p.buf) && rem != 0 {
		return 0, io.EOF
	}

	bits := p.buf[i]

	if rem != 0 {
		next := p.buf[i+1]
		bits = bits << rem
		next = next >> (8 - rem)
		bits = bits | next
	}

	return bits, nil
}

// nBits reads between 1 and 7 bits of data from this packet, beginning at
// index i.  The bits will fill the least significant bits of the returned byte,
// and the other bits will be zero'd out.  Returns io.EOF if the underlying
// buffer is not long enough to read these bits.
func (p Packet) nBits(bitIndex int, numBits byte) (byte, error) {
	bits, err := p.eightBits(bitIndex)
	if err != nil {
		return 0, err
	}

	shift := 8 - numBits
	mask := byte(2<<numBits - 1)

	return bits >> shift & mask, nil
}
