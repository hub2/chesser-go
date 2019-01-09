package main

import (
	"errors"
	"fmt"
	"math"
	"math/bits"
	"sort"
	"time"

	dt "github.com/dylhunn/dragontoothmg"
)

const MAXVALUE int = 100000
const MINVALUE int = -100000

var pieceTypesnoking []int
var pieceTypes []int

var pieceAttackUnits map[int]int
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

var pieceVal map[int]int
var attackSquareVal map[int]int

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

	pieceAttackUnits = map[int]int{
		dt.Knight: 2,
		dt.Bishop: 2,
		dt.Queen:  5,
		dt.Rook:   3,
		dt.Pawn:   1,
		dt.King:   0}

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

	pawns = []int{0, 0, 0, 0, 0, 0, 0, 0,
		50, 50, 50, 50, 50, 50, 50, 50,
		10, 10, 20, 30, 30, 20, 10, 10,
		5, 5, 10, 25, 25, 10, 5, 5,
		0, 0, 0, 20, 20, 0, 0, 0,
		5, -5, -10, 0, 0, -10, -5, 5,
		5, 10, 10, -20, -20, 10, 10, 5,
		0, 0, 0, 0, 0, 0, 0, 0}

	pawnsBlack = reverse(pawns)

	knights = []int{-50, -40, -30, -30, -30, -30, -40, -50,
		-40, -20, 0, 0, 0, 0, -20, -40,
		-30, 0, 10, 15, 15, 10, 0, -30,
		-30, 5, 15, 20, 20, 15, 5, -30,
		-30, 0, 15, 20, 20, 15, 0, -30,
		-30, 5, 10, 15, 15, 10, 5, -30,
		-40, -20, 0, 5, 5, 0, -20, -40,
		-50, -40, -30, -30, -30, -30, -40, -50}

	knightsBlack = reverse(knights)

	bishops = []int{-20, -10, -10, -10, -10, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 10, 10, 5, 0, -10,
		-10, 5, 5, 10, 10, 5, 5, -10,
		-10, 0, 10, 10, 10, 10, 0, -10,
		-10, 10, 10, 10, 10, 10, 10, -10,
		-10, 5, 0, 0, 0, 0, 5, -10,
		-20, -10, -10, -10, -10, -10, -10, -20}

	bishopsBlack = reverse(bishops)

	rooks = []int{0, 0, 0, 0, 0, 0, 0, 0,
		5, 10, 10, 10, 10, 10, 10, 5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		-5, 0, 0, 0, 0, 0, 0, -5,
		0, 0, 0, 5, 5, 0, 0, 0}

	rooksBlack = reverse(rooks)

	queens = []int{-20, -10, -10, -5, -5, -10, -10, -20,
		-10, 0, 0, 0, 0, 0, 0, -10,
		-10, 0, 5, 5, 5, 5, 0, -10,
		-5, 0, 5, 5, 5, 5, 0, -5,
		0, 0, 5, 5, 5, 5, 0, -5,
		-10, 5, 5, 5, 5, 5, 0, -10,
		-10, 0, 5, 0, 0, 0, 0, -10,
		-20, -10, -10, -5, -5, -10, -10, -20}

	queensBlack = reverse(queens)

	kingMiddlegame = []int{-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-30, -40, -40, -50, -50, -40, -40, -30,
		-20, -30, -30, -40, -40, -30, -30, -20,
		-10, -20, -20, -20, -20, -20, -20, -10,
		20, 20, 0, 0, 0, 0, 20, 20,
		20, 30, 10, 0, 0, 10, 30, 20}

	kingMiddlegameBlack = reverse(kingMiddlegame)

	kingEndgame = []int{-50, -40, -30, -20, -20, -30, -40, -50,
		-30, -20, -10, 0, 0, -10, -20, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 30, 40, 40, 30, -10, -30,
		-30, -10, 20, 30, 30, 20, -10, -30,
		-30, -30, 0, 0, 0, 0, -30, -30,
		-50, -30, -30, -30, -30, -30, -30, -50}

	kingEndgameBlack = reverse(kingEndgame)
	pieceVal = map[int]int{
		dt.Pawn:   100,
		dt.Knight: 320,
		dt.Bishop: 330,
		dt.Rook:   500,
		dt.Queen:  935,
		dt.King:   0}

	attackSquareVal = map[int]int{
		dt.Pawn:   1,
		dt.Knight: 4,
		dt.Bishop: 2,
		dt.Rook:   2,
		dt.Queen:  2,
		dt.King:   0}

}

