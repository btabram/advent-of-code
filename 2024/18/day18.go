package main

import (
	utils "aoc2024"
	"fmt"
)

type vec struct {
	x, y int
}

// Max index value within the memory space
const limit = 70

func main() {
	lines := utils.ReadLines("input.txt")

	// Map of byte locations to the order they fall in. The first byte to fall has value 0 and so on
	fallingBytes := make(map[vec]int)
	for i, line := range lines {
		ints := utils.Ints(line)
		fallingBytes[vec{x: ints[0], y: ints[1]}] = i
	}

	start := vec{0, 0}

	isFinish := func(node vec) bool {
		return node.x == limit && node.y == limit
	}

	var fallenByteCount int

	getNeighbours := func(node vec) map[vec]int {
		x := node.x
		y := node.y

		possibleMoves := []vec{
			{x + 1, y},
			{x - 1, y},
			{x, y + 1},
			{x, y - 1},
		}

		neighbours := make(map[vec]int)
		for _, move := range possibleMoves {
			if move.x < 0 || move.x > limit || move.y < 0 || move.y > limit {
				continue // Out of bounds
			}

			fallIndex, ok := fallingBytes[move]
			if ok && fallIndex < fallenByteCount {
				continue // Corrupted memory, can't go there
			}

			neighbours[move] = 1 // Just count steps, every move has same cost
		}

		return neighbours
	}

	heuristic := func(node vec) int {
		return utils.Abs(limit-node.x) + utils.Abs(limit-node.y)
	}

	// Tries to solve the maze after a given number of bytes have fallen. Returns move count or -1
	tryReachExit := func(fallenBytes int) int {
		// `fallenByteCount` has been captured by `getNeighbours` so we just need to update this and
		// then we can keep calling the pathfinding util with the same arguments.
		fallenByteCount = fallenBytes

		moveCount, _ := utils.AStarPathfinding(start, isFinish, getNeighbours, heuristic)

		return moveCount
	}

	part1 := tryReachExit(1024)

	// In part 2 we want to find the coordinates of the byte which cuts of the exist when it falls
	var part2 string
	for i := 1025; i < len(lines); i++ { // Just brute force... it takes <10s
		moveCount := tryReachExit(i)
		if moveCount == -1 {
			// Maze wasn't solvable - exit has just been cut off by byte `i`` falling
			part2 = lines[i-1]
			break
		}
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
