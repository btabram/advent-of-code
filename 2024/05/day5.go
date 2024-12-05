package main

import (
	utils "aoc2024"
	"fmt"
	"slices"
)

func isValid(rulesMap map[int][]int, pages []int) bool {
	printedPages := make(map[int]bool)

	for _, page := range pages {
		// Go through pages which must have been printed before `page` according to the rules
		for _, mustBePrintedPage := range rulesMap[page] {
			if !slices.Contains(pages, mustBePrintedPage) {
				continue // Rule not applicable
			}
			if !printedPages[mustBePrintedPage] {
				return false // Invalid
			}
		}
		printedPages[page] = true
	}

	return true
}

func makeValid(rulesMap map[int][]int, pages []int) []int {
	// Start with no pages printed, and all pages remaining
	printedPages := make([]int, 0, len(pages))
	remainingPages := make(map[int]bool)
	for _, page := range pages {
		remainingPages[page] = true
	}

	// Keep looping until all pages have been printed
	for len(printedPages) != len(pages) {

		// Loop through all remaining pages until we find one that is allowed to be printed next
	pageLoop:
		for page := range remainingPages {
			for _, mustBePrintedPage := range rulesMap[page] {
				if remainingPages[mustBePrintedPage] {
					// `page` is blocked and can't be printed yet, continue and try another page
					continue pageLoop
				}
			}

			// If we've reached this point then `page` is unblocked and can be printed
			printedPages = append(printedPages, page)
			delete(remainingPages, page)

			// Start a new iteration now we've printed another page (and modified `remainingPages`)
			break pageLoop
		}
	}

	return printedPages
}

// Assuming all lists have an odd length and therefore a well-defined middle value
func middle(pages []int) int {
	return pages[len(pages)/2]
}

func main() {
	rawLines := utils.ReadLines("input.txt")

	lines := make([][]int, len(rawLines))
	for i, line := range rawLines {
		lines[i] = utils.Ints(line)
	}

	split := slices.IndexFunc(lines, func(l []int) bool { return len(l) == 0 })
	// Two separate parts to the input, separated by a blank line
	orderingRules := lines[:split]
	pagesToPrint := lines[split+1:]

	// map of page number to all the page numbers that must come before it
	rulesMap := make(map[int][]int)
	for _, rule := range orderingRules {
		before := rule[0]
		after := rule[1]
		//rulesMap[before] = append(rulesMap[before], after)
		rulesMap[after] = append(rulesMap[after], before)
	}

	part1 := 0
	var invalidPagesToPrint [][]int
	for _, pages := range pagesToPrint {
		if isValid(rulesMap, pages) {
			part1 += middle(pages)
		} else {
			invalidPagesToPrint = append(invalidPagesToPrint, pages)
		}
	}

	part2 := 0
	for _, invalidPages := range invalidPagesToPrint {
		part2 += middle(makeValid(rulesMap, invalidPages))
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
