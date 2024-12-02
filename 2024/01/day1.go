package main

import (
	utils "aoc2024"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
)

func main() {
	lines := utils.ReadLines("input.txt")

	// Parse ids lists
	ids := [2][]int{make([]int, 0, len(lines)), make([]int, 0, len(lines))}
	for _, line := range lines {
		stringVals := strings.Fields(line)

		for i := 0; i < 2; i++ {
			val, err := strconv.Atoi(stringVals[i])
			if err != nil {
				log.Fatal(err)
			}
			ids[i] = append(ids[i], val)
		}
	}

	// Sort both lists
	sort.Ints(ids[0])
	sort.Ints(ids[1])

	// Calculate total distance
	var part1 int
	for i, leftVal := range ids[0] {
		rightVal := ids[1][i]
		part1 += utils.Abs(leftVal - rightVal)
	}

	rightCounts := make(map[int]int)
	for _, rightVal := range ids[1] {
		rightCounts[rightVal]++
	}

	// Calculate similarity score
	var part2 int
	for _, leftVal := range ids[0] {
		part2 += leftVal * rightCounts[leftVal]
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
