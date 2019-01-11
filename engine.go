package main

import (
	"container/heap"
	"errors"
	"fmt"
	"math"
	"math/bits"
	"os"
	"sort"
	"time"

	dt "github.com/dylhunn/dragontoothmg"
)

const MAXVALUE int = 100000
const MINVALUE int = -100000
const LMR_LIMIT = 6

var pieceTypesnoking []int
var pieceTypes []int

var pieceAttackUnits []int
var kingSafety []int

var pawns []int
var pawnsBlack []int

var knights []int
var knightsBlack []int

var bishops []int
var bishopsBlack []int

var rooks []int
var rooksBlack []int

var queens []int
var queensBlack []int

var kings []int
var kingsBlack []int

var kingMiddlegame []int
var kingMiddlegameBlack []int

var kingEndgame []int
var kingEndgameBlack []int

var pieceVal []int
var attackSquareVal []int

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

func init() {
	pieceTypesnoking = []int{dt.Pawn, dt.Knight, dt.Bishop, dt.Rook, dt.Queen}
	pieceTypes = append(pieceTypesnoking, dt.King)

	pieceAttackUnits = []int{0, 2, 2, 5, 3, 1, 0}

	kingSafety = []int{
		0, 0, 1, 2, 3, 5, 7, 9, 12, 15,
		18, 22, 26, 30, 35, 39, 44, 50, 56, 62,
		68, 75, 82, 85, 89, 97, 105, 113, 122, 131,
		140, 150, 169, 180, 191, 202, 213, 225, 237, 248,
		260, 272, 283, 295, 307, 319, 330, 342, 354, 366,
		377, 389, 401, 412, 424, 436, 448, 459, 471, 483,
		494, 500, 500, 500, 500, 500, 500, 500, 500, 500,
		500, 500, 500, 500, 500, 500, 500, 500, 500, 500,
		500, 500, 500, 500, 500, 500, 500, 500, 500, 500,
		500, 500, 500, 500, 500, 500, 500, 500, 500, 500,
	}

	pawnsBlack = []int{0, 0, 0, 0, 0, 0, 0, 0,
		50, 50, 50, 50, 50, 50, 50, 50,
		10, 10, 20, 30, 30, 20, 10, 10,
		5, 5, 10, 25, 25, 10, 5, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		5, -5, -10, 0, 0, -10, -5, 5,
		5, 10, 10, -20, -20, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0}

	pawns = reverse(pawnsBlack)

	knightsBlack = []int{-50, -40, -30, -30, -30, -30, -40, -50,
		-40, -20, 0, 0, 0, 0, -20, -40,
		-30, 0, 10, 15, 15, 10, 0, -30,
		-30, 5, 15, 20, 20, 15, 5, -30,
		-30, 0, 15, 20, 20, 15, 0, -30,
		-30, 5, 10, 15, 15, 10, 5, -30,
		-40, -20, 0, 5, 5, 0, -20, -40,
		-50, -40, -30, -30, -30, -30, -40, -50}

	knights = reverse(knightsBlack)

	bishopsBlack = []int{-20, -10, -10, -10, -10, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 5, 5, 10, 10, 5, 5, -10,
		-10, 0, 10, 10, 10, 10, 0, -10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		-10, 5, 0, 0, 0, 0, 5, -10,
		-20, -10, -10, -10, -10, -10, -10, -20}

	bishops = reverse(bishopsBlack)

	rooksBlack = []int{0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, 10, 10, 10, 10, 5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		0, 0, 0, 5, 5, 0, 0, 0}

	rooks = reverse(rooksBlack)

	queensBlack = []int{-20, -10, -10, -5, -5, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 5, 5, 5, 0, -10,
		-5, 0, 5, 5, 5, 5, 0, -5,
		0, 0, 5, 5, 5, 5, 0, -5,
		-10, 5, 5, 5, 5, 5, 0, -10,
		-10, 0, 5, 0, 0, 0, 0, -10,
		-20, -10, -10, -5, -5, -10, -10, -20}

	queens = reverse(queensBlack)

	kingMiddlegameBlack = []int{-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-20, -30, -30, -40, -40, -30, -30, -20,
		-10, -20, -20, -20, -20, -20, -20, -10,
		20, 20, 0, 0, 0, 0, 20, 20,
		20, 30, 10, 0, 0, 10, 30, 20}

	kingMiddlegame = reverse(kingMiddlegameBlack)

	kingEndgameBlack = []int{-50, -40, -30, -20, -20, -30, -40, -50,
		-30, -20, -10, 0, 0, -10, -20, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -30, 0, 0, 0, 0, -30, -30,
		-50, -30, -30, -30, -30, -30, -30, -50}

	kingEndgame = reverse(kingEndgameBlack)

	pieceVal = []int{0, 100, 320, 330, 500, 935, 0}

	attackSquareVal = []int{0, 1, 4, 2, 2, 2, 0}
}

