package main

import (
	"fmt"
	"math"
	"os"
	"strconv"

	"AoC/utils"
)

type Coord struct {
	x, y int
}

func (c Coord) add(other Coord) Coord {
	return Coord{c.x + other.x, c.y + other.y}
}

func (c Coord) sub(other Coord) Coord {
	return Coord{c.x - other.x, c.y - other.y}
}

type RopeBase struct {
	tailVisited map[Coord]bool
}

func (r *RopeBase) registerTailVisit(c Coord) {
	r.tailVisited[c] = true
}

func (r *RopeBase) countTailVisited() int {
	return len(r.tailVisited)
}

func (r *RopeBase) parseMoveStr(moveStr string) (Coord, int) {
	var directionDelta Coord
	switch moveStr[0] {
	case 'U':
		directionDelta = Coord{1, 0}
	case 'D':
		directionDelta = Coord{-1, 0}
	case 'R':
		directionDelta = Coord{0, 1}
	case 'L':
		directionDelta = Coord{0, -1}
	}
	distance := utils.CheckErr(strconv.Atoi(moveStr[2:]))
	return directionDelta, distance
}

func newRopeBase() *RopeBase {
	return &RopeBase{make(map[Coord]bool)}
}

type ShortRope struct {
	*RopeBase
	head, tail Coord
}

func (r *ShortRope) catchUpTail() {
	gap := r.head.sub(r.tail)
	if dist := math.Sqrt(math.Pow(float64(gap.x), 2) + math.Pow(float64(gap.y), 2)); dist < 2 {
		return // Nothing to do, head and tail are already touching
	}

	sign := func(a int) int {
		switch {
		case a > 0:
			return 1
		case a < 0:
			return -1
		default: // a == 0
			return 0
		}
	}

	r.tail = r.tail.add(Coord{sign(gap.x), sign(gap.y)})
}

func (r *ShortRope) move(moveStr string) {
	directionDelta, distance := r.parseMoveStr(moveStr)
	for step := 0; step < distance; step++ {
		r.head = r.head.add(directionDelta)
		r.catchUpTail()
		r.registerTailVisit(r.tail)
	}
}

func newShortRope() *ShortRope {
	return &ShortRope{
		RopeBase: newRopeBase(),
		head:     Coord{0, 0},
		tail:     Coord{0, 0},
	}
}

type LongRope struct {
	*RopeBase
	segments [9]*ShortRope
}

// We can model our long rope as a series of short rope segments where the head of any given
// segment is the tail of the previous segment.
func (r *LongRope) move(moveStr string) {
	directionDelta, distance := r.parseMoveStr(moveStr)
	for step := 0; step < distance; step++ {
		r.segments[0].head = r.segments[0].head.add(directionDelta)
		r.segments[0].catchUpTail()
		for i := 1; i < 9; i++ {
			r.segments[i].head = r.segments[i-1].tail
			r.segments[i].catchUpTail()
		}
		r.registerTailVisit(r.segments[8].tail)
	}
}

func newLongRope() *LongRope {
	r := new(LongRope)
	r.RopeBase = newRopeBase()
	for i := range r.segments {
		r.segments[i] = newShortRope()
	}
	return r
}

type Rope interface {
	move(moveStr string)
	countTailVisited() int
}

func main() {
	inputLines := utils.Lines(string(utils.CheckErr(os.ReadFile("input.txt"))))
	ropes := []Rope{newShortRope(), newLongRope()}

	for _, inputLine := range inputLines {
		for _, rope := range ropes {
			rope.move(inputLine)
		}
	}

	fmt.Printf("The answer to Part 1 is %v.\n", ropes[0].countTailVisited())
	fmt.Printf("The answer to Part 2 is %v.\n", ropes[1].countTailVisited())
}
