package main

import (
	utils "aoc2024"
	"fmt"
	"strings"
)

func main() {
	lines := utils.ReadLines("input.txt")

	// The input is a list of available towel patterns then a list of desired designs
	availableTowels := strings.Split(lines[0], ", ")
	desiredDesigns := make([]string, len(lines)-2)
	copy(desiredDesigns, lines[2:])

	cache := make(map[string]int)

	// Recursively work out how many ways we can solve a given design with the towels available
	var solve func(string) int
	solve = func(design string) int {
		cachedValue, ok := cache[design]
		if ok {
			// Memoise since there's a LOT of different combinations to consider
			return cachedValue
		}

		solutionCount := 0
		for _, towel := range availableTowels {
			if strings.HasPrefix(design, towel) {
				if len(design) == len(towel) {
					solutionCount += 1 // Solved
				} else {
					solutionCount += solve(design[len(towel):])
				}
			}
		}

		cache[design] = solutionCount

		return solutionCount
	}

	part1 := 0 // The number of designs that can be solved
	part2 := 0 // The total number of possible solutions for all solvable designs

	for _, design := range desiredDesigns {
		solutionCount := solve(design)
		if solutionCount > 0 {
			part1++
		}
		part2 += solutionCount
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
