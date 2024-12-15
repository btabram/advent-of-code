package main

import (
	utils "aoc2024"
	"fmt"
	"strings"
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

type warehouse struct {
	robot        vec
	grid         map[vec]rune
	instructions string
	cursor       int
}

func loadWarehouse(lines []string) *warehouse {
	w := warehouse{grid: make(map[vec]rune)}

	for y, line := range lines {
		// The instructions comes after the grid, separated by a blank line
		if len(line) == 0 {
			// They may be spread over multiple lines but the linebreaks don't matter
			for _, l := range lines[y+1:] {
				w.instructions += l
			}
			break
		}
		// Track all features on the grid, don't care about empty squares
		for x, char := range line {
			if char == '@' {
				w.robot = vec{x, y}
			} else if char != '.' {
				w.grid[vec{x, y}] = char
			}
		}
	}

	return &w
}

func (w *warehouse) popNextMove() vec {
	next := rune(w.instructions[w.cursor])
	w.cursor++

	switch next {
	case '^':
		return vec{x: 0, y: -1} // Positive y is down
	case 'v':
		return vec{x: 0, y: 1}
	case '<':
		return vec{x: -1, y: 0}
	case '>':
		return vec{x: 1, y: 0}
	}

	panic("Bad instruction")
}

func (w *warehouse) doMove() {
	move := w.popNextMove()

	nextRobotPos := w.robot.add(move)

	occupied := w.grid[nextRobotPos]
	if occupied == 0 {
		// Moving into empty square
		w.robot = nextRobotPos
		return
	}

	if occupied == '#' {
		// Trying to move directly into a wall. Nothing happens
		return
	}

	// Must be moving into a box. Need to work out if we can move it, and which other boxes move too
	boxesToMove := make(map[vec]rune)

	// Track the square that the boxes will need to move into, but which we haven't resolved yet
	squaresToMoveInto := []vec{nextRobotPos}
	for {
		nextPos := squaresToMoveInto[0]
		squaresToMoveInto = squaresToMoveInto[1:]

		if boxesToMove[nextPos] != 0 {
			continue // Already resolved as a (part of a) box to move, skip
		}

		nextOccupied := w.grid[nextPos]
		if nextOccupied == 0 {
			// Empty square, can move into here and doesn't trigger any more moves
			if len(squaresToMoveInto) == 0 {
				break
			} else {
				continue
			}
		}

		if nextOccupied == '#' {
			// One of the boxes is up against a wall and can't move. Nothing happens
			return
		}

		// Must be another box, add it to the line of boxes we're trying to move
		boxesToMove[nextPos] = nextOccupied
		squaresToMoveInto = append(squaresToMoveInto, nextPos.add(move))

		// In part 2 the boxes are two squares wide and we need to track both halves
		if nextOccupied == '[' {
			otherHalf := nextPos.add(vec{1, 0})
			boxesToMove[otherHalf] = ']'
			squaresToMoveInto = append(squaresToMoveInto, otherHalf.add(move))
		} else if nextOccupied == ']' {
			otherHalf := nextPos.add(vec{-1, 0})
			boxesToMove[otherHalf] = '['
			squaresToMoveInto = append(squaresToMoveInto, otherHalf.add(move))
		}
	}

	// Move the whole stack of boxes (and the robot)
	for pos := range boxesToMove {
		delete(w.grid, pos)
	}
	for pos, box := range boxesToMove {
		w.grid[pos.add(move)] = box
	}
	w.robot = nextRobotPos
}

// Executes all robot instructions and then calculates the sum of the box GPS coordinates
func (w *warehouse) run() int {
	for w.cursor < len(w.instructions) {
		w.doMove()
	}

	gpsSum := 0
	for pos, char := range w.grid {
		if char == 'O' || char == '[' {
			gpsSum += 100*pos.y + pos.x
		}
	}
	return gpsSum
}

func (w *warehouse) print() {
	var width, height int
	for p := range w.grid {
		if p.x >= width {
			width = p.x + 1
		}
		if p.y >= height {
			height = p.y + 1
		}
	}

	str := ""
	for y := range height {
		for x := range width {
			obj := w.grid[vec{x, y}]
			if x == w.robot.x && y == w.robot.y {
				str += "@"
			} else if obj != 0 {
				str += string(obj)
			} else {
				str += "."
			}
		}
		str += "\n"
	}
	fmt.Println(str)
}

func main() {
	lines := utils.ReadLines("input.txt")

	warehouse1 := loadWarehouse(lines)
	part1 := warehouse1.run()

	// In the part 2 everything is twice as wide, except the robot doing the pushing. We can simply
	// replace characters in the input to get the wider warehouse we need to model.
	part2Input := make([]string, len(lines))
	for i, line := range lines {
		part2Input[i] =
			strings.ReplaceAll(
				strings.ReplaceAll(
					strings.ReplaceAll(
						strings.ReplaceAll(line, "#", "##"),
						"O", "[]"),
					".", ".."),
				"@", "@.")
	}
	warehouse2 := loadWarehouse(part2Input)
	part2 := warehouse2.run()
	warehouse2.print()

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
