package main

import (
	utils "aoc2024"
	"fmt"
)

const (
	north = iota
	east
	south
	west
)

type vec struct {
	x, y int
}

type node struct {
	position  vec
	direction int
}

func main() {
	lines := utils.ReadLines("input.txt")

	var start, end vec
	walls := make(map[vec]bool)
	for y, line := range lines {
		for x, char := range line {
			switch char {
			case 'S':
				start = vec{x, y}
			case 'E':
				end = vec{x, y}
			case '#':
				walls[vec{x, y}] = true
			}
		}
	}

	initialNode := node{position: start, direction: east}

	isFinish := func(n node) bool {
		return n.position == end
	}

	getNeighbours := func(n node) map[node]int {
		moves := make(map[node]int)

		// Can either turn at a cost of 1000
		turnOptions := []node{
			{n.position, north},
			{n.position, east},
			{n.position, south},
			{n.position, west},
		}
		for _, opt := range turnOptions {
			if opt.direction != n.direction {
				moves[opt] = 1000
			}
		}

		// Or move forwards at a cost of 1 (unless there's a wall!)
		var newPosition vec
		switch n.direction {
		case north:
			newPosition = vec{n.position.x, n.position.y - 1}
		case east:
			newPosition = vec{n.position.x + 1, n.position.y}
		case south:
			newPosition = vec{n.position.x, n.position.y + 1}
		case west:
			newPosition = vec{n.position.x - 1, n.position.y}
		}
		if !walls[newPosition] {
			moves[node{newPosition, n.direction}] = 1
		}

		return moves
	}

	// Use manhattan distance as a simple lower-bound
	heuristic := func(n node) int {
		return utils.Abs(n.position.x-end.x) + utils.Abs(n.position.y-end.y)
	}

	minCost, prevNodes := utils.AStarPathfinding(initialNode, isFinish, getNeighbours, heuristic)

	// Part 1 is simply the cost of the best path through the maze
	part1 := minCost

	// In part 2 we need to work out the number of tiles which are included in any of the best paths
	visited := make(map[vec]bool)
	queue := []node{
		{end, north},
		{end, east},
		{end, south},
		{end, west},
	}
	// Follow the paths backwards using the `prevNodes` map which our pathfinding produced
	for len(queue) != 0 {
		current := queue[0]
		queue = queue[1:]

		visited[current.position] = true

		queue = append(queue, prevNodes[current]...)
	}
	part2 := len(visited)

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
