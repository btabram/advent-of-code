package utils

import (
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Int tries to parse the sting as an int.
func Int(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

// Ints splits a string on common separators and then tries to parse all fields as integers.
func Ints(str string) []int {
	fields := strings.FieldsFunc(
		str, func(r rune) bool { return r == ' ' || r == ',' || r == '|' },
	)
	ints := make([]int, len(fields))

	for i, strVal := range fields {
		ints[i] = Int(strVal)
	}

	return ints
}

func Read(filename string) string {
	contents, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(contents)
}

func ReadLines(filename string) []string {
	return strings.Split(Read(filename), "\n")
}

// See https://en.wikipedia.org/wiki/A*_search_algorithm
// Returns the minimum cost to get to the finish and the `prevNodes` map
func AStarPathfinding[node comparable](
	start node,
	isFinish func(node) bool,
	// Should return a map of possible neighbours -> the cost of moving to that neighbour
	getNeighbours func(node) map[node]int,
	// A heuristic to guide the search, should provide a lower-bound (never overestimate) for the
	// cost to get to the finish from the given node
	heuristic func(node) int,
) (int, map[node][]node) {
	// Set of nodes to consider
	queue := make(map[node]bool)
	queue[start] = true

	// Known (best) costs to visit particular nodes
	knownCosts := make(map[node]int)
	knownCosts[start] = 0

	// Predicted total costs to get to the finish via this particular node
	predicatedTotalCosts := make(map[node]int)
	predicatedTotalCosts[start] = heuristic(start)

	// Track the previous nodes(s) for the best routes to particular nodes
	prevNodes := make(map[node][]node)

	for len(queue) != 0 {
		// Pop the lowest predicted cost item out of the queue. This is the node we're visiting now
		var current node
		min := math.MaxInt
		for n := range queue {
			ptc := predicatedTotalCosts[n]
			if ptc < min {
				current = n
				min = ptc
			}
		}
		delete(queue, current)

		costToCurrent := knownCosts[current]

		// Check if we've reached the goal
		if isFinish(current) {
			return costToCurrent, prevNodes
		}

		// Consider neighbours of the current node we're visiting
		for neighbour, costToMove := range getNeighbours(current) {
			costToNeighbour := costToCurrent + costToMove

			knownCostToNeighbour, ok := knownCosts[neighbour]
			if ok && costToNeighbour > knownCostToNeighbour {
				// We've already found a better route to `neighbour`
				continue
			}

			if ok && costToNeighbour == knownCostToNeighbour {
				// Register a new joint best route (we want to know all possibilities)
				prevNodes[neighbour] = append(prevNodes[neighbour], current)
				continue
			}

			// At this point we've found a new / better route which needs consideration
			queue[neighbour] = true
			knownCosts[neighbour] = costToNeighbour
			predicatedTotalCosts[neighbour] = costToNeighbour + heuristic(neighbour)
			prevNodes[neighbour] = []node{current}
		}
	}

	panic("Failed to find a path")
}
