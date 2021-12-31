package alu

import (
	"bytes"
	"fmt"
	"io"
	"strconv"

	"github.com/pkg/errors"
)

// Compile interprets the given program code, checks it for syntax errors,
// and generates a sequence of instructions for the ALU to execute when it runs.
func Compile(code []byte) (Program, error) {
	buf := bytes.NewBuffer(code)
	instructions := make([]instruction, 0, 128)

	for line := 1; true; line++ {
		var (
			text []byte
			inst instruction
			err  error
		)
		inst.line = line

		if text, err = buf.ReadBytes('\n'); err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.WithMessagef(err, "line %d", line)
		}

		// make it so I can comment out irrelevant code, and add blank lines
		if text[0] == '#' || text[0] == '\n' {
			continue
		}

		if inst.r1, err = parseR1(text[4]); err != nil {
			return nil, errors.WithMessagef(err, "line %d", line)
		}

		if inst.p2, err = parseP2(text[5:]); err != nil {
			return nil, errors.WithMessagef(err, "line %d", line)
		}

		if inst.fn, err = chooseOperation(text[:3], inst.p2); err != nil {
			return nil, errors.WithMessagef(err, "line %d", line)
		}

		instructions = append(instructions, inst)
	}

	return instructions, nil
}

// opcode is the identifier for an arithmetic operation
type opcode string

const (
	opInput    opcode = "inp"
	opAdd      opcode = "add"
	opMultiply opcode = "mul"
	opDivide   opcode = "div"
	opModulo   opcode = "mod"
	opEquals   opcode = "eql"
)

func chooseOperation(code []byte, p2 interface{}) (operation, error) {
	_, useRegister := p2.(registerID)

	switch opcode(code) {
	case opInput:
		return (*ALU).input, nil

	case opAdd:
		if useRegister {
			return (*ALU).addRegister, nil
		}
		return (*ALU).addImmediate, nil

	case opMultiply:
		if useRegister {
			return (*ALU).mulRegister, nil
		}
		return (*ALU).mulImmediate, nil

	case opDivide:
		if useRegister {
			return (*ALU).divRegister, nil
		}
		if n := p2.(int); n == 0 {
			return nil, errors.New("divide by 0")
		}
		return (*ALU).divImmediate, nil

	case opModulo:
		if useRegister {
			return (*ALU).modRegister, nil
		}
		if n := p2.(int); n == 0 {
			return nil, errors.New("divide by 0")
		}
		return (*ALU).modImmediate, nil

	case opEquals:
		if useRegister {
			return (*ALU).eqRegister, nil
		}
		return (*ALU).eqImmediate, nil

	default:
		return nil, fmt.Errorf("unrecognised opcode %q", code)
	}
}

// parseR1 checks to ensure that b is a register ID,
// and returns an error if it isn't.
func parseR1(b byte) (registerID, error) {
	if !isRegisterID(b) {
		return 0, errors.New("not a register ID")
	}
	return registerID(b), nil
}

// parseP2 parses the second parameter to an instruction, returning either
// a registerID or an integer.  If the value cannot be parsed as either of
// these, then an error is returned.
func parseP2(b []byte) (interface{}, error) {
	b = bytes.TrimSpace(b)
	if len(b) == 0 {
		return nil, nil
	}

	if len(b) == 1 && isRegisterID(b[0]) {
		return registerID(b[0]), nil
	}

	n, err := strconv.Atoi(string(b))
	if err != nil {
		return nil, errors.Wrap(err, "cannot parse second parameter")
	}
	return n, nil
}

// isRegisterID ensures the given byte is a valid register ID. Used during compilation.
func isRegisterID(b byte) bool {
	switch registerID(b) {
	case regW, regX, regY, regZ:
		return true
	default:
		return false
	}
}
