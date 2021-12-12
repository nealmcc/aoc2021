package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// main solves the both part 1 and part 2, reading from input.txt
func main() {
	in, err := os.Open("input.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer in.Close()

	caves, err := read(in)
	if err != nil {
		log.Fatal(err)
	}

	p1 := part1(caves)
	fmt.Println("part1", p1)
}

type node struct {
	id  string
	out map[string]*node
}

func read(r io.Reader) (map[string]*node, error) {
	s := bufio.NewScanner(r)

	caves := make(map[string]*node, 24)

	getOrCreate := func(id string) *node {
		n, ok := caves[id]
		if !ok {
			n = &node{id: id, out: make(map[string]*node, 2)}
			caves[id] = n
		}
		return n
	}

	for s.Scan() {
		parts := strings.Split(s.Text(), "-")
		left, right := getOrCreate(parts[0]), getOrCreate(parts[1])
		left.out[right.id] = right
		right.out[left.id] = left
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return caves, nil
}

// part1 solves part1 of the puzzle
func part1(caves map[string]*node) int {
	visited := map[string]bool{}
	paths := pathsFrom(caves, caves["start"], caves["end"], visited)
	return len(paths)
}

// pathsFrom lists all the possible paths through the cave system from
// start to end, such that small caves (lowercase) are visited at most once.
func pathsFrom(caves map[string]*node, start, end *node, visited map[string]bool) []string {
	if start.id == end.id {
		return []string{end.id}
	}

	canVisit := func(id string) bool {
		if id[0] < 'a' {
			return true
		}
		return !visited[id]
	}

	visited[start.id] = true
	paths := make([]string, 0, 2)
	for key, child := range start.out {
		if canVisit(key) {
			childPaths := pathsFrom(caves, child, end, copyMap(visited))
			for _, p := range childPaths {
				paths = append(paths, start.id+","+p)
			}
		}
	}

	return paths
}

func copyMap(m1 map[string]bool) map[string]bool {
	m2 := make(map[string]bool, len(m1))
	for key, val := range m1 {
		m2[key] = val
	}
	return m2
}
