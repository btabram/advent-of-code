package main

import (
	"container/ring"
	"fmt"
	"os"
	"strconv"

	"AoC/utils"
)

type RingItem struct {
	originalPosition int
	value            int
}

func doRound(ring *ring.Ring, size int) {
	for i := 0; i < size; i++ {
		// Find the right item to mix.
		for ring.Value.(RingItem).originalPosition != i {
			ring = ring.Next()
		}
		// We move the item around a number of times given by its value. We ignore the item itself
		// for these moves hence the mod by size-1 rather than size.
		toMoveBy := (ring.Value.(RingItem).value) % (size - 1)
		// Pop out item to move.
		ring = ring.Prev()
		itemToMove := ring.Unlink(1)
		// Do the move.
		ring = ring.Move(toMoveBy)
		// Put the item back.
		ring.Link(itemToMove)
	}
}

func sumGroveCoordinates(ring *ring.Ring) int {
	for ring.Value.(RingItem).value != 0 {
		ring = ring.Next()
	}
	ans := 0
	for i := 0; i < 3; i++ {
		ring = ring.Move(1000)
		ans += ring.Value.(RingItem).value
	}
	return ans
}

func main() {
	inputLines := utils.Lines(string(utils.CheckErr(os.ReadFile("input.txt"))))
	size := len(inputLines)

	part1Ring := ring.New(size)
	part2Ring := ring.New(size)
	for i, line := range inputLines {
		part1Ring.Value = RingItem{i, utils.CheckErr(strconv.Atoi(line))}
		part1Ring = part1Ring.Next()
		// In part 2 the values must be multiplied by the special decryption key.
		part2Ring.Value = RingItem{i, utils.CheckErr(strconv.Atoi(line)) * 811589153}
		part2Ring = part2Ring.Next()
	}

	doRound(part1Ring, size)

	for i := 0; i < 10; i++ {
		doRound(part2Ring, size)
	}

	fmt.Printf("The answer to Part 1 is %v.\n", sumGroveCoordinates(part1Ring))
	fmt.Printf("The answer to Part 2 is %v.\n", sumGroveCoordinates(part2Ring))
}
