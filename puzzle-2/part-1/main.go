package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const (
	inputFile = "input.txt"

	redCount   = 12
	greenCount = 13
	blueCount  = 14
)

type color string

const (
	red   color = "red"
	green color = "green"
	blue  color = "blue"
)

type cubes struct {
	count int
	color color
}

type set []cubes

func (set *set) countColors() (int, int, int) {
	r, g, b := 0, 0, 0

	for _, cubes := range []cubes(*set) {
		switch cubes.color {
		case red:
			r += cubes.count
		case green:
			g += cubes.count
		case blue:
			b += cubes.count
		}
	}

	return r, g, b
}

type game struct {
	id   int
	sets []set
}

func (game *game) isPossible() bool {
	for _, set := range game.sets {
		r, g, b := set.countColors()
		if r > redCount || g > greenCount || b > blueCount {
			return false
		}
	}

	return true
}

func parseGameLine(line string) (*game, error) {
	// Splitting the line into parts
	parts := strings.Split(line, ": ")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid format")
	}

	// Parsing the game id
	gameID, err := strconv.Atoi(parts[0][5:])
	if err != nil {
		return nil, fmt.Errorf("invalid game id: %v", err)
	}

	// Splitting into sets
	rawSets := strings.Split(parts[1], "; ")
	var sets []set

	for _, rawSet := range rawSets {
		var cubesSet set

		// Splitting into cubes
		rawCubes := strings.Split(rawSet, ", ")
		for _, rawCube := range rawCubes {
			cubeParts := strings.Split(rawCube, " ")
			if len(cubeParts) != 2 {
				return nil, fmt.Errorf("invalid cubes format")
			}

			count, err := strconv.Atoi(cubeParts[0])
			if err != nil {
				return nil, fmt.Errorf("invalid count: %v", err)
			}

			cubesSet = append(cubesSet, cubes{count: count, color: color(cubeParts[1])})
		}

		sets = append(sets, cubesSet)
	}

	return &game{id: gameID, sets: sets}, nil
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

	idsSum := 0
	for _, line := range lines {
		game, errParse := parseGameLine(line)
		if errParse != nil {
			log.Fatal(errParse)
		}

		if game.isPossible() {
			idsSum += game.id
		}
	}

	fmt.Println(idsSum)
}
