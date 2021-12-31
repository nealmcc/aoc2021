// ALU is an arithmetic logic unit with 4 registers.
package alu

import (
	"bufio"
	"io"
	"reflect"
	"runtime"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// ALU is an Arithmetic Logic Unit
type ALU struct {
	mu   sync.Mutex
	in   *bufio.Scanner
	reg  [4]int
	code Program
	log  *zap.SugaredLogger
}

// Program is the sequence of instructions that the ALU will execute.
type Program []instruction

// instruction is a command for the ALU to perform.
type instruction struct {
	fn   operation
	r1   registerID
	p2   interface{}
	line int
}

type operation func(*ALU, registerID, interface{}) error

// registerID is the identifier for each of the ALU's registers
type registerID byte

const (
	regW registerID = 'w'
	regX registerID = 'x'
	regY registerID = 'y'
	regZ registerID = 'z'
)

// MarshalJSON implements json.Marshaller, so that register IDs will print
// nicely in traces.
func (reg registerID) MarshalJSON() ([]byte, error) {
	return []byte{'"', byte(reg), '"'}, nil
}

// New initializes a new ALU with the given pre-compiled program.
func New(code Program) *ALU {
	return &ALU{code: code}
}

// EnableTrace makes this ALU write trace output to the given logger.
func (a *ALU) EnableTrace(log *zap.SugaredLogger) {
	a.log = log
}

// Reset this ALU's registers to 0.
func (a *ALU) Reset() {
	for i := range a.reg {
		a.reg[i] = 0
	}
}

// Run executes this ALU's program, reading its input from r.
// An ALU will only execute one Run() command at a time.
func (a *ALU) Run(r io.Reader) (int, error) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.Reset()

	a.in = bufio.NewScanner(r)
	a.in.Split(bufio.ScanBytes)

	for i, inst := range a.code {
		if a.log != nil {
			a.log.Debugw("begin instruction",
				"line", inst.line,
				"seq", i+1,
				"fn", getFunctionName(inst.fn),
				"r1", inst.r1,
				"p2", inst.p2)
		}
		err := inst.fn(a, inst.r1, inst.p2)
		if err != nil {
			return 1, errors.Wrapf(err, "execution failed on line %d", inst.line)
		}
		if a.log != nil {
			a.log.Debugw("instruction complete", "line", inst.line, "registers", a.reg)
		}
	}

	return a.get(regZ), nil
}

// input reads a value from the ALU's input and saves it in the given register.
func (a *ALU) input(r1 registerID, _ interface{}) error {
	if ok := a.in.Scan(); !ok {
		return errors.New("input requires an input value")
	}
	n := a.in.Bytes()[0]
	a.set(r1, int(n))
	return nil
}

// addRegister adds the value stored in r2 to r1. Always returns a nil error.
func (a *ALU) addRegister(r1 registerID, p2 interface{}) error {
	r2 := p2.(registerID)

	sum := a.get(r1) + a.get(r2)
	a.set(r1, sum)
	return nil
}

// addImmediate adds n to r1. Always returns a nil error.
func (a *ALU) addImmediate(r1 registerID, p2 interface{}) error {
	n := p2.(int)

	sum := a.get(r1) + n
	a.set(r1, sum)
	return nil
}

// mulRegister multiples r1 by the value stored in r2. Always returns a nil error.
func (a *ALU) mulRegister(r1 registerID, p2 interface{}) error {
	r2 := p2.(registerID)

	prod := a.get(r1) * a.get(r2)
	a.set(r1, prod)
	return nil
}

// mulImmediate multiplies r1 by n. Always returns a nil error.
func (a *ALU) mulImmediate(r1 registerID, p2 interface{}) error {
	n := p2.(int)

	prod := a.get(r1) * n
	a.set(r1, prod)
	return nil
}

// divRegister divides r1 by the value stored in r2.
func (a *ALU) divRegister(r1 registerID, p2 interface{}) error {
	r2 := p2.(registerID)
	n := a.get(r2)
	if n == 0 {
		return errors.New("divide by 0")
	}

	quo := a.get(r1) / n
	a.set(r1, quo)
	return nil
}

// divImmediate divides r1 by n.
func (a *ALU) divImmediate(r1 registerID, p2 interface{}) error {
	// the compiler will prevent n == 0
	n := p2.(int)

	quo := a.get(r1) / n
	a.set(r1, quo)
	return nil
}

// modRegister stores r1 % r2 in r1.
func (a *ALU) modRegister(r1 registerID, p2 interface{}) error {
	r2 := p2.(registerID)
	n := a.get(r2)
	if n == 0 {
		return errors.New("modulo of 0 is undefined")
	}

	mod := a.get(r1) % n
	a.set(r1, mod)
	return nil
}

// modImmediate stores r1 % n in r1. Returns an error at initialisation time if n is 0
func (a *ALU) modImmediate(r1 registerID, p2 interface{}) error {
	// the compiler will prevent n == 0
	n := p2.(int)

	mod := a.get(r1) % n
	a.set(r1, mod)
	return nil
}

// eqRegister stores 1 or 0 in r1 depending on if the value in
// r1 equals the value in r2. (1 for true, 0 for false)
func (a *ALU) eqRegister(r1 registerID, p2 interface{}) error {
	r2 := p2.(registerID)

	if a.get(r1) == a.get(r2) {
		a.set(r1, 1)
		return nil
	}

	a.set(r1, 0)
	return nil
}

// eqImmediate stores 1 in r1 iff the value in r1 == n.
func (a *ALU) eqImmediate(r1 registerID, p2 interface{}) error {
	n := p2.(int)

	if a.get(r1) == n {
		a.set(r1, 1)
		return nil
	}

	a.set(r1, 0)
	return nil
}

// get returns the value of the given register.
func (a *ALU) get(reg registerID) int {
	return a.reg[reg-regW]
}

// set the given register to the given value.
func (a *ALU) set(reg registerID, n int) {
	a.reg[reg-regW] = n
}

// getFunctionName is used during tracing
func getFunctionName(i interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
	return strings.TrimPrefix(name, "github.com/nealmcc/aoc2021/pkg/alu.(*ALU).")
}
