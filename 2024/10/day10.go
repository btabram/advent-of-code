package main

import (
	utils "aoc2024"
	"fmt"
)

type vec struct {
	x, y int
}

type TopographicMap map[vec]int

func (tm TopographicMap) getMoves(currentPos vec, requiredHeight int) []vec {
	x := currentPos.x
	y := currentPos.y

	neighbours := []vec{{x + 1, y}, {x - 1, y}, {x, y + 1}, {x, y - 1}}

	// We don't need to worry about neighbours being out of bounds here, they'll
	// end up as height 0 and we're always looking for height >= 1
	var moves []vec
	for _, neighbour := range neighbours {
		if tm[neighbour] == requiredHeight {
			moves = append(moves, neighbour)
		}
	}

	return moves
}

func (tm TopographicMap) evaluateTrailhead(pos vec) (int, int) {

	reachedPeaks := make(map[vec]bool)

	// Recursively try every possible valid move to find all trails.
	var countTrails func(vec, int) int
	countTrails = func(currentPos vec, currentHeight int) int {
		if currentHeight == 9 {
			// We've found a valid trail!
			reachedPeaks[currentPos] = true
			return 1
		}

		// Valid trails gain one unit of height per step
		requiredHeight := currentHeight + 1
		nextMoves := tm.getMoves(currentPos, requiredHeight)

		count := 0
		for _, move := range nextMoves {
			count += countTrails(move, requiredHeight)
		}

		return count
	}

	uniqueTrails := countTrails(pos, 0)

	// Return "score" (number of reachable 9s, part1) and "rating" (number of unique trails, part2)
	return len(reachedPeaks), uniqueTrails
}

func main() {
	lines := utils.ReadLines("input.txt")

	topoMap := make(TopographicMap)
	var possibleTrailheads []vec
	for y, line := range lines {
		for x, char := range line {
			topoMap[vec{x, y}] = int(char - '0')
			if char == '0' {
				possibleTrailheads = append(possibleTrailheads, vec{x, y})
			}
		}
	}

	part1, part2 := 0, 0
	for _, p := range possibleTrailheads {
		score, rating := topoMap.evaluateTrailhead(p)
		part1 += score
		part2 += rating
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
