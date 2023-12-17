package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"slices"
)

const inputFile = "input.txt"

// reflectionDirection represents the direction of reflection in the pattern.
// It can either be vertical or horizontal.
type reflectionDirection int

const (
	vertical   reflectionDirection = iota // vertical represents vertical reflection.
	horizontal                            // horizontal represents horizontal reflection.
)

// reflection stores the details of a specific reflection in a pattern.
// It includes the direction of the reflection and the index at which it occurs.
type reflection struct {
	direction reflectionDirection
	index     int
}

// value calculates a specific value based on the reflection's properties.
// If the reflection direction is vertical, it returns the index.
// If the direction is horizontal, it returns the index multiplied by 100.
func (r reflection) value() int {
	if r.direction == vertical {
		return r.index
	}

	return r.index * 100
}

// pattern represents a 2D pattern, defined as a slice of rune slices.
type pattern [][]rune

// hasReflection checks if a given pattern has a reflection specified by the argument 'r'.
// It returns true if the pattern contains the specified reflection, otherwise false.
func (p pattern) hasReflection(r reflection) bool {
	if r.direction == vertical {
		for i := 0; i+r.index < len(p[0]) && r.index-i > 0; i++ {
			for y := 0; y < len(p); y++ {
				if p[y][i+r.index] != p[y][r.index-i-1] {
					return false
				}
			}
		}
	} else {
		for i := 0; i+r.index < len(p) && r.index-i > 0; i++ {
			for x := 0; x < len(p[0]); x++ {
				if p[i+r.index][x] != p[r.index-i-1][x] {
					return false
				}
			}
		}
	}

	return true
}

// getAllPossibleCleanPatterns generates all possible variations of the current pattern
// by toggling each rune between '.' and '#'.
// It returns a slice of these modified patterns.
func (p pattern) getAllPossibleCleanPatterns() []pattern {
	cleanPatterns := []pattern{}
	for y := range p {
		for x := range p[y] {
			cleanPattern := pattern{}
			for i := range p {
				var row []rune
				for j := range p[i] {
					row = append(row, p[i][j])
				}
				cleanPattern = append(cleanPattern, row)
			}
			if cleanPattern[y][x] == '.' {
				cleanPattern[y][x] = '#'
			} else {
				cleanPattern[y][x] = '.'
			}
			cleanPatterns = append(cleanPatterns, cleanPattern)
		}
	}

	return cleanPatterns
}

// findCleanReflection identifies a reflection in the pattern after modifying it
// using getAllPossibleCleanPatterns. It returns the found reflection and an error, if any.
func (p pattern) findCleanReflection() (reflection, error) {
	ogReflection, errOgR := p.findReflection()
	if errOgR != nil {
		return reflection{}, errOgR
	}
	cleanPatterns := p.getAllPossibleCleanPatterns()
	for _, cleanPattern := range cleanPatterns {
		r, errR := cleanPattern.findReflection(ogReflection)
		if errR == nil {
			return r, nil
		}
	}

	return reflection{}, fmt.Errorf("no clean reflection in pattern")
}

// findReflection searches for a reflection in the pattern, optionally skipping some reflections.
// It returns the first found reflection and an error if no reflection is found.
func (p pattern) findReflection(skip ...reflection) (reflection, error) {
	possibleReflections := []reflection{}
	for y := 1; y < len(p); y++ {
		possibleReflections = append(possibleReflections, reflection{
			index:     y,
			direction: horizontal,
		})
	}
	for x := 1; x < len(p[0]); x++ {
		possibleReflections = append(possibleReflections, reflection{
			index:     x,
			direction: vertical,
		})
	}

	for _, r := range possibleReflections {
		if p.hasReflection(r) && !slices.Contains(skip, r) {
			return r, nil
		}
	}

	return reflection{}, fmt.Errorf("no reflection in pattern")
}

// parsePatterns parses a slice of strings into a slice of patterns.
// Each pattern is separated by an empty string in the slice.
func parsePatterns(lines []string) []pattern {
	currentPattern := pattern{}
	patterns := []pattern{}

	for _, line := range lines {
		if line == "" {
			patterns = append(patterns, currentPattern)
			currentPattern = pattern{}
		} else {
			currentPattern = append(currentPattern, []rune(line))
		}
	}

	if len(currentPattern) != 0 {
		patterns = append(patterns, currentPattern)
	}

	return patterns
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

	patterns := parsePatterns(lines)

	sum := 0
	for _, p := range patterns {
		r, errR := p.findCleanReflection()
		if errR != nil {
			log.Fatal(errR)
		}

		sum += r.value()
	}

	fmt.Println(sum)
}
