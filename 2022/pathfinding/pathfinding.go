package pathfinding

import "sort"

type queueItem[T any] struct {
	item T
	cost int
}

func Dijkstra[Node comparable](
	start Node,
	isFinish func(Node) bool,
	getNeighbours func(Node) map[Node]int,
) int {
	queue := make([]queueItem[Node], 1)
	visited := map[Node]bool{}

	// Add initial node to the queue with zero cost.
	queue[0] = queueItem[Node]{start, 0}

	for {
		// Pop the lowest cost item out of the queue.
		next := queue[len(queue)-1]
		queue = queue[:len(queue)-1]

		toVisit, costSoFar := next.item, next.cost

		// The queue can have multiple entires for the same node so check for already visited.
		if visited[toVisit] {
			continue
		}

		if isFinish(toVisit) {
			return costSoFar
		}

		// Add neighbours of the node we're about to visit to our queue (and keep it sorted!).
		for neighbour, cost := range getNeighbours(toVisit) {
			if !visited[neighbour] {
				queue = append(queue, queueItem[Node]{neighbour, costSoFar + cost})
			}
		}
		sort.Slice(queue, func(i, j int) bool { return queue[i].cost > queue[j].cost })

		// Finally visit the node.
		visited[toVisit] = true
	}
}
