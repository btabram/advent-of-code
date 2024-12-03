package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"AoC/utils"
)

type CratePile struct {
	craneModelNumber int
	labels           []int
	stacks           map[int][]byte
}

func (cp *CratePile) String() string {
	s := fmt.Sprintf("CraneMover %v\n", cp.craneModelNumber)
	for _, label := range cp.labels {
		s += fmt.Sprintf("%d: %s\n", label, cp.stacks[label])
	}
	return s
}

func (cp *CratePile) move(moveCount, fromLabel, toLabel int) {
	to := cp.stacks[toLabel]
	from := cp.stacks[fromLabel]

	// The CrateMover 9001 can move mutliple crates at once!
	if cp.craneModelNumber > 9000 {
		to = append(to, from[len(from)-moveCount:]...)
		from = from[:len(from)-moveCount]
	} else {
		for i := 0; i < moveCount; i++ {
			to = append(to, from[len(from)-1:]...)
			from = from[:len(from)-1]
		}
	}

	cp.stacks[toLabel] = to
	cp.stacks[fromLabel] = from
}

func (cp *CratePile) topMessage() string {
	message := ""
	for _, label := range cp.labels {
		stack := cp.stacks[label]
		message += string(stack[len(stack)-1])
	}
	return message
}

// This is an example representation of a crate pile as a string:
// > [D]
// > [N] [C]
// > [Z] [M] [P]
// >  1   2   3
func newCratePile(stringEncoding []string, craneModelNumber int) *CratePile {
	cp := new(CratePile)
	cp.craneModelNumber = craneModelNumber
	cp.labels = make([]int, 0)
	cp.stacks = make(map[int][]byte)

	cratesStr := stringEncoding[:len(stringEncoding)-1]
	labelsStr := stringEncoding[len(stringEncoding)-1]

	for x := 1; x < len(labelsStr); x += 4 {
		stack := make([]byte, 0)
		for y := len(cratesStr) - 1; y >= 0; y-- {
			if value := cratesStr[y][x]; value != ' ' {
				stack = append(stack, value)
			}
		}
		label := utils.CheckErr(strconv.Atoi(string(labelsStr[x])))
		cp.labels = append(cp.labels, label)
		cp.stacks[label] = stack
	}

	return cp
}

func main() {
	lines := utils.Lines(string(utils.CheckErr(os.ReadFile("input.txt"))))

	cratePileStr, instructions := utils.SplitAt(lines, "")
	part1Crates := newCratePile(cratePileStr, 9000)
	part2Crates := newCratePile(cratePileStr, 9001)

	instructionsRegex := regexp.MustCompile("^move ([0-9]+) from ([0-9]+) to ([0-9]+)$")
	for _, instruction := range instructions {
		groups := instructionsRegex.FindStringSubmatch(instruction)
		moveCount := utils.CheckErr(strconv.Atoi(groups[1]))
		fromLabel := utils.CheckErr(strconv.Atoi(groups[2]))
		toLabel := utils.CheckErr(strconv.Atoi(groups[3]))
		part1Crates.move(moveCount, fromLabel, toLabel)
		part2Crates.move(moveCount, fromLabel, toLabel)
	}

	fmt.Printf("The answer to Part 1 is %v.\n", part1Crates.topMessage())
	fmt.Printf("The answer to Part 2 is %v.\n", part2Crates.topMessage())
}
