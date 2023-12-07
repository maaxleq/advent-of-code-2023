package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

// race struct represents a racing scenario with a given time and record distance.
const inputFile = "input.txt"

type race struct {
	timeMs     int
	distanceMm int
}

// countWaysOfWinning calculates the number of ways the race can be won.
// This is determined by the number of times the distance reached is greater than the race's record distance,
// within the given time limit.
func (r *race) countWaysOfWinning() int {
	count := 0
	for i := 0; i <= r.timeMs; i++ {
		if r.calculateDistanceReached(i) > r.distanceMm {
			count++
		}
	}

	return count
}

// calculateDistanceReached calculates the distance reached based on a given press time.
// It returns the product of the press time and the remaining travel time, but returns 0 if
// the press time is not within the valid range.
func (r *race) calculateDistanceReached(pressTimeMs int) int {
	if pressTimeMs >= r.timeMs || pressTimeMs <= 0 {
		return 0
	}

	travelTimeMs := r.timeMs - pressTimeMs
	return pressTimeMs * travelTimeMs
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

// parseRaces takes an array of strings representing the race time and distance, and parses them into a race slice.
// It splits the input strings to extract time and distance, converts them to integers,
// and handles any format errors in the input data.
func parseRaces(lines [2]string) ([]race, error) {
	times := strings.Fields(strings.Split(lines[0], ":")[1])
	distances := strings.Fields(strings.Split(lines[1], ":")[1])

	if len(times) != len(distances) {
		return []race{}, fmt.Errorf("cannot parse races: different number of times (%d) and distances (%d)", len(times), len(distances))
	}

	races := []race{}
	for i := range times {
		timeInt, errTime := strconv.Atoi(times[i])
		if errTime != nil {
			return []race{}, fmt.Errorf("cannot parse races: bad time format: %s", times[i])
		}

		distanceInt, errDistance := strconv.Atoi(distances[i])
		if errDistance != nil {
			return []race{}, fmt.Errorf("cannot parse races: bad distance format: %s", distances[i])
		}

		races = append(races, race{
			timeMs:     timeInt,
			distanceMm: distanceInt,
		})
	}

	return races, nil
}

func main() {
	lines, errRead := readFileLines(inputFile)
	if errRead != nil {
		log.Fatal(errRead)
	}

	races, errRaces := parseRaces(*(*[2]string)(lines[0:2])) // Call parseRaces with the 2 first lines of the input.
	if errRaces != nil {
		log.Fatal(errRaces)
	}

	prod := 0
	if len(races) > 0 {
		prod = races[0].countWaysOfWinning()
		for i := 1; i < len(races); i++ {
			prod *= races[i].countWaysOfWinning()
		}
	}

	fmt.Println(prod)
}
