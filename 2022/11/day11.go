package main

import (
	"fmt"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"AoC/utils"
)

type Monkey struct {
	id              int
	items           []int
	operation       func(int) int
	testDivisor     int
	trueThrowTo     int
	falseThrowTo    int
	inspectionCount int
}

func parseMonkeys(lines []string) []*Monkey {
	monkeys := []*Monkey{}
	for i := 0; i < len(lines); i++ {
		idStr := regexp.MustCompile("^Monkey ([0-9]+):$").FindStringSubmatch(lines[i])[1]
		monkey := Monkey{id: utils.CheckErr(strconv.Atoi(idStr))}
		for {
			i++
			fields := strings.Fields(lines[i])
			if len(fields) == 0 { // Blank line. End of this monkey
				break
			}
			switch fields[0] {
			case "Starting": // "Starting items: 79, 98"
				itemStrs := fields[2:]
				monkey.items = make([]int, len(itemStrs))
				for j, itemStr := range itemStrs {
					monkey.items[j] = utils.CheckErr(strconv.Atoi(strings.TrimRight(itemStr, ",")))
				}
			case "Operation:": // "Operation: new = old + 19" or "Operation: new = old * old"
				value, err := strconv.Atoi(fields[5])
				if err != nil {
					monkey.operation = func(w int) int { return w * w } // Bit of a hack...
				} else {
					switch fields[4] {
					case "+":
						monkey.operation = func(w int) int { return w + value }
					case "*":
						monkey.operation = func(w int) int { return w * value }
					}
				}
			case "Test:": // "Test: divisible by 23"
				monkey.testDivisor = utils.CheckErr(strconv.Atoi(fields[3]))
			case "If": // "If true: throw to monkey 2" or "If false: throw to monkey 3
				value := utils.CheckErr(strconv.Atoi(fields[5]))
				if fields[1] == "true:" {
					monkey.trueThrowTo = value
				} else {
					monkey.falseThrowTo = value
				}
			default:
				panic(fmt.Sprintf("Unexpected line: %v", lines[i]))
			}
		}
		monkeys = append(monkeys, &monkey)
	}
	return monkeys
}

func main() {
	lines := utils.Lines(string(utils.CheckErr(os.ReadFile("input.txt"))))
	monkeys := parseMonkeys(lines)

	doRound := func(worryReducer func(int) int) {
		// Each monkey inspects all of its items, changing their worry level and throwing them on.
		for _, monkey := range monkeys {
			for _, worryLevel := range monkey.items {
				newWorryLevel := worryReducer(monkey.operation(worryLevel))
				var throwTo int
				if (newWorryLevel % monkey.testDivisor) == 0 {
					throwTo = monkey.trueThrowTo
				} else {
					throwTo = monkey.falseThrowTo
				}
				monkeys[throwTo].items = append(monkeys[throwTo].items, newWorryLevel)
				monkey.inspectionCount++
			}
			monkey.items = []int{}
		}
	}

	calculateMonkeyBusiness := func() int {
		biggestInspectionCounts := make([]int, 2)
		for _, monkey := range monkeys {
			if monkey.inspectionCount > biggestInspectionCounts[0] {
				biggestInspectionCounts[0] = monkey.inspectionCount
			}
			sort.Ints(biggestInspectionCounts)
		}
		return utils.Product(biggestInspectionCounts)
	}

	// In part 1 we get less worried whenever a monkey loses interested and throws an item away.
	for i := 0; i < 20; i++ {
		doRound(func(w int) int { return w / 3 })
	}
	part1 := calculateMonkeyBusiness()

	// Reset the monkeys back to their initial states for part 2.
	monkeys = parseMonkeys(lines)

	// In part 2 we don't get any less worried so the worry levels become unmanageably high! The
	// monkeys look at whether the worry level is divisble by the test divisor when deciding where
	// to throw an item. To simulate this we check the remainder (X % Y) is zero. If we replace X
	// with X modulo a multiple of Y then this doesn't change our answer. Therefore we can safely
	// replace the real worry level with itself modulo a number which is a multiple of all the test
	// divisors.
	divisorProduct := 1
	for _, m := range monkeys {
		divisorProduct *= m.testDivisor
	}
	for i := 0; i < 10000; i++ {
		doRound(func(w int) int { return w % divisorProduct })
	}
	part2 := calculateMonkeyBusiness()

	fmt.Printf("The answer to Part 1 is %v.\n", part1)
	fmt.Printf("The answer to Part 2 is %v.\n", part2)
}
