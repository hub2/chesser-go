package dragontoothmg

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
)

func recomputeBoardHash(b *Board) uint64 {
	var hash uint64 = 0
	if b.Wtomove {
		hash ^= whiteToMoveZobristC
	}
	if b.whiteCanCastleKingside() {
		hash ^= castleRightsZobristC[0]
	}
	if b.whiteCanCastleQueenside() {
		hash ^= castleRightsZobristC[1]
	}
	if b.blackCanCastleKingside() {
		hash ^= castleRightsZobristC[2]
	}
	if b.blackCanCastleQueenside() {
		hash ^= castleRightsZobristC[3]
	}
	hash ^= uint64(b.enpassant)
	for i := uint8(0); i < 64; i++ {
		whitePiece, _ := determinePieceType(&(b.White), uint64(1)<<i)
		blackPiece, _ := determinePieceType(&(b.Black), uint64(1)<<i)
		if whitePiece != Nothing {
			hash ^= pieceSquareZobristC[whitePiece-1][i]
		}
		if blackPiece != Nothing {
			hash ^= pieceSquareZobristC[blackPiece+5][i]
		}
	}
	return hash
}

func IsCapture(m Move, b *Board) bool {
	toBitboard := (uint64(1) << m.To())
	if (toBitboard&b.White.All != 0) || (toBitboard&b.Black.All != 0) {
		return true
	}
	// Is it an en passant capture?
	fromBitboard := (uint64(1) << m.From())
	originIsPawn := fromBitboard&b.White.Pawns != 0 || fromBitboard&b.Black.Pawns != 0
	return originIsPawn && (toBitboard&(uint64(1) << b.enpassant) != 0)
}

// A testing-use function that ignores the error output
func parseMove(movestr string) Move {
	res, _ := ParseMove(movestr)
	return res
}

func (b *Bitboards) sanityCheck() {
	if b.All != b.Pawns|b.Knights|b.Bishops|b.Rooks|b.Kings|b.Queens {
		fmt.Println("Bitboard sanity check problem.")
	}
	if ((((((b.All ^ b.Pawns) ^ b.Knights) ^ b.Bishops) ^ b.Rooks) ^ b.Kings) ^ b.Queens) != 0 {
		fmt.Println("Bitboard sanity check problem.")
	}
}

// Some example valid move strings:
// e1e2 b4d6 e7e8q a2a1n
// TODO(dylhunn): Make the parser more forgiving. Eg: 0-0, O-O-O, a2-a3, D3D4
func ParseMove(movestr string) (Move, error) {
	if movestr == "0000" {
		return 0, nil
	}
	var mv Move
	if len(movestr) < 4 || len(movestr) > 5 {
		return mv, errors.New("Invalid move to parse.")
	}
	from, errf := AlgebraicToIndex(movestr[0:2])
	to, errto := AlgebraicToIndex(movestr[2:4])
	if errf != nil || errto != nil {
		return mv, errors.New("Invalid move to parse.")
	}
	mv.Setto(Square(to)).Setfrom(Square(from))
	if len(movestr) == 5 {
		switch movestr[4] {
		case 'n':
			mv.Setpromote(Knight)
		case 'b':
			mv.Setpromote(Bishop)
		case 'q':
			mv.Setpromote(Queen)
		case 'r':
			mv.Setpromote(Rook)
		default:
			return mv, errors.New("Invalid promotion symbol in move.")
		}
	}
	return mv, nil
}

func printBitboard(bitboard uint64) {
	for i := 63; i >= 0; i-- {
		j := (i/8)*8 + (7 - (i % 8))
		if bitboard&(uint64(1)<<uint8(j)) == 0 {
			fmt.Print("-")
		} else {
			fmt.Print("X")
		}
		if i%8 == 0 {
			fmt.Println()
		}
	}
	fmt.Println()
}

func printMoves(moves []Move) {
	fmt.Println("Moves:")
	for _, v := range moves {
		fmt.Println(&v)
	}
}

// Used for in-place algtoindex parsing where the result is guaranteed to be correct
func algebraicToIndexFatal(alg string) uint8 {
	res, err := AlgebraicToIndex(alg)
	if err != nil {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatal("Could not parse algebraic: ", alg)
	}
	return res
}

// Accepts an algebraic notation chess square, and converts it to a square ID
// as used by Dragontooth (in both the board and move types).
func AlgebraicToIndex(alg string) (uint8, error) {
	firstchar := strings.ToLower(alg)[0]
	if firstchar < 'a' || firstchar > 'h' || alg[1] < '1' || alg[1] > '8' {
		return 64, errors.New("Invalid algebraic " + alg)
	}
	return (firstchar - 'a') + ((alg[1] - '1') * 8), nil
}

// Accepts a Dragontooth Square ID, and converts it to an algebraic square.
func IndexToAlgebraic(id Square) string {
	if id < 0 || id > 63 {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		log.Fatal("Could not parse index: ", id)
	}
	rune := rune((uint8(id) % 8) + 'a')
	return fmt.Sprintf("%c", rune) + strconv.Itoa((int(id)/8)+1)
}

