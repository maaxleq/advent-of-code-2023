package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const inputFile = "input.txt"

// node represents a single node in the network.
type node string

// crossing represents a crossing in the network with left and right nodes.
type crossing struct {
	left, right node
}

// instructions stores the directions to follow and the current cursor position.
type instructions struct {
	directions []rune
	cursor     uint
}

// next returns the current direction from the instructions and advances the cursor.
func (i *instructions) next() rune {
	direction := i.directions[i.cursor]
	i.cursor = (i.cursor + 1) % uint(len(i.directions))

	return direction
}

// parseLines takes an array of string lines and parses them into instructions and a network map.
// It returns a set of instructions, a map of nodes to their corresponding crossings, and an error if any.
func parseLines(lines []string) (instructions, map[node]crossing, error) {
	directions := []rune{}
	network := make(map[node]crossing)

	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)

		if trimmedLine == "" {
			continue
		}

		if strings.Contains(trimmedLine, "=") {
			splitEquals := strings.Split(trimmedLine, " = ")
			n := splitEquals[0]
			splitComma := strings.Split(splitEquals[1], ", ")
			left, cutLeft := strings.CutPrefix(splitComma[0], "(")
			right, cutRight := strings.CutSuffix(splitComma[1], ")")

			if !(cutLeft && cutRight) {
				return instructions{}, nil, fmt.Errorf("bad network line format: %s", line)
			}

			network[node(n)] = crossing{
				left:  node(left),
				right: node(right),
			}
		} else {
			directions = []rune(trimmedLine)
		}
	}

	return instructions{
		directions: directions,
		cursor:     0,
	}, network, nil
}

// readFileLines reads a file and returns its contents as an array of strings.
func readFileLines(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func main() {
	lines, errRead := readFileLines(inputFile)
	if errRead != nil {
		log.Fatal(errRead)
	}

	instructions, network, errParse := parseLines(lines)
	if errParse != nil {
		log.Fatal(errParse)
	}

	currentNode := node("AAA")
	steps := 0
	for currentNode != node("ZZZ") {
		steps++
		direction := instructions.next()
		crossing := network[currentNode]
		if direction == 'L' {
			currentNode = crossing.left
		} else if direction == 'R' {
			currentNode = crossing.right
		}
	}

	fmt.Println(steps)
}
