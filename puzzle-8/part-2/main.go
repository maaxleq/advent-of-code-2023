package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
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

// getStartingNodes extracts and returns all nodes from the network map that end with the character 'A'.
func getStartingNodes(network map[node]crossing) []node {
	nodes := []node{}

	for node := range network {
		if strings.HasSuffix(string(node), "A") {
			nodes = append(nodes, node)
		}
	}

	return nodes
}

// gcd computes the Greatest Common Divisor using the Euclidean algorithm
func gcd(a, b uint) uint {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

// lcm computes the Least Common Multiple of two numbers
func lcm(a, b uint) uint {
	return a / gcd(a, b) * b
}

// lcmSlice computes the LCM of a slice of uint
func lcmSlice(numbers []uint) uint {
	if len(numbers) == 0 {
		return 0 // No LCM for empty slice
	}

	result := numbers[0]
	for _, number := range numbers[1:] {
		result = lcm(result, number)
	}
	return result
}

func main() {
	lines, errRead := readFileLines(inputFile)
	if errRead != nil {
		log.Fatal(errRead)
	}

	instr, network, errParse := parseLines(lines)
	if errParse != nil {
		log.Fatal(errParse)
	}

	startingNodes := getStartingNodes(network)
	steps := make(chan uint, len(startingNodes))

	wg := sync.WaitGroup{}
	for _, n := range startingNodes {
		wg.Add(1)
		go func(n node, ch chan uint, instr instructions) {
			defer wg.Done()
			currentNode := n
			var stepCount uint = 0
			for !strings.HasSuffix(string(currentNode), "Z") {
				stepCount++
				direction := instr.next()
				crossing := network[currentNode]
				if direction == 'L' {
					currentNode = crossing.left
				} else if direction == 'R' {
					currentNode = crossing.right
				}
			}

			ch <- stepCount
		}(n, steps, instr)
	}
	wg.Wait()
	close(steps)

	stepCounts := []uint{}
	for stepCount := range steps {
		stepCounts = append(stepCounts, stepCount)
	}

	result := lcmSlice(stepCounts)

	fmt.Println(result)
}
