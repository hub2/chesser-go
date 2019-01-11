package main

import (
	"time"

	dt "github.com/dylhunn/dragontoothmg"
)

func reverse(array []int) []int {
	newArray := make([]int, len(array))
	for i, Val := range array {
		newArray[len(array)-i-1] = Val
	}
	return newArray
}

func reverseMove(array []dt.Move) []dt.Move {
	newArray := make([]dt.Move, len(array))
	for i, Val := range array {
		newArray[len(array)-i-1] = Val
	}
	return newArray
}

func getHalfMoveCount(board *dt.Board) int {
	halfMove := 0
	if board.Wtomove == false {
		halfMove = 1
	}
	return int(board.Fullmoveno) + halfMove
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func isValidMove(move dt.Move, moveList []dt.Move) bool {
	for _, mv := range moveList {
		if move == mv {
			return true
		}
	}
	return false
}
func updateTimer() {
	timeCheckCounter--
	if timeCheckCounter == 0 {
		if time.Now().After(endTime) {
			searching = false
		}
		timeCheckCounter = TIMECHECK_FREQ
	}
}

func getColorMutliplier(color bool) int {
	if color {
		return 1
	}
	return -1
}

func nortFill(gen uint64) uint64 {
	gen |= (gen << 8)
	gen |= (gen << 16)
	gen |= (gen << 32)
	return gen
}

func soutFill(gen uint64) uint64 {
	gen |= (gen >> 8)
	gen |= (gen >> 16)
	gen |= (gen >> 32)
	return gen
}
