package main

import (
	"fmt"
	"math"
	"math/bits"
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

func search(board *dt.Board, depth int) {
	// check if endgame and set appproproeirpoeporiylu
	// isEndgame =...
	nodes = 0
	valf := 0.0
	for i := 1; i < depth; i++ {
		t := time.Now()
		val, _, tpv := negaMax(board, i, math.MinInt32, math.MaxInt32)
		timeElapsed := time.Since(t)

		valf = float64(val) / 100.0
		pv := reverseMove(tpv)
		outMoves := ""
		for _, mv := range pv {
			outMoves += mv.String() + " "
		}
		fmt.Printf("depth %d val %.2f time %v nodes %d\n", i, valf, timeElapsed, nodes)
		fmt.Println(outMoves)
	}
}

func negaMax(board *dt.Board, depth int, alpha, beta int) (int, dt.Move, []dt.Move) {
	nodes++
	moveList := board.GenerateLegalMoves()

	if depth == 0 || len(moveList) == 0 {
		return evalBoard(board), 0, []dt.Move{} // kurwa co
	}

	vMax := MINVALUE
	var bestMove dt.Move
	var tpv []dt.Move

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
