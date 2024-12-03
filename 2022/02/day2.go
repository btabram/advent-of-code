package main

import (
	"fmt"
	"os"
	"strings"

	"AoC/utils"
)

var (
	decryptionMap = map[string]string{
		"A": "rock",
		"B": "paper",
		"C": "scissors",
		"X": "rock",
		"Y": "paper",
		"Z": "scissors",
	}
	moveScoreMap = map[string]int{
		"rock":     1,
		"paper":    2,
		"scissors": 3,
	}
	winningMoveMap = map[string]string{
		"rock":     "paper",
		"paper":    "scissors",
		"scissors": "rock",
	}
	losingMoveMap = map[string]string{
		"rock":     "scissors",
		"paper":    "rock",
		"scissors": "paper",
	}
)

func scoreMoves(ourMove, theirMove string) int {
	moveScore := moveScoreMap[ourMove]
	if ourMove == winningMoveMap[theirMove] {
		return moveScore + 6 // we won
	} else if ourMove == theirMove {
		return moveScore + 3 // we drew
	} else {
		return moveScore // we lost
	}
}

type RPSGame struct {
	EncryptedThem string
	EncryptedUs   string
}

func (g RPSGame) scorePart1() int {
	us := decryptionMap[g.EncryptedUs]
	them := decryptionMap[g.EncryptedThem]
	return scoreMoves(us, them)
}

func (g RPSGame) scorePart2() int {
	them := decryptionMap[g.EncryptedThem]
	us := ""
	switch g.EncryptedUs {
	case "X": // we need to lose
		us = losingMoveMap[them]
	case "Y": // we need to draw
		us = them
	case "Z": // we need to win
		us = winningMoveMap[them]
	}
	return scoreMoves(us, them)
}

func main() {
	input := utils.CheckErr(os.ReadFile("input.txt"))

	games := []RPSGame{}
	for _, line := range utils.Lines(string(input)) {
		encryptedMoves := strings.Fields(line)
		games = append(games, RPSGame{
			EncryptedThem: encryptedMoves[0],
			EncryptedUs:   encryptedMoves[1],
		})
	}

	part1 := utils.Sum(utils.Transform(games, RPSGame.scorePart1))
	part2 := utils.Sum(utils.Transform(games, RPSGame.scorePart2))

	fmt.Printf("The answer to Part 1 is %v.\n", part1)
	fmt.Printf("The answer to Part 2 is %v.\n", part2)
}
