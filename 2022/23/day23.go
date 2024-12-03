package main

import (
	"fmt"
	"os"
	"sort"

	"AoC/utils"
)

const (
	north = iota
	south
	west
	east
)

type Coord struct {
	x, y int
}

func (c Coord) move(direction int) Coord {
	switch direction {
	case north:
		return Coord{c.x, c.y - 1}
	case east:
		return Coord{c.x + 1, c.y}
	case south:
		return Coord{c.x, c.y + 1}
	case west:
		return Coord{c.x - 1, c.y}
	default:
		panic(fmt.Sprintf("Invalid direction: %v", direction))
	}
}

type Elf struct {
	pos     Coord
	nextPos *Coord
}

func (elf *Elf) getNeighbourInfo(elfPositions map[Coord]bool) uint8 {
	e := elf.pos
	neighbourCoords := [8]Coord{
		{e.x - 1, e.y - 1}, // NW
		{e.x, e.y - 1},     // N
		{e.x + 1, e.y - 1}, // NE
		{e.x + 1, e.y},     // E
		{e.x + 1, e.y + 1}, // SE
		{e.x, e.y + 1},     // S
		{e.x - 1, e.y + 1}, // SW
		{e.x - 1, e.y},     // W
	}
	info := uint8(0)
	for i, neighbourCoord := range neighbourCoords {
		if _, ok := elfPositions[neighbourCoord]; ok {
			info |= 1 << i
		}
	}
	return info
}

func main() {
	input := utils.CheckErr(os.ReadFile("input.txt"))

	elves := []*Elf{}
	for y, line := range utils.Lines(string(input)) {
		for x, char := range line {
			if char == '#' {
				elves = append(elves, &Elf{pos: Coord{x, y}})
			}
		}
	}

	firstDirection := north
	n := 0
	for {
		elfPositions := map[Coord]bool{}
		for _, elf := range elves {
			elfPositions[elf.pos] = true
		}

		// Elves choose moves.
		for _, elf := range elves {
			neighbourInfo := elf.getNeighbourInfo(elfPositions)
			if neighbourInfo == 0 {
				continue // No neighbours means no move
			}
			direction := firstDirection
			for i := 0; i < 4; i++ {
				var neighbourMask uint8
				switch direction {
				case north:
					neighbourMask = 0b00000111
				case east:
					neighbourMask = 0b00011100
				case south:
					neighbourMask = 0b01110000
				case west:
					neighbourMask = 0b11000001
				}
				// Elves only move in directions without neighbours.
				if (neighbourInfo & neighbourMask) == 0 {
					nextPos := elf.pos.move(direction)
					elf.nextPos = &nextPos
					break
				}
				direction = (direction + 1) % 4
			}
		}

		// Elves only do their moves if no other elf tries to move into the same square.
		nextPositions := map[Coord]*Elf{}
		for _, elf := range elves {
			if elf.nextPos != nil {
				if other, ok := nextPositions[*elf.nextPos]; ok {
					elf.nextPos = nil
					other.nextPos = nil
				} else {
					nextPositions[*elf.nextPos] = elf
				}
			}
		}

		// Do the valid moves.
		thisRoundHasChanges := false
		for _, elf := range elves {
			if elf.nextPos != nil {
				thisRoundHasChanges = true
				elf.pos = *elf.nextPos
				elf.nextPos = nil
			}
		}

		n++
		firstDirection = (firstDirection + 1) % 4 // The first direction changes every round

		if n == 10 {
			xs, ys := []int{}, []int{}
			for _, elf := range elves {
				xs = append(xs, elf.pos.x)
				ys = append(ys, elf.pos.y)
			}
			sort.Ints(xs)
			sort.Ints(ys)
			containingRectangleArea := (xs[len(xs)-1] - xs[0] + 1) * (ys[len(ys)-1] - ys[0] + 1)
			fmt.Printf("The answer to Part 1 is %v.\n", containingRectangleArea-len(elves))
		}

		if !thisRoundHasChanges {
			fmt.Printf("The answer to Part 2 is %v.\n", n)
			break
		}
	}
}
