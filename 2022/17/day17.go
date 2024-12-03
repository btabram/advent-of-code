package main

import (
	"fmt"
	"os"

	"AoC/utils"
)

type Generator struct {
	pattern string
	index   int
}

func (g *Generator) next() byte {
	next := g.pattern[g.index]
	g.index = (g.index + 1) % len(g.pattern)
	return next
}

func newGenerator(pattern string) *Generator {
	return &Generator{pattern, 0}
}

type Coord struct {
	x, y int
}

type Chamber struct {
	highestRockPos int
	occupied       map[Coord]bool
}

func (c *Chamber) addRock(rockPoints []Coord) {
	for _, point := range rockPoints {
		c.occupied[point] = true
		if point.y > c.highestRockPos {
			c.highestRockPos = point.y
		}
	}
}

func (c *Chamber) canMove(rockPoints []Coord) bool {
	for _, point := range rockPoints {
		if point.x < 1 || point.y < 1 || point.x > 7 {
			return false // Out of bounds
		}
		if c.occupied[point] {
			return false // Another rock is already there
		}
	}
	return true
}

func (c *Chamber) String() string {
	str := ""
	for y := c.highestRockPos; y >= 1; y-- {
		str += "|"
		for x := 1; x <= 7; x++ {
			if c.occupied[Coord{x, y}] {
				str += "#"
			} else {
				str += "."
			}
		}
		str += "|\n"
	}
	str += "+-------+\n"
	return str
}

func newChamber() *Chamber {
	return &Chamber{0, make(map[Coord]bool)}
}

// There's 5 rock shapes.
const (
	rsHorizontalLine   byte   = 'H'
	rsPlus             byte   = '+'
	rsBackwardsL       byte   = 'L'
	rsVerticaLine      byte   = 'V'
	rsSquare           byte   = 'S'
	rsGeneratorPattern string = "H+LVS"
)

func spwanRock(rockShape byte, maxY int) []Coord {
	// Rocks spawn with a 2 unit gap to the left wall & 3 units gap up from the highest rock so far.
	switch rockShape {
	case rsHorizontalLine:
		return []Coord{
			{3, maxY + 4}, {4, maxY + 4}, {5, maxY + 4}, {6, maxY + 4},
		}
	case rsPlus:
		return []Coord{
			{4, maxY + 4}, {3, maxY + 5}, {4, maxY + 5}, {5, maxY + 5}, {4, maxY + 6},
		}
	case rsBackwardsL:
		return []Coord{
			{3, maxY + 4}, {4, maxY + 4}, {5, maxY + 4}, {5, maxY + 5}, {5, maxY + 6},
		}
	case rsVerticaLine:
		return []Coord{
			{3, maxY + 4}, {3, maxY + 5}, {3, maxY + 6}, {3, maxY + 7},
		}
	case rsSquare:
		return []Coord{
			{3, maxY + 4}, {4, maxY + 4}, {3, maxY + 5}, {4, maxY + 5},
		}
	default:
		panic(fmt.Sprintf("Invalid rock shape: %s", string(rockShape)))
	}
}

func moveRockDown(rockPoints []Coord) []Coord {
	return utils.Transform(rockPoints, func(c Coord) Coord { return Coord{c.x, c.y - 1} })
}

func moveRockLeft(rockPoints []Coord) []Coord {
	return utils.Transform(rockPoints, func(c Coord) Coord { return Coord{c.x - 1, c.y} })
}

func moveRockRight(rockPoints []Coord) []Coord {
	return utils.Transform(rockPoints, func(c Coord) Coord { return Coord{c.x + 1, c.y} })
}

func doMainLoop(chamber *Chamber, rockShapeGenerator, jetPattern *Generator) {
	rock := spwanRock(rockShapeGenerator.next(), chamber.highestRockPos)
	for {
		// First the rock is blown around by the jets of hot gas.
		var maybeNewRockPos []Coord
		switch jetPattern.next() {
		case '<':
			maybeNewRockPos = moveRockLeft(rock)
		case '>':
			maybeNewRockPos = moveRockRight(rock)
		}
		if chamber.canMove(maybeNewRockPos) {
			rock = maybeNewRockPos
		}
		// Then it falls a bit.
		maybeNewRockPos = moveRockDown(rock)
		if chamber.canMove(maybeNewRockPos) {
			rock = maybeNewRockPos
		} else {
			break // We've hit the floor or a rock that fell earlier
		}
	}
	chamber.addRock(rock)
}

func getTowerHeightAfterNRocks(input string, n int) int {
	chamber := newChamber()
	jetPattern := newGenerator(input)
	rockShapeGenerator := newGenerator(rsGeneratorPattern)
	for i := 0; i < n; i++ {
		doMainLoop(chamber, rockShapeGenerator, jetPattern)
	}
	return chamber.highestRockPos
}

// For part 2 we need the tower height after a very large number of rocks. Since the sequence of
// rocks and jet directions both repeat there must be a periodicity to the tower of rocks that is
// produced. If we can find that out we can work out how tall the tower will be after a number of
// rocks without actually simulating all the rocks.
func findPeriodicity(input string) (int, int) {
	chamber := newChamber()
	jetPattern := newGenerator(input)
	rockShapeGenerator := newGenerator(rsGeneratorPattern)

	stateCache := map[string]int{}

	i := 0
	for {
		// We include the last 50 rows of the rock tower in our state. It's theortically possible,
		// but very unlikely, that the shape of the tower >50 rows down will affect future moves.
		towerState := ""
		for j := 0; j < 50; j++ {
			for x := 1; x <= 7; x++ {
				if chamber.occupied[Coord{x, chamber.highestRockPos - j}] {
					towerState += "#"
				} else {
					towerState += "."
				}
			}
		}
		state := fmt.Sprintf("%d,%d,%s", rockShapeGenerator.index, jetPattern.index, towerState)
		if prevI, ok := stateCache[state]; ok {
			return prevI, i - prevI // Return info about the first repeated state
		}
		stateCache[state] = i

		doMainLoop(chamber, rockShapeGenerator, jetPattern)

		i++
	}
}

func main() {
	input := string(utils.CheckErr(os.ReadFile("input.txt")))

	part1 := getTowerHeightAfterNRocks(input, 2022)

	// Part 2 involves very large numbers of rocks so we need to use the fact that the rock tower
	// has a repeating pattern (since the rock shape and jet patterns repeat). Note that there's an
	// initial region of the tower which is different to the periodic bit later on.
	prePeriodicRockCount, periodicity := findPeriodicity(input)
	prePeriodicHeight := getTowerHeightAfterNRocks(input, prePeriodicRockCount)

	repeats := (1000000000000 - prePeriodicRockCount) / periodicity
	remainder := (1000000000000 - prePeriodicRockCount) % periodicity

	// If we ignore the pre-periodic region, where things are still settling down, then how much
	// height do |n| rocks add?
	getHeightInPeriodicRegion := func(n int) int {
		return getTowerHeightAfterNRocks(input, prePeriodicRockCount+n) - prePeriodicHeight
	}
	periodHeight := getHeightInPeriodicRegion(periodicity)
	remainderHeight := getHeightInPeriodicRegion(remainder)

	part2 := prePeriodicHeight + (periodHeight * repeats) + remainderHeight

	fmt.Printf("The answer to Part 1 is %v.\n", part1)
	fmt.Printf("The answer to Part 2 is %v.\n", part2)
}
