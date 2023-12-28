package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const inputFile = "input.txt"

// hashAlgorithm calculates a simple hash of a given string.
// It iterates through each character in the string, converting it to an integer
// and performing a series of calculations to produce a hash value.
// The function returns an integer representing the hash.
func hashAlgorithm(s string) int {
	val := 0

	for _, r := range s {
		val += int(r)
		val *= 17
		val %= 256
	}

	return val
}

// parseSteps splits a string by commas and returns a slice of the substrings.
// It is used to process a line of text into individual steps or commands.
func parseSteps(line string) []string {
	return strings.Split(line, ",")
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

	if len(lines) < 1 {
		log.Fatal(fmt.Errorf("no line in file"))
	}

	steps := parseSteps(lines[0])

	sum := 0
	for _, step := range steps {
		sum += hashAlgorithm(step)
	}

	fmt.Println(sum)
}
