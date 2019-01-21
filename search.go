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
	transpositionTable = make(transpositionMapping, 500000)
	hashMoveTable = make([]dt.Move, 512)
	killerOneTable = make([]dt.Move, 512)
	killerTwoTable = make([]dt.Move, 512)
	searching = true

	if movetime != -1 {
		endTime = time.Now().Add(time.Millisecond * time.Duration(movetime))
	}

	for i := 1; i < depth; i += 2 {
		if time.Now().Add(time.Duration(lastTime) * time.Millisecond).After(endTime) {
			break
		}
		maxDepth = i
		deepestQuiescence = 0
		t := time.Now()
		moveList := board.GenerateLegalMoves()
		sortMoves(moveList, board)

		val, bmv := negaMax(board, i, math.MinInt32, math.MaxInt32, moveList, true)
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
			pv = recoverPv(board, bestMove)
			for i, mv := range pv {
				hashMoveTable[getHalfMoveCount(board)+i] = mv
				outMoves += mv.String() + " "
			}
		} else {
			searching = false
		}
		lastTime = int(timeElapsed.Nanoseconds() / 1000000)
		fmt.Printf("info depth %d score cp %d time %d nodes %d\n", i, val, timeElapsed.Nanoseconds()/1000000, nodes)
		fmt.Fprintf(os.Stderr, "info depth %d score cp %d time %d nodes %d\n", i, val, timeElapsed.Nanoseconds()/1000000, nodes)
		fmt.Fprintln(os.Stderr, outMoves)
	}
	return valf, bestMove
}

func recoverPv(board *dt.Board, move dt.Move) []dt.Move {
	pvArray := []dt.Move{move}
	antiRepeat := map[uint64]int{board.Hash(): 1}
	copyBoard := *board
	board.ApplyNoFunc(move)
	for {
		entry, ok := transpositionTable.get(board)
		moveList := board.GenerateLegalMoves()
		_, ok2 := antiRepeat[board.Hash()]
		if ok2 {
			break
		}
		if ok == nil && isValidMove(entry.move, moveList) {
			pvArray = append(pvArray, entry.move)
			antiRepeat[board.Hash()] = 1
			board.ApplyNoFunc(entry.move)
		} else {
			break
		}
	}
	*board = copyBoard
	return pvArray
}

func pickReduction(remainingDepth int, moveCount int) int {
	if maxDepth-remainingDepth > 3 { // if we are at depth >=5
		if moveCount > 6 {
			return min(remainingDepth-1, max(remainingDepth/3, 1))
		}
		return min(remainingDepth-1, 1)

	}
	return 0
}

func negaMax(board *dt.Board, depth int, alpha, beta int, moveList []dt.Move, doNull bool) (int, dt.Move) {
	var bestMove dt.Move
	var v int
	var inCheck bool
	var ourPieces dt.Bitboards
	alphaOriginal := alpha
	trEntry, err := transpositionTable.get(board)

	nodes++

	if err == nil && trEntry.depth >= depth && isValidMove(trEntry.move, moveList) {
		switch trEntry.flag {
		case EXACT:
			return trEntry.value, trEntry.move
		case LOWERBOUND:
			alpha = max(alpha, trEntry.value)
		case UPPERBOUND:
			beta = min(beta, trEntry.value)
		}
		if alpha >= beta {
			return trEntry.value, trEntry.move
		}
	}

	updateTimer()
	if !searching {
		return -evalBoard(board, nil), 0
	}
	if board.Halfmoveclock >= 100 {
		return 0, 0
	}
	if board.Halfmoveclock > 1 {
		// Check for 3fold
		for i := 0; i < 4; i++ {
			if board.Last4Hashes[i] == board.Hash() {
				return 0, 0
			}
		}
	}
	if depth == 0 || len(moveList) == 0 {
		val, move := quiescenceSearch(board, alpha, beta, depth)
		return val, move
	}
	if board.OurKingInCheck() {
		depth++
		inCheck = true
	}
	if board.Wtomove {
		ourPieces = board.White
	} else {
		ourPieces = board.Black
	}

	if doNull && maxDepth != depth && depth >= 3 && !inCheck && (ourPieces.All^ourPieces.Pawns^ourPieces.Kings) > 0 {
		boardCopy := *board
		board.MakeNullMove()
		moveList := board.GenerateLegalMoves()

		val, move := negaMax(board, depth-1-2, -beta, -beta+1, moveList, false)
		val = -val

		*board = boardCopy
		if val >= beta {
			trEntry.value = beta
			trEntry.move = move
			trEntry.depth = depth
			trEntry.flag = LOWERBOUND
			transpositionTable.put(board, trEntry)
			return beta, move
		}
	}

	bSearchPv := true
	sortMoves(moveList, board)
	for moveCount, currMove := range moveList {
		boardCopy := *board
		board.ApplyNoFunc(currMove)
		moveList := board.GenerateLegalMoves()
		R := 0
		if !(moveCount < LMR_LIMIT || isInteresting(currMove, &boardCopy, board)) {
			R = pickReduction(depth, moveCount)
		}
		if bSearchPv {
			v, _ = negaMax(board, depth-1, -beta, -alpha, moveList, false)
		} else {
			v, _ = negaMax(board, depth-1-R, -alpha-1, -alpha, moveList, true)
			if -v > alpha {
				v, _ = negaMax(board, depth-1, -beta, -alpha, moveList, true)
			}
		}

		v = -v
		if v > alpha {
			alpha = v
			bestMove = currMove
			bSearchPv = false
		}
		*board = boardCopy

		if alpha >= beta {
			break
		}
	}

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

	return alpha, bestMove
}

func quiescenceSearch(board *dt.Board, alpha, beta, depth int) (int, dt.Move) {
	var val int
	var bestMove dt.Move
	var alphaOriginal int = alpha

	updateTimer()
	if !searching {
		return -evalBoard(board, nil), 0
	}
	if board.Halfmoveclock >= 100 {
		return 0, 0
	}
	if board.Halfmoveclock > 1 {
		// Check for 3fold
		for i := 0; i < 4; i++ {
			if board.Last4Hashes[i] == board.Hash() {
				return 0, 0
			}
		}
	}

	deepestQuiescence = min(depth, deepestQuiescence)
	isCheck := board.OurKingInCheck()

	if !isCheck {
		val = evalBoard(board, nil)
		if val >= beta {
			return val, 0
		}
		queenValue := pieceVal[dt.Queen]
		if val < alpha-queenValue {
			return alpha, 0
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
			return -MAXVALUE, 0
		}
		for _, move := range moves {
			heap.Push(&pq, &moveValPair{val: 0, move: move})
		}
	} else {
		for _, move := range moves {
			if dt.IsCapture(move, board) {
				captureValue := getCaptureValue(board, move)
				if val+captureValue+200 >= alphaOriginal {
					heap.Push(&pq, &moveValPair{val: captureValue, move: move})
				}
			}
		}
	}

	for pq.Len() > 0 {
		nodes++
		mvP := heap.Pop(&pq).(*moveValPair)
		copyBoard := *board
		board.ApplyNoFunc(mvP.move)

		if !isCheck {
			if board.OurKingInCheck() {
				*board = copyBoard
				continue
			}
		}

		val, _ := quiescenceSearch(board, -beta, -alpha, depth-1)
		val = -val
		*board = copyBoard

		if val >= beta {
			return beta, mvP.move
		}
		if val > alpha {
			alpha = val
			bestMove = mvP.move
		}
	}
	return alpha, bestMove

}
