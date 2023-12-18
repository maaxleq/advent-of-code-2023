package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

// inputFile defines the name of the file to be read.
const inputFile = "input.txt"

// cacheItem stores a snapshot of the platform's state and the iteration count at which this state was recorded.
type cacheItem struct {
	runes [][]rune // The 2D array of runes representing the platform's current state.
	i     int      // The iteration count when this state was recorded.
}

// cache maps a platform's hash string to its cached state.
type cache map[string]cacheItem

// platform represents a 2D grid of runes where 'O' represents a rounded rock, '#' represents a cube-shaped rock and '.' represents an empty space.
type platform [][]rune

// getLoad calculates the total load on the platform. It adds up the vertical positions of all 'O's, with the top row being the highest value.
func (p platform) getLoad() int {
	height := len(p)
	load := 0

	for y := range p {
		for x := range p[y] {
			if p[y][x] == 'O' {
				load += height - y
			}
		}
	}

	return load
}

// hash generates a string representation of the platform's current state.
func (p platform) hash() string {
	h := ""

	// Concatenate each line of the platform to form the hash.
	for _, line := range p {
		h += string(line)
	}

	return h
}

// getNorthboundPosition finds the next position to the north where an 'O' can move.
func (p platform) getNorthboundPosition(x, y int) int {
	// Check each position north of the current one.
	for i := y - 1; i >= 0; i-- {
		if p[i][x] != '.' {
			return i + 1
		}
	}

	return 0
}

// getSouthboundPosition finds the next position to the south where an 'O' can move.
func (p platform) getSouthboundPosition(x, y int) int {
	// Check each position south of the current one.
	for i := y + 1; i < len(p); i++ {
		if p[i][x] != '.' {
			return i - 1
		}
	}

	return len(p) - 1
}

// getWestboundPosition finds the next position to the west where an 'O' can move.
func (p platform) getWestboundPosition(x, y int) int {
	// Check each position west of the current one.
	for i := x - 1; i >= 0; i-- {
		if p[y][i] != '.' {
			return i + 1
		}
	}

	return 0
}

// getEastboundPosition finds the next position to the east where an 'O' can move.
func (p platform) getEastboundPosition(x, y int) int {
	// Check each position east of the current one.
	for i := x + 1; i < len(p[0]); i++ {
		if p[y][i] != '.' {
			return i - 1
		}
	}

	return len(p[0]) - 1
}

// tiltNorth tilts the platform north, moving all 'O's upwards.
func (p *platform) tiltNorth() {
	// Move each 'O' to its northbound position.
	for y := range *p {
		for x := range (*p)[y] {
			if (*p)[y][x] == 'O' {
				nPos := p.getNorthboundPosition(x, y)
				(*p)[y][x] = '.'
				(*p)[nPos][x] = 'O'
			}
		}
	}
}

// tiltSouth tilts the platform south, moving all 'O's downwards.
func (p *platform) tiltSouth() {
	// Move each 'O' to its southbound position.
	for y := len(*p) - 1; y >= 0; y-- {
		for x := range (*p)[y] {
			if (*p)[y][x] == 'O' {
				sPos := p.getSouthboundPosition(x, y)
				(*p)[y][x] = '.'
				(*p)[sPos][x] = 'O'
			}
		}
	}
}

// tiltWest tilts the platform west, moving all 'O's to the left.
func (p *platform) tiltWest() {
	// Move each 'O' to its westbound position.
	for x := range (*p)[0] {
		for y := range *p {
			if (*p)[y][x] == 'O' {
				wPos := p.getWestboundPosition(x, y)
				(*p)[y][x] = '.'
				(*p)[y][wPos] = 'O'
			}
		}
	}
}

// tiltEast tilts the platform east, moving all 'O's to the right.
func (p *platform) tiltEast() {
	// Move each 'O' to its eastbound position.
	for x := len((*p)[0]) - 1; x >= 0; x-- {
		for y := range *p {
			if (*p)[y][x] == 'O' {
				ePos := p.getEastboundPosition(x, y)
				(*p)[y][x] = '.'
				(*p)[y][ePos] = 'O'
			}
		}
	}
}

// rotate performs a full rotation of the platform, tilting it in all four cardinal directions.
func (p *platform) rotate() {
	p.tiltNorth()
	p.tiltWest()
	p.tiltSouth()
	p.tiltEast()
}

// detectCyclePeriod detects the cycle period of the platform's movement.
func (p *platform) detectCyclePeriod(c cache, maxSearch int) (int, int, error) {
	// Iteratively check for cycle in platform states.
	for i := 0; i < maxSearch; i++ {
		h := p.hash()
		item, exists := c[h]
		if exists {
			// Return cycle start and end if found.
			return item.i, i, nil
		} else {
			// Cache current state and rotate.
			c[h] = cacheItem{
				runes: *p,
				i:     i,
			}
			p.rotate()
		}
	}

	// Return an error if no cycle is found within the specified limit.
	return 0, 0, fmt.Errorf("no cycle found within %d rotations", maxSearch)
}

// parsePlatform converts an array of strings into a platform structure.
func parsePlatform(lines []string) platform {
	p := platform{}

	for _, line := range lines {
		p = append(p, []rune(line))
	}

	return p
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
	// Read each line of the file.
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

	pDetect := parsePlatform(lines)
	c := cache(make(map[string]cacheItem))

	// Detect the cycle period of the platform.
	start, end, errPeriod := pDetect.detectCyclePeriod(c, 1_000_000_000)
	if errPeriod != nil {
		log.Fatal(errPeriod)
	}

	period := end - start

	p := parsePlatform(lines)

	// Perform rotations to simulate the final state.
	for i := 0; i < start+(1_000_000_000-start)%period; i++ {
		p.rotate()
	}

	// Calculate and print the final load on the platform.
	load := p.getLoad()

	fmt.Println(load)
}
