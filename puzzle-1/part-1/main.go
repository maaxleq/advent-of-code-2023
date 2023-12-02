package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

const inputFile = "input.txt"

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

func findFirstLastDigits(line string) [2]int {
	runes := []rune(line)

	var firstDigit int
	for i := 0; i < len(runes); i++ {
		if digit, err := strconv.Atoi(line[i : i+1]); err == nil {
			firstDigit = digit
			break
		}
	}

	var lastDigit int
	for i := len(runes) - 1; i >= 0; i-- {
		if digit, err := strconv.Atoi(line[i : i+1]); err == nil {
			lastDigit = digit
			break
		}
	}

	return [2]int{firstDigit, lastDigit}
}

func computeCalibrationValue(digits [2]int) int {
	return digits[0]*10 + digits[1]
}

func main() {
	lines, errRead := readFileLines(inputFile)
	if errRead != nil {
		log.Fatal(errRead)
	}

	sum := 0
	for _, line := range lines {
		calibrationValue := computeCalibrationValue(findFirstLastDigits(line))
		sum += calibrationValue
	}

	fmt.Println(sum)
}
