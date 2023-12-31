package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const inputFile = "input.txt"

// beam represents a ray with a starting position (x, y) and a direction (dx, dy).
type beam struct {
	x, y, dx, dy int
}

// memBeam is a map that tracks if a beam has been visited in the simulation.
type memBeam map[beam]bool

// grid represents a 2D grid of runes, where each rune corresponds to a tile type.
type grid [][]rune

// simulate projects a beam through the grid, altering its path based on tile types,
// and records the energized tiles in eGrid. It uses memBeam to avoid revisiting the same path.
func (g grid) simulate(eGrid [][]bool, mem memBeam, startX, startY, dx, dy int) {
	b := beam{
		x:  startX,
		y:  startY,
		dx: dx,
		dy: dy,
	}

	if visited, found := mem[b]; found && visited {
		return
	}

	mem[b] = true

	x, y := startX, startY
	cont := true

	for cont {
		if x < 0 || y < 0 || y >= len(g) || x >= len(g[y]) {
			break
		}

		eGrid[y][x] = true

		switch g[y][x] {
		case '.':
			x += dx
			y += dy
		case '\\':
			prevDy := dy
			dy = dx
			dx = prevDy

			x += dx
			y += dy
		case '/':
			prevDy := dy
			dy = -dx
			dx = -prevDy

			x += dx
			y += dy
		case '-':
			if dy != 0 {
				g.simulate(eGrid, mem, x-1, y, -1, 0)
				g.simulate(eGrid, mem, x+1, y, 1, 0)
				cont = false
			} else {
				x += dx
				y += dy
			}
		case '|':
			if dx != 0 {
				g.simulate(eGrid, mem, x, y-1, 0, -1)
				g.simulate(eGrid, mem, x, y+1, 0, 1)
				cont = false
			} else {
				x += dx
				y += dy
			}
		}
	}
}

// countEnergizedTiles runs the simulation for the entire grid and counts the number of energized tiles.
func (g grid) countEnergizedTiles() int {
	eGrid := [][]bool{}

	for _, gLine := range g {
		eLine := []bool{}
		for range gLine {
			eLine = append(eLine, false)
		}

		eGrid = append(eGrid, eLine)
	}

	mem := make(memBeam)
	g.simulate(eGrid, mem, 0, 0, 1, 0)

	energizedCount := 0

	for _, eLine := range eGrid {
		for _, eTile := range eLine {
			if eTile {
				energizedCount++
			}
		}
	}

	return energizedCount
}

// parseGrid converts an array of strings into a grid structure.
func parseGrid(lines []string) grid {
	g := grid{}

	for _, line := range lines {
		g = append(g, []rune(line))
	}

	return g
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

	g := parseGrid(lines)

	energizedTiles := g.countEnergizedTiles()

	fmt.Println(energizedTiles)
}
