package main

import (
	"bytes"
	"fmt"
	"io"
	logPkg "log"
	"os"
	"time"

	"go.uber.org/zap"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	logNoSugar, err := zap.NewDevelopment(zap.AddCaller(), zap.AddCallerSkip(1))
	if err != nil {
		logPkg.Fatal(err)
	}

	log := logNoSugar.Sugar()

	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	start := time.Now()

	packet, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1, err := part1(*packet, nil /* optional logger for debugging */)
	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()

	fmt.Println("part1:", p1)
	fmt.Printf("time taken: %s\n", end.Sub(start))
}

// read a packet from the given reader
func read(r io.Reader) (*Packet, error) {
	buf := bytes.Buffer{}
	_, err := buf.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	p := new(Packet)
	err = p.UnmarshalText(buf.Bytes())
	return p, err
}

func part1(p Packet, log *zap.SugaredLogger) (int, error) {
	s := PacketStack{}
	if log != nil {
		p.SetLogger(log)
	}
	s.Push(p)

	sum := 0
	for s.Length() > 0 {
		p = s.Pop()

		ver, err := p.Version()
		if err != nil {
			return 0, err
		}

		sum += ver

		children, err := p.Children()
		if err != nil {
			return sum, err
		}

		for _, child := range children {
			s.Push(child)
		}
	}

	return sum, nil
}
