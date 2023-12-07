package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

const inputFile = "input.txt"

// handType is an enumeration of different types of poker hands.
type handType int

// Enumeration constants representing different types of poker hands.
const (
	fiveOfAKind  handType = 7
	fourOfAKind           = 6
	fullHouse             = 5
	threeOfAKind          = 4
	twoPair               = 3
	onePair               = 2
	highCard              = 1
)

// cardOrder defines the order of cards from lowest to highest.
const cardOrder = "23456789TJQKA"

// hand represents a poker hand with its 5 cards and a bid.
type hand struct {
	cards [5]rune
	bid   int
}

// getHandType evaluates and returns the type of the poker hand.
func (h *hand) getHandType() handType {
	rankCount := make(map[rune]int)
	for _, card := range h.cards {
		rankCount[card]++
	}

	pairCounter := 0
	var threes, pairCount bool
	for _, count := range rankCount {
		switch count {
		case 2:
			pairCounter++
			if pairCounter == 2 {
				pairCount = true
			}
		case 3:
			threes = true
		case 4:
			return fourOfAKind
		case 5:
			return fiveOfAKind
		}
	}

	if threes {
		if pairCounter > 0 {
			return fullHouse
		}
		return threeOfAKind
	}

	if pairCount {
		return twoPair
	}

	if pairCounter == 1 {
		return onePair
	}

	return highCard
}

// Less compares two hands and determines if the first hand is "less" than the other based on hand type and card ranks.
func (h *hand) Less(other hand) bool {
	handType, otherHandType := h.getHandType(), other.getHandType()
	if handType < otherHandType {
		return true
	} else if otherHandType < handType {
		return false
	} else {
		for i := 0; i < len(h.cards); i++ {
			value, otherValue := strings.IndexRune(cardOrder, h.cards[i]), strings.IndexRune(cardOrder, other.cards[i])
			if value < otherValue {
				return true
			} else if otherValue < value {
				return false
			}
		}

		return false
	}
}

// parseHands takes a slice of string lines, each representing a hand, and returns a slice of hand structs.
func parseHands(lines []string) ([]hand, error) {
	hands := []hand{}

	for _, line := range lines {
		fields := strings.Fields(line)
		bid, errBid := strconv.Atoi(fields[1])
		if errBid != nil {
			return []hand{}, fmt.Errorf("cannot parse hands: %w", errBid)
		}

		cards := (*[5]rune)(([]rune)(fields[0]))

		hands = append(hands, hand{
			cards: *cards,
			bid:   bid,
		})
	}

	return hands, nil
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

	hands, errHands := parseHands(lines)
	if errHands != nil {
		log.Fatal(errHands)
	}

	sort.Slice(hands, func(i, j int) bool {
		return hands[i].Less(hands[j])
	})

	winnings := 0
	for i, hand := range hands {
		winnings += (i + 1) * hand.bid
	}

	fmt.Println(winnings)
}
