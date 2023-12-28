package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

const inputFile = "input.txt"

// lens represents a lens with a label and focal length.
type lens struct {
	label       string
	focalLength int
}

// focusingPower calculates the focusing power of a lens based on its position in a box and slot.
// It uses the formula: (boxN + 1) * (slotN + 1) * focalLength.
func (l lens) focusingPower(boxN int, slotN int) int {
	return (boxN + 1) * (slotN + 1) * l.focalLength
}

// box represents a container for multiple lenses.
type box struct {
	lenses []lens
}

// removeLens removes a lens from the box based on its label.
func (b *box) removeLens(label string) {
	lenses := []lens{}

	for _, l := range b.lenses {
		if l.label != label {
			lenses = append(lenses, l)
		}
	}

	b.lenses = lenses
}

// putLens adds a new lens to the box or replaces an existing one with the same label.
func (b *box) putLens(newLens lens) {
	replaced := false

	for i, l := range b.lenses {
		if l.label == newLens.label {
			b.lenses[i] = newLens
			replaced = true
		}
	}

	if !replaced {
		b.lenses = append(b.lenses, newLens)
	}
}

// boxes represents a collection of box objects.
type boxes []box

// totalFocusingPower calculates the total focusing power of all lenses in all boxes.
func (bs boxes) totalFocusingPower() int {
	sum := 0

	for boxN, box := range bs {
		for slotN, lens := range box.lenses {
			sum += lens.focusingPower(boxN, slotN)
		}
	}

	return sum
}

// hashAlgorithm calculates a simple hash of a given string.
// It iterates through each character in the string, converting it to an integer
// and performing a series of calculations to produce a hash value.
// The function returns an integer representing the hash.
func hashAlgorithm(s string) int {
	val := 0

	for _, r := range s {
		val += int(r)
		val *= 17
		val %= 256
	}

	return val
}

// parseSteps splits a string by commas and returns a slice of the substrings.
// It is used to process a line of text into individual steps or commands.
func parseSteps(line string) []string {
	return strings.Split(line, ",")
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

	if len(lines) < 1 {
		log.Fatal(fmt.Errorf("no line in file"))
	}

	bs := make([]box, 256)

	steps := parseSteps(lines[0])

	for _, step := range steps {
		stepSplit := strings.Split(strings.Split(step, "-")[0], "=")
		boxN := hashAlgorithm(stepSplit[0])

		if len(stepSplit) == 2 {
			focalLength, errConv := strconv.Atoi(stepSplit[1])
			if errConv != nil {
				log.Fatal(errConv)
			}

			bs[boxN].putLens(lens{
				label:       stepSplit[0],
				focalLength: focalLength,
			})
		} else {
			bs[boxN].removeLens(stepSplit[0])
		}
	}

	totalFocusingPower := boxes(bs).totalFocusingPower()

	fmt.Println(totalFocusingPower)
}
