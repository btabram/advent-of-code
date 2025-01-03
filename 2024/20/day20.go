package main

import (
	utils "aoc2024"
	"fmt"
)

type vec struct {
	x, y int
}

func getNeighbours(node vec) []vec {
	return []vec{
		{node.x + 1, node.y},
		{node.x - 1, node.y},
		{node.x, node.y + 1},
		{node.x, node.y - 1},
	}
}

func distance(a, b vec) int {
	return utils.Abs(b.x-a.x) + utils.Abs(b.y-a.y)
}

func getGoodShortcutCount(costs map[vec]int, maxShortcutDuration int) int {
	// Count of shortcuts with a given saving
	shortcutCounts := make(map[int]int)

	for node, cost := range costs {
		for otherNode, otherCost := range costs {
			if node == otherNode {
				continue
			}

			distance := distance(node, otherNode)
			if distance > maxShortcutDuration {
				continue // Too far apart
			}

			saving := cost - (otherCost + distance)
			if saving <= 0 {
				continue // Not a shortcut
			}

			shortcutCounts[saving]++
		}
	}

	goodShortcutCount := 0
	for saving, count := range shortcutCounts {
		if saving >= 100 {
			goodShortcutCount += count
		}
	}

	return goodShortcutCount
}

func main() {
	lines := utils.ReadLines("input.txt")

	walls := make(map[vec]bool)
	var end vec

	for y, line := range lines {
		for x, char := range line {
			switch char {
			case '#':
				walls[vec{x, y}] = true
			case 'E':
				end = vec{x, y}
			}
		}
	}

	// The number of moves required to get to the finish from every location on the race track
	costs := make(map[vec]int)
	costs[end] = 0

	// Flood fill from the finish to populate the costs map
	currentNodes := []vec{end}
	for i := 1; true; i++ {
		nextNodes := []vec{}

		for _, node := range currentNodes {
			for _, neighbour := range getNeighbours(node) {
				if walls[neighbour] {
					continue
				}

				_, visited := costs[neighbour]
				if visited {
					continue
				}

				nextNodes = append(nextNodes, neighbour)
				costs[neighbour] = i
			}
		}

		if len(nextNodes) == 0 {
			break
		}

		currentNodes = nextNodes
	}

	fmt.Printf("The answer to Part 1 is %v\n", getGoodShortcutCount(costs, 2))
	fmt.Printf("The answer to Part 2 is %v\n", getGoodShortcutCount(costs, 20))
}
