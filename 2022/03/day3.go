package main

import (
	"fmt"
	"os"

	"AoC/utils"
)

type itemSet = map[rune]bool

func scoreItem(item rune) int {
	if item > 90 { // lowercase letter
		return int(item) - int('a') + 1
	} else { // uppercase letter
		return int(item) - int('A') + 27
	}
}

func main() {
	lines := utils.Lines(string(utils.CheckErr(os.ReadFile("input.txt"))))

	errorItems := []rune{}
	for _, line := range lines {
		leftCompartment := make(itemSet)
		rightCompartment := make(itemSet)
		for i, item := range line {
			if i < len(line)/2 {
				leftCompartment[item] = true
			} else {
				rightCompartment[item] = true
			}
		}
		for x := range utils.Intersection(leftCompartment, rightCompartment) {
			errorItems = append(errorItems, x)
		}
	}

	badgeItems := []rune{}
	for i := 0; i < len(lines); i += 3 {
		team := []itemSet{}
		for j := 0; j < 3; j++ {
			backpack := make(itemSet)
			for _, item := range lines[i+j] {
				backpack[item] = true
			}
			team = append(team, backpack)
		}
		for x := range utils.Intersection(team[0], utils.Intersection(team[1], team[2])) {
			badgeItems = append(badgeItems, x)
		}
	}

	fmt.Printf("The answer to Part 1 is %v.\n", utils.Sum(utils.Transform(errorItems, scoreItem)))
	fmt.Printf("The answer to Part 2 is %v.\n", utils.Sum(utils.Transform(badgeItems, scoreItem)))
}
