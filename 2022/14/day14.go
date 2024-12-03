package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"AoC/utils"
)

type Coord struct {
	x, y int
}

func makeCoord(str string) Coord {
	split := strings.Split(str, ",")
	return Coord{
		x: utils.CheckErr(strconv.Atoi(split[0])),
		y: utils.CheckErr(strconv.Atoi(split[1])),
	}
}

var sandSpawnLocation Coord = Coord{500, 0}

func getPossibleSandMoves(currentPos Coord) []Coord {
	// Sand tries to move down. If it can't, then it moves down & left, or finally down & right.
	return []Coord{
		{currentPos.x, currentPos.y + 1},
		{currentPos.x - 1, currentPos.y + 1},
		{currentPos.x + 1, currentPos.y + 1},
	}
}

func dropSand(occupied map[Coord]bool, fallthroughHeight int) *Coord {
	fallingSand := sandSpawnLocation
	for {
		cantMove := true
		for _, move := range getPossibleSandMoves(fallingSand) {
			if !occupied[move] {
				cantMove = false
				fallingSand = move
				break
			}
		}
		if cantMove {
			return &fallingSand // Sand has come to rest, return final position
		}
		if fallingSand.y > fallthroughHeight {
			return nil // Fallen into the abyss and not coming to rest, no final position
		}
	}
}

func part1(rocks map[Coord]bool) int {
	occupiedPoints := map[Coord]bool{}
	lowestRockHeight := 0
	for rock := range rocks {
		occupiedPoints[rock] = true
		if rock.y > lowestRockHeight {
			lowestRockHeight = rock.y
		}
	}

	sandComingToRestCount := 0
	for {
		restLocation := dropSand(occupiedPoints, lowestRockHeight)
		if restLocation == nil {
			return sandComingToRestCount
		}
		occupiedPoints[*restLocation] = true
		sandComingToRestCount++
	}
}

func part2(rocks map[Coord]bool) int {
	occupiedPoints := map[Coord]bool{}
	lowestRockHeight := 0
	for rock := range rocks {
		occupiedPoints[rock] = true
		if rock.y > lowestRockHeight {
			lowestRockHeight = rock.y
		}
	}

	// Part 2 has a solid floor rather than an endless void.
	floorHeight := lowestRockHeight + 2
	for x := -10000; x < 10000; x++ { // Not exactly infinite but this should be enough
		occupiedPoints[Coord{x, floorHeight}] = true
	}

	sandComingToRestCount := 0
	for {
		restLocation := dropSand(occupiedPoints, floorHeight+1)
		occupiedPoints[*restLocation] = true
		sandComingToRestCount++
		if *restLocation == sandSpawnLocation {
			return sandComingToRestCount
		}
	}
}

func main() {
	input := utils.CheckErr(os.ReadFile("input.txt"))

	rocks := map[Coord]bool{}
	for _, line := range utils.Lines(string(input)) {
		pathPoints := utils.Transform(strings.Split(line, " -> "), makeCoord)
		prevPathPoint := pathPoints[0]
		for i := 1; i < len(pathPoints); i++ {
			nextPathPoint := pathPoints[i]
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
			// Note that one of these will always be zero because lines are never diagonal.
			unitDeltaX := sign(nextPathPoint.x - prevPathPoint.x)
			unitDeltaY := sign(nextPathPoint.y - prevPathPoint.y)
			point := prevPathPoint
			for point != nextPathPoint {
				rocks[point] = true
				point = Coord{point.x + unitDeltaX, point.y + unitDeltaY}
			}
			rocks[nextPathPoint] = true
			prevPathPoint = nextPathPoint
		}
	}

	fmt.Printf("The answer to Part 1 is %v.\n", part1(rocks))
	fmt.Printf("The answer to Part 2 is %v.\n", part2(rocks))
}
