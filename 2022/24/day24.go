package main

import (
	"fmt"
	"os"

	"AoC/pathfinding"
	"AoC/utils"
)

type Coord struct {
	x, y int
}

type Blizzards = map[Coord][]byte

type State struct {
	pos  Coord
	time int
}

func main() {
	inputLines := utils.Lines(string(utils.CheckErr(os.ReadFile("input.txt"))))

	// Width and height of the valley, ignoring the walls.
	width := len(inputLines[0]) - 2
	height := len(inputLines) - 2

	// Our coodinate system has (0, 0) at the top left corner of the valley ground.
	start := Coord{0, -1}
	finish := Coord{width - 1, height}

	initialBlizzards := Blizzards{}
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			if char := inputLines[y+1][x+1]; char != '.' {
				initialBlizzards[Coord{x, y}] = []byte{char}
			}
		}
	}

	// Store the blizzard positions for a given timestep.
	blizzardsMemoMap := map[int]Blizzards{}
	blizzardsMemoMap[0] = initialBlizzards

	getNextMoves := func(s State) map[State]int {
		possibleMoves := []State{
			{s.pos, s.time + 1},                       // Wait in place
			{Coord{s.pos.x + 1, s.pos.y}, s.time + 1}, // Move down
			{Coord{s.pos.x - 1, s.pos.y}, s.time + 1}, // Move up
			{Coord{s.pos.x, s.pos.y + 1}, s.time + 1}, // Move left
			{Coord{s.pos.x, s.pos.y - 1}, s.time + 1}, // Move right
		}

		nextBlizzards := Blizzards{}
		if b, ok := blizzardsMemoMap[s.time+1]; ok {
			nextBlizzards = b
		} else {
			// We need to calculate the position of the blizzards next step. Note that an entry for
			// the current blizzards must always exist in the map.
			for pos, dirs := range blizzardsMemoMap[s.time] {
				for _, dir := range dirs {
					var newPos Coord
					switch dir {
					case '>':
						newPos = Coord{(pos.x + 1) % width, pos.y}
					case 'v':
						newPos = Coord{pos.x, (pos.y + 1) % height}
					case '<':
						newPos = Coord{(pos.x - 1 + width) % width, pos.y}
					case '^':
						newPos = Coord{pos.x, (pos.y - 1 + height) % height}
					}
					nextBlizzards[newPos] = append(nextBlizzards[newPos], dir)
				}
			}
			blizzardsMemoMap[s.time+1] = nextBlizzards
		}

		validMoves := map[State]int{}
		for _, move := range possibleMoves {
			// The start and finish are special gaps in the wall which never have blizzards.
			isAlwaysAllowed := move.pos == start || move.pos == finish
			if !isAlwaysAllowed {
				// The valley is surrounded by a wall that block movement.
				if move.pos.x < 0 || move.pos.x >= width {
					continue
				}
				if move.pos.y < 0 || move.pos.y >= height {
					continue
				}
				// We can't move to positions not occupied by blizzards.
				if _, ok := nextBlizzards[move.pos]; ok {
					continue
				}
			}
			validMoves[move] = 1 // All moves take 1 minute
		}
		return validMoves
	}

	isStart := func(s State) bool { return s.pos == start }
	isFinish := func(s State) bool { return s.pos == finish }

	// For part 1 we need to go from the start to the finish.
	part1 := pathfinding.Dijkstra(State{start, 0}, isFinish, getNextMoves)
	fmt.Printf("The answer to Part 1 is %v.\n", part1)

	// For part 2 we need to go start -> finish -> start -> finish.
	backToStart := pathfinding.Dijkstra(State{finish, part1}, isStart, getNextMoves)
	toFinishAgain := pathfinding.Dijkstra(State{start, part1 + backToStart}, isFinish, getNextMoves)
	part2 := part1 + backToStart + toFinishAgain
	fmt.Printf("The answer to Part 2 is %v.\n", part2)
}
