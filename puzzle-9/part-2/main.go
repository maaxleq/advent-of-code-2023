package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const inputFile = "input.txt"

// composedOfZeros checks if all elements in the given slice are zero.
// It returns true if all elements are zero, and false otherwise.
func composedOfZeros(nums []int) bool {
	for _, num := range nums {
		if num != 0 {
			return false
		}
	}

	return true
}

// getDiffList calculates the difference between adjacent elements in the given slice.
// It returns a new slice containing these differences
func getDiffList(nums []int) []int {
	diffs := []int{}

	for i := range nums[1:] {
		diffs = append(diffs, nums[i+1]-nums[i])
	}

	return diffs
}

// getDiffSequences generates a sequence of differences from the provided slice.
// It creates a series of slices where each subsequent slice is the difference
// of its preceding slice, until a slice of all zeros is reached.
// It returns a slice of slices containing these difference sequences.
func getDiffSequences(nums []int) [][]int {
	firstLine := make([]int, len(nums))
	copy(firstLine, nums)

	res := [][]int{
		firstLine,
	}

	depth := 0
	for !composedOfZeros(res[depth]) {
		res = append(res, getDiffList(res[depth]))
		depth++
	}

	return res
}

// extrapolate uses the difference sequences to extrapolate the first number in the series.
// It calculates the next number based on the difference sequences generated from the input slice.
// Returns the extrapolated value.
func extrapolate(nums []int) int {
	diffSeqs := getDiffSequences(nums)

	for i := len(diffSeqs) - 1; i >= 0; i-- {
		if i == len(diffSeqs)-1 {
			diffSeqs[i] = append([]int{0}, diffSeqs[i]...)
			continue
		}

		a := diffSeqs[i+1][0]
		b := diffSeqs[i][0]
		diffSeqs[i] = append([]int{b - a}, diffSeqs[i]...)
	}

	return diffSeqs[0][0]
}

// parseData converts a slice of string lines into a slice of slices of integers.
// Each line is expected to contain space-separated integers.
// Returns the parsed data or an error if the parsing fails.
func parseData(lines []string) ([][]int, error) {
	data := [][]int{}

	for _, line := range lines {
		dataLine := []int{}

		fields := strings.Fields(line)
		for _, field := range fields {
			num, errConv := strconv.Atoi(field)
			if errConv != nil {
				return nil, fmt.Errorf("bad data line format: %s", line)
			}

			dataLine = append(dataLine, num)
		}

		data = append(data, dataLine)
	}

	return data, nil
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

	data, errParse := parseData(lines)
	if errParse != nil {
		log.Fatal(errParse)
	}

	sum := 0
	for _, dataLine := range data {
		sum += extrapolate(dataLine)
	}

	fmt.Println(sum)
}
