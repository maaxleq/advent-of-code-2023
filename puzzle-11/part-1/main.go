package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"slices"
)

const inputFile = "input.txt"

// universe represents a 2D grid of tiles.
type universe [][]tile

// getEmptyRowsCols identifies the indices of entirely empty rows and columns in the universe.
// It returns two slices of integers, the first representing row indices and the second column indices.
func (u universe) getEmptyRowsCols() ([]int, []int) {
	emptyRows := []int{}
	emptyCols := []int{}

	for i := 0; i < len(u); i++ {
		emptySegment := true
		for j := 0; j < len(u[i]); j++ {
			if u[i][j] != empty {
				emptySegment = false
				break
			}
		}
		if emptySegment {
			emptyRows = append(emptyRows, i)
		}
	}

	for i := 0; i < len(u[0]); i++ {
		emptySegment := true
		for j := 0; j < len(u); j++ {
			if u[j][i] != empty {
				emptySegment = false
				break
			}
		}
		if emptySegment {
			emptyCols = append(emptyCols, i)
		}
	}

	return emptyRows, emptyCols
}

// getGalaxiesCoordinates returns a slice of coordinate pairs representing the positions of galaxies in the universe.
// Each coordinate pair is an array of two integers [x, y].
func (u universe) getGalaxiesCoordinates() [][2]int {
	galaxies := [][2]int{}

	for y := range u {
		for x := range u[y] {
			if u[y][x] == galaxy {
				galaxies = append(galaxies, [2]int{x, y})
			}
		}
	}

	return galaxies
}

// getDistinctCoordinatePairs computes all distinct pairs of galaxy coordinates in the universe.
// It returns a slice of pairs of coordinate pairs, each represented as [2][2]int.
func (u universe) getDistinctCoordinatePairs() [][2][2]int {
	galaxies := u.getGalaxiesCoordinates()
	pairs := [][2][2]int{}

	for i := 0; i < len(galaxies); i++ {
		for j := i + 1; j < len(galaxies); j++ {
			pairs = append(pairs, [2][2]int{
				galaxies[i], galaxies[j],
			})
		}
	}

	return pairs
}

// tile represents the state of a single cell in the universe.
type tile int

// Constants representing the different types of tiles in the universe.
const (
	galaxy tile = iota
	empty
)

// expandUniverse duplicates tiles corresponding to empty rows and columns in the universe,
// effectively expanding it. It returns a new, expanded universe.
func expandUniverse(u universe) universe {
	expanded := universe{}
	emptyRows, emptyCols := u.getEmptyRowsCols()

	for y := 0; y < len(u); y++ {
		line := []tile{}
		for x := 0; x < len(u[y]); x++ {
			line = append(line, u[y][x])
			if slices.Contains(emptyCols, x) {
				line = append(line, u[y][x])
			}
		}
		expanded = append(expanded, line)
		if slices.Contains(emptyRows, y) {
			expanded = append(expanded, line)
		}
	}

	return expanded
}

// parseUniverse converts a slice of string lines into a universe.
// It returns the parsed universe and an error if the parsing fails.
func parseUniverse(lines []string) (universe, error) {
	universe := universe{}

	for _, line := range lines {
		universeLine := []tile{}
		for _, r := range line {
			switch r {
			case '#':
				universeLine = append(universeLine, galaxy)
			case '.':
				universeLine = append(universeLine, empty)
			default:
				return nil, fmt.Errorf("cannot parse universe: invalid tile %c", r)
			}
		}
		universe = append(universe, universeLine)
	}

	return universe, nil
}

// pairDistance calculates the Manhattan distance between two coordinate pairs.
// It returns the distance as an integer.
func pairDistance(pair [2][2]int) int {
	p1, p2 := pair[0], pair[1]
	return int(math.Abs(float64(p1[0]-p2[0])) + math.Abs(float64(p1[1]-p2[1])))
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

	u, errParse := parseUniverse(lines)
	if errParse != nil {
		log.Fatal(errParse)
	}

	uExp := expandUniverse(u)

	distinctCoordPairs := uExp.getDistinctCoordinatePairs()

	sum := 0
	for _, pair := range distinctCoordPairs {
		sum += pairDistance(pair)
	}

	fmt.Println(sum)
}
