package main

import (
	utils "aoc2024"
	"fmt"
)

// Value used to represent EMPTY blocks in memory, file IDs are always non-negative
const EMPTY = -1

func checkSum(memory []int) int {
	checksum := 0
	for i, fileId := range memory {
		if fileId != EMPTY {
			checksum += i * fileId
		}
	}
	return checksum
}

func main() {
	diskMap := utils.Read("input.txt")

	var initialMemory, fileSizes []int
	for i, char := range diskMap {
		length := int(char - '0')
		if i&1 == 1 {
			// Odd-indexed values in the disk map are free space
			for range length {
				initialMemory = append(initialMemory, EMPTY)
			}
		} else {
			// Even-indexed values in the disk map are files
			fileId := len(fileSizes) // The first file is 0, then 1, 2 etc.
			for range length {
				initialMemory = append(initialMemory, fileId)
			}
			fileSizes = append(fileSizes, length)
		}
	}

	var part1CondensedMemory []int
	remaining := initialMemory // Not mutating underlying array, so just need a new slice (no copy)
	for len(remaining) != 0 {
		next := remaining[0]
		if next != EMPTY {
			// Next block is a file, nothing special to do
			part1CondensedMemory = append(part1CondensedMemory, next)
			remaining = remaining[1:]
		} else {
			// Next block is empty space, move the last file block here instead
			lastFileBlock := remaining[len(remaining)-1]
			for lastFileBlock == EMPTY {
				remaining = remaining[:len(remaining)-1]
				lastFileBlock = remaining[len(remaining)-1]
			}

			part1CondensedMemory = append(part1CondensedMemory, lastFileBlock)
			remaining = remaining[1 : len(remaining)-1]
		}
	}

	// In part 2 we try to move each file into the first suitable run of free space
	// We move the files in order of decreasing file ID, and only try each file once
	part2CondensedMemory := initialMemory
	for fileId := len(fileSizes) - 1; fileId >= 0; fileId-- {
		fileSize := fileSizes[fileId]

		// This nested looping over all memory isn't very efficient, but is fine enough for my input
		emptyRunLength := 0
		for i, block := range part2CondensedMemory {
			if block == fileId {
				break // We've searched all empty space before the file, can't move this one
			}

			if block == EMPTY {
				emptyRunLength++
			} else {
				emptyRunLength = 0
			}

			if emptyRunLength == fileSize {
				// We've found suitable empty space to move the file to!
				for j, block := range part2CondensedMemory {
					if j > i-fileSize && j <= i {
						part2CondensedMemory[j] = fileId // Copy file to new location
					}
					if block == fileId {
						part2CondensedMemory[j] = EMPTY // Empty the old location
					}
				}
				break
			}
		}
	}

	fmt.Printf("The answer to Part 1 is %v\n", checkSum(part1CondensedMemory))
	fmt.Printf("The answer to Part 2 is %v\n", checkSum(part2CondensedMemory))
}
