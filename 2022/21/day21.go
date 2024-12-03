package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"AoC/utils"
)

type Monkey struct {
	name  string
	value *int
	arg1  string
	op    string
	arg2  string
}

func parseMonkeys(input []byte) map[string]*Monkey {
	monkeys := map[string]*Monkey{}
	for _, line := range utils.Lines(string(input)) {
		fields := strings.Fields(line)
		name := fields[0][:4]
		switch len(fields) {
		case 2: // e.g. "dbpl: 5"
			value := utils.CheckErr(strconv.Atoi(fields[1]))
			monkeys[name] = &Monkey{
				name:  name,
				value: &value,
			}
		case 4: // e.g. "cczh: sllz + lgvd"
			monkeys[name] = &Monkey{
				name: name,
				arg1: fields[1],
				op:   fields[2],
				arg2: fields[3],
			}
		default:
			panic(fmt.Sprintf("Unexpected input line: %s", line))
		}
	}
	return monkeys
}

func solveMonkeys(monkeys map[string]*Monkey) {
	for {
		var toSolve *Monkey
		for _, m := range monkeys {
			if m.value == nil {
				hasValue := func(argName string) bool {
					if argMonkey, ok := monkeys[argName]; ok {
						return argMonkey.value != nil
					}
					return false
				}
				if hasValue(m.arg1) && hasValue(m.arg2) {
					toSolve = m
					break
				}
			}
		}
		if toSolve == nil {
			return
		}
		arg1 := *monkeys[toSolve.arg1].value
		arg2 := *monkeys[toSolve.arg2].value
		value := 0
		switch toSolve.op {
		case "+":
			value = arg1 + arg2
		case "-":
			value = arg1 - arg2
		case "*":
			value = arg1 * arg2
		case "/":
			value = arg1 / arg2
		default:
			panic(fmt.Sprintf("Invalid operator: %s", toSolve.op))
		}
		toSolve.value = &value
	}
}

func main() {
	input := utils.CheckErr(os.ReadFile("input.txt"))

	monkeys := parseMonkeys(input)
	solveMonkeys(monkeys)
	part1 := *monkeys["root"].value

	monkeys = parseMonkeys(input)
	monkeys["humn"].value = nil
	solveMonkeys(monkeys)
	var part2 int
	// We work backwards, undoing operations to work out the required monkey values so that the
	// "root" equality is met, until we have the "humn" value we need.
	currentMonkey := monkeys[monkeys["root"].arg1]
	requiredValue := *monkeys[monkeys["root"].arg2].value
	for {
		if currentMonkey.name == "humn" {
			part2 = requiredValue
			break
		}

		arg1Monkey := monkeys[currentMonkey.arg1]
		arg2Monkey := monkeys[currentMonkey.arg2]

		var arg1HasValue bool
		var knownArgValue int
		if arg1Monkey.value != nil {
			arg1HasValue = true
			knownArgValue = *arg1Monkey.value
		} else if arg2Monkey.value != nil {
			arg1HasValue = false
			knownArgValue = *arg2Monkey.value
		} else {
			panic("Expected one argument monkey to have a value")
		}

		switch currentMonkey.op {
		case "+":
			requiredValue = requiredValue - knownArgValue
		case "-":
			if arg1HasValue {
				requiredValue = knownArgValue - requiredValue
			} else {
				requiredValue = requiredValue + knownArgValue
			}
		case "*":
			requiredValue = requiredValue / knownArgValue
		case "/":
			if arg1HasValue {
				requiredValue = knownArgValue / requiredValue
			} else {
				requiredValue = requiredValue * knownArgValue
			}
		}

		if arg1HasValue {
			currentMonkey = arg2Monkey
		} else {
			currentMonkey = arg1Monkey
		}
	}

	fmt.Printf("The answer to Part 1 is %v.\n", part1)
	fmt.Printf("The answer to Part 2 is %v.\n", part2)
}
