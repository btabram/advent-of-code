package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"AoC/utils"
)

const (
	screenWidth  = 40
	screenHeight = 6
)

func main() {
	input := utils.CheckErr(os.ReadFile("input.txt"))

	totalSignalStrength := 0
	pixels := [screenWidth * screenHeight]byte{}

	cycleNumber, registerX := 1, 1
	doCycle := func() {
		// For part 1 we want the signal strength during cycles 20, 60, 100 etc.
		if cycleNumber%40 == 20 {
			totalSignalStrength += cycleNumber * registerX
		}
		// For part 2 we draw a pixel every cycle. The pixel is lit if the sprite (whose position is
		// given by the X register value) overlaps with the pixel being drawn, otherwise it is dark.
		horizontalDrawingPosition := (cycleNumber % screenWidth) - 1
		if utils.Abs(registerX-horizontalDrawingPosition) <= 1 { // The sprite is 3 pixels wide
			pixels[cycleNumber-1] = '#'
		} else {
			pixels[cycleNumber-1] = '.'
		}
		// Finally, increment the cycle number.
		cycleNumber += 1
	}

	for _, instruction := range utils.Lines(string(input)) {
		fields := strings.Fields(instruction)
		switch fields[0] {
		case "noop": // Takes one cycle and does nothing
			doCycle()
		case "addx": // Takes two cycles and afterwards increments the X register
			doCycle()
			doCycle()
			registerX += utils.CheckErr(strconv.Atoi(fields[1]))
		}
	}

	fmt.Printf("The answer to Part 1 is %v.\n", totalSignalStrength)
	fmt.Println("The answer is Part 2 is:")
	for i := 0; i < screenHeight; i++ {
		fmt.Printf("%s\n", pixels[i*screenWidth:((i+1)*screenWidth)-1])
	}
}
