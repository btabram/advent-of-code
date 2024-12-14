package main

import (
	utils "aoc2024"
	"fmt"
	"regexp"
	"strconv"
)

const (
	width     = 101
	height    = 103
	midWidth  = width / 2
	midHeight = height / 2
)

type vec struct {
	x, y int
}

type guard struct {
	position, velocity vec
}

func (g *guard) move(time int) {
	g.position.x = (g.position.x + (g.velocity.x * time)) % width
	g.position.y = (g.position.y + (g.velocity.y * time)) % height

	if g.position.x < 0 {
		g.position.x += width
	}
	if g.position.y < 0 {
		g.position.y += height
	}
}

// Multiply the number of guards in each quadrant of the tile floor
func safetyFactor(guards []*guard) int {
	quadrantCounts := [4]int{0, 0, 0, 0}
	for _, g := range guards {
		if g.position.x == midWidth || g.position.y == midHeight {
			continue // In the middle, not in a quadrant
		}

		// Just need a unique and consistent key for each quadrant
		key := 0
		if g.position.x > midWidth {
			key = key | 1
		}
		if g.position.y > midHeight {
			key = key | 2
		}

		quadrantCounts[key]++
	}

	safetyFactor := 1
	for _, count := range quadrantCounts {
		safetyFactor *= count
	}
	return safetyFactor
}

func print(guards []*guard) {
	guardPositions := make(map[vec]int)
	for _, g := range guards {
		guardPositions[g.position]++
	}

	str := ""
	for y := range height {
		for x := range width {
			count := guardPositions[vec{x, y}]
			if count != 0 {
				str += strconv.Itoa(count)
			} else {
				str += "."
			}
		}
		str += "\n"
	}
	fmt.Println(str)
}

func main() {
	lines := utils.ReadLines("input.txt")

	numberRegex := regexp.MustCompile(`-?\d+`)

	guards := make([]*guard, len(lines))
	for i, line := range lines {
		// Each line is of the form "p=0,4 v=3,-3"
		numbers := numberRegex.FindAllString(line, 4)
		guards[i] = &guard{
			position: vec{x: utils.Int(numbers[0]), y: utils.Int(numbers[1])},
			velocity: vec{x: utils.Int(numbers[2]), y: utils.Int(numbers[3])},
		}
	}

	// For part 1 we need the safety factor after 100 moves
	for _, g := range guards {
		g.move(100)
	}
	part1 := safetyFactor(guards)

	// For part 2 we need to find the first time that the robots arrange themselves in a Christmas
	// tree pattern! To solve this I guessed that the "safety factor" from part 1 must be relevant
	// and simulated a larger number of iterations while tracking the min and max safety factor
	// values, and printing the state of the guards at each new min/max. I quickly found the tree as
	// the state with the lowest safety factor (out of all possible states, the system is periodic).
	part2 := 6516
	for _, g := range guards {
		g.move(part2 - 100) // Already done 100 moves above
	}
	print(guards)

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
