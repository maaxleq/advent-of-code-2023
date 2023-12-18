package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

const inputFile = "input.txt"

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

// getNorthboundPosition finds the position directly above a given coordinate (x, y) where the next 'O' can be placed. If no obstruction is found, it returns 0.
func (p platform) getNorthboundPosition(x, y int) int {
	for i := y - 1; i >= 0; i-- {
		if p[i][x] != '.' {
			return i + 1
		}
	}

	return 0
}

// tiltNorth simulates tilting the platform northward, causing all 'O's to move up as far as possible without overlapping.
func (p *platform) tiltNorth() {
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

// parsePlatform converts a slice of strings into a platform. Each string represents a row in the platform.
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

	p := parsePlatform(lines)
	p.tiltNorth()

	load := p.getLoad()

	fmt.Println(load)
}
