package main

import (
	"math/bits"

	dt "github.com/dylhunn/dragontoothmg"
)

func evalBoard(board *dt.Board, moveList []dt.Move) int {
	if board.OurKingInCheck() {
		if moveList != nil && len(moveList) == 0 {
			return MINVALUE
		}
	}
	allPieces := board.White.All | board.Black.All
	whitePieces := board.White.All
	blackPieces := board.Black.All
	// Give bonus for having right to move
	// This should help with comparing same positions
	// but with us to move is probably better (if not zugzwang)
	v := 0
	if !isEndgame {
		v += 25
	}

	// Count material
	v += (bits.OnesCount64(board.White.Pawns) - bits.OnesCount64(board.Black.Pawns)) * pieceVal[dt.Pawn]
	v += (bits.OnesCount64(board.White.Knights) - bits.OnesCount64(board.Black.Knights)) * pieceVal[dt.Knight]
	v += (bits.OnesCount64(board.White.Bishops) - bits.OnesCount64(board.Black.Bishops)) * pieceVal[dt.Bishop]
	v += (bits.OnesCount64(board.White.Rooks) - bits.OnesCount64(board.Black.Rooks)) * pieceVal[dt.Rook]
	v += (bits.OnesCount64(board.White.Queens) - bits.OnesCount64(board.Black.Queens)) * pieceVal[dt.Queen]

	mobility := 0
	// Piece square tables
	// TODO: consider lerping between early and end game
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
		bishopTargets := dt.CalculateBishopMoveBitboard(uint8(idx), allPieces) & (^whitePieces)
		mobility += bits.OnesCount64(bishopTargets) * BISHOP_MOBILITY
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.Black.Bishops
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v -= bishopsBlack[idx]
		bishopTargets := dt.CalculateBishopMoveBitboard(uint8(idx), allPieces) & (^blackPieces)
		mobility -= bits.OnesCount64(bishopTargets) * BISHOP_MOBILITY
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.White.Knights
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v += knights[idx]
		knightTargets := dt.KnightMasks[idx] & (^whitePieces)
		mobility += bits.OnesCount64(knightTargets) * KNIGHT_MOBILITY
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.Black.Knights
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v -= knightsBlack[idx]
		knightTargets := dt.KnightMasks[idx] & (^blackPieces)
		mobility -= bits.OnesCount64(knightTargets) * KNIGHT_MOBILITY
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.White.Rooks
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v += rooks[idx]
		rookTargets := dt.CalculateRookMoveBitboard(uint8(idx), allPieces) & (^whitePieces)
		mobility += bits.OnesCount64(rookTargets) * ROOK_MOBILITY
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.Black.Rooks
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v -= rooksBlack[idx]
		rookTargets := dt.CalculateRookMoveBitboard(uint8(idx), allPieces) & (^blackPieces)
		mobility -= bits.OnesCount64(rookTargets) * ROOK_MOBILITY
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.White.Queens
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v += queens[idx]
		diagTargets := dt.CalculateBishopMoveBitboard(uint8(idx), allPieces) & (^whitePieces)
		orthoTargets := dt.CalculateRookMoveBitboard(uint8(idx), allPieces) & (^whitePieces)
		mobility += (bits.OnesCount64(diagTargets) + bits.OnesCount64(orthoTargets)) * QUEEN_MOBILITY
		tmp &= ^(1 << uint(idx))
	}

	tmp = board.Black.Queens
	for tmp != 0 {
		idx := bits.TrailingZeros64(tmp)
		v -= queensBlack[idx]
		diagTargets := dt.CalculateBishopMoveBitboard(uint8(idx), allPieces) & (^blackPieces)
		orthoTargets := dt.CalculateRookMoveBitboard(uint8(idx), allPieces) & (^blackPieces)
		mobility -= (bits.OnesCount64(diagTargets) + bits.OnesCount64(orthoTargets)) * QUEEN_MOBILITY
		tmp &= ^(1 << uint(idx))
	}

	whiteKing := board.White.Kings
	whiteKingIdx := bits.TrailingZeros64(whiteKing)

	blackKing := board.Black.Kings
	blackKingIdx := bits.TrailingZeros64(blackKing)
	isEndgame := (board.Black.Queens | board.White.Queens) == 0
	// TODO: instead of flag implement phases
	if isEndgame {
		v += kingEndgame[whiteKingIdx]
		v -= kingEndgameBlack[blackKingIdx]
	} else {
		v += kingMiddlegame[whiteKingIdx]
		v -= kingMiddlegameBlack[blackKingIdx]
	}
	// King Safety
	numberOfPiecesAroundKingSafetyPoints := (bits.OnesCount64(dt.KingMasks[whiteKingIdx]&(board.White.Pawns)) - bits.OnesCount64(dt.KingMasks[blackKingIdx]&(board.Black.Pawns))) * KING_SAFETY_SQUARE
	v += numberOfPiecesAroundKingSafetyPoints

	// Space
	// Counting rearfill
	v += (bits.OnesCount64(nortFill(board.White.Pawns)) - bits.OnesCount64(soutFill(board.Black.Pawns))) * SPACE_PER_FRONTSPAN

	// stacked pawns & Isolated pawns
	doublePawnsWhite := 0
	doublePawnsBlack := 0

	isolatedPawnsWhite := 0
	isolatedPawnsBlack := 0
	for i := 0; i < 8; i++ {
		// Double pawns
		doublePawnsWhite += bits.OnesCount64(onlyFile[i] & board.White.Pawns)
		doublePawnsBlack += bits.OnesCount64(onlyFile[i] & board.Black.Pawns)
		if doublePawnsWhite > 1 {
			doublePawnsWhite++
		}
		if doublePawnsBlack > 1 {
			doublePawnsBlack++
		}

		// Isolated pawns
		if board.White.Pawns&onlyFile[i] > 0 {
			if board.White.Pawns & ^onlyFile[i] & isolatedPawnTable[i] == 0 {
				isolatedPawnsWhite++
			}
		}
		if board.Black.Pawns&onlyFile[i] > 0 {
			if board.Black.Pawns & ^onlyFile[i] & isolatedPawnTable[i] == 0 {
				isolatedPawnsBlack++
			}
		}

	}
	v += (doublePawnsWhite - doublePawnsBlack) * DOUBLE_PAWNS_PENALTY
	v += (isolatedPawnsWhite - isolatedPawnsBlack) * ISOLATED_PAWNS_PENALTY

	// Mobility

	v += mobility

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
		isTsqLastLine = toBitboard&onlyRank[7] != 0
	} else {
		ourBitboard = &board.Black
		theirBitboard = &board.White
		isTsqLastLine = toBitboard&onlyRank[0] != 0
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
	}
	return swaplist[0]
}
