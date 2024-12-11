package main

import (
	utils "aoc2024"
	"fmt"
	"strconv"
)

type stones map[int]int

// The stones change every time you blink! The number of stones increases dramatically over time so
// it's not practical to track them all individually but there's no interaction between stones so we
// can just all group stones with the same value together and do calculations at the group level.
func (s *stones) blink() {
	newStones := make(stones)

	for stone, count := range *s {
		// First rule - a stone engraved with 0 becomes 1
		if stone == 0 {
			newStones[1] += count
			continue
		}

		// Second rule - a stone with an even number of digits is split into two stones
		stoneStr := strconv.Itoa(stone)
		if len(stoneStr)&1 == 0 {

			leftStone := utils.Int(stoneStr[:len(stoneStr)/2])
			rightStone := utils.Int(stoneStr[len(stoneStr)/2:])

			newStones[leftStone] += count
			newStones[rightStone] += count
			continue
		}

		// Third rule - any other stone has its value multiplied by 2024
		newStones[stone*2024] += count
	}

	// Change what `s` is pointing to
	*s = newStones
}

func (s stones) total() int {
	total := 0
	for _, count := range s {
		total += count
	}
	return total
}

func main() {
	input := utils.Read("input.txt")

	stones := make(stones)
	for _, initialStone := range utils.Ints(input) {
		stones[initialStone]++
	}

	// For part 1 we want the number of stones after 25 blinks
	for range 25 {
		stones.blink()
	}
	fmt.Printf("The answer to Part 1 is %v\n", stones.total())

	// For part 2 we need to go up to 75 blinks
	for range 50 {
		stones.blink()
	}
	fmt.Printf("The answer to Part 2 is %v\n", stones.total())
}
