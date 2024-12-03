package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"AoC/pathfinding"
	"AoC/utils"
)

type Valve struct {
	name      string
	flowRate  int
	tunnelsTo []string
}

func main() {
	input := utils.CheckErr(os.ReadFile("input.txt"))
	inputRegex := regexp.MustCompile(
		`^Valve (\w+) has flow rate=([0-9]+); tunnels? leads? to valves? (.+)$`,
	)
	valves := map[string]*Valve{}
	for _, line := range utils.Lines(string(input)) {
		matches := inputRegex.FindStringSubmatch(line)
		valves[matches[1]] = &Valve{
			name:      matches[1],
			flowRate:  utils.CheckErr(strconv.Atoi(matches[2])),
			tunnelsTo: strings.Split(matches[3], ", "),
		}
	}

	// We can mostly ignore valves with zero flow rate since there's no point turning them on. The
	// only exception is "AA" which we can't ignore because we always start there.
	usefulValves := map[string]*Valve{}
	nonZeroValves := []string{}
	for name, valve := range valves {
		if name == "AA" {
			usefulValves[name] = valve
		}
		if valve.flowRate != 0 {
			usefulValves[name] = valve
			nonZeroValves = append(nonZeroValves, name)
		}
	}

	// Calculate the distance between all the useful valves to make things quicker later on.
	distances := map[string]int{}
	for _, a := range usefulValves {
		for _, b := range usefulValves {
			distance := pathfinding.Dijkstra(
				a,
				func(v *Valve) bool { return v.name == b.name },
				func(v *Valve) map[*Valve]int {
					neighbours := map[*Valve]int{}
					for _, name := range v.tunnelsTo {
						neighbours[valves[name]] = 1 // All moves only take 1 minute
					}
					return neighbours
				},
			)
			distances[a.name+","+b.name] = distance
			distances[b.name+","+a.name] = distance
		}
	}

	// Use a recursive function to find the best pressure released value we can achieve in a given
	// amount of time by opening a given set of valves.
	var getBestPressureReleased func(pos string, timeLeft, pressueReleased int, valves []string) int
	getBestPressureReleased = func(pos string, timeLeft, pressueReleased int, valves []string) int {
		var bestMoveOutcome *int
		for _, valveToOpen := range valves {
			timeToOpen := distances[pos+","+valveToOpen] + 1
			newPressureReleased := (timeLeft - timeToOpen) * usefulValves[valveToOpen].flowRate
			if newPressureReleased <= 0 {
				continue // Not a move worth doing (we might be running out of time)
			}
			newValves := make([]string, 0, len(valves)-1)
			for _, name := range valves {
				if name != valveToOpen {
					newValves = append(newValves, name)
				}
			}
			pressureReleased := getBestPressureReleased(
				valveToOpen,
				timeLeft-timeToOpen,
				pressueReleased+newPressureReleased,
				newValves,
			)
			if bestMoveOutcome == nil || pressureReleased > *bestMoveOutcome {
				bestMoveOutcome = &pressureReleased
			}
		}
		if bestMoveOutcome == nil {
			// There's no good moves left, return the pressure released by the moves already made.
			return pressueReleased
		} else {
			return *bestMoveOutcome
		}
	}

	part1 := getBestPressureReleased("AA", 30, 0, nonZeroValves)

	// For part2 we still want to maximise the pressure released but we now have two 'people'
	// opening valves independently. We want to try every possibly way of splitting the valves
	// between the two openers and generate this list of possible combinations using a bitmask.
	valveCombinations := [][][]string{}
	for bitmask := 0; bitmask < 1<<len(nonZeroValves); bitmask++ {
		valvesForMe := []string{}
		valvesForElephant := []string{}
		for i, name := range nonZeroValves {
			if bitmask&(1<<i) == 0 {
				valvesForMe = append(valvesForMe, name)
			} else {
				valvesForElephant = append(valvesForElephant, name)
			}
		}
		// A little hack to speed things up, the optimal solution probably gives a similar number
		// of values to each opener otherwise one will run out of useful moves to make. There's 15
		// non-zero valves in my puzzle input.
		if len(valvesForMe) <= 5 || len(valvesForElephant) <= 5 {
			continue
		}
		valveCombinations = append(valveCombinations, [][]string{valvesForMe, valvesForElephant})
	}

	// With the hack above and this bit of parallelisation it takes 15-20s to run.
	doPartOfPart2 := func(combinations [][][]string, resultChannel chan int) {
		result := 0
		for _, combination := range combinations {
			myValves, elephantValves := combination[0], combination[1]
			myPressureReleased := getBestPressureReleased("AA", 26, 0, myValves)
			elephantPressureReleased := getBestPressureReleased("AA", 26, 0, elephantValves)
			if myPressureReleased+elephantPressureReleased > result {
				result = myPressureReleased + elephantPressureReleased
			}
		}
		resultChannel <- result
	}
	resChan := make(chan int)
	go doPartOfPart2(valveCombinations[:len(valveCombinations)/2], resChan)
	go doPartOfPart2(valveCombinations[len(valveCombinations)/2:], resChan)
	part2 := <-resChan
	if otherResult := <-resChan; otherResult > part2 {
		part2 = otherResult
	}

	fmt.Printf("The answer to Part 1 is %v.\n", part1)
	fmt.Printf("The answer to Part 2 is %v.\n", part2)
}
