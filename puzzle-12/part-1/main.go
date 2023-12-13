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

// countArrangements serves as a wrapper function for countArrangementsInner. It initializes
// the memoization map and calls the inner function to compute the arrangements.
func countArrangements(condition []rune, groups []int) int {
	return countArrangementsInner(condition, groups, make(map[string]int))
}

// Key function to generate unique keys for memoization.
func key(condition string, groups []int) string {
	return condition + ":" + fmt.Sprint(groups)
}

// countArrangementsInner recursively counts the number of valid arrangements of damaged springs.
// It uses memoization to optimize repeated state calculations.
func countArrangementsInner(condition []rune, groups []int, memoizationMap map[string]int) int {
	// Generate a unique key for the current state.
	memKey := key(string(condition), groups)

	// Check if result is already computed.
	if val, exists := memoizationMap[memKey]; exists {
		return val
	}

	// Base case: empty condition.
	if len(condition) == 0 {
		if len(groups) == 0 {
			return 1
		}
		return 0
	}

	firstChar := condition[0]
	var permutations int

	switch firstChar {
	case '.':
		// Operational spring, skip it.
		permutations = countArrangementsInner(condition[1:], groups, memoizationMap)
	case '?':
		// Unknown status, count both possibilities.
		permutations = countArrangementsInner(append([]rune{'.'}, condition[1:]...), groups, memoizationMap) +
			countArrangementsInner(append([]rune{'#'}, condition[1:]...), groups, memoizationMap)
	case '#':
		// Damaged spring.
		if len(groups) == 0 {
			permutations = 0
		} else {
			nrDamaged := groups[0]
			if nrDamaged <= len(condition) {
				valid := true
				for i := 0; i < nrDamaged; i++ {
					if condition[i] == '.' {
						valid = false
						break
					}
				}

				if valid {
					newGroups := groups[1:]
					if nrDamaged == len(condition) {
						if len(newGroups) == 0 {
							permutations = 1
						}
					} else if condition[nrDamaged] == '.' {
						permutations = countArrangementsInner(condition[nrDamaged+1:], newGroups, memoizationMap)
					} else if condition[nrDamaged] == '?' {
						permutations = countArrangementsInner(append([]rune{'.'}, condition[nrDamaged+1:]...), newGroups, memoizationMap)
					}
				}
			}
		}
	}

	// Memoize the result.
	memoizationMap[memKey] = permutations
	return permutations
}

// parseConditionAndGroups parses a line of input into a condition (array of runes) and a slice of groups (integers).
// Returns an error if any part of the input cannot be parsed correctly.
func parseConditionAndGroups(line string) ([]rune, []int, error) {
	fields := strings.Fields(line)
	groupsStr := strings.Split(fields[1], ",")
	groups := []int{}

	for _, groupStr := range groupsStr {
		group, errConv := strconv.Atoi(groupStr)
		if errConv != nil {
			return nil, nil, fmt.Errorf("cannot parse condition and groups: %w", errConv)
		}
		groups = append(groups, group)
	}

	return []rune(fields[0]), groups, nil
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
		condition, groups, errParse := parseConditionAndGroups(line)
		if errParse != nil {
			log.Fatal(errParse)
		}

		sum += countArrangements(condition, groups)
	}

	fmt.Println(sum)
}
