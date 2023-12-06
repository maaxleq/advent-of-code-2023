package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
	"strconv"
	"strings"
)

const inputFile = "input.txt"

// evaluateCard takes a card line and returns the amount of points won with that card
func evaluateCard(line string) (int, error) {
	splitColon := strings.Split(line, ":")
	if len(splitColon) != 2 {
		return 0, fmt.Errorf("invalid format: %s", line)
	}

	splitBar := strings.Split(splitColon[1], "|")
	if len(splitBar) != 2 {
		return 0, fmt.Errorf("invalid format: %s", line)
	}

	var winningNums, ourNums []int

	for _, s := range strings.Fields(splitBar[0]) {
		n, errConv := strconv.Atoi(s)
		if errConv != nil {
			return 0, fmt.Errorf("invalid number: %s", s)
		}

		winningNums = append(winningNums, n)
	}

	for _, s := range strings.Fields(splitBar[1]) {
		n, errConv := strconv.Atoi(s)
		if errConv != nil {
			return 0, fmt.Errorf("invalid number: %s", s)
		}

		ourNums = append(ourNums, n)
	}

	winCount := 0
	for _, n := range ourNums {
		if slices.Contains(winningNums, n) {
			winCount++
		}
	}

	if winCount == 0 {
		return 0, nil
	}

	return int(math.Pow(2, float64(winCount-1))), nil
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

	sum := 0
	for _, line := range lines {
		winCount, errEval := evaluateCard(line)
		if errEval != nil {
			log.Fatal(errEval)
		}

		sum += winCount
	}

	fmt.Println(sum)
}
