package main

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"

	"AoC/utils"
)

type PacketData interface {
}

type PacketPair struct {
	left, right PacketData
}

func isCorrectOrder(left, right PacketData) *bool {
	// Left and right packets can be correctly ordered, incorrectly ordered, or equal (nil).
	correct := true
	incorrect := false

	leftInt, leftIsInt := left.(int)
	rightInt, rightIsInt := right.(int)

	// Both are ints, simply compare.
	if leftIsInt && rightIsInt {
		if leftInt == rightInt {
			return nil
		} else if leftInt < rightInt {
			return &correct
		} else { // leftInt > rightInt
			return &incorrect
		}
	}

	// At least one is a list. Compare as lists, making a list from an int if necessary. Note that
	// a |PacketData| instance must be either a list or an int.
	leftList, ok := left.([]PacketData)
	if !ok {
		leftList = []PacketData{leftInt}
	}
	rightList, ok := right.([]PacketData)
	if !ok {
		rightList = []PacketData{rightInt}
	}

	for i := 0; i < utils.Min(len(leftList), len(rightList)); i++ {
		if isCorrect := isCorrectOrder(leftList[i], rightList[i]); isCorrect != nil {
			return isCorrect
		}
	}

	if len(leftList) == len(rightList) {
		return nil
	} else if len(leftList) < len(rightList) {
		return &correct
	} else { // len(leftList) > len(rightList)
		return &incorrect
	}
}

func makePacketData(str string) PacketData {
	data := []PacketData{}
	for i := 1; i < len(str)-1; i++ { // Skip leading "[" and trailing "]"
		switch char := str[i]; char {
		case '[':
			// We've hit another list! Find the matching closing bracket and recurse.
			start := i
			depth := 0
			for {
				i++
				if str[i] == '[' {
					depth++
				}
				if str[i] == ']' {
					if depth == 0 {
						break
					}
					depth--
				}
			}
			data = append(data, makePacketData(str[start:i+1]))
		case ',':
			// no op
		default:
			// It's an integer! Read it all in and then parse.
			integerChars := []byte{char}
			for i < len(str)-2 {
				i++
				char = str[i]
				if char == '[' || char == ',' {
					i-- // Go back, we've overshot
					break
				}
				integerChars = append(integerChars, char)
			}
			data = append(data, utils.CheckErr(strconv.Atoi((string(integerChars)))))
		}
	}
	return data
}

func main() {
	inputLines := utils.Lines(string(utils.CheckErr(os.ReadFile("input.txt"))))

	pairs := []PacketPair{}
	for i := 0; i+2 < len(inputLines); i += 3 {
		pairs = append(pairs, PacketPair{
			left:  makePacketData(inputLines[i]),
			right: makePacketData(inputLines[i+1]),
		})
	}

	// For part 1 we need to sum the indicies of the pairs which are in the correct oder.
	part1 := 0
	for i, pair := range pairs {
		if *isCorrectOrder(pair.left, pair.right) {
			part1 += i + 1
		}
	}

	// For part 2 we need to order all packets, including two additional divider packets.
	divider1 := makePacketData("[[2]]")
	divider2 := makePacketData("[[6]]")
	packets := []PacketData{divider1, divider2}
	for _, pair := range pairs {
		packets = append(packets, pair.left)
		packets = append(packets, pair.right)
	}

	sort.Slice(packets, func(i, j int) bool {
		return *isCorrectOrder(packets[i], packets[j])
	})

	div1Idx, div2Idx := 0, 0
	for i, packet := range packets {
		if reflect.DeepEqual(packet, divider1) {
			div1Idx = i + 1
		}
		if reflect.DeepEqual(packet, divider2) {
			div2Idx = i + 1
		}
	}

	fmt.Printf("The answer to Part 1 is %v.\n", part1)
	fmt.Printf("The answer to Part 1 is %v.\n", div1Idx*div2Idx)
}
