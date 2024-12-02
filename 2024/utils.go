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

func Ints(str string) []int {
	fields := strings.Fields(str)

	ints := make([]int, len(fields))

	for i, strVal := range fields {
		val, err := strconv.Atoi(strVal)
		if err != nil {
			log.Fatal(err)
		}
		ints[i] = val
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
