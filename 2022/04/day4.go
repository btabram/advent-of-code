package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"AoC/utils"
)

type Range struct {
	min int
	max int
}

func (r Range) includes(val int) bool {
	return val >= r.min && val <= r.max
}

func (r Range) contains(other Range) bool {
	return r.includes(other.min) && r.includes(other.max)
}

func (r Range) overlaps(other Range) bool {
	return r.includes(other.min) || r.includes(other.max)
}

func makeRange(rangeStr string) Range {
	ns := strings.Split(rangeStr, "-")
	min := utils.CheckErr(strconv.Atoi(ns[0]))
	max := utils.CheckErr(strconv.Atoi(ns[1]))
	return Range{min, max}
}

func main() {
	input := utils.CheckErr(os.ReadFile("input.txt"))

	var part1, part2 int
	for _, line := range utils.Lines(string(input)) {
		pair := strings.Split(line, ",")
		r1 := makeRange(pair[0])
		r2 := makeRange(pair[1])
		if r1.contains(r2) || r2.contains(r1) {
			part1 += 1
		}
		if r1.overlaps(r2) || r2.overlaps(r1) {
			part2 += 1
		}
	}

	fmt.Printf("The answer to Part 1 is %v.\n", part1)
	fmt.Printf("The answer to Part 2 is %v.\n", part2)
}
