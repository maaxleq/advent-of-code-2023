package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

const inputFile = "input.txt"

type grid [][]rune

func (g grid) simulate(eGrid [][]bool, startX, startY, dx, dy int) {
	x, y := startX, startY
	cont := true

	for cont {
		fmt.Println(x, y, dx, dy)

		if x < 0 || y < 0 || x >= len(g[y]) || y >= len(g) {
			break
		}

		if math.Abs(float64(dx))-math.Abs(float64(dy)) == 0 {
			log.Fatalf("illegal movement")
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
				g.simulate(eGrid, x-1, y, -1, 0)
				g.simulate(eGrid, x+1, y, 1, 0)
				cont = false
			} else {
				x += dx
				y += dy
			}
		case '|':
			if dx != 0 {
				g.simulate(eGrid, x, y-1, 0, -1)
				g.simulate(eGrid, x, y+1, 0, 1)
				cont = false
			} else {
				x += dx
				y += dy
			}
		}
	}
}

func (g grid) countEnergizedTiles() int {
	eGrid := [][]bool{}

	for _, gLine := range g {
		eLine := []bool{}
		for range gLine {
			eLine = append(eLine, false)
		}

		eGrid = append(eGrid, eLine)
	}

	g.simulate(eGrid, 0, 0, 1, 0)

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
