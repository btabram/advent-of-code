package main

import (
	utils "aoc2024"
	"fmt"
	"slices"
	"strings"
)

func main() {
	lines := utils.ReadLines("input.txt")

	networkMap := map[string][]string{}

	for _, line := range lines {
		parts := strings.Split(line, "-")
		computerA := parts[0]
		computerB := parts[1]

		// Links are bidirectional
		networkMap[computerA] = append(networkMap[computerA], computerB)
		networkMap[computerB] = append(networkMap[computerB], computerA)
	}

	// For part 1 we find all cliques of 3 computers that are all connected. Map avoids duplicates
	triples := map[string][]string{}

	for computerA, connectionsA := range networkMap {
		for _, computerB := range connectionsA {
			connectionsB := networkMap[computerB]
			for _, computerC := range connectionsB {
				if slices.Contains(connectionsA, computerC) {
					group := []string{computerA, computerB, computerC}
					slices.Sort(group)
					triples[strings.Join(group, ",")] = group
				}
			}
		}
	}

	// Count the number of triples that contain a computer starting with "t"
	part1 := 0
	for _, group := range triples {
		for _, computer := range group {
			if strings.HasPrefix(computer, "t") {
				part1++
				break
			}
		}
	}

	// For part 2 we need to find the largest clique which is potentially a lot of work. Looking at
	// the input every computer is connected to 13 others so the largest possible clique is 14. In
	// the example every computer is connected to 4 others and the largest clique is 4 so let's try
	// to find a clique of size 13 (hopefully there's only one!)
	part2 := ""
	for computer, connections := range networkMap {
		for _, connectionToExclude := range connections {
			// Try a possible size 13 clique by excluding one connection (and including `computer`)
			cliqueToTry := map[string]bool{}
			for _, c := range connections {
				if c != connectionToExclude {
					cliqueToTry[c] = true
				}
			}

			valid := true
			for c := range cliqueToTry {
				matchCount := 0
				for _, connection := range networkMap[c] {
					if cliqueToTry[connection] || connection == computer {
						matchCount++
					}
				}
				if matchCount != 12 { // each computer in the clique is connected to the n-1 others
					valid = false
					break
				}
			}

			if valid {
				fullClique := []string{computer}
				for c := range cliqueToTry {
					fullClique = append(fullClique, c)
				}
				slices.Sort(fullClique)
				clique := strings.Join(fullClique, ",")
				if part2 != "" && part2 != clique {
					panic("found multiple cliques of size 13")
				}
				part2 = clique
			}
		}
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
