package main

import (
	"fmt"
	"os"
	"sort"
	"strconv"

	"AoC/utils"
)

func main() {
	input := utils.CheckErr(os.ReadFile("input.txt"))

	elfCalorieLists := [][]int{{}}
	for _, line := range utils.Lines(string(input)) {
		if line != "" {
			val := utils.CheckErr(strconv.Atoi(line))
			currentList := &elfCalorieLists[len(elfCalorieLists)-1]
			*currentList = append(*currentList, val)
		} else {
			elfCalorieLists = append(elfCalorieLists, []int{})
		}
	}

	elfCalorieCounts := utils.Transform(elfCalorieLists, utils.Sum)
	sort.Ints(elfCalorieCounts)

	size := len(elfCalorieCounts)
	fmt.Printf("The answer to Part 1 is %v.\n", elfCalorieCounts[size-1])
	fmt.Printf("The answer to Part 2 is %v.\n", utils.Sum(elfCalorieCounts[size-3:]))
}
