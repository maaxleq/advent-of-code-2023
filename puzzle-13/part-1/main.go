package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const inputFile = "input.txt"

type pattern [][]rune

func parsePatterns(lines []string) []pattern {
	currentPattern := pattern{}
	patterns := []pattern{}

	for _, line := range lines {
		if line == "" {
			patterns = append(patterns, currentPattern)
		}

		currentPattern = append(currentPattern, []rune(line))
	}

	if len(currentPattern) != 0 {
		patterns = append(patterns, currentPattern)
	}

	return patterns
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

	fmt.Println(len(parsePatterns(lines)))
}
