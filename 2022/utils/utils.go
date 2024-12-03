package utils

import (
	"regexp"
)

func Abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func Max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// We don't expect any errors, treat them all as fatal.
func CheckErr[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// There's no built-in set type but we can use a map with bool values and always insert true.
func Intersection[T comparable](a, b map[T]bool) map[T]bool {
	intersection := make(map[T]bool)
	for val := range a {
		if b[val] {
			intersection[val] = true
		}
	}
	return intersection
}

func Lines(str string) []string {
	// Use a regex to make this work for both Windows and Unix line endings.
	return regexp.MustCompile("\r?\n").Split(str, -1)
}

func SplitAt[T comparable](slice []T, splitAt T) ([]T, []T) {
	for i, item := range slice {
		if item == splitAt {
			return slice[:i], slice[i+1:]
		}
	}
	return slice, nil
}

func Reduce[T any](slice []T, reducer func(T, T) T) T {
	switch len(slice) {
	case 0:
		panic("Can't reduce an empty slice")
	case 1:
		return slice[0]
	case 2:
		return reducer(slice[0], slice[1])
	default:
		reduction := reducer(slice[0], slice[1])
		for _, item := range slice[2:] {
			reduction = reducer(reduction, item)
		}
		return reduction
	}
}

func Product(slice []int) int {
	return Reduce(slice, func(a, b int) int { return a * b })
}

func Sum(slice []int) int {
	return Reduce(slice, func(a, b int) int { return a + b })
}

func Transform[A, B any](slice []A, transformFn func(A) B) []B {
	newSlice := make([]B, len(slice))
	for i := range slice {
		newSlice[i] = transformFn(slice[i])
	}
	return newSlice
}
