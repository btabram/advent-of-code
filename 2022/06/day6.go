package main

import (
	"fmt"
	"os"

	"AoC/utils"
)

// Markers are groups of consecutive unique characters of a particular length.
func findMarker(datastream string, markerLength int) *int {
	for endPos := markerLength; endPos < len(datastream); endPos++ {
		possibleMarker := datastream[endPos-markerLength : endPos]
		set := make(map[rune]interface{})
		for _, char := range possibleMarker {
			set[char] = nil
		}
		if len(set) == markerLength {
			return &endPos
		}
	}
	return nil
}

func main() {
	input := string(utils.CheckErr(os.ReadFile("input.txt")))

	fmt.Printf("The answer to Part 1 is %v.\n", *findMarker(input, 4))
	fmt.Printf("The answer to Part 2 is %v.\n", *findMarker(input, 14))
}
