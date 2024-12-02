package main

import (
	utils "aoc2024"
	"fmt"
)

type Levels []int

func (levels Levels) check() bool {
	reportShouldIncrease := levels[1] > levels[0]

	prevLevel := levels[0]
	for _, nextLevel := range levels[1:] {
		diff := nextLevel - prevLevel

		if utils.Abs(diff) < 1 || utils.Abs(diff) > 3 {
			return false // Adjacent levels must change by at least one and at most three
		}

		isIncreasing := diff > 0
		if isIncreasing != reportShouldIncrease {
			return false // Adjacent levels must all be increasing or all be decreasing
		}

		prevLevel = nextLevel
	}

	return true
}

func (levels Levels) checkWithProblemDampener() bool {
	for i := range levels {
		// Construct a slice like `levels` but omitting the `i`th element.
		withOmission := append(append(Levels{}, levels[:i]...), levels[i+1:]...)

		if withOmission.check() {
			return true
		}
	}

	return false
}

func main() {
	lines := utils.ReadLines("input.txt")

	reports := make([]Levels, len(lines))
	for i, line := range lines {
		reports[i] = utils.Ints(line)
	}

	var part1, part2 int
	for _, report := range reports {
		if report.check() {
			part1++
			part2++
			continue
		}
		// In part 2 we can also tolerate one bad level (bad levels can be completely ignored).
		if report.checkWithProblemDampener() {
			part2++
		}
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
