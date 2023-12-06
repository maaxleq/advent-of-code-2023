package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

const inputFile = "input.txt"

// mapOrder defines the order of mappings to be applied to each seed.
var mapOrder = []string{
	"seed-to-soil",
	"soil-to-fertilizer",
	"fertilizer-to-water",
	"water-to-light",
	"light-to-temperature",
	"temperature-to-humidity",
	"humidity-to-location",
}

// intervalMap is a map that associates a string key with a slice of uint64 triples.
type intervalMap = map[string][][3]uint64

// getSeedsAndIntervals parses the input lines and extracts seeds and interval mappings.
// It returns a slice of seeds, a mapping of intervals, and an error if any occurs.
func getSeedsAndIntervals(lines []string) ([]uint64, intervalMap, error) {
	intervals := make(intervalMap)
	seeds := []uint64{}

	var currentMap string
	for _, line := range lines {
		if len(line) == 0 { // Skip line if it is empty.
			continue
		}

		if strings.HasPrefix(line, "seeds") { // Get seeds
			seedsStr := strings.Fields(strings.Split(line, ":")[1])
			for _, seedStr := range seedsStr {
				seed, errParse := strconv.ParseUint(seedStr, 10, 64)
				if errParse != nil {
					return nil, nil, fmt.Errorf("invalid seed: %s", seedStr)
				}

				seeds = append(seeds, seed)
			}
		} else if strings.Contains(line, "map:") { // Get current map
			currentMap = strings.Fields(line)[0]
		} else if _, errParse := strconv.Atoi(line[0:1]); errParse == nil { // Get mapping rule
			fields := strings.Fields(line)
			ruleNumbers := [3]uint64{}

			if len(fields) != 3 {
				return nil, nil, fmt.Errorf("invalid rule: %s", line)
			}

			for i := 0; i < 3; i++ {
				number, errParse := strconv.ParseUint(fields[i], 10, 64)
				if errParse != nil {
					return nil, nil, fmt.Errorf("invalid number: %s", fields[i])
				}

				ruleNumbers[i] = number
			}

			intervals[currentMap] = append(intervals[currentMap], ruleNumbers)
		}
	}

	return seeds, intervals, nil
}

// getMapping applies the mapping rules to a given number for a specified mapName.
// It returns the transformed number according to the mapping rules, or an error if the map does not exist.
func getMapping(mapName string, intervals intervalMap, num uint64) (uint64, error) {
	rules, exists := intervals[mapName]
	if !exists {
		return 0, fmt.Errorf("map does not exist: %s", mapName)
	}

	for _, rule := range rules {
		if num >= rule[1] && num < rule[1]+rule[2] {
			return rule[0] + (num - rule[1]), nil
		}
	}

	return num, nil
}

// getLocationForSeed computes the final location value for a given seed by applying the sequence of mappings.
// It returns the computed location or an error if any mapping fails.
func getLocationForSeed(intervals intervalMap, seed uint64) (uint64, error) {
	currentNum := seed
	for _, mapName := range mapOrder {
		newNum, errMapping := getMapping(mapName, intervals, currentNum)
		if errMapping != nil {
			return 0, fmt.Errorf("cannot get location for seed %d: %w", seed, errMapping)
		}

		currentNum = newNum
	}

	return currentNum, nil
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

	seeds, intervals, errGet := getSeedsAndIntervals(lines)
	if errGet != nil {
		log.Fatal(errGet)
	}

	var lowestLoc uint64 = math.MaxUint64

	for _, seed := range seeds {
		location, errLoc := getLocationForSeed(intervals, seed)
		if errLoc != nil {
			log.Fatal(errLoc)
		}

		if location < lowestLoc {
			lowestLoc = location
		}
	}

	fmt.Println(lowestLoc)
}
