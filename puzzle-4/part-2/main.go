package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

const inputFile = "input.txt"

// evaluateCard takes a line number and a card line and returns an array containing the card copy numbers won
func evaluateCard(lineNum int, line string) ([]int, error) {
	splitColon := strings.Split(line, ":")
	if len(splitColon) != 2 {
		return nil, fmt.Errorf("invalid format: %s", line)
	}

	splitBar := strings.Split(splitColon[1], "|")
	if len(splitBar) != 2 {
		return nil, fmt.Errorf("invalid format: %s", line)
	}

	var winningNums, ourNums []int

	for _, s := range strings.Fields(splitBar[0]) {
		n, errConv := strconv.Atoi(s)
		if errConv != nil {
			return nil, fmt.Errorf("invalid number: %s", s)
		}

		winningNums = append(winningNums, n)
	}

	for _, s := range strings.Fields(splitBar[1]) {
		n, errConv := strconv.Atoi(s)
		if errConv != nil {
			return nil, fmt.Errorf("invalid number: %s", s)
		}

		ourNums = append(ourNums, n)
	}

	i := 1
	var copiesWon []int
	for _, n := range ourNums {
		if slices.Contains(winningNums, n) {
			copiesWon = append(copiesWon, lineNum+i)
			i++
		}
	}

	return copiesWon, nil
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

	exemplarsCount := make(map[int]int)
	for i := 0; i < len(lines); i++ {
		exemplarsCount[i+1] = 1
	}

	for i, line := range lines {
		copiesWon, errEval := evaluateCard(i+1, line)
		if errEval != nil {
			log.Fatal(errEval)
		}

		for _, copy := range copiesWon {
			exemplarsCount[copy] += exemplarsCount[i+1]
		}
	}

	sum := 0

	for _, count := range exemplarsCount {
		sum += count
	}

	fmt.Println(sum)
}
