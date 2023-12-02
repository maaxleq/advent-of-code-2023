package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
)

const inputFile = "input.txt"

// digitMap maps spelled out digits to their numeric equivalents.
var digitMap = map[string]int{
	"one":   1,
	"two":   2,
	"three": 3,
	"four":  4,
	"five":  5,
	"six":   6,
	"seven": 7,
	"eight": 8,
	"nine":  9,
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

// findDigits finds and returns the first and last digit (numerical or spelled out) in a string.
func findFirstLastDigits(line string) ([2]int, error) {
	var matches []string

	for i := 0; i < len(line); i++ {
		// Check for each digit word
		for word := range digitMap {
			if strings.HasPrefix(line[i:], word) {
				matches = append(matches, word)
				break
			}
		}

		// Check if the character is a digit
		if len(line) > i && line[i] >= '0' && line[i] <= '9' {
			matches = append(matches, line[i:i+1])
		}
	}

	if len(matches) < 1 {
		return [2]int{}, fmt.Errorf("not enough digits found")
	}

	firstDigit, errFirst := getDigit(matches[0])
	if errFirst != nil {
		return [2]int{}, errFirst
	}

	lastDigit, errLast := getDigit(matches[len(matches)-1])
	if errLast != nil {
		return [2]int{}, errLast
	}

	return [2]int{firstDigit, lastDigit}, nil
}

// getDigit converts a spelled-out digit or a single digit to its numerical equivalent.
func getDigit(s string) (int, error) {
	if val, exists := digitMap[s]; exists {
		return val, nil
	}

	if len(s) == 1 && s[0] >= '0' && s[0] <= '9' {
		return int(s[0] - '0'), nil
	}

	return 0, fmt.Errorf("invalid digit: %s", s)
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
		digits, errDigits := findFirstLastDigits(line)
		if errDigits != nil {
			log.Fatal(errDigits)
		}

		calibrationValue := computeCalibrationValue(digits)
		sum += calibrationValue
	}

	fmt.Println(sum)
}
