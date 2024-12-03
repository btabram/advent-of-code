package main

import (
	"fmt"
	"os"
	"strconv"

	"AoC/utils"
)

type Coord struct {
	i, j int
}

func (c Coord) add(other Coord) Coord {
	return Coord{c.i + other.i, c.j + other.j}
}

type TreeGrid struct {
	trees [][]int
	size  int // Assume square
}

func (tg *TreeGrid) getHeight(c Coord) int {
	return tg.trees[c.i][c.j]
}

func (tg *TreeGrid) outOfBounds(c Coord) bool {
	return c.i < 0 || c.j < 0 || c.i >= tg.size || c.j >= tg.size
}

func newTreeGrid(input string) *TreeGrid {
	lines := utils.Lines(input)
	size := len(lines) // Assume square
	trees := make([][]int, size)
	for i, line := range lines {
		trees[i] = make([]int, size)
		for j, tree := range line {
			trees[i][j] = utils.CheckErr(strconv.Atoi(string(tree)))
		}
	}
	return &TreeGrid{trees, size}
}

func main() {
	grid := newTreeGrid(string(utils.CheckErr(os.ReadFile("input.txt"))))

	visibleTrees := make(map[Coord]bool)
	bestScenicScore := 0
	for i := 0; i < grid.size; i++ {
		for j := 0; j < grid.size; j++ {
			self := Coord{i, j}
			selfHeight := grid.getHeight(self)

			visible := false
			viewingDistances := make([]int, 0, 4)
			// Explore as far as we can see in each grid direction. We can work out the viewing
			// distances from the current tree and whether it is visible from outside the grid.
			for _, direction := range []Coord{{1, 0}, {-1, 0}, {0, 1}, {0, -1}} {
				pos := self
				seenTrees := 0
				for {
					pos = pos.add(direction)
					if grid.outOfBounds(pos) {
						visible = true
						break
					}
					seenTrees += 1
					if grid.getHeight(pos) >= selfHeight {
						break
					}
				}
				viewingDistances = append(viewingDistances, seenTrees)
			}

			if visible {
				visibleTrees[self] = true
			}
			if scenicScore := utils.Product(viewingDistances); scenicScore > bestScenicScore {
				bestScenicScore = scenicScore
			}
		}
	}

	fmt.Printf("The answer to Part 1 is %v.\n", len(visibleTrees))
	fmt.Printf("The answer to Part 2 is %v.\n", bestScenicScore)
}
