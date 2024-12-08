package main

import (
	utils "aoc2024"
	"fmt"
)

type vec struct {
	x, y int
}

func main() {
	lines := utils.ReadLines("input.txt")

	width := len(lines[0])
	height := len(lines)
	outOfBounds := func(pos vec) bool {
		return pos.x < 0 || pos.x >= width || pos.y < 0 || pos.y >= height
	}

	antennaeByFrequency := make(map[rune][]vec)
	for y, line := range lines {
		for x, char := range line {
			if char != '.' {
				antennaeByFrequency[char] = append(antennaeByFrequency[char], vec{x, y})
			}
		}
	}

	antinodesP1, antinodesP2 := make(map[vec]bool), make(map[vec]bool)

	for _, positions := range antennaeByFrequency {
		// Consider all possible antenna pairs (within a given frequency)
		for i := range positions {
			for j := range positions {
				if i == j {
					continue
				}

				dx := positions[i].x - positions[j].x
				dy := positions[i].y - positions[j].y

				// Follow the line forwards, past antenna `i`, until we're out of bounds
				antinode := positions[i]
				for k := 0; true; k++ {
					antinode = vec{antinode.x + dx, antinode.y + dy}
					if outOfBounds(antinode) {
						break
					}

					// In part 1 only the first position is an antinode, in part 2 the whole line is
					if k == 0 {
						antinodesP1[antinode] = true
					}
					antinodesP2[antinode] = true
				}

				// Follow the line backwards, past antenna `j`, until we're out of bounds
				antinode = positions[j]
				for k := 0; true; k++ {
					antinode = vec{antinode.x - dx, antinode.y - dy}
					if outOfBounds(antinode) {
						break
					}

					// In part 1 only the first position is an antinode, in part 2 the whole line is
					if k == 0 {
						antinodesP1[antinode] = true
					}
					antinodesP2[antinode] = true
				}

				// In part 2 every point along the line is an antinode, which includes the antennae!
				antinodesP2[positions[i]] = true
				antinodesP2[positions[j]] = true
			}
		}
	}

	fmt.Printf("The answer to Part 1 is %v\n", len(antinodesP1))
	fmt.Printf("The answer to Part 2 is %v\n", len(antinodesP2))
}
