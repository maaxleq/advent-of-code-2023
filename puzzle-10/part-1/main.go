package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
)

const inputFile = "input.txt"

// It returns the parsed network and an error if any line contains invalid tiles.
type tile = int

// Tile types are enumerated here.
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

// getMaxTravelDistance computes the maximum distance that can be navigated from the start tile.
// It returns the maximum distance and an error if the start tile is not found or navigation is not possible.
func (n network) getMaxTravelDistance() (int, error) {
	x, y, errStart := n.findStart()
	if errStart != nil {
		return 0, fmt.Errorf("cannot navigate network: %w", errStart)
	}

	width, height := n.size()
	visited := make([][]bool, height)
	for i := range visited {
		visited[i] = make([]bool, width)
	}

	followPipe := func(x, y, dx, dy, distance int) int {
		maxDistance := distance

		for {
			x += dx
			y += dy

			// Check if out of bounds
			if x < 0 || x >= width || y < 0 || y >= height {
				break
			}

			currentTile := n[y][x]

			// Check if it's a looping tile or an empty tile
			if currentTile == empty || visited[y][x] {
				break
			}

			visited[y][x] = true
			distance++

			// Update maxDistance
			if distance > maxDistance {
				maxDistance = distance
			}

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

		return maxDistance
	}

	// Check possible initial directions from the start
	maxDist := 0
	if y+1 < height && !slices.Contains([]tile{horizontal, northEast, northWest, empty}, n[y+1][x]) { // Down
		maxDist = max(maxDist, followPipe(x, y, 0, 1, 0))
	}
	if y-1 >= 0 && !slices.Contains([]tile{horizontal, northEast, northWest, empty}, n[y-1][x]) { // Up
		maxDist = max(maxDist, followPipe(x, y, 0, -1, 0))
	}
	if x+1 < width && !slices.Contains([]tile{vertical, northWest, southWest, empty}, n[y][x+1]) { // Right
		maxDist = max(maxDist, followPipe(x, y, 1, 0, 0))
	}
	if x-1 >= 0 && !slices.Contains([]tile{vertical, northEast, southEast, empty}, n[y][x-1]) { // Left
		maxDist = max(maxDist, followPipe(x, y, -1, 0, 0))
	}

	return maxDist / 2, nil
}

// max returns the larger of two integers.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
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

	distance, errNav := n.getMaxTravelDistance()
	if errNav != nil {
		log.Fatal(errNav)
	}

	fmt.Println(distance)
}