var isEndgame = false
var nodes int = 0

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

func search(board *dt.Board, depth int) {
	// check if endgame and set appproproeirpoeporiylu
	// isEndgame =...
	nodes = 0
	valf := 0.0
	transpositionTable = make(transpositionMapping, 5000000)
	hashMoveTable = make([]dt.Move, 512)

	for i := 1; i < depth; i++ {
		t := time.Now()
		val, _, tpv := negaMax(board, i, math.MinInt32, math.MaxInt32)
		timeElapsed := time.Since(t)

		valf = float64(val) / 100.0
		pv := reverseMove(tpv)
		outMoves := ""
		halfMove := 0
		if board.Wtomove == false {
			halfMove = 1
		}
		for i, mv := range pv {
			hashMoveTable[int(board.Fullmoveno)+i+halfMove] = mv
			outMoves += mv.String() + " "
		}
		fmt.Printf("depth %d val %.2f time %v nodes %d\n", i, valf, timeElapsed, nodes)
		fmt.Println(outMoves)
	}
}

type moveValue struct {
	val  int
	move dt.Move
}

func getMoveValue(move dt.Move, board *dt.Board) int {
	halfMove := 0
	if board.Wtomove == false {
		halfMove = 1
	}

	if hashMoveTable[int(board.Fullmoveno)+halfMove] == move {
		return MAXVALUE
	}
	if dt.IsCapture(move, board) {
		return 10
	}
	return 1
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
func negaMax(board *dt.Board, depth int, alpha, beta int) (int, dt.Move, []dt.Move) {
	nodes++
	alphaOriginal := alpha

	trEntry, err := transpositionTable.get(board)
	if err == nil && trEntry.depth >= depth {
		unApply := board.Apply(trEntry.move)
		switch trEntry.flag {
		case EXACT:
			unApply()
			return trEntry.value, trEntry.move, []dt.Move{}
		case LOWERBOUND:
			alpha = max(alpha, trEntry.value)
		case UPPERBOUND:
			beta = min(beta, trEntry.value)
		}
		unApply()
		if alpha >= beta {
			return trEntry.value, trEntry.move, []dt.Move{}
		}
	}

	moveList := board.GenerateLegalMoves()

	if depth == 0 || len(moveList) == 0 {
		return evalBoard(board), 0, []dt.Move{} // kurwa co
	}

	vMax := MINVALUE
	var bestMove dt.Move
	var tpv []dt.Move

	sortMoves(moveList, board)

	for _, currMove := range moveList {
		unapplyFunc := board.Apply(currMove)

		v, _, ttpv := negaMax(board, depth-1, -beta, -alpha)
		v = -v

		v = int(math.Max(float64(alpha), float64(v)))
		alpha = v

		if v > vMax {
			vMax = v
			bestMove = currMove
			tpv = append(ttpv, currMove)
		}
		unapplyFunc()

		if alpha >= beta {
			break
		}
	}

	trEntry.value = vMax
	trEntry.move = bestMove
	trEntry.depth = depth
	if vMax <= alphaOriginal {
		trEntry.flag = UPPERBOUND
	} else if vMax >= beta {
		trEntry.flag = LOWERBOUND
	} else {
		trEntry.flag = EXACT
	}
	transpositionTable.put(board, trEntry)

	return vMax, bestMove, tpv

}

func evalBoard(board *dt.Board) int {
	if board.OurKingInCheck() && len(board.GenerateLegalMoves()) == 0 {
		return MINVALUE
	}
	v := 0

	v += bits.OnesCount64(board.White.Pawns) * pieceVal[dt.Pawn]
	v += bits.OnesCount64(board.White.Bishops) * pieceVal[dt.Bishop]
	v += bits.OnesCount64(board.White.Knights) * pieceVal[dt.Knight]
	v += bits.OnesCount64(board.White.Rooks) * pieceVal[dt.Rook]
	v += bits.OnesCount64(board.White.Queens) * pieceVal[dt.Queen]

	v -= bits.OnesCount64(board.White.Pawns) * pieceVal[dt.Pawn]
	v -= bits.OnesCount64(board.White.Bishops) * pieceVal[dt.Bishop]
	v -= bits.OnesCount64(board.White.Knights) * pieceVal[dt.Knight]
	v -= bits.OnesCount64(board.White.Rooks) * pieceVal[dt.Rook]
	v -= bits.OnesCount64(board.White.Queens) * pieceVal[dt.Queen]

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

	blackKing := board.White.Kings
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

func getColorMutliplier(color bool) int {
	if color {
		return 1
	}
	return -1
}
