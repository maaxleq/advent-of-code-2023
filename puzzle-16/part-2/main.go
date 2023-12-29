package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sync"
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

// countEnergizedTiles calculates the number of energized tiles in the grid
// by simulating the path of a beam from a given starting point and direction.
func (g grid) countEnergizedTiles(startX, startY, dx, dy int) int {
	eGrid := [][]bool{}

	for _, gLine := range g {
		eLine := []bool{}
		for range gLine {
			eLine = append(eLine, false)
		}

		eGrid = append(eGrid, eLine)
	}

	mem := make(memBeam)
	g.simulate(eGrid, mem, startX, startY, dx, dy)

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

// findMaxEnergizedTiles finds the maximum number of tiles that can be energized
// by a beam emitted from any edge of the grid.
func (g grid) findMaxEnergizedTiles() int {
	maxEnergizedTiles := 0

	bs := []beam{}

	maxX := len(g[0]) - 1
	maxY := len(g) - 1

	for y := range g {
		bs = append(bs, beam{
			x:  0,
			y:  y,
			dx: 1,
			dy: 0,
		}, beam{
			x:  maxX,
			y:  y,
			dx: -1,
			dy: 0,
		})
	}

	for x := range g[0] {
		bs = append(bs, beam{
			x:  x,
			y:  0,
			dx: 0,
			dy: 1,
		}, beam{
			x:  x,
			y:  maxY,
			dx: 0,
			dy: -1,
		})
	}

	ch := make(chan int, len(bs))
	wg := sync.WaitGroup{}

	for _, b := range bs {
		wg.Add(1)
		go func(b beam) {
			defer wg.Done()
			energizedTiles := g.countEnergizedTiles(b.x, b.y, b.dx, b.dy)
			ch <- energizedTiles
		}(b)
	}
	wg.Wait()
	close(ch)

	for count := range ch {
		if count > maxEnergizedTiles {
			maxEnergizedTiles = count
		}
	}

	return maxEnergizedTiles
}

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

	energizedTiles := g.findMaxEnergizedTiles()

	fmt.Println(energizedTiles)
}
