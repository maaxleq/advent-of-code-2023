package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

const inputFile = "input.txt"

// expansionMultiplier is used to calculate the expanded distance between galaxy pairs.
const expansionMultiplier int = 1000000

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

// pairDistanceWithExpansion calculates the expanded distance between a pair of galaxy coordinates.
// It takes into account empty rows and columns that expand the distance.
func pairDistanceWithExpansion(pair [2][2]int, emptyRows, emptyCols []int) int {
	rowExpCount, colExpCount := 0, 0

	p1, p2 := pair[0], pair[1]

	for _, row := range emptyRows {
		if row <= int(math.Max(float64(p1[1]), float64(p2[1]))) && row >= int(math.Min(float64(p1[1]), float64(p2[1]))) {
			rowExpCount++
		}
	}

	for _, col := range emptyCols {
		if col <= int(math.Max(float64(p1[0]), float64(p2[0]))) && col >= int(math.Min(float64(p1[0]), float64(p2[0]))) {
			colExpCount++
		}
	}

	return int(math.Abs(float64(p1[0]-p2[0]))+math.Abs(float64(p1[1]-p2[1]))) + (expansionMultiplier-1)*(rowExpCount+colExpCount)
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

	distinctCoordPairs := u.getDistinctCoordinatePairs()
	emptyRows, emptyCols := u.getEmptyRowsCols()

	sum := 0
	for _, pair := range distinctCoordPairs {
		sum += pairDistanceWithExpansion(pair, emptyRows, emptyCols)
	}

	fmt.Println(sum)
}
