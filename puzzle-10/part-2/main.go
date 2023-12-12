package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
)

const inputFile = "input.txt"

// tile represents a type of tile in the network.
type tile = int

const (
	empty tile = iota
	start
	vertical
	horizontal
	northEast
	northWest
	southEast
	southWest
)

// network represents a grid of tiles.
type network [][]tile

// findStart locates the starting point in the network.
// It returns the x and y coordinates of the start tile and an error if the start tile is not found.
func (n network) findStart() (int, int, error) {
	for y, line := range n {
		for x, t := range line {
			if t == start {
				return x, y, nil
			}
		}
	}

	return 0, 0, fmt.Errorf("unable to find start")
}

// size returns the width and height of the network.
func (n network) size() (int, int) {
	if len(n) == 0 {
		return 0, 0
	}

	return len(n[0]), len(n)
}

// getArea calculates the area enclosed by the pipes in the network.
// It returns the calculated area and an error if the network cannot be navigated.
func (n network) getArea() (int, error) {
	x, y, errStart := n.findStart()
	if errStart != nil {
		return 0, fmt.Errorf("cannot navigate network: %w", errStart)
	}

	width, height := n.size()
	loopBoundary := make([][]tile, height)
	for i := range loopBoundary {
		loopBoundary[i] = make([]tile, width)
		for j := range loopBoundary[i] {
			loopBoundary[i][j] = empty
		}
	}

	followPipe := func(x, y, dx, dy int) {
		for {
			x += dx
			y += dy

			// Check if out of bounds
			if x < 0 || x >= width || y < 0 || y >= height {
				break
			}

			currentTile := n[y][x]

			// Check if it's a looping tile or an empty tile
			if currentTile == empty || loopBoundary[y][x] != empty {
				break
			}

			loopBoundary[y][x] = currentTile

			// Change direction based on the type of pipe
			switch currentTile {
			case northEast:
				if dx == 0 {
					dx = 1
					dy = 0
				} else {
					dx = 0
					dy = -1
				}
			case northWest:
				if dx == 0 {
					dx = -1
					dy = 0
				} else {
					dx = 0
					dy = -1
				}
			case southEast:
				if dx == 0 {
					dx = 1
					dy = 0
				} else {
					dx = 0
					dy = 1
				}
			case southWest:
				if dx == 0 {
					dx = -1
					dy = 0
				} else {
					dx = 0
					dy = 1
				}
			}
		}
	}

	// Check possible initial directions from the start
	if y+1 < height && !slices.Contains([]tile{horizontal, northEast, northWest, empty}, n[y+1][x]) { // Down
		followPipe(x, y, 0, 1)
	}
	if y-1 >= 0 && !slices.Contains([]tile{horizontal, northEast, northWest, empty}, n[y-1][x]) { // Up
		followPipe(x, y, 0, -1)
	}
	if x+1 < width && !slices.Contains([]tile{vertical, northWest, southWest, empty}, n[y][x+1]) { // Right
		followPipe(x, y, 1, 0)
	}
	if x-1 >= 0 && !slices.Contains([]tile{vertical, northEast, southEast, empty}, n[y][x-1]) { // Left
		followPipe(x, y, -1, 0)
	}

	// Applying the even-odd rule to count the surface area
	surfaceArea := 0

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			inside := false
			for k := 0; k+j < width && k+i < height; k++ {
				if slices.Contains([]tile{horizontal, vertical, northWest, southEast}, loopBoundary[i+k][j+k]) {
					inside = !inside
				}
			}
			if inside && loopBoundary[i][j] == empty {
				surfaceArea++
			}
		}
	}

	return surfaceArea, nil
}

// tileFromRune converts a rune to a corresponding tile type.
// It returns the tile type and an error if the rune does not correspond to a valid tile.
func tileFromRune(r rune) (tile, error) {
	switch r {
	case '.':
		return empty, nil
	case '|':
		return vertical, nil
	case '-':
		return horizontal, nil
	case 'L':
		return northEast, nil
	case 'J':
		return northWest, nil
	case '7':
		return southWest, nil
	case 'F':
		return southEast, nil
	case 'S':
		return start, nil
	default:
		return tile(0), fmt.Errorf("invalid tile: %c", r)
	}
}

// parseNetwork converts a slice of string lines into a network.
// It returns the parsed network and an error if any line contains invalid tiles.
func parseNetwork(lines []string) (network, error) {
	net := network{}

	for _, line := range lines {
		netLine := []tile{}
		for _, r := range line {
			tile, errTile := tileFromRune(r)
			if errTile != nil {
				return nil, fmt.Errorf("cannot parse network: %w", errTile)
			}

			netLine = append(netLine, tile)
		}
		net = append(net, netLine)
	}

	return net, nil
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

	n, _ := parseNetwork(lines)

	area, errNav := n.getArea()
	if errNav != nil {
		log.Fatal(errNav)
	}

	fmt.Println(area)
}
