package main

import (
	utils "aoc2024"
	"fmt"
)

func mix(a, b int) int {
	return a ^ b // Bitwise XOR
}

func prune(a int) int {
	return a % 16777216 // 2^24
}

func evolve(secret int) int {
	secret = prune(mix(secret, secret*64))
	secret = prune(mix(secret, secret/32))
	secret = prune(mix(secret, secret*2048))

	return secret
}

func evolveMany(secret int, n int) []int {
	values := []int{}
	value := secret
	for range n {
		value = evolve(value)
		values = append(values, value)
	}
	return values
}

// Number of new prices each seller generates (and so number of price diffs we have in part 2)
const N = 2000

func main() {
	lines := utils.ReadLines("input.txt")

	part1 := 0 // Sum of the 2000th secret values

	sellerPrices := make([][]int, len(lines))
	for i, sellerSecret := range lines {
		initialSecret := utils.Int(sellerSecret)
		secrets := evolveMany(initialSecret, N+1)

		part1 += secrets[1999]

		prices := make([]int, N+1)
		for j := range N + 1 {
			prices[j] = secrets[j] % 10
		}
		sellerPrices[i] = prices
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)

	sellerDiffs := make([][]int, len(sellerPrices))
	for i, prices := range sellerPrices {
		sellerDiffs[i] = make([]int, N)
		for j := range N {
			sellerDiffs[i][j] = prices[j+1] - prices[j]
		}
	}

	// There's a lot of possible sequences so limit the amount of brute force work a bit by only
	// considering price diff sequences that are actually present in the data.
	knownSequences := make(map[string][]int)
	for _, diffs := range sellerDiffs {
		for i := range N - 4 {
			sequence := diffs[i : i+4]
			knownSequences[fmt.Sprintf("%v", sequence)] = sequence
		}
	}

	// Just brute force part 2 to find the bestPrice total price we can get. It takes ~3 minutes to run
	// but it does the job and running it is quicker than optimising the code.
	bestPrice := 0
	for _, sequence := range knownSequences {
		price := 0
		for i, diffs := range sellerDiffs {
			for j := range N - 4 {
				isMatch := diffs[j] == sequence[0] &&
					diffs[j+1] == sequence[1] &&
					diffs[j+2] == sequence[2] &&
					diffs[j+3] == sequence[3]

				if isMatch {
					price += sellerPrices[i][j+4]
					break
				}
			}
		}
		if price > bestPrice {
			bestPrice = price
		}
	}

	fmt.Printf("The answer to Part 2 is %v\n", bestPrice)
}
