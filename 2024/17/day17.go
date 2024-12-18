package main

import (
	"fmt"
	"math"
	"slices"
	"strings"
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
	// Integer division by a power of 2 is equivalent to bit shifting downwards
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

// Didn't bother writing an input parser, these are just copied from 'input.txt'
var inputA = 60589763
var program = []int{2, 4, 1, 5, 7, 5, 1, 6, 4, 1, 5, 5, 0, 3, 3, 0}

func main() {
	part1Computer := computer{
		regA:    inputA,
		program: program,
	}

	// For part 1 we just need to run the program as-is and format the output a bit
	part1Computer.run()
	output := fmt.Sprintf("%#v", part1Computer.output)
	part1 := strings.ReplaceAll(output[6:len(output)-1], " ", "")

	// For part 2 we've got to reverse engineer the program and work out the minimum initial value
	// for register A that will make the program output itself. Trying to understand the program:
	//
	// 2, 4, 1, 5, 7, 5, 1, 6, 4, 1, 5, 5, 0, 3, 3, 0
	//
	// Split into pairs of operations and operands:
	//
	// 2,4 -> set B to A mod 8 (set B to the last 3 bits of A)
	//
	// 1,5 -> bitwise xor B with 5
	//
	// 7,5 -> set C to A / 2**B (set C to some bits from the middle of A)
	//
	// 1,6 -> bitwise xor B with 6
	//
	// 4,1 -> bitwise xor B with C
	//
	// 5,5 -> output B mod 8
	//
	// 0,3 -> set A to A / 2**3 (bit shift A 3 places to the right, dropping last 3 bits)
	//
	// 3,0 -> jump back to start if A != 0
	//
	// So we see that the program loops, consuming the 3 least significant bits from the register A
	// value on each iteration. There's no state carried over between iterations except for the A
	// value (B and C are assigned before use). The value that gets outputted depends on the last
	// three bits of A but also some of the other bits (via register C).
	//
	// It's simplest to look at the last value we need to output, by this point our register A value
	// will have been entirely consumed (since we want the minimum working input) so it'll just be
	// three bits (value 0-7) and we don't need to worry about higher bits getting involved via
	// register C. Once we have the last value, we can look at the next last value and so on. Since
	// we've already worked out the later values we know what the higher bits are and can evaluate
	// the register C stuff. Note that there's multiple possible solutions in some cases and it's
	// hard to work out in advance which solutions will turn out to be dead-ends because of the
	// register C stuff so we just track them all at once, rather than bothering with backtracking.
	toOutput := make([]int, len(program))
	copy(toOutput, program)
	slices.Reverse(toOutput)

	findNextBits := func(bitsSoFar, nextTargetOutput int) []int {
		var nextBits []int

		for v := range 8 {
			c := computer{
				regA:    (bitsSoFar << 3) + v,
				program: program,
			}
			c.run()
			newOutput := c.output[0]
			if newOutput == nextTargetOutput {
				nextBits = append(nextBits, v)
			}
		}

		return nextBits
	}

	solutions := []int{0}
	for _, targetOutput := range toOutput {
		var nextSolutions []int
		for _, solutionSoFar := range solutions {
			nextBits := findNextBits(solutionSoFar, targetOutput)
			for _, value := range nextBits {
				nextSolutions = append(nextSolutions, (solutionSoFar<<3)+value)
			}
		}
		solutions = nextSolutions
	}

	part2 := math.MaxInt
	for _, s := range solutions {
		if s < part2 {
			part2 = s
		}
	}

	fmt.Printf("The answer to Part 1 is %v\n", part1)
	fmt.Printf("The answer to Part 2 is %v\n", part2)
}
