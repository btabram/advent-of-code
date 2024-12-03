package main

import (
	"fmt"
	"os"

	"AoC/pathfinding"
	"AoC/utils"
)

type Coord struct {
	i, j int
}

func main() {
	inputLines := utils.Lines(string(utils.CheckErr(os.ReadFile("input.txt"))))

	var start, end Coord
	heightGrid := make([][]byte, len(inputLines))
	for i, line := range inputLines {
		row := make([]byte, len(line))
		for j := range line {
			// Note any special positions but always put the correct height in the grid.
			switch val := line[j]; val {
			case 'S':
				start = Coord{i, j}
				row[j] = 'a'
			case 'E':
				end = Coord{i, j}
				row[j] = 'z'
			default:
				row[j] = val
			}
		}
		heightGrid[i] = row
	}

	getNeighbours := func(c Coord, isValidMove func(from, to Coord) bool) map[Coord]int {
		validNeighbours := map[Coord]int{}
		for _, n := range []Coord{{c.i + 1, c.j}, {c.i - 1, c.j}, {c.i, c.j + 1}, {c.i, c.j - 1}} {
			if n.i < 0 || n.j < 0 || n.i >= len(heightGrid) || n.j >= len(heightGrid[0]) {
				continue // Out of bounds
			}
			if !isValidMove(c, n) {
				continue
			}
			validNeighbours[n] = 1 // We're just counting moves so all costs are 1
		}
		return validNeighbours
	}

	// We can't go up more than one unit of height per move.
	isValidMovePart1 := func(from, to Coord) bool {
		return heightGrid[to.i][to.j] <= (heightGrid[from.i][from.j] + 1)
	}
	part1 := pathfinding.Dijkstra(
		start,
		func(c Coord) bool { return c == end },
		func(c Coord) map[Coord]int { return getNeighbours(c, isValidMovePart1) })

	// In part 2 we're working backwards from the summit, trying to find the quickest route down.
	isValidMovePart2 := func(a, b Coord) bool { return isValidMovePart1(b, a) }
	part2 := pathfinding.Dijkstra(
		end,
		func(c Coord) bool { return heightGrid[c.i][c.j] == 'a' },
		func(c Coord) map[Coord]int { return getNeighbours(c, isValidMovePart2) })

	fmt.Printf("The answer to Part 1 is %v.\n", part1)
	fmt.Printf("The answer to Part 2 is %v.\n", part2)
}
