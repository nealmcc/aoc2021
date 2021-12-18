package main

import (
	"encoding"
	"fmt"
	"io"
	"strconv"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// Packet is a compact slice of bytes, with a predefined schema.
// All packets have a header, of varying length
// (either 22 bits or 18 bits, depending on the length type)
//
// The first three bits [0:3] are the version number of the Packet.
// The next three bits [3:6] are the type of the Packet.
//
// The 7th bit (index = 6) is the sizeTypeID, and defines how to measure
// the inner size (size of the children) of the Packet as follows:
//
// - If the sizeTypeID is 0, then the following 15 bits define the total length
// in bits of the sub-packets contained by this Packet.  In this case, the
// header of this Packet will be 22 bits long.
//
// - If the sizeTypeID is a 1, then the following 11 bits define the number
// of sub-packets immediately contained by this Packet.  In this case, the
// header of this Packet will be 18 bits long.
type Packet struct {
	buf []byte
	// firstbit is an offset from 0, defining where this packet begins
	// within the bufffer.
	firstBit int

	log *zap.SugaredLogger
}

// compile-time interface checks
var (
	_ encoding.TextUnmarshaler = &Packet{}
	_ encoding.TextMarshaler   = Packet{}
	_ fmt.Formatter            = Packet{}
)

type (
	// fieldIndex specifies the starting index for each field.
	fieldIndex = int

	// fieldWidth specifies the number of bits used for each field.
	fieldWidth uint8
)

const (
	_versionIx   fieldIndex = 0
	_versionBits fieldWidth = 3

	_packetTypeIx   fieldIndex = 3
	_packetTypeBits fieldWidth = 3

	_sizeTypeIx   fieldIndex = 6
	_sizeTypeBits fieldWidth = 1

	_sizeIx             fieldIndex = 7
	_sizeImmediateUpper fieldWidth = 7
	_sizeImmediateLower fieldWidth = 8
	_sizePacketsUpper   fieldWidth = 3
	_sizePacketsLower   fieldWidth = 8

	_valueIx    fieldIndex = 6 // (only for the first chunk)
	_valueChunk fieldWidth = 5
)

// Version number for this packet.
func (p Packet) Version() (ver int, err error) {
	p.Debugw("enter p.Version()", "p", p, "firstBit", p.firstBit)
	defer func() {
		p.Debugw(" exit p.Version()", "ver", ver, "err", err)
	}()

	const bitIndex = 0
	bits, err := p.nBits(bitIndex, _versionBits)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return int(bits), nil
}

// Value returns the cumulative value of this packet including its children.
func (p Packet) Value() (val int, err error) {
	p.Debugw("enter p.Value()", "p", p, "firstBit", p.firstBit)
	defer func() {
		p.Debugw(" exit p.Value()", "val", val, "err", err)
	}()

	pt, err := p.packetType()
	if err != nil {
		return 0, err
	}

	switch pt {
	case _literalValue:
		const (
			mask_val      = 0x0F // read the lower 4 bits of the value
			mask_continue = 0x10 // check the 5th least significant bit to see if we'll continue
		)

		var sum int
		index := _valueIx
		for {
			data, err := p.nBits(index, _valueChunk)
			if err != nil {
				return 0, err
			}

			sum = sum<<4 + int(data&mask_val)

			index += int(_valueChunk)
			if data&mask_continue == 0 {
				break
			}
		}
		return sum, nil

	default:
		return 0, errors.New("operator values not implemented")
	}
}

// Children returns a slice of the immediate children of this packet.
// Returns an empty slice if this packet has no children.
func (p Packet) Children() ([]Packet, error) {
	p.Debugw("enter p.Children()", "p", p, "firstBit", p.firstBit)
	pt, err := p.packetType()
	if err != nil {
		return nil, err
	}
	if pt == _literalValue {
		p.Debugw("value packet", "pt", pt)
		return nil, nil
	}

	p.Debugw("operator packet", "pt", pt)
	children := make([]Packet, 0, 2)

	st, inner, next, err := p.innerSize()
	if err != nil {
		return nil, err
	}

	if st == _sizeInBits {
		p.Debugw("size in bits", "st", st, "inner", inner)
		// inner is the total number of bits that the children will occupy
		sum := 0
		for sum < inner {
			child := p.childFrom(next)
			p.Debugw("child found", "bitIndex", next,
				"child", child, "child.firstBit", child.firstBit)
			children = append(children, child)
			length, err := child.bitLength()
			if err != nil {
				return nil, err
			}
			next += length
			sum += length
		}
		return children, nil
	}

	p.Debugw("size in packets", "st", st, "inner", inner)
	// inner is the number of child packets
	return p.nSiblingsAt(next, inner)
}

func (p Packet) nSiblingsAt(bitIndex, n int) (children []Packet, err error) {
	p.Debugw("enter p.nSiblingsAt()", "p", p, "firstBit", p.firstBit,
		"bitIndex", bitIndex, "n", n)
	defer func() {
		p.Debugw(" exit p.nSiblingsAt()", "len(children)", len(children), "err", err)
	}()

	children = make([]Packet, 0, 3)
	for i := 0; i < n; i++ {
		child := p.childFrom(bitIndex)
		p.Debugw("\tchild found", "bitIndex", bitIndex,
			"child", child, "child.firstBit", child.firstBit)
		children = append(children, child)
		length, err := child.bitLength()
		if err != nil {
			return nil, err
		}
		bitIndex += length
	}

	return children, nil
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
func (p Packet) packetType() (pt packetType, err error) {
	p.Debugw("enter p.packetType()", "p", p, "firstBit", p.firstBit)
	defer func() {
		p.Debugw(" exit p.packetType()", "pt", pt, "err", err)
	}()

	const bitIndex = 3

	bits, err := p.nBits(bitIndex, _packetTypeBits)
	if err != nil {
		return 0, err
	}

	return packetType(bits), nil
}

// bitLength is the size in bits of this packet including its children, but
// excluding the initial bit offset if any, and also excluding any trailing bits.
func (p Packet) bitLength() (length int, err error) {
	p.Debugw("enter p.bitLength()", "p", p, "firstBit", p.firstBit)
	defer func() {
		p.Debugw(" exit p.bitLength()", "length", length, "err", err)
	}()

	pt, err := p.packetType()
	if err != nil {
		return 0, err
	}
	p.Debugw("packet type", "pt", pt)

	if pt == _literalValue {
		_, sz, err := p.literal()
		return sz, err
	}

	st, inner, next, err := p.innerSize()
	if err != nil {
		return 0, err
	}

	if st == _sizeInBits {
		// inner is the sum total of the size of the inner packets in bits
		return next + inner, nil
	}

	// inner is the number of child packets
	p.Debugw("summing child lengths", "numChildren", inner)
	children, err := p.nSiblingsAt(next, inner)
	if err != nil {
		return 0, err
	}

	var sum int
	for _, c := range children {
		l, err := c.bitLength()
		if err != nil {
			return 0, err
		}
		sum += l
	}

	return next + sum, nil
}

// literal returns the integer value of this Packet and the number of bits
// that are used to store it.
func (p Packet) literal() (val, width int, err error) {
	pt, err := p.packetType()
	if err != nil {
		return 0, 0, err
	}
	if pt != _literalValue {
		return 0, 0, errors.New("invalid operation")
	}

	const chunk_size = 5
	i, sum := 6, 0
	for {
		data, err := p.nBits(i, chunk_size)
		if err != nil {
			return 0, 0, err
		}
		i += chunk_size

		const maskVal = 0b_0000_1111
		sum = sum<<4 + int(data&maskVal)

		const maskContinue = 0b_0001_0000
		if data&maskContinue == 0 {
			return sum, i, nil
		}
	}
}

type sizeTypeID byte

const (
	_sizeInBits    sizeTypeID = 0 // inner bit length defined in the following 15 bits
	_sizeInPackets sizeTypeID = 1 // inner packet size defined in the following 11 bits
)

// innerSize returns the size type of the inner portion of this container packet,
// which could have two different meanings. Either it is the sum total
// of the bit lengths of its children, or it is the number of packets contained
// inside this one.  InnerSize() also returns the next index to read child bits from.
func (p Packet) innerSize() (st sizeTypeID, size int, next int, err error) {
	p.Debugw("enter p.innerSize()", "p", p, "firstBit", p.firstBit)
	defer func() {
		p.Debugw(" exit p.innerSize()", "st", st, "size", size, "next", next, "err", err)
	}()

	b, err := p.nBits(_sizeTypeIx, _sizeTypeBits)
	if err != nil {
		return 0, 0, 0, err
	}

	st = sizeTypeID(b)
	var upper, lower byte
	switch st {
	case _sizeInBits:
		upper, err = p.nBits(_sizeIx, _sizeImmediateUpper)
		if err != nil {
			return 0, 0, 0, err
		}
		lowerBitsIndex := _sizeIx + int(_sizeImmediateUpper)
		lower, err = p.nBits(lowerBitsIndex, _sizeImmediateLower)
		if err != nil {
			return 0, 0, 0, err
		}
		next = 22

	case _sizeInPackets:
		upper, err = p.nBits(_sizeIx, _sizePacketsUpper)
		if err != nil {
			return 0, 0, 0, err
		}
		lowerPacketsIndex := _sizeIx + int(_sizePacketsUpper)
		lower, err = p.nBits(lowerPacketsIndex, _sizePacketsLower)
		if err != nil {
			return 0, 0, 0, err
		}
		next = 18

	default:
		return 0, 0, 0, errors.New("unexpected size type")
	}

	size = int(upper)<<8 + int(lower)
	return st, size, next, nil
}

// nBits reads between 1 and 7 bits of data from this packet, beginning at
// bitIndex.  The bits will fill the least significant bits of the returned byte,
// and the other bits will be zero'd out.  Returns io.EOF if the underlying
// buffer is not long enough to read these bits.
func (p Packet) nBits(bitIndex fieldIndex, numBits fieldWidth) (res byte, err error) {
	// uncomment for more verbose tracing
	p.Debugw("enter p.nBits()", "p", p, "firstBit", p.firstBit,
		"bitIndex", bitIndex, "numBits", numBits)
	defer func() {
		p.Debugw(" exit p.nBits()", "res", res, "err", err)
	}()

	byteIndex := (p.firstBit + bitIndex) / 8
	rem := (p.firstBit + bitIndex) % 8
	if byteIndex >= len(p.buf) {
		return 0, io.EOF
	}

	var (
		left  uint8
		right uint8
	)

	left = p.buf[byteIndex]
	if rem == 0 {
		left = left >> (8 - numBits)
		return left, nil
	}

	if byteIndex+1 < len(p.buf) {
		right = p.buf[byteIndex+1]
	}

	left <<= rem
	right >>= (8 - rem)
	data := (left | right) >> (8 - numBits)

	// uncomment for more verbose tracing
	// if p.log != nil {
	// 	p.log.Debugf("shifted and masked 8 bits: %0.8b aka %0.2x", data, data)
	// }

	return data, nil
}

// childFrom re-frames the current packet at the given index, returning that
// as a new Packet, but re-using the underlying buffer.
func (p Packet) childFrom(bitIndex int) (child Packet) {
	p.Debugw("enter p.childFrom()", "p", p, "firstBit", p.firstBit,
		"bitIndex", bitIndex)
	defer func() {
		p.Debugw(" exit p.childFrom()", "child", child, "firstBit", child.firstBit)
	}()

	byteIndex := (p.firstBit + bitIndex) / 8
	firstBit := (p.firstBit + bitIndex) % 8
	return Packet{
		buf:      p.buf[byteIndex:],
		firstBit: firstBit,
		log:      p.log,
	}
}

// SetLogger assigns a logger for this packet to use for debugging
func (p *Packet) SetLogger(log *zap.SugaredLogger) {
	p.log = log.Desugar().
		WithOptions(zap.AddCaller(), zap.AddCallerSkip(1)).
		Sugar()
}

// DebugWith adds keys and values to the logger
func (p *Packet) DebugWith(keysAndValues ...interface{}) {
	if p.log == nil {
		return
	}
	p.log = p.log.With(keysAndValues...)
}

// Debugw logs a message with some additional context. The variadic key-value
// pairs are treated as they are in zap.SugaredLogger.With.
// If this packet has no logger configured, this is a no-op.
func (p Packet) Debugw(msg string, keysAndValues ...interface{}) {
	if p.log == nil {
		return
	}
	p.log.Debugw(msg, keysAndValues...)
}
