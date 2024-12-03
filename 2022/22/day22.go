package main

import (
	"fmt"
	"os"
	"strconv"

	"AoC/utils"
)

type Coord struct {
	x, y int
}

func (c Coord) add(delta Coord) Coord {
	return Coord{c.x + delta.x, c.y + delta.y}
}

const (
	east = iota
	south
	west
	north
)

var facingVectors = [4]Coord{
	{1, 0},
	{0, 1},
	{-1, 0},
	{0, -1},
}

type Position struct {
	coord  Coord
	facing int
}

type Move interface {
	do(pos Position) Position
}

type Turn struct {
	clockwise bool
}

func (t Turn) do(p Position) Position {
	var newFacing int
	if t.clockwise {
		newFacing = (p.facing + 1) % 4
	} else {
		newFacing = (p.facing - 1 + 4) % 4
	}
	return Position{p.coord, newFacing}
}

type Forwards struct {
}

func (f Forwards) do(p Position) Position {
	change := facingVectors[p.facing]
	return Position{p.coord.add(change), p.facing}
}

func main() {
	inputLines := utils.Lines(string(utils.CheckErr(os.ReadFile("input.txt"))))

	boardStrs, movesStrs := utils.SplitAt(inputLines, "")
	movesStr := movesStrs[0]

	board := map[Coord]byte{}
	var initialPosition *Coord
	for y := 1; y <= len(boardStrs); y++ {
		boardStr := boardStrs[y-1]
		for x := 1; x <= len(boardStr); x++ {
			if char := boardStr[x-1]; char != ' ' {
				board[Coord{x, y}] = char
				if initialPosition == nil {
					initialPosition = &Coord{x, y}
				}
			}
		}
	}

	moves := []Move{}
	numberStr := ""
	for i := range movesStr {
		switch char := string(movesStr[i]); char {
		case "L":
			moves = append(moves, Turn{false})
		case "R":
			moves = append(moves, Turn{true})
		default:
			numberStr += char
			if i == len(movesStr)-1 || movesStr[i+1] == 'L' || movesStr[i+1] == 'R' {
				tilesToMove := utils.CheckErr(strconv.Atoi(numberStr))
				for t := 0; t < tilesToMove; t++ {
					moves = append(moves, Forwards{})
				}
				numberStr = ""
			}
		}
	}

	/*
		For part 2 we conceptually fold the board up into a cube and move around the faces of the
		cube. I can't think of a clever general way to do the edge transitions so I'm just doing a
		solution for my board shape by working out the edge pairs and creating a map of points to
		wrap.

		Sketch of my input net with edges labelled:
		          A   B
		        11112222
		      C 11112222 D
		        11112222
		        11112222
		        3333  E
		      F 3333 G
		        3333
		     H  3333
		    44445555
		  I 44445555 J
		    44445555
		    44445555
		    6666  K
		  L 6666 M
		    6666
		    6666
		     N

		Edge pairs:
		A <-> L
		B <-> N
		C <-> I
		D <-> J
		E <-> G
		F <-> H
		K <-> M
	*/
	part2EdgeWrappings := map[Position]Position{}
	for i := 1; i <= 50; i++ {
		// A -> L
		part2EdgeWrappings[Position{Coord{50 + i, 0}, north}] = Position{Coord{1, 150 + i}, east}
		// L -> A
		part2EdgeWrappings[Position{Coord{0, 150 + i}, west}] = Position{Coord{50 + i, 1}, south}

		// B -> N
		part2EdgeWrappings[Position{Coord{100 + i, 0}, north}] = Position{Coord{i, 200}, north}
		// N -> B
		part2EdgeWrappings[Position{Coord{i, 200 + 1}, south}] = Position{Coord{100 + i, 1}, south}

		// C -> I
		part2EdgeWrappings[Position{Coord{50, i}, west}] = Position{Coord{1, 151 - i}, east}
		// I -> C
		part2EdgeWrappings[Position{Coord{0, 151 - i}, west}] = Position{Coord{50 + 1, i}, east}

		// D -> J
		part2EdgeWrappings[Position{Coord{151, i}, east}] = Position{Coord{100, 151 - i}, west}
		// J -> D
		part2EdgeWrappings[Position{Coord{101, 151 - i}, east}] = Position{Coord{150, i}, west}

		// E -> G
		part2EdgeWrappings[Position{Coord{100 + i, 50 + 1}, south}] = Position{Coord{100, 50 + i}, west}
		// G -> E
		part2EdgeWrappings[Position{Coord{101, 50 + i}, east}] = Position{Coord{100 + i, 50}, north}

		// F -> H
		part2EdgeWrappings[Position{Coord{50, 50 + i}, west}] = Position{Coord{i, 101}, south}
		// H -> F
		part2EdgeWrappings[Position{Coord{i, 100}, north}] = Position{Coord{50 + 1, 50 + i}, east}

		// K -> M
		part2EdgeWrappings[Position{Coord{50 + i, 151}, south}] = Position{Coord{50, 150 + i}, west}
		// M -> J
		part2EdgeWrappings[Position{Coord{50 + 1, 150 + i}, east}] = Position{Coord{50 + i, 150}, north}
	}

	part1Pos := Position{*initialPosition, east}
	part2Pos := Position{*initialPosition, east}
	for _, move := range moves {
		newPart1Pos := move.do(part1Pos)
		newPart2Pos := move.do(part2Pos)

		// In part 1 we simply wrap around to the opposite side of the board when we go off the edge.
		if _, ok := board[newPart1Pos.coord]; !ok {
			// New position isn't on the board. Turn around and go backwards...
			newPart1Pos = Turn{true}.do(Turn{true}.do(newPart1Pos))
			for {
				newPart1Pos = Forwards{}.do(newPart1Pos)
				if _, ok := board[newPart1Pos.coord]; !ok {
					break
				}
			}
			// ...until we go off the other side. Then turn back and go forward onto the board.
			newPart1Pos = Forwards{}.do(Turn{true}.do(Turn{true}.do(newPart1Pos)))
		}
		// In part 2 the board is the net of a cube and when we go off the edge we move as if we
		// travelling around the faces on the cube.
		if wrappedPos, ok := part2EdgeWrappings[newPart2Pos]; ok {
			newPart2Pos = wrappedPos
		}

		// Only make the move if we're moving into an empty space (not a wall).
		if board[newPart1Pos.coord] == '.' {
			part1Pos = newPart1Pos
		}
		if board[newPart2Pos.coord] == '.' {
			part2Pos = newPart2Pos
		}
	}

	part1Password := (1000 * part1Pos.coord.y) + (4 * part1Pos.coord.x) + part1Pos.facing
	part2Password := (1000 * part2Pos.coord.y) + (4 * part2Pos.coord.x) + part2Pos.facing

	fmt.Printf("The answer to Part 1 is %v.\n", part1Password)
	fmt.Printf("The answer to Part 2 is %v.\n", part2Password)
}
