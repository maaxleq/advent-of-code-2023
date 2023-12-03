package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
)

const inputFile = "input.txt"

// point represents a 2d point in a plan.
type point struct {
	x int
	y int
}

// partNumber is a number which is potentially part of the engine plan.
type partNumber struct {
	num    int
	points []point
}

// isIncluded returns true if the part number is adjacent to at least one symbol.
func (pn *partNumber) isIncluded(pss []partSymbol) bool {
	for _, ps := range pss {
		for _, p := range pn.points {
			if areAdjacent(ps.point, p) {
				return true
			}
		}
	}

	return false
}

// partSymbol represents a symbol in the engine plan.
type partSymbol struct {
	point point
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

// areAdjacent returns true if two points are adjacent, even diagonally.
func areAdjacent(p1, p2 point) bool {
	return math.Abs(float64(p1.x)-float64(p2.x)) <= 1 && math.Abs(float64(p1.y)-float64(p2.y)) <= 1
}

func main() {
	lines, errRead := readFileLines(inputFile)
	if errRead != nil {
		log.Fatal(errRead)
	}

	pns := []partNumber{}
	pss := []partSymbol{}

	currentPn := partNumber{}
	for y, line := range lines {
		for x, r := range line {
			if r >= '0' && r <= '9' {
				currentPn.points = append(currentPn.points, point{
					x: x, y: y,
				})
				currentPn.num = currentPn.num*10 + int(r-'0')
			} else {
				if currentPn.num != 0 {
					pns = append(pns, currentPn)
					currentPn = partNumber{}
				}

				if r != '.' {
					pss = append(pss, partSymbol{
						point{
							x: x, y: y,
						},
					})
				}
			}
		}

		if currentPn.num != 0 {
			pns = append(pns, currentPn)
			currentPn = partNumber{}
		}
	}

	sum := 0
	for _, pn := range pns {
		if pn.isIncluded(pss) {
			sum += pn.num
		}
	}

	fmt.Println(sum)
}
