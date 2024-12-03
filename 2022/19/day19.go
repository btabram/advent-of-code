package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"AoC/utils"
)

const (
	ore = iota
	clay
	obsidian
	geode
	typesLen
)

var types []int = []int{ore, clay, obsidian, geode}

type Blueprint struct {
	id         int
	robotCosts [typesLen][typesLen]int // costs[robot type][resource type]
}

type State struct {
	ore, clay, obsidian, geodes             int
	oreBots, clayBots, obsidBots, geodeBots int
	toBuild                                 int
	time                                    int
}

func findMaxGeodeProduction(bp Blueprint, timeLimit int) int {
	// Work out the maximum number of robots of a given type that it's useful to have. For example,
	// if everything costs less than 5 ore then a 5th ore bot is useless because you can only build
	// one robot per turn and you already had enough ore production to build every turn. Any branch
	// where we build more robots than is useful must be sub-optimal so we stop simulating it.
	var oreCosts, clayCosts, obsidianCosts []int
	for _, costs := range bp.robotCosts {
		oreCosts = append(oreCosts, costs[ore])
		clayCosts = append(clayCosts, costs[clay])
		obsidianCosts = append(obsidianCosts, costs[obsidian])
	}
	maxOre := utils.Reduce(oreCosts, utils.Max)
	maxClay := utils.Reduce(clayCosts, utils.Max)
	maxObsidian := utils.Reduce(obsidianCosts, utils.Max)

	// Have a stack of states to simulate. For now filled only with the possible initial states.
	statesToTry := []State{
		{oreBots: 1, toBuild: ore},
		{oreBots: 1, toBuild: clay},
	}

	// Take a state and simulate it, working out the number of geodes produced. Whenever there's a
	// branching point we add the possible branch states to the |statesToTry| stack.
	simulate := func(s State) int {
		// Identify doomed states which are trying to build a robot they can't ever afford.
		if (s.toBuild == obsidian && s.clayBots == 0) || (s.toBuild == geode && s.obsidBots == 0) {
			return 0
		}

		for {
			// Work out whether we can build a new robot during this round.
			building := false
			costs := bp.robotCosts[s.toBuild]
			if s.ore >= costs[ore] && s.clay >= costs[clay] && s.obsidian >= costs[obsidian] {
				s.ore -= costs[ore]
				s.clay -= costs[clay]
				s.obsidian -= costs[obsidian]
				building = true
			}

			// Gather resources.
			s.ore += s.oreBots
			s.clay += s.clayBots
			s.obsidian += s.obsidBots
			s.geodes += s.geodeBots

			s.time++
			if s.time == timeLimit {
				return s.geodes
			}

			// Maybe build a robot and decide what to build next, with branching and some pruning.
			if building {
				switch s.toBuild {
				case ore:
					s.oreBots++
					if s.oreBots > maxOre {
						return 0
					}
				case clay:
					s.clayBots++
					if s.clayBots > maxClay {
						return 0
					}
				case obsidian:
					s.obsidBots++
					if s.obsidBots > maxObsidian {
						return 0
					}
				case geode:
					s.geodeBots++
				}
				for _, robotType := range types {
					if robotType == s.toBuild {
						continue
					}
					newState := s
					newState.toBuild = robotType
					statesToTry = append(statesToTry, newState)
				}
			}
		}
	}

	best := 0
	for len(statesToTry) != 0 {
		poppedState := statesToTry[len(statesToTry)-1]
		statesToTry = statesToTry[:len(statesToTry)-1]
		if result := simulate(poppedState); result > best {
			best = result
		}
	}
	return best
}

func main() {
	input := utils.CheckErr(os.ReadFile("input.txt"))
	inputRegex := regexp.MustCompile(`^Blueprint ([0-9]+): ` +
		`Each ore robot costs ([0-9]+) ore. ` +
		`Each clay robot costs ([0-9]+) ore. ` +
		`Each obsidian robot costs ([0-9]+) ore and ([0-9]+) clay. ` +
		`Each geode robot costs ([0-9]+) ore and ([0-9]+) obsidian.$`,
	)
	blueprints := []Blueprint{}
	for _, line := range utils.Lines(string(input)) {
		matches := inputRegex.FindStringSubmatch(line)
		bp := Blueprint{id: utils.CheckErr(strconv.Atoi(matches[1]))}
		bp.robotCosts[ore][ore] = utils.CheckErr(strconv.Atoi(matches[2]))
		bp.robotCosts[clay][ore] = utils.CheckErr(strconv.Atoi(matches[3]))
		bp.robotCosts[obsidian][ore] = utils.CheckErr(strconv.Atoi(matches[4]))
		bp.robotCosts[obsidian][clay] = utils.CheckErr(strconv.Atoi(matches[5]))
		bp.robotCosts[geode][ore] = utils.CheckErr(strconv.Atoi(matches[6]))
		bp.robotCosts[geode][obsidian] = utils.CheckErr(strconv.Atoi(matches[7]))
		blueprints = append(blueprints, bp)
	}

	part1 := 0
	for _, bp := range blueprints {
		mostGeodes := findMaxGeodeProduction(bp, 24)
		part1 += mostGeodes * bp.id
	}
	fmt.Printf("The answer to Part 1 is %v.\n", part1)

	resChan := make(chan int)
	for _, bp := range blueprints[:3] {
		bp := bp // Avoid loop variable
		go func() {
			resChan <- findMaxGeodeProduction(bp, 32)
		}()
	}
	part2 := 1
	for i := 0; i < 3; i++ {
		part2 *= <-resChan
	}
	fmt.Printf("The answer to Part 2 is %v.\n", part2)
}
