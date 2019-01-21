package main

import (
	"sort"

	dt "github.com/dylhunn/dragontoothmg"
)

type moveValue struct {
	val  int
	move dt.Move
}

func getMoveValue(move dt.Move, board *dt.Board) int {
	if hashMoveTable[getHalfMoveCount(board)] == move {
		return MAXVALUE
	}
	if dt.IsCapture(move, board) {
		return getCaptureValue(board, move)
	}
	piece := move.Promote()
	if piece != dt.Nothing {
		if piece == dt.Queen {
			return 11
		}
		return 6
	}
	if killerOneTable[getHalfMoveCount(board)] == move {
		return 10
	}
	if killerTwoTable[getHalfMoveCount(board)] == move {
		return 8
	}
	if dt.IsCastle(move, board) {
		return 7
	}

	return 0
}

func isInteresting(move dt.Move, board *dt.Board, newBoard *dt.Board) bool {
	if board.OurKingInCheck() {
		return true
	}

	if newBoard.OurKingInCheck() {
		return true
	}

	return getMoveValue(move, board) > 10
}

func sortMoves(moveList []dt.Move, board *dt.Board) {
	tuples := make([]moveValue, len(moveList))

	for i, mv := range moveList {
		val := getMoveValue(mv, board)
		tuples[i] = moveValue{val: val, move: mv}
	}

	sort.Slice(tuples, func(i, j int) bool { // czy i < j
		return tuples[i].val > tuples[j].val
	})

	for i := range moveList {
		moveList[i] = tuples[i].move
	}
}
