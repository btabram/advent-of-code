package main

import (
	utils "aoc2024"
	"fmt"
)

func get(charMap [][]rune, x, y int) rune {
	if x < 0 || x >= len(charMap[0]) || y < 0 || y >= len(charMap) {
		return 0
	}
	return charMap[x][y]
}

func main() {
	lines := utils.ReadLines("input.txt")

	charMap := make([][]rune, len(lines[0]))
	for i, line := range lines {
		charMap[i] = []rune(line)
	}

	var part1 int
	var part2 int
	for x, row := range charMap {
		for y, c := range row {
			// Part 1 - look for "XMAS" in all directions
			if c == 'X' {
				for dx := -1; dx <= 1; dx++ {
					for dy := -1; dy <= 1; dy++ {
						if dx == 0 && dy == 0 {
							continue
						}

						// Check the next characters in the chosen direction
						spellsXmas := get(charMap, x+dx, y+dy) == 'M' &&
							get(charMap, x+(2*dx), y+(2*dy)) == 'A' &&
							get(charMap, x+(3*dx), y+(3*dy)) == 'S'

						if spellsXmas {
							part1++
						}
					}
				}
			}

			// Part 2 - look for "MAS" in an X shape
			if c == 'A' {
				ne := get(charMap, x+1, y+1)
				se := get(charMap, x+1, y-1)
				nw := get(charMap, x-1, y+1)
				sw := get(charMap, x-1, y-1)

				firstDiagonalSpellsMas :=
					(ne == 'M' && sw == 'S') || (ne == 'S' && sw == 'M')
				secondDiagonalSpellsMas :=
					(nw == 'M' && se == 'S') || (nw == 'S' && se == 'M')

				if firstDiagonalSpellsMas && secondDiagonalSpellsMas {
					part2++
				}
			}
		}
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
