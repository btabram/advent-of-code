package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"AoC/utils"
)

type File struct {
	name string
	size int
}

type Directory struct {
	name     string
	files    []File
	children map[string]*Directory
	parent   *Directory
	size     *int // Pointer because it's optional and may not be set
}

func (d *Directory) getSize() int {
	if d.size != nil {
		return *d.size
	}

	size := 0
	for _, file := range d.files {
		size += file.size
	}
	for _, child := range d.children {
		size += child.getSize()
	}
	d.size = &size
	return size
}

func newDirectory(name string, parent *Directory) *Directory {
	return &Directory{
		name:     name,
		files:    make([]File, 0),
		children: make(map[string]*Directory),
		parent:   parent,
	}
}

// Sum the size of all directories smaller than 100000.
func part1(cwd *Directory) int {
	ans := 0
	if cwd.getSize() < 100000 {
		ans += cwd.getSize()
	}
	for _, child := range cwd.children {
		ans += part1(child)
	}
	return ans
}

// Return the size of the smallest directory which is greater than |requiredSize|.
func part2(cwd *Directory, requiredSize int) *int {
	// We don't need to consider the children of directories that are too small because they must
	// also be too small.
	if cwd.getSize() >= requiredSize {
		ans := cwd.getSize()
		for _, child := range cwd.children {
			s := part2(child, requiredSize)
			if s != nil && *s >= requiredSize && *s < ans {
				ans = *s
			}
		}
		return &ans
	}
	return nil
}

func main() {
	inputLines := utils.Lines(string(utils.CheckErr(os.ReadFile("input.txt"))))

	// Work through the input, building up our directory structure.
	var cwd *Directory
	for i := 0; i < len(inputLines); i++ {
		line := inputLines[i]
		if line == "$ ls" {
			// Loop over the output of the ls command.
			for {
				i += 1
				if i == len(inputLines) {
					// End of input.
					break
				}
				line = inputLines[i]
				if line[0] == '$' {
					// We've overshot and are on to the next command, go back.
					i -= 1
					break
				}
				words := strings.Fields(line)
				if words[0] == "dir" { // line == "dir foo"
					name := words[1]
					cwd.children[name] = newDirectory(name, cwd)
				} else { // line == "1234 bar.txt"
					size, name := words[0], words[1]
					cwd.files = append(cwd.files, File{
						name: name,
						size: utils.CheckErr(strconv.Atoi(size)),
					})
				}
			}
		} else { // line == "$ cd baz"
			newDir := strings.Fields(line)[2]
			if cwd == nil {
				// Create the root directory.
				cwd = newDirectory(newDir, nil)
			} else if newDir == ".." {
				cwd = cwd.parent
			} else {
				cwd = cwd.children[newDir]
			}
		}
	}

	// Find the root directory by working our way back up.
	root := cwd
	for root.parent != nil {
		root = root.parent
	}

	totalDiskSpace := 70000000
	requiredDiskSpace := 30000000
	currentFreeSpace := totalDiskSpace - root.getSize()

	fmt.Printf("The answer to Part 1 is %v.\n", part1(root))
	fmt.Printf("The answer to Part 2 is %v.\n", *part2(root, requiredDiskSpace-currentFreeSpace))
}
