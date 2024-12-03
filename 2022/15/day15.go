package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"AoC/utils"
)

type Coord struct {
	x, y int
}

func (a Coord) dist(b Coord) int {
	return utils.Abs(b.x-a.x) + utils.Abs(b.y-a.y)
}

type Sensor struct {
	pos        Coord
	beacon     Coord
	beaconDist int
}

func (s Sensor) cantBeBeacon(c Coord) bool {
	return s.pos.dist(c) <= s.beaconDist && c != s.beacon
}

func (s Sensor) getOuterBorder() []Coord {
	outerBorderDist := s.beaconDist + 1

	borderLines := []struct {
		start  Coord
		dX, dY int
	}{
		{start: Coord{s.pos.x + outerBorderDist, s.pos.y}, dX: -1, dY: -1},
		{start: Coord{s.pos.x, s.pos.y + outerBorderDist}, dX: -1, dY: -1},
		{start: Coord{s.pos.x - outerBorderDist, s.pos.y}, dX: 1, dY: -1},
		{start: Coord{s.pos.x, s.pos.y - outerBorderDist}, dX: 1, dY: 1},
	}

	borderPoints := make([]Coord, 0, 4*outerBorderDist)
	for _, line := range borderLines {
		point := line.start
		for i := 0; i < outerBorderDist; i++ {
			borderPoints = append(borderPoints, point)
			point = Coord{point.x + line.dX, point.y + line.dY}
		}
	}

	return borderPoints
}

func main() {
	input := utils.CheckErr(os.ReadFile("input.txt"))
	inputRegex := regexp.MustCompile(
		"^Sensor at x=(-?[0-9]+), y=(-?[0-9]+): closest beacon is at x=(-?[0-9]+), y=(-?[0-9]+)$",
	)
	sensors := []Sensor{}
	for _, line := range utils.Lines(string(input)) {
		matches := inputRegex.FindStringSubmatch(line)
		pos := Coord{
			x: utils.CheckErr(strconv.Atoi(matches[1])),
			y: utils.CheckErr(strconv.Atoi(matches[2])),
		}
		beacon := Coord{
			x: utils.CheckErr(strconv.Atoi(matches[3])),
			y: utils.CheckErr(strconv.Atoi(matches[4])),
		}
		sensors = append(sensors, Sensor{pos, beacon, pos.dist(beacon)})
	}

	minX, maxX := 0, 0
	for _, s := range sensors {
		if s.pos.x-s.beaconDist < minX {
			minX = s.pos.x - s.beaconDist
		}
		if s.pos.x+s.beaconDist > maxX {
			maxX = s.pos.x + s.beaconDist
		}
	}

	// For part 1 we need the number of positions along y=2000000 where a beacon cannot be.
	part1 := 0
	for x := minX; x <= maxX; x++ {
		for _, s := range sensors {
			if s.cantBeBeacon(Coord{x, 2000000}) {
				part1++
				break
			}
		}
	}

	// For part 2 we have to find the one point with x and y between 0 and 4000000 which could be a
	// beacon. There's too many points to test them all but we can reason that the one valid point
	// must be just outside the region where a beacon cannot be around a sensor. We can calculate
	// all of these points at distances |Sensor.beacondist + 1| from the sensors and it is then
	// a manageable number of points to test them all.
	var part2 *Coord
	for _, s := range sensors {
		for _, c := range s.getOuterBorder() {
			if c.x < 0 || c.y < 0 || c.x > 4000000 || c.y > 4000000 {
				continue
			}
			couldBeBeacon := true
			for _, s := range sensors {
				if s.cantBeBeacon(c) {
					couldBeBeacon = false
					break
				}
			}
			if couldBeBeacon {
				// Avoid taking a reference to the loop variable because its value will change.
				cInstance := c
				part2 = &cInstance
				break
			}
		}
		if part2 != nil {
			break
		}
	}

	fmt.Printf("The answer to Part 1 is %v.\n", part1)
	fmt.Printf("The answer to Part 2 is %v.\n", 4000000*part2.x+part2.y)
}
