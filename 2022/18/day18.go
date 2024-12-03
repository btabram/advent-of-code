package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"AoC/utils"
)

type Coord struct {
	x, y, z int
}

func (c Coord) getNeighbours() []Coord {
	return []Coord{
		{c.x + 1, c.y, c.z}, {c.x - 1, c.y, c.z},
		{c.x, c.y + 1, c.z}, {c.x, c.y - 1, c.z},
		{c.x, c.y, c.z + 1}, {c.x, c.y, c.z - 1},
	}
}

func makeCoord(str string) Coord {
	split := strings.Split(str, ",")
	return Coord{
		x: utils.CheckErr(strconv.Atoi(split[0])),
		y: utils.CheckErr(strconv.Atoi(split[1])),
		z: utils.CheckErr(strconv.Atoi(split[2])),
	}
}

func main() {
	input := utils.CheckErr(os.ReadFile("input.txt"))

	lavaCubes := []Coord{}
	for _, line := range utils.Lines(string(input)) {
		lavaCubes = append(lavaCubes, makeCoord(line))
	}

	occupied := map[Coord]bool{}
	for _, lavaCube := range lavaCubes {
		occupied[lavaCube] = true
	}

	// For part 1 we count the number of exposed lava cube faces.
	part1 := 0
	for _, lavaCube := range lavaCubes {
		for _, neighbouringCube := range lavaCube.getNeighbours() {
			if !occupied[neighbouringCube] {
				part1 += 1
			}
		}
	}

	xs, ys, zs := []int{}, []int{}, []int{}
	for _, lavaCube := range lavaCubes {
		xs = append(xs, lavaCube.x)
		ys = append(ys, lavaCube.y)
		zs = append(zs, lavaCube.z)
	}
	sort.Ints(xs)
	sort.Ints(ys)
	sort.Ints(zs)

	getBoundedExteriorNeighbours := func(c Coord) []Coord {
		bens := []Coord{}
		for _, n := range c.getNeighbours() {
			if occupied[n] {
				continue
			}
			// Only consider neighbours within a bounding cube slightly bigger than the lava blob.
			if n.x < xs[0]-1 || n.x > xs[len(xs)-1]+1 {
				continue
			}
			if n.y < ys[0]-1 || n.y > ys[len(ys)-1]+1 {
				continue
			}
			if n.z < zs[0]-1 || n.z > zs[len(zs)-1]+1 {
				continue
			}
			bens = append(bens, n)
		}
		return bens
	}

	// For part 2 we count the number of externally exposed lava cube faces. To do this we
	// identify the exterior region outside the lava blob using a recurisve filling function.
	exteriorRegion := map[Coord]bool{}
	var fillExteriorRegion func(c Coord)
	fillExteriorRegion = func(c Coord) {
		exteriorRegion[c] = true
		for _, neighbour := range getBoundedExteriorNeighbours(c) {
			if !exteriorRegion[neighbour] {
				fillExteriorRegion(neighbour)
			}
		}
	}
	fillExteriorRegion(Coord{xs[0], ys[0], zs[0]})

	part2 := 0
	for _, lavaCube := range lavaCubes {
		for _, neighbouringCube := range lavaCube.getNeighbours() {
			if !occupied[neighbouringCube] && exteriorRegion[neighbouringCube] {
				part2 += 1
			}
		}
	}

	fmt.Printf("The answer to Part 1 is %v.\n", part1)
	fmt.Printf("The answer to Part 2 is %v.\n", part2)
}
