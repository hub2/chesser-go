package main

import (
	"container/heap"
	"fmt"
	"math"
	"os"
	"sync/atomic"
	"time"

	dt "github.com/dylhunn/dragontoothmg"
)

func search(board *dt.Board, depth int, movetime int) (float64, dt.Move) {
	// check if endgame and set appproproeirpoeporiylu
	// isEndgame =...
	var outMoves string
	var pv []dt.Move
	var bestMove dt.Move

	valf := 0.0
	searching = true

	if movetime != -1 {
		endTime = time.Now().Add(time.Millisecond * time.Duration(movetime))
	}
	hits := 0
	for i := 1; i < depth; i++ {
		maxDepth = i
		deepestQuiescence = 0
		t := time.Now()
		moveList := board.GenerateLegalMoves()
		sortMoves(moveList, board)

		val, bmv, tpv := negaMax(board, i, math.MinInt32, math.MaxInt32, moveList, &hits)
		timeElapsed := time.Since(t)

		// dont return not fully searched tree
		// force check
		timeCheckCounter = 1
		updateTimer()
		if !searching {
			break
		}
		valf = float64(getColorMutliplier(board.Wtomove)*val) / 100.0
		if bmv != 0 {
			outMoves = ""
			bestMove = bmv
			pv = reverseMove(tpv)
			for i, mv := range pv {
				// if mv != 0{
				hashMoveTable[getHalfMoveCount(board)+i] = mv
				outMoves += mv.String() + " "
				// }
			}
		} else {
			searching = false
		}
		fmt.Printf("info depth %d score cp %d time %d nodes %d hits %d\n", i, val, timeElapsed.Nanoseconds()/1000000, nodes, hits)
		fmt.Fprintf(os.Stderr, "info depth %d score cp %d time %d nodes %d hits %d\n", i, val, timeElapsed.Nanoseconds()/1000000, nodes, hits)
		fmt.Fprintln(os.Stderr, outMoves)
	}
	return valf, bestMove
}

func pickReduction(remainingDepth int, moveCount int) int {
	if maxDepth-remainingDepth > 3 { // if we are at depth >=5
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

func negaMax(board *dt.Board, depth int, alpha, beta int, moveList []dt.Move, hits *int) (int, dt.Move, []dt.Move) {
	atomic.AddUint64(&nodes, 1)
	var bestMove dt.Move
	var tpv []dt.Move
	var bestTtpv []dt.Move
	var v int
	var ttpv []dt.Move

	alphaOriginal := alpha
	trEntry, err := transpositionTable.get(board)

	if err == nil && trEntry.depth >= depth && isValidMove(trEntry.move, moveList) {
		*hits++
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
		val, move, tpv := quiescenceSearch(board, alpha, beta, depth, moveList, hits)
		return val, move, tpv
	}
	bSearchPv := true
	sortMoves(moveList, board)
	for moveCount, currMove := range moveList {
		boardCopy := *board
		board.ApplyNoFunc(currMove)
		moveList := board.GenerateLegalMoves()
		kingCheckDepthBonus := 0
		if board.OurKingInCheck() {
			kingCheckDepthBonus = 1
		}

		if moveCount < LMR_LIMIT || isInteresting(currMove, &boardCopy, board) {
			if bSearchPv {
				v, _, ttpv = negaMax(board, depth-1+kingCheckDepthBonus, -beta, -alpha, moveList, hits)
			} else {
				v, _, ttpv = negaMax(board, depth-1+kingCheckDepthBonus, -alpha-1, -alpha, moveList, hits)
				if -v > alpha {
					v, _, ttpv = negaMax(board, depth-1+kingCheckDepthBonus, -beta, -alpha, moveList, hits)
				}
			}
		} else {
			R := pickReduction(depth, moveCount)
			v, _, ttpv = negaMax(board, depth-1-R, -alpha-1, -alpha, moveList, hits)
			if -v > alpha {
				v, _, ttpv = negaMax(board, depth-1, -beta, -alpha, moveList, hits)
			}
		}

		v = -v
		if v > alpha {
			alpha = v
			bestMove = currMove
			bestTtpv = ttpv
			bSearchPv = false
		}
		*board = boardCopy

		if alpha >= beta {
			break
		}
	}
	tpv = append(bestTtpv, bestMove)

	trEntry.value = alpha
	trEntry.move = bestMove
	trEntry.depth = depth
	if alpha <= alphaOriginal {
		trEntry.flag = UPPERBOUND
	} else if alpha >= beta {
		trEntry.flag = LOWERBOUND
		if !dt.IsCapture(bestMove, board) && bestMove.Promote() == dt.Nothing {
			addKiller(bestMove, getHalfMoveCount(board))
		}
	} else {
		trEntry.flag = EXACT
	}
	transpositionTable.put(board, trEntry)

	return alpha, bestMove, tpv
}

func quiescenceSearch(board *dt.Board, alpha, beta, depth int, moves []dt.Move, hits *int) (int, dt.Move, []dt.Move) {
	var val int
	var bestTpv []dt.Move
	var bestMove dt.Move

	alphaOriginal := alpha
	trEntry, err := transpositionTable.get(board)

	if err == nil && trEntry.depth >= -depth && isValidMove(trEntry.move, moves) {
		*hits++
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

	if board.Halfmoveclock >= 100 {
		return 0, 0, []dt.Move{}
	}
	if board.Halfmoveclock > 1 {
		// Check for 3fold
		for i := 0; i < 4; i++ {
			if board.Last4Hashes[i] == board.Hash() {
				return 0, 0, []dt.Move{}
			}
		}
	}

	updateTimer()
	if !searching {
		return -evalBoard(board, nil), 0, []dt.Move{}
	}

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
		atomic.AddUint64(&nodes, 1)
		mvP := heap.Pop(&pq).(*moveValPair)
		copyBoard := *board
		board.ApplyNoFunc(mvP.move)

		if !isCheck {
			if board.OurKingInCheck() {
				*board = copyBoard
				continue
			}
		}

		newMoves := board.GenerateLegalMoves()
		val, _, tpv := quiescenceSearch(board, -beta, -alpha, depth-1, newMoves, hits)
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

	trEntry.value = alpha
	trEntry.move = bestMove
	trEntry.depth = -depth
	if alpha <= alphaOriginal {
		trEntry.flag = UPPERBOUND
	} else if alpha >= beta {
		trEntry.flag = LOWERBOUND
		if !dt.IsCapture(bestMove, board) && bestMove.Promote() == dt.Nothing {
			addKiller(bestMove, getHalfMoveCount(board))
		}
	} else {
		trEntry.flag = EXACT
	}
	transpositionTable.put(board, trEntry)

	if bestTpv != nil {
		return alpha, bestMove, append(bestTpv, bestMove)
	}
	return alpha, 0, []dt.Move{}

}