var TIMECHECK_FREQ int = 5000
var isEndgame = false
var nodes int
var deepestQuiescence int
var timeCheckCounter = TIMECHECK_FREQ
var endTime = time.Now().AddDate(1000, 10, 10)
var searching = true

type transpositionFlag int

const (
	// EXACT value from search
	EXACT transpositionFlag = iota
	// LOWERBOUND alpha from search
	LOWERBOUND
	// UPPERBOUND beta from search
	UPPERBOUND
)

type transpositionEntry struct {
	value int
	depth int
	move  dt.Move
	flag  transpositionFlag
}

type transpositionMapping map[uint64]transpositionEntry

var transpositionTable transpositionMapping

var errNoTranspositionEntry = errors.New("No entry")
var hashMoveTable []dt.Move
var killerOneTable []dt.Move
var killerTwoTable []dt.Move

func (t transpositionMapping) put(board *dt.Board, trEntry transpositionEntry) {
	h := board.Hash()
	t[h] = trEntry
}

func (t transpositionMapping) get(board *dt.Board) (transpositionEntry, error) {
	h := board.Hash()
	entry, ok := t[h]

	if !ok {
		return transpositionEntry{}, errNoTranspositionEntry
	}
	return entry, nil
}

var maxDepth int

func search(board *dt.Board, depth int, movetime int) (float64, dt.Move) {
	// check if endgame and set appproproeirpoeporiylu
	// isEndgame =...
	nodes = 0
	valf := 0.0
	transpositionTable = make(transpositionMapping, 5000000)
	hashMoveTable = make([]dt.Move, 512)
	killerOneTable = make([]dt.Move, 512)
	killerTwoTable = make([]dt.Move, 512)
	searching = true
	maxDepth = depth
	var outMoves string
	var pv []dt.Move
	if movetime != -1 {
		endTime = time.Now().Add(time.Millisecond * time.Duration(movetime))
	}

	var bestMove dt.Move

	for i := 1; i < depth; i++ {
		deepestQuiescence = 0
		t := time.Now()
		moveList := board.GenerateLegalMoves()
		sortMoves(moveList, board)

		//fmt.Fprintf(os.Stderr, "depth %d\n", i)
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
func addKiller(move dt.Move, depth int) {
	if killerOneTable[depth] == 0 {
		killerOneTable[depth] = move
	} else if move != killerOneTable[depth] {
		killerTwoTable[depth] = move
	}
}

type moveValue struct {
	val  int
	move dt.Move
}

func getHalfMoveCount(board *dt.Board) int {
	halfMove := 0
	if board.Wtomove == false {
		halfMove = 1
	}
	return int(board.Fullmoveno) + halfMove
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
		return pieceVal[dt.Pawn] - pieceVal[piece] + 500
	}
	if killerOneTable[getHalfMoveCount(board)] == move {
		return 10
	}
	if killerTwoTable[getHalfMoveCount(board)] == move {
		return 8
	}
	return 0
}
func isInteresting(move dt.Move, board *dt.Board) bool {
	if board.OurKingInCheck() {
		return true
	}
	unApply := board.Apply(move)
	if board.OurKingInCheck() {
		unApply()
		return true
	}
	unApply()
	return getMoveValue(move, board) > 0
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
	nodes++
	alphaOriginal := alpha
	trEntry, err := transpositionTable.get(board)
	if err == nil && trEntry.depth >= depth && isValidMove(trEntry.move, moveList) {
		unApply := board.Apply(trEntry.move)
		switch trEntry.flag {
		case EXACT:
			unApply()
			return trEntry.value, trEntry.move, []dt.Move{trEntry.move}
		case LOWERBOUND:
			alpha = max(alpha, trEntry.value)
		case UPPERBOUND:
			beta = min(beta, trEntry.value)
		}
		unApply()
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
	var bestMove dt.Move
	var tpv []dt.Move
	var bestTtpv []dt.Move
	var v int
	var ttpv []dt.Move

	sortMoves(moveList, board)
	for moveCount, currMove := range moveList {
		boardCopy := *board
		board.Apply(currMove)
		moveList := board.GenerateLegalMoves()

		if moveCount < LMR_LIMIT || isInteresting(currMove, &boardCopy) {
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

func evalBoard(board *dt.Board, moveList []dt.Move) int {
	if board.OurKingInCheck() {
		if moveList != nil && len(moveList) == 0 {
			return MINVALUE
		}
	}
	v := 0

	v += (bits.OnesCount64(board.White.Pawns) - bits.OnesCount64(board.Black.Pawns)) * pieceVal[dt.Pawn]
	v += (bits.OnesCount64(board.White.Knights) - bits.OnesCount64(board.Black.Knights)) * pieceVal[dt.Knight]
	v += (bits.OnesCount64(board.White.Bishops) - bits.OnesCount64(board.Black.Bishops)) * pieceVal[dt.Bishop]
	v += (bits.OnesCount64(board.White.Rooks) - bits.OnesCount64(board.Black.Rooks)) * pieceVal[dt.Rook]
	v += (bits.OnesCount64(board.White.Queens) - bits.OnesCount64(board.Black.Queens)) * pieceVal[dt.Queen]

	tmp := board.White.Pawns
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v += pawns[idx]
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.Black.Pawns
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v -= pawnsBlack[idx]
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.White.Bishops
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v += bishops[idx]
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.Black.Bishops
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v -= bishopsBlack[idx]
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.White.Knights
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v += knights[idx]
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.Black.Knights
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v -= knightsBlack[idx]
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.White.Rooks
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v += rooks[idx]
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.Black.Rooks
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v -= rooksBlack[idx]
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.White.Queens
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v += queens[idx]
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.Black.Queens
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v -= queensBlack[idx]
		tmp &= ^(1 << uint(idx))
	}

	whiteKing := board.White.Kings
	whiteKingIdx := bits.TrailingZeros64(whiteKing)

	blackKing := board.Black.Kings
	blackKingIdx := bits.TrailingZeros64(blackKing)

	if isEndgame {
		v += kingEndgame[whiteKingIdx]
		v -= kingEndgameBlack[blackKingIdx]
	} else {
		v += kingMiddlegame[whiteKingIdx]
		v -= kingMiddlegameBlack[blackKingIdx]
	}
	return v * getColorMutliplier(board.Wtomove)

}

func getCaptureValue(board *dt.Board, move dt.Move) int {
	var ourBitboard *dt.Bitboards
	var theirBitboard *dt.Bitboards
	var theirVal int
	var candidateMove dt.Move
	var isTsqLastLine bool
	ourColor := board.Wtomove
	fsq := move.From()
	tsq := move.To()

	fromBitboard := (uint64(1) << fsq)
	toBitboard := (uint64(1) << tsq)

	if board.Wtomove {
		ourBitboard = &board.White
		theirBitboard = &board.Black
		isTsqLastLine = toBitboard&dt.OnlyRank[7] != 0
	} else {
		ourBitboard = &board.Black
		theirBitboard = &board.White
		isTsqLastLine = toBitboard&dt.OnlyRank[0] != 0
	}

	ourPieceType, _ := dt.DeterminePieceType(ourBitboard, fromBitboard)
	theirPieceType, _ := dt.DeterminePieceType(theirBitboard, toBitboard)
	ourVal := pieceVal[ourPieceType]

	if theirPieceType == dt.Nothing { // handle en passant -> capture on empty square
		theirVal = pieceVal[dt.Pawn]
	} else {
		theirVal = pieceVal[theirPieceType]
	}

	if theirVal > ourVal {
		return theirVal - ourVal
	}

	// Create copy of the board instead of unapplying everything in sequence
	boardCopy := *board
	board.ApplyNoGoingBackBadHash(move)

	swaplist := make([]int, 0, 10)
	swaplist = append(swaplist, theirVal)

	swapCount := 0
	lastVal := ourVal

	for {
		attackers := board.GetAttackersForSquare(!board.Wtomove, tsq)
		if attackers.All == 0 {
			break
		}
		// We want to capture with lowest valued pieces first
		attackerPieceType, lowestValueAttackerBitboard := dt.LowestValuePiece(&attackers)
		lowestValueAttackerSquare := dt.Square(bits.TrailingZeros64(*lowestValueAttackerBitboard)) // 0 = A1 ... 63 - H8

		// Create move to use in capturing
		candidateMove.Setfrom(lowestValueAttackerSquare)
		candidateMove.Setto(dt.Square(tsq))

		if attackerPieceType == dt.Pawn && isTsqLastLine { // set promotion if applicable
			candidateMove.Setpromote(dt.Queen)
		} else {
			candidateMove.Setpromote(dt.Nothing)
		}
		board.ApplyNoGoingBackBadHash(candidateMove)
		swapCount++

		swaplist = append(swaplist, swaplist[len(swaplist)-1]+lastVal)
		lastVal = pieceVal[attackerPieceType] * getColorMutliplier(ourColor == board.Wtomove)
	}
	// restore board
	*board = boardCopy

	for i := len(swaplist) - 1; i > 0; i-- {
		if i&1 != 0 {
			if swaplist[i] <= swaplist[i-1] {
				swaplist[i-1] = swaplist[i]
			}
		} else {
			if swaplist[i] >= swaplist[i-1] {
				swaplist[i-1] = swaplist[i]
			}
		}
	}
	if swaplist[0] < 0 {
		return MINVALUE
	} else {
		return swaplist[0]
	}

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
func quiescenceSearch(board *dt.Board, alpha, beta, depth int) (int, dt.Move, []dt.Move) {
	updateTimer()
	if !searching {
		return -evalBoard(board, nil), 0, []dt.Move{}
	}
	deepestQuiescence = min(depth, deepestQuiescence)
	isCheck := board.OurKingInCheck()
	var val int
	//var unApplyFunc func()
	var bestTpv []dt.Move
	var bestMove dt.Move

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

	moves := board.GenerateLegalMoves()
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

func getColorMutliplier(color bool) int {
	if color {
		return 1
	}
	return -1
}
