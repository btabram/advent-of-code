package main

import (
	utils "aoc2024"
	"fmt"
	"strings"
)

type equation struct {
	answer int
	inputs []int
}

type operator func(int, int) int

func (eq *equation) canSolve(operators []operator) bool {

	// Use recursion to try all possible operator combinations
	var trySolve func(int, []int) bool
	trySolve = func(value int, inputs []int) bool {
		// Early return to save some work, already overshot
		if value > eq.answer {
			return false
		}

		// End case, time to check the answer
		if len(inputs) == 0 {
			return value == eq.answer
		}

		// Branch, trying all available operators
		for _, op := range operators {
			if trySolve(op(value, inputs[0]), inputs[1:]) {
				return true
			}
		}
		return false
	}

	return trySolve(eq.inputs[0], eq.inputs[1:])
}

func add(a, b int) int {
	return a + b
}

func mul(a, b int) int {
	return a * b
}

func concat(a, b int) int {
	return utils.Int(fmt.Sprintf("%v%v", a, b))
}

func main() {
	lines := utils.ReadLines("input.txt")

	equations := make([]equation, len(lines))
	for i, line := range lines {
		// Equation lines are of the form "292: 11 6 16 20"
		split := strings.Split(line, ":")
		equations[i] = equation{
			answer: utils.Int(split[0]),
			inputs: utils.Ints(split[1]),
		}
	}

	var part1, part2 int
	for _, eq := range equations {
		if eq.canSolve([]operator{add, mul}) {
			part1 += eq.answer
		}
		// In part2 we have an extra operator to try
		if eq.canSolve([]operator{add, mul, concat}) {
			part2 += eq.answer
		}
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
