package main

import (
	utils "aoc2024"
	"fmt"
)

type Vec struct {
	x, y int
}

type Guard struct {
	position, direction Vec
}

var (
	up    = Vec{0, -1}
	right = Vec{1, 0}
	down  = Vec{0, 1}
	left  = Vec{-1, 0}
)

// turn 90 degrees to the right
func (g *Guard) turn() {
	switch g.direction {
	case up:
		g.direction = right
	case right:
		g.direction = down
	case down:
		g.direction = left
	case left:
		g.direction = up
	}
}

// walk the guard forwards, returns true on success and false if blocked
func (g *Guard) walk(obstacles map[Vec]bool) bool {
	nextPos := Vec{
		x: g.position.x + g.direction.x,
		y: g.position.y + g.direction.y,
	}

	if obstacles[nextPos] {
		return false
	}

	g.position = nextPos
	return true
}

// Let the guard follow their route and count how many squares they enter before leaving the area
func part1(guard Guard, obstacles map[Vec]bool, width, height int) int {
	visited := make(map[Vec]bool)
	visited[guard.position] = true

	for {
		success := guard.walk(obstacles)

		if !success {
			guard.turn()
			continue
		}

		x := guard.position.x
		y := guard.position.y
		if x < 0 || x >= width || y < 0 || y >= height {
			return len(visited)
		}

		visited[guard.position] = true
	}
}

// Count how many different ways you can get the guard stuck in a loop by adding one extra obstacle
func part2(initialGuard Guard, initialObstacles map[Vec]bool, width, height int) int {
	loops := 0

	for x := range width {
		for y := range height {
			if initialObstacles[Vec{x, y}] {
				continue // Already an obstacle here
			}

			// Try adding a new obstacle at (x, y)
			obstacles := make(map[Vec]bool)
			for key := range initialObstacles {
				obstacles[key] = true
			}
			obstacles[Vec{x, y}] = true

			guard := initialGuard

			seen := make(map[Guard]bool)
			seen[guard] = true

			for {
				success := guard.walk(obstacles)

				if !success {
					guard.turn()
					continue
				}

				x := guard.position.x
				y := guard.position.y
				if x < 0 || x >= width || y < 0 || y >= height {
					break // Going out of the area - haven't found a loop
				}

				if seen[guard] {
					loops++ // We've been in this exact state before - we've found a loop!
					break
				}

				seen[guard] = true
			}
		}
	}

	return loops
}

func main() {
	lines := utils.ReadLines("input.txt")

	obstacles := make(map[Vec]bool)
	var guard *Guard

	width := len(lines[0])
	height := len(lines)

	for y, line := range lines {
		for x, char := range line {
			if char == '#' {
				obstacles[Vec{x, y}] = true
			} else if char == '^' {
				guard = &Guard{position: Vec{x, y}, direction: up}
			}
		}
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1(*guard, obstacles, width, height))
	fmt.Printf("The answer to Part 2 is %v\n", part2(*guard, obstacles, width, height))
}
