package utils

import (
	"log"
	"os"
	"strconv"
	"strings"
)

func Abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

// Int tries to parse the sting as an int.
func Int(str string) int {
	val, err := strconv.Atoi(str)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

// Ints splits a string and then tries to parse all fields as integers.
func Ints(str string) []int {
	fields := strings.Fields(str)
	ints := make([]int, len(fields))

	for i, strVal := range fields {
		ints[i] = Int(strVal)
	}

	return ints
}

func Read(filename string) string {
	contents, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return string(contents)
}

func ReadLines(filename string) []string {
	return strings.Split(Read(filename), "\n")
}
