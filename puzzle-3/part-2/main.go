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

// gear represents a geat in the engine plan.
type gear struct {
	point point
}

// ratio returns the gear ratio if the gear is adjacent to exactly 2 part numbers (multiplying these 2 numbers), and 0 otherwise.
func (g *gear) ratio(pns []partNumber) int {
	adjacent := []partNumber{}

	for _, pn := range pns {
		for _, p := range pn.points {
			if areAdjacent(g.point, p) {
				adjacent = append(adjacent, pn)
				break
			}
		}
	}

	if len(adjacent) == 2 {
		return adjacent[0].num * adjacent[1].num
	}

	return 0
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
	gs := []gear{}

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

				if r == '*' {
					gs = append(gs, gear{
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
	for _, g := range gs {
		sum += g.ratio(pns)
	}

	fmt.Println(sum)
}
