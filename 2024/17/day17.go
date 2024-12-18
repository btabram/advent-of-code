package main

import (
	"fmt"
	"math"
	"slices"
)

// A 3-bit computer a described in the puzzle input
type computer struct {
	regA, regB, regC   int // Registers
	instructionPointer int
	program            []int
	output             []int
}

// Operands are either used as literally or as "combo operands" which have more complex values
func (c *computer) comboValue(operand int) int {
	if operand >= 0 && operand <= 3 {
		return operand // Act as literals
	}
	if operand == 4 {
		return c.regA // Act as value of register A
	}
	if operand == 5 {
		return c.regB // Act as value of register b
	}
	if operand == 6 {
		return c.regC // Act as value of register C
	}

	panic("Invalid combo operand")
}

// Divide register A by 2 to the power of the combo operand. Use by multiple instructions
func (c *computer) divisionImpl(operand int) int {
	//return c.regA / (1 << c.comboValue(operand))
	return c.regA >> c.comboValue(operand)
}

func (c *computer) run() {
	for {
		if c.instructionPointer < 0 || c.instructionPointer >= len(c.program) {
			return // Halt
		}

		opcode := c.program[c.instructionPointer]
		operand := c.program[c.instructionPointer+1]

		c.instructionPointer += 2

		switch opcode {
		case 0:
			// adv (division stored in register A)
			c.regA = c.divisionImpl(operand)
		case 1:
			// bxl (bitwise XOR register B with the literal operand)
			c.regB ^= operand
		case 2:
			// bst (set register B to the combo operator modulo 8)
			c.regB = c.comboValue(operand) % 8
		case 3:
			// jnz (conditionally jump instruction pointer)
			if c.regA != 0 {
				c.instructionPointer = operand
			}
		case 4:
			// bxc (bitwise XOR register B with register C)
			c.regB ^= c.regC
		case 5:
			// out (output the combo operand value modulo 8)
			c.output = append(c.output, c.comboValue(operand)%8)
		case 6:
			// bdv (division stored in register B)
			c.regB = c.divisionImpl(operand)
		case 7:
			// cdv (division stored in register C)
			c.regC = c.divisionImpl(operand)
		default:
			panic("Invalid instruction")
		}
	}
}

func main() {
	test := computer{
		regA:    729,
		program: []int{0, 1, 5, 4, 3, 0},
	}

	test.run()
	fmt.Println(test)

	p1 := computer{
		regA:    100032, // 0b011_000_011,
		program: []int{2, 4, 1, 5, 7, 5, 1, 6, 4, 1, 5, 5, 0, 3, 3, 0},
	}
	p1.run()
	fmt.Println(p1)

	toFind := []int{2, 4, 1, 5, 7, 5, 1, 6, 4, 1, 5, 5, 0, 3, 3, 0}
	slices.Reverse(toFind)

	findNexts := func(soFar, next int) []int {
		ns := []int{}

		for v := range 8 {
			c := computer{
				regA:    (soFar << 3) + v,
				program: []int{2, 4, 1, 5, 7, 5, 1, 6, 4, 1, 5, 5, 0, 3, 3, 0},
			}
			c.run()
			done := c.output[0]
			if done == next {
				ns = append(ns, v)
			}
		}

		return ns
	}

	soFars := []int{0}
	for _, f := range toFind {
		for i, sf := range soFars {
			if sf == -1 {
				continue // skip deadends
			}

			nexts := findNexts(sf, f)

			if len(nexts) == 0 {
				soFars[i] = -1
			}

			for j, nx := range nexts {
				if j == 0 {
					soFars[i] = (sf << 3) + nx
				} else {
					soFars = append(soFars, (sf<<3)+nx)
				}
			}
		}
	}

	min := math.MaxInt
	for _, v := range soFars {
		if v != -1 && v < min {
			min = v
		}
	}

	p := computer{
		regA:    min,
		program: []int{2, 4, 1, 5, 7, 5, 1, 6, 4, 1, 5, 5, 0, 3, 3, 0},
	}
	p.run()
	fmt.Println(p)
	fmt.Println(min)

	//fmt.Printf("The answer to Part 1 is %v\n", part1)
	//fmt.Printf("The answer to Part 2 is %v\n", part2)
}

/*
	2, 4, 1, 5, 7, 5, 1, 6, 4, 1, 5, 5, 0, 3, 3, 0

	2,4 -> set B to A%8

	1,5 -> bitwise xor B with 5

	7,5 -> set C to A / 2**B

	1,6 -> bitwise xor B with 6

	4,1 -> bitwise xor B with C

	5,5 -> output B%8

	0,3 -> set A to A / 2**3

	3,0 -> jump back to start if A != 0


~~~~~~~~~~~~~~~~~~~~~~~~~~
	B becomes (A % 8) ^ 5 -> only depends on last three bts of A

	C = A / (2 ** ((A % 8) ^ 5)) -> A bit shifted down by B

	/6 at the end effectively bit shifts 3 places to the right (dropping last 3 bits)


	111 -> 7
	110 -> 6
	101 -> 5
	100 -> 4
	011 -> 3
	010 -> 2
	001 -> 1
	000 -> 0


	~~~~~~~~~~~~~~~~~~

	* last output needs to be 0, so most significant octet must be 3
	* prev output needs to be 3, smallest fitting octet is 0 (C = 0)
	* next needs to be 3 again... what about C!?
*/