// Serializes a board position to a Fen string.
func (b *Board) ToFen() string {
	b.White.sanityCheck()
	b.Black.sanityCheck()
	var position string
	var empty int // empty slots
	for i := 63; i >= 0; i-- {
		// Loop file A to H, within ranks 8 to 1
		currIdx := (i/8)*8 + (7 - (i % 8))
		var currMask uint64
		currMask = 1 << uint64(currIdx)

		toprint := ""
		if b.White.Pawns&currMask != 0 {
			toprint += "P"
		} else if b.White.Knights&currMask != 0 {
			toprint += "N"
		} else if b.White.Bishops&currMask != 0 {
			toprint += "B"
		} else if b.White.Rooks&currMask != 0 {
			toprint += "R"
		} else if b.White.Queens&currMask != 0 {
			toprint += "Q"
		} else if b.White.Kings&currMask != 0 {
			toprint += "K"
		} else if b.Black.Pawns&currMask != 0 {
			toprint += "p"
		} else if b.Black.Knights&currMask != 0 {
			toprint += "n"
		} else if b.Black.Bishops&currMask != 0 {
			toprint += "b"
		} else if b.Black.Rooks&currMask != 0 {
			toprint += "r"
		} else if b.Black.Queens&currMask != 0 {
			toprint += "q"
		} else if b.Black.Kings&currMask != 0 {
			toprint += "k"
		} else {
			empty++
		}
		if toprint != "" {
			if empty != 0 {
				position += strconv.Itoa(empty)
				empty = 0
			}
			position += toprint
		}

		if i%8 == 0 {
			if empty != 0 {
				position += strconv.Itoa(empty)
				empty = 0
			}
			if i != 0 {
				position += "/"
			}
		}
	}
	if b.Wtomove {
		position += " w"
	} else {
		position += " b"
	}
	position += " "
	castleCount := 0
	if b.whiteCanCastleKingside() {
		position += "K"
		castleCount++
	}
	if b.whiteCanCastleQueenside() {
		position += "Q"
		castleCount++
	}
	if b.blackCanCastleKingside() {
		position += "k"
		castleCount++
	}
	if b.blackCanCastleQueenside() {
		position += "q"
		castleCount++
	}
	if castleCount == 0 {
		position += "-"
	}
	position += " "
	if b.enpassant != 0 {
		position += IndexToAlgebraic(Square(b.enpassant))
	} else {
		position += "-"
	}
	position = position + " " + strconv.Itoa(int(b.Halfmoveclock)) + " " + strconv.Itoa(int(b.Fullmoveno))
	return position
}

// Parse a board from a FEN string.
func ParseFen(fen string) Board {
	// BUG(dylhunn): This FEN parsing implementation doesn't handle malformed inputs.
	tokens := strings.Fields(fen)
	var b Board
	// replace digits with the appropriate number of dashes
	for i := 1; i <= 8; i++ {
		var replacement string
		for j := 0; j < i; j++ {
			replacement += "-"
		}
		tokens[0] = strings.Replace(tokens[0], strconv.Itoa(i), replacement, -1)
	}
	// reverse the order of the ranks, removing slashes
	ranks := strings.Split(tokens[0], "/")
	for i := 0; i < len(ranks)/2; i++ {
		j := len(ranks) - i - 1
		ranks[i], ranks[j] = ranks[j], ranks[i]
	}
	tokens[0] = ranks[0]
	for i := 1; i < len(ranks); i++ {
		tokens[0] += ranks[i]
	}
	// add every piece to the board
	for i := uint8(0); i < 64; i++ {
		switch tokens[0][i] {
		case 'p':
			b.Black.Pawns |= 1 << i
		case 'n':
			b.Black.Knights |= 1 << i
		case 'b':
			b.Black.Bishops |= 1 << i
		case 'r':
			b.Black.Rooks |= 1 << i
		case 'q':
			b.Black.Queens |= 1 << i
		case 'k':
			b.Black.Kings |= 1 << i
		case 'P':
			b.White.Pawns |= 1 << i
		case 'N':
			b.White.Knights |= 1 << i
		case 'B':
			b.White.Bishops |= 1 << i
		case 'R':
			b.White.Rooks |= 1 << i
		case 'Q':
			b.White.Queens |= 1 << i
		case 'K':
			b.White.Kings |= 1 << i
		}
	}
	b.White.All = b.White.Pawns | b.White.Knights | b.White.Bishops | b.White.Rooks | b.White.Queens | b.White.Kings
	b.Black.All = b.Black.Pawns | b.Black.Knights | b.Black.Bishops | b.Black.Rooks | b.Black.Queens | b.Black.Kings

	b.Wtomove = tokens[1] == "w" || tokens[1] == "W"
	if strings.Contains(tokens[2], "K") {
		b.flipWhiteKingsideCastle()
	}
	if strings.Contains(tokens[2], "Q") {
		b.flipWhiteQueensideCastle()
	}
	if strings.Contains(tokens[2], "k") {
		b.flipBlackKingsideCastle()
	}
	if strings.Contains(tokens[2], "q") {
		b.flipBlackQueensideCastle()
	}
	if tokens[3] != "-" {
		res, err := AlgebraicToIndex(tokens[3])
		if err != nil {
			var b2 Board
			return b2 // TODO(dylhunn): return error instead of blank board
		}
		b.enpassant = res
	}

	if len(tokens) > 4 {
		result, _ := strconv.Atoi(tokens[4])
		b.Halfmoveclock = uint8(result)
	}

	if len(tokens) > 5 {
		result, _ := strconv.Atoi(tokens[5])
		b.Fullmoveno = uint16(result)
	}
	b.hash = recomputeBoardHash(&b)
	return b
}
