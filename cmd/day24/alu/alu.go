// ALU is an arithmetic logic unit with 4 registers.
package alu

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strconv"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ALU is an Arithmetic Logic Unit
type ALU struct {
	mu   sync.Mutex
	in   *bufio.Scanner
	reg  [4]int
	code []instruction
}

// instruction is a command for the ALU to perform.
type instruction struct {
	fn operation
	r1 registerID
	p2 interface{}
}

type operation func(*ALU, registerID, ...interface{}) error

// New initializes a new ALU with its program.
// Returns an error if the program syntax is invalid.
func New(program []byte) (*ALU, error) {
	a := new(ALU)
	var err error
	a.code, err = compile(program)
	return a, err
}

// Run executes this ALU's program, reading its input from r.
// An ALU will only execute one Run() command at a time.
func (a *ALU) Run(r io.Reader) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.in = bufio.NewScanner(r)
	a.in.Split(bufio.ScanBytes)

	for i, inst := range a.code {
		err := inst.fn(a, inst.r1, inst.p2)
		if err != nil {
			return 1, errors.Wrapf(err, "execution failed on line %d", i)
		}
	}

	return a.get(regZ), nil
}

// get returns the value of the given register.
func (a *ALU) get(reg registerID) int {
	return a.reg[reg-regW]
}

// set the given register to the given value.
func (a *ALU) set(reg registerID, n int) {
	a.reg[reg-regW] = n
}

// compile interprets the given program code, checks it for syntax errors,
// and generates a sequence of instructions for the ALU to execute when it runs.
func compile(code []byte) ([]instruction, error) {
	buf := bytes.NewBuffer(code)
	instructions := make([]instruction, 0, 128)

	for line := 1; true; line++ {
		var (
			text []byte
			inst instruction
			err  error
		)
		if text, err = buf.ReadBytes('\n'); err != nil {
			if err == io.EOF {
				break
			}
			return nil, errors.WithMessagef(err, "line %d", line)
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

// input reads a value from the ALU's input and saves it in the given register.
func (a *ALU) input(r1 registerID, _ ...interface{}) error {
	if ok := a.in.Scan(); !ok {
		return errors.New("input requires an input value")
	}
	n := a.in.Bytes()[0]
	a.set(r1, int(n))
	return nil
}

// addRegister adds the value stored in r2 to r1. Always returns a nil error.
func (a *ALU) addRegister(r1 registerID, args ...interface{}) error {
	r2 := args[0].(registerID)

	sum := a.get(r1) + a.get(r2)
	a.set(r1, sum)
	return nil
}

// addImmediate adds n to r1. Always returns a nil error.
func (a *ALU) addImmediate(r1 registerID, args ...interface{}) error {
	n := args[0].(int)

	sum := a.get(r1) + n
	a.set(r1, sum)
	return nil
}

// mulRegister multiples r1 by the value stored in r2. Always returns a nil error.
func (a *ALU) mulRegister(r1 registerID, args ...interface{}) error {
	r2 := args[0].(registerID)

	prod := a.get(r1) * a.get(r2)
	a.set(r1, prod)
	return nil
}

// mulImmediate multiplies r1 by n. Always returns a nil error.
func (a *ALU) mulImmediate(r1 registerID, args ...interface{}) error {
	n := args[0].(int)

	sum := a.get(r1) + n
	a.set(r1, sum)
	return nil
}

// divRegister divides r1 by the value stored in r2.
func (a *ALU) divRegister(r1 registerID, args ...interface{}) error {
	r2 := args[0].(registerID)
	n := a.get(r2)
	if n == 0 {
		return errors.New("divide by 0")
	}

	quo := a.get(r1) / n
	a.set(r1, quo)
	return nil
}

// divImmediate divides r1 by n.
func (a *ALU) divImmediate(r1 registerID, args ...interface{}) error {
	n := args[0].(int)
	if n == 0 {
		return errors.New("divide by 0")
	}

	quo := a.get(r1) / n
	a.set(r1, quo)
	return nil
}

// modRegister stores r1 % r2 in r1.
func (a *ALU) modRegister(r1 registerID, args ...interface{}) error {
	r2 := args[0].(registerID)
	n := a.get(r2)
	if n == 0 {
		return errors.New("modulo of 0 is undefined")
	}

	mod := a.get(r1) % n
	a.set(r1, mod)
	return nil
}

// modImmediate stores r1 % n in r1. Returns an error at initialisation time if n is 0
func (a *ALU) modImmediate(r1 registerID, args ...interface{}) error {
	n := args[0].(int)
	if n == 0 {
		return errors.New("modulo of 0 is undefined")
	}

	mod := a.get(r1) % n
	a.set(r1, mod)
	return nil
}

// eqRegister stores 1 or 0 in r1 depending on if the value in
// r1 equals the value in r2. (1 for true, 0 for false)
func (a *ALU) eqRegister(r1 registerID, args ...interface{}) error {
	r2 := args[0].(registerID)

	if a.get(r1) == a.get(r2) {
		a.set(r1, 1)
		return nil
	}

	a.set(r1, 0)
	return nil
}

// eqImmediate stores 1 in r1 iff the value in r1 == n.
func (a *ALU) eqImmediate(r1 registerID, args ...interface{}) error {
	n := args[0].(int)
	var (
		left  = a.reg[r1-regW]
		right = n
		out   int
	)
	if left == right {
		out = 1
	}
	a.reg[r1-regW] = out

	return nil
}

func trace(logs []*zap.SugaredLogger, msg string, keysAndValues ...interface{}) {
	for _, l := range logs {
		l.Debugw(msg, keysAndValues...)
	}
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

// registerID is the identifier for each of the ALU's registers
type registerID byte

const (
	regW registerID = 'w'
	regX registerID = 'x'
	regY registerID = 'y'
	regZ registerID = 'z'
)

func isRegisterID(b byte) bool {
	switch registerID(b) {
	case regW, regX, regY, regZ:
		return true
	default:
		return false
	}
}
