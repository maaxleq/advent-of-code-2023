package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const inputFile = "input.txt"

func hashAlgorithm(s string) int {
	val := 0

	for _, r := range s {
		val += int(r)
		val *= 17
		val %= 256
	}

	return val
}

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
