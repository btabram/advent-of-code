package main

import (
	utils "aoc2024"
	"fmt"
	"math"
	"regexp"
)

type vec struct {
	x, y int
}

var buttonRegex = regexp.MustCompile(`Button [AB]: X\+(\d+), Y\+(\d+)`)
var prizeRegex = regexp.MustCompile(`Prize: X=(\d+), Y=(\d+)`)

func vecFromMatch(match []string) vec {
	return vec{
		x: utils.Int(match[1]),
		y: utils.Int(match[2]),
	}
}

type clawMachine struct {
	buttonA, buttonB vec
	prize            vec
}

/*
The problem is just simultaneous linear equations, we have two equations and two unknowns. There's
only one solution for every machine, but we're only interested in machines that can be solved by
integer values of `a` and `b`.

Here's one of the examples as equations. We can rearrange to solve:

	94a + 22b = 8400 [x equation]
	34a + 67b = 5400 [y equation]

	a = (8400 - 22b)/94

	34*(8400 - 22b)/94 + 67b = 5400

	(67 - (22*34/94))b = (5400 - (8400*34/94))

	b = (5400 - (8400*34/94)) / (67 - (22*34/94))

From this example we can write out a general solution:

	(ax * a) + (bx * b) = px
	(ay * a) + (bx * b) = px

	a := (px - (bx * b)) / ax
	b := (py - (px * ay / ax)) / (by - (bx * ay / ax))
*/
func (m clawMachine) solve() (float64, float64) {
	ax := float64(m.buttonA.x)
	ay := float64(m.buttonA.y)
	bx := float64(m.buttonB.x)
	by := float64(m.buttonB.y)
	px := float64(m.prize.x)
	py := float64(m.prize.y)

	b := (py - (px * ay / ax)) / (by - (bx * ay / ax))
	a := (px - (bx * b)) / ax

	return a, b
}

func (m clawMachine) test(a, b int) bool {
	return ((a*m.buttonA.x)+(b*m.buttonB.x) == m.prize.x) &&
		((a*m.buttonA.y)+(b*m.buttonB.y) == m.prize.y)
}

func main() {
	lines := utils.ReadLines("input.txt")

	var machines []clawMachine
	for i := 0; i < len(lines); i += 4 {
		a := buttonRegex.FindStringSubmatch(lines[i])
		b := buttonRegex.FindStringSubmatch(lines[i+1])
		prize := prizeRegex.FindStringSubmatch(lines[i+2])

		machines = append(machines, clawMachine{
			buttonA: vecFromMatch(a),
			buttonB: vecFromMatch(b),
			prize:   vecFromMatch(prize),
		})
	}

	tokensToGetAllPrizes := func(offset int) int {
		tokens := 0
		for _, m := range machines {
			// Note that `m` is a copy so this doesn't mutate `machines`
			m.prize.x += offset
			m.prize.y += offset

			aFloat, bFloat := m.solve()

			// Due to floating point errors it's likely that `aFloat` and `bFloat` are not exactly
			// integers so we round to the nearest integer and then see if these integers are the
			// correct solution. We're relying on the floating point errors staying small here, but
			// it does work out as long as we use `float64`. In part 2 we have to deal with a large
			// range of magnitudes which can lead to significant floating point errors unless the
			// calculation is done carefully (which it isn't!) - so this code fails with `float32`
			a := int(math.Round(aFloat))
			b := int(math.Round(bFloat))

			if m.test(a, b) {
				tokens += 3*int(a) + int(b) // Pressing button A costs 3 tokens and button B costs 1
			}
		}
		return tokens
	}

	part1 := tokensToGetAllPrizes(0)
	part2 := tokensToGetAllPrizes(10_000_000_000_000)

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
