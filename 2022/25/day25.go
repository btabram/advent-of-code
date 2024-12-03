package main

import (
	"fmt"
	"os"
	"strings"

	"AoC/utils"
)

var fromSNAFUDigit = map[byte]int{
	'2': 2,
	'1': 1,
	'0': 0,
	'-': -1,
	'=': -2,
}

var toSNAFUDigit = map[int]string{
	2:  "2",
	1:  "1",
	0:  "0",
	-1: "-",
	-2: "=",
}

func addSNAFU(a, b string) string {
	// Pad with zeros so that both numbers are the same length.
	lenA, lenB := len(a), len(b)
	if lenA != lenB {
		if lenA > lenB {
			for i := 0; i < (lenA - lenB); i++ {
				b = "0" + b
			}
		} else {
			for i := 0; i < (lenB - lenA); i++ {
				a = "0" + a
			}
		}
	}

	ans := ""
	toCarry := 0
	for i := len(a) - 1; i >= 0; i-- {
		summedDigit := fromSNAFUDigit[a[i]] + fromSNAFUDigit[b[i]] + toCarry
		if summedDigit > 2 {
			summedDigit -= 5
			toCarry = 1
		} else if summedDigit < -2 {
			summedDigit += 5
			toCarry = -1
		} else {
			toCarry = 0
		}
		ans = toSNAFUDigit[summedDigit] + ans
	}
	ans = toSNAFUDigit[toCarry] + ans

	// Trim any leading zeros.
	return strings.TrimLeft(ans, "0")
}

func main() {
	fuelRequirements := utils.Lines(string(utils.CheckErr(os.ReadFile("input.txt"))))

	sum := "0"
	for _, fr := range fuelRequirements {
		sum = addSNAFU(sum, fr)
	}

	fmt.Printf("The answer to Part 1 is %v.\n", sum)
}
