package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
	"time"

	dt "github.com/dylhunn/dragontoothmg"
)

func search(board *dt.Board, depth int, movetime int) (float64, dt.Move) {
	// check if endgame and set appproproeirpoeporiylu
	// isEndgame =...
	var outMoves string
	var pv []dt.Move
	var bestMove dt.Move

	nodes = 0
	valf := 0.0
	transpositionTable = make(transpositionMapping, 5000000)
	hashMoveTable = make([]dt.Move, 512)
	killerOneTable = make([]dt.Move, 512)
	killerTwoTable = make([]dt.Move, 512)
	searching = true
	maxDepth = depth

	if movetime != -1 {
		endTime = time.Now().Add(time.Millisecond * time.Duration(movetime))
	}

	for i := 1; i < depth; i++ {
		deepestQuiescence = 0
		t := time.Now()
		moveList := board.GenerateLegalMoves()
		sortMoves(moveList, board)

		val, bmv, tpv := negaMax(board, i, math.MinInt32, math.MaxInt32, moveList)
		timeElapsed := time.Since(t)

		// dont return not fully searched tree
		// force check
		timeCheckCounter = 1
		updateTimer()
		if !searching {
			break
		}
		valf = float64(val) / 100.0
		if bmv != 0 {
			outMoves = ""
			bestMove = bmv
			pv = reverseMove(tpv)
			for i, mv := range pv {
				hashMoveTable[getHalfMoveCount(board)+i] = mv
				outMoves += mv.String() + " "
			}
		} else {
			searching = false
		}

		fmt.Printf("info depth %d score %.2f time %d nodes %d\n", i, valf, timeElapsed.Nanoseconds()/1000000, nodes)
		fmt.Fprintf(os.Stderr, "info depth %d/%d score %.2f time %d nodes %d\n", i, depth-deepestQuiescence, valf, timeElapsed.Nanoseconds()/1000000, nodes)
		fmt.Fprintln(os.Stderr, outMoves)
	}
	return valf, bestMove
}

func pickReduction(remainingDepth int, moveCount int) int {
	if maxDepth-remainingDepth > 4 { // if we are at depth >=5
		if moveCount > 20 {
			return (4 * remainingDepth) / 5
		}
		if moveCount > 10 {
			return (2 * remainingDepth) / 3
		}
		return remainingDepth / 3

	}
	return 0
}

func negaMax(board *dt.Board, depth int, alpha, beta int, moveList []dt.Move) (int, dt.Move, []dt.Move) {
	var bestMove dt.Move
	var tpv []dt.Move
	var bestTtpv []dt.Move
	var v int
	var ttpv []dt.Move

	alphaOriginal := alpha
	trEntry, err := transpositionTable.get(board)

	nodes++

	if err == nil && trEntry.depth >= depth && isValidMove(trEntry.move, moveList) {
		switch trEntry.flag {
		case EXACT:
			return trEntry.value, trEntry.move, []dt.Move{trEntry.move}
		case LOWERBOUND:
			alpha = max(alpha, trEntry.value)
		case UPPERBOUND:
			beta = min(beta, trEntry.value)
		}
		if alpha >= beta {
			return trEntry.value, trEntry.move, []dt.Move{trEntry.move}
		}
	}

	updateTimer()
	if !searching {
		return -evalBoard(board, nil), 0, []dt.Move{}
	}

	if depth == 0 || len(moveList) == 0 {
		val, move, tpv := quiescenceSearch(board, alpha, beta, depth)
		return val, move, tpv
	}

	vMax := MINVALUE

	sortMoves(moveList, board)
	for moveCount, currMove := range moveList {
		boardCopy := *board
		board.ApplyNoFunc(currMove)
		moveList := board.GenerateLegalMoves()

		if moveCount < LMR_LIMIT || isInteresting(currMove, &boardCopy, board) {
			v, _, ttpv = negaMax(board, depth-1, -beta, -alpha, moveList)
		} else {
			R := pickReduction(depth, moveCount)
			v, _, ttpv = negaMax(board, depth-1-R, -beta, -alpha, moveList)
			if -v > alpha {
				v, _, ttpv = negaMax(board, depth-1, -beta, -alpha, moveList)
			}
		}

		v = -v

		v = max(alpha, v)
		alpha = v

		if v > vMax {
			vMax = v
			bestMove = currMove
			bestTtpv = ttpv
		}
		*board = boardCopy

		if alpha >= beta {
			break
		}
	}
	tpv = append(bestTtpv, bestMove)

	trEntry.value = vMax
	trEntry.move = bestMove
	trEntry.depth = depth
	if vMax <= alphaOriginal {
		trEntry.flag = UPPERBOUND
	} else if vMax >= beta {
		trEntry.flag = LOWERBOUND
		if !dt.IsCapture(bestMove, board) && bestMove.Promote() == dt.Nothing {
			addKiller(bestMove, getHalfMoveCount(board))
		}
	} else {
		trEntry.flag = EXACT
	}
	transpositionTable.put(board, trEntry)

	return vMax, bestMove, tpv

}

func quiescenceSearch(board *dt.Board, alpha, beta, depth int) (int, dt.Move, []dt.Move) {
	var val int
	var bestTpv []dt.Move
	var bestMove dt.Move

	updateTimer()
	if !searching {
		return -evalBoard(board, nil), 0, []dt.Move{}
	}

	deepestQuiescence = min(depth, deepestQuiescence)
	isCheck := board.OurKingInCheck()

	if !isCheck {
		val = evalBoard(board, nil)
		if val >= beta {
			return val, 0, []dt.Move{}
		}
		if val > alpha {
			alpha = val
		}
	}

	moves := board.GenerateLegalMoves()
	pq := make(PriorityQueue, 0, 40)
	heap.Init(&pq)

	if isCheck {
		if len(moves) == 0 {
			return -MAXVALUE, 0, []dt.Move{}
		}
		for _, move := range moves {
			heap.Push(&pq, &moveValPair{val: 0, move: move})
		}
	} else {
		for _, move := range moves {
			if dt.IsCapture(move, board) {
				heap.Push(&pq, &moveValPair{val: getCaptureValue(board, move), move: move})
			}
		}
	}

	for pq.Len() > 0 {
		nodes++
		mvP := heap.Pop(&pq).(*moveValPair)
		copyBoard := *board
		board.ApplyNoGoingBackBadHash(mvP.move)

		if !isCheck {
			if board.OurKingInCheck() {
				*board = copyBoard
				continue
			}
		}

		val, _, tpv := quiescenceSearch(board, -beta, -alpha, depth-1)
		val = -val
		*board = copyBoard

		if val >= beta {
			return beta, mvP.move, append(tpv, mvP.move)
		}
		if val > alpha {
			alpha = val
			bestTpv = tpv
			bestMove = mvP.move
		}
	}
	if bestTpv != nil {
		return alpha, bestMove, append(bestTpv, bestMove)
	} else {
		return alpha, 0, []dt.Move{}
	}

}
