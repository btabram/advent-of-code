package main

import (
	utils "aoc2024"
	"fmt"
	"regexp"
	"strings"
)

func main() {
	input := utils.Read("input.txt")

	// Three types - `mul`, `do` and `don't`
	instructionRegex := regexp.MustCompile(`mul\((\d{1,3}),(\d{1,3})\)|do\(\)|don't\(\)`)

	// -1 means return all matches (no limit)
	validInstructions := instructionRegex.FindAllStringSubmatch(input, -1)

	var part1, part2 int
	enabled := true
	for _, instruction := range validInstructions {
		if strings.HasPrefix(instruction[0], "mul") {
			value := utils.Int(instruction[1]) * utils.Int(instruction[2])

			// Part 1 just sums all `mul` values but part 2 needs to consider the `do` / `don't` state
			part1 += value
			if enabled {
				part2 += value
			}

			continue
		}

		// The instruction must either be a `do` or `don't` at this point
		enabled = instruction[0] == "do()"
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
