package main

import (
	"time"

	dt "github.com/dylhunn/dragontoothmg"
)

const MAXVALUE int16 = 32767 - 7000
const MINVALUE int16 = -32767 + 7000
const LMR_LIMIT = 5
const DOUBLE_PAWNS_PENALTY = -15
const ISOLATED_PAWNS_PENALTY = -20
const KNIGHT_MOBILITY = 6
const BISHOP_MOBILITY = 4
const ROOK_MOBILITY = 3
const QUEEN_MOBILITY = 3
const SPACE_PER_FRONTSPAN = 1

const KING_SAFETY_SQUARE = 8

var TIMECHECK_FREQ int = 5000
var isEndgame = false
var nodes int
var deepestQuiescence int16
var timeCheckCounter = TIMECHECK_FREQ
var endTime = time.Now().AddDate(1000, 10, 10)
var searching = true
var maxDepth int16
var lastTime int
var transpositionTableSize uint64

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

var pieceVal []int16
var attackSquareVal []int

var isolatedPawnTable = [8]uint64{
	0x303030303030303, 0x707070707070707, 0xe0e0e0e0e0e0e0e, 0x1c1c1c1c1c1c1c1c,
	0x3838383838383838, 0x7070707070707070, 0xe0e0e0e0e0e0e0e0, 0xc0c0c0c0c0c0c0c0,
}

var onlyFile = [8]uint64{
	0x0101010101010101, 0x0202020202020202, 0x0404040404040404, 0x0808080808080808,
	0x1010101010101010, 0x2020202020202020, 0x4040404040404040, 0x8080808080808080}

var onlyRank = [8]uint64{
	0xFF, 0xFF00, 0xFF0000, 0xFF000000,
	0xFF00000000, 0xFF0000000000, 0xFF000000000000, 0xFF00000000000000}

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

	pieceVal = []int16{0, 100, 320, 330, 500, 935, 0}

	attackSquareVal = []int{0, 1, 4, 2, 2, 2, 0}
}
