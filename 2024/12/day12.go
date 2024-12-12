package main

import (
	utils "aoc2024"
	"fmt"
	"slices"
)

type vec struct {
	x, y int
}

func (a vec) add(b vec) vec {
	return vec{
		x: a.x + b.x,
		y: a.y + b.y,
	}
}

func (a vec) sub(b vec) vec {
	return vec{
		x: a.x - b.x,
		y: a.y - b.y,
	}
}

type neighbour struct {
	pos       vec
	direction rune // direction of the neighbour from us, will be 'N', 'E', 'S' or 'W'
}

func neighbours(position vec) []neighbour {
	x := position.x
	y := position.y
	return []neighbour{
		{vec{x + 1, y}, 'E'},
		{vec{x - 1, y}, 'W'},
		{vec{x, y + 1}, 'S'},
		{vec{x, y - 1}, 'N'},
	}
}

func numberOfSides(border []neighbour) int {
	// Not the most efficient but from a given position on the border we expand a side / line as far
	// as we can until we have a complete side. This means we'll get the same side multiple times so
	// to deduplicate we use a map with a stringified key.
	sides := make(map[string]bool)
	toKey := func(direction rune, positions []vec) string {
		// This mutates the input which isn't great, but it's fine in this case because calling this
		// function is the last thing we do with our slices.
		slices.SortFunc(positions, func(a, b vec) int {
			if a.x == b.x {
				return a.y - b.y
			}
			return a.x - b.x
		})
		return fmt.Sprintf("%v_%v", direction, positions)
	}

	// The direction of the border relative to the plot is important when working out the number of
	// sides, we may have sides sharing common positions because of complex plot shapes so it is
	// not sufficient to just track the border positions alone.
	borderByDirection := make(map[rune][]vec)
	for _, b := range border {
		borderByDirection[b.direction] = append(borderByDirection[b.direction], b.pos)
	}

	for direction, positions := range borderByDirection {
		// Set for easy checking whether a position is on the border
		positionSet := make(map[vec]bool)
		for _, p := range positions {
			positionSet[p] = true
		}

		for _, pos := range positions {
			// Work out if any neighbours are also in the border (with the same direction) and
			// therefore on the same side as us
			var sideNeighbours []vec
			for _, n := range neighbours(pos) {
				if positionSet[n.pos] {
					sideNeighbours = append(sideNeighbours, n.pos)
				}
			}

			// Handle side of length 1 case
			if len(sideNeighbours) == 0 {
				sides[fmt.Sprintf("%v_%v", direction, pos)] = true
				continue
			}

			// Could be cleverer, but just expand every neighbour which we're on a side with out
			// into the longest line we can (while staying on the border) to record all sides.
			for _, sideNeighbour := range sideNeighbours {
				var sidePositions []vec

				diff := sideNeighbour.sub(pos)

				// Expand side forwards
				next := sideNeighbour
				for positionSet[next] {
					sidePositions = append(sidePositions, next)
					next = next.add(diff)
				}

				// Expand side backwards
				prev := pos
				for positionSet[prev] {
					sidePositions = append(sidePositions, prev)
					prev = prev.sub(diff)
				}

				sides[toKey(direction, sidePositions)] = true
			}
		}
	}

	return len(sides)
}

func main() {
	lines := utils.ReadLines("input.txt")

	// A map is nice because out-of-bounds access will just work and return 0, which will correctly
	// be considered a border for any of our plants (no NULL characters in the input)
	garden := make(map[vec]rune)
	for y, line := range lines {
		for x, char := range line {
			garden[vec{x, y}] = char
		}
	}

	part1, part2 := 0, 0

	plotted := make(map[vec]bool)
	for pos, plant := range garden {
		if plotted[pos] {
			continue
		}

		plot := make(map[vec]bool)
		plot[pos] = true

		var border []neighbour

		// Work through plants in the plot one at a time, considering all their neighbours.
		queue := []vec{pos}
		for len(queue) > 0 {
			for _, n := range neighbours(queue[0]) {
				// It's a different plant, not part of this plot so add it to the border. Note that
				// the same position could be counted twice in the perimeter calculation if the plot
				// does not have straight sides (but it should be a different direction ever time!)
				if garden[n.pos] != plant {
					border = append(border, n)
					continue
				}

				// It's a plant in the plot! If we've not seen it before then register it and add
				// it to the queue for processing.
				if !plot[n.pos] {
					plot[n.pos] = true
					queue = append(queue, n.pos)
				}
			}

			queue = queue[1:]
		}

		part1 += len(plot) * len(border)
		part2 += len(plot) * numberOfSides(border)

		for p := range plot {
			plotted[p] = true
		}
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
