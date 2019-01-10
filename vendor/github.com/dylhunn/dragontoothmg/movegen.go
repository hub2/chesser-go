// dragontoothmg is a fast chess legal move generator library based on magic bitboards.
package dragontoothmg

// The main Dragontooth move generator file.
// Functions are in this file if (and only if) they are performance-critical
// move generator components, called while actually generating moves in-game.
// (The exception is a few one-line helpers for Move and Board in types.go)

import (
	//"fmt"
	"math/bits"
)

// The main API entrypoint. Generates all legal moves for a given board.
func (b *Board) GenerateLegalMoves() []Move {
	moves := make([]Move, 0, kDefaultMoveListLength)
	// First, see if we are currently in check. If we are, invoke a special check-
	// evasion move generator.
	var kingLocation uint8
	var ourPiecesPtr *Bitboards
	if b.Wtomove { // assumes only one king
		kingLocation = uint8(bits.TrailingZeros64(b.White.Kings))
		ourPiecesPtr = &(b.White)
	} else {
		kingLocation = uint8(bits.TrailingZeros64(b.Black.Kings))
		ourPiecesPtr = &(b.Black)
	}
	kingAttackers, blockerDestinations := b.countAttacks(b.Wtomove, kingLocation, 2)
	if kingAttackers >= 2 { // Under multiple attack, we must move the king.
		b.kingPushes(&moves, ourPiecesPtr)
		return moves
	}

	// Several move types can work in single check, but we must block the check
	if kingAttackers == 1 {
		// calculate pinned pieces
		pinnedPieces := b.generatePinnedMoves(&moves, blockerDestinations)
		nonpinnedPieces := ^pinnedPieces
		// TODO
		b.pawnPushes(&moves, nonpinnedPieces, blockerDestinations)
		b.pawnCaptures(&moves, nonpinnedPieces, blockerDestinations)
		b.knightMoves(&moves, nonpinnedPieces, blockerDestinations)
		b.rookMoves(&moves, nonpinnedPieces, blockerDestinations)
		b.bishopMoves(&moves, nonpinnedPieces, blockerDestinations)
		b.queenMoves(&moves, nonpinnedPieces, blockerDestinations)
		b.kingPushes(&moves, ourPiecesPtr)
		return moves
	}

	// Then, calculate all the absolutely pinned pieces, and compute their moves.
	// If we are in check, we can only move to squares that block the check.
	pinnedPieces := b.generatePinnedMoves(&moves, everything)
	nonpinnedPieces := ^pinnedPieces

	// Finally, compute ordinary moves, ignoring absolutely pinned pieces on the board.
	b.pawnPushes(&moves, nonpinnedPieces, everything)
	b.pawnCaptures(&moves, nonpinnedPieces, everything)
	b.knightMoves(&moves, nonpinnedPieces, everything)
	b.rookMoves(&moves, nonpinnedPieces, everything)
	b.bishopMoves(&moves, nonpinnedPieces, everything)
	b.queenMoves(&moves, nonpinnedPieces, everything)
	b.kingMoves(&moves)
	return moves
}

// Calculate the available moves for absolutely pinned pieces (pinned to the king).
// We are only allowed to move to squares in allowDest, to block checks.
// Return a bitboard of all pieces that are pinned.
func (b *Board) generatePinnedMoves(moveList *[]Move, allowDest uint64) uint64 {
	var ourKingIdx uint8
	var ourPieces, oppPieces *Bitboards
	var allPinnedPieces uint64 = 0
	var pawnPushDirection int
	var doublePushRank, ourPromotionRank uint64
	if b.Wtomove { // Assumes only one king on the board
		ourKingIdx = uint8(bits.TrailingZeros64(b.White.Kings))
		ourPieces = &(b.White)
		oppPieces = &(b.Black)
		pawnPushDirection = 1
		doublePushRank = onlyRank[3]
		ourPromotionRank = onlyRank[7]
	} else {
		ourKingIdx = uint8(bits.TrailingZeros64(b.Black.Kings))
		ourPieces = &(b.Black)
		oppPieces = &(b.White)
		pawnPushDirection = -1
		doublePushRank = onlyRank[4]
		ourPromotionRank = onlyRank[0]
	}
	allPieces := oppPieces.All | ourPieces.All

	// Calculate king moves as if it was a rook.
	// "king targets" includes our own friendly pieces, for the purpose of identifying pins.
	kingOrthoTargets := CalculateRookMoveBitboard(ourKingIdx, allPieces)
	oppRooks := oppPieces.Rooks | oppPieces.Queens
	for oppRooks != 0 { // For each opponent ortho slider
		currRookIdx := uint8(bits.TrailingZeros64(oppRooks))
		oppRooks &= oppRooks - 1
		rookTargets := CalculateRookMoveBitboard(currRookIdx, allPieces) & (^(oppPieces.All))
		// A piece is pinned iff it falls along both attack rays.
		pinnedPiece := rookTargets & kingOrthoTargets & ourPieces.All
		if pinnedPiece == 0 { // there is no pin
			continue
		}
		pinnedPieceIdx := uint8(bits.TrailingZeros64(pinnedPiece))
		sameRank := pinnedPieceIdx/8 == ourKingIdx/8 && pinnedPieceIdx/8 == currRookIdx/8
		sameFile := pinnedPieceIdx%8 == ourKingIdx%8 && pinnedPieceIdx%8 == currRookIdx%8
		if !sameRank && !sameFile {
			continue // it's just an intersection, not a pin
		}
		allPinnedPieces |= pinnedPiece        // store the pinned piece location
		if pinnedPiece&ourPieces.Pawns != 0 { // it's a pawn; we might be able to push it
			if sameFile { // push the pawn
				var pawnTargets uint64 = 0
				pawnTargets |= (1 << uint8(int(pinnedPieceIdx)+8*pawnPushDirection)) & ^allPieces
				if pawnTargets != 0 { // single push worked; try double
					pawnTargets |= (1 << uint8(int(pinnedPieceIdx)+16*pawnPushDirection)) & ^allPieces & doublePushRank
				}
				pawnTargets &= allowDest // TODO this might be a promotion. Is that possible?
				genMovesFromTargets(moveList, Square(pinnedPieceIdx), pawnTargets)
			}
			continue
		}
		// If it's not a rook or queen, it can't move
		if pinnedPiece&ourPieces.Rooks == 0 && pinnedPiece&ourPieces.Queens == 0 {
			continue
		}
		// all ortho moves, as if it was not pinned
		pinnedPieceAllMoves := CalculateRookMoveBitboard(pinnedPieceIdx, allPieces) & (^(ourPieces.All))
		// actually available moves
		pinnedTargets := pinnedPieceAllMoves & (rookTargets | kingOrthoTargets | (uint64(1) << currRookIdx))
		pinnedTargets &= allowDest
		genMovesFromTargets(moveList, Square(pinnedPieceIdx), pinnedTargets)
	}

	// Calculate king moves as if it was a bishop.
	// "king targets" includes our own friendly pieces, for the purpose of identifying pins.
	kingDiagTargets := CalculateBishopMoveBitboard(ourKingIdx, allPieces)
	oppBishops := oppPieces.Bishops | oppPieces.Queens
	for oppBishops != 0 {
		currBishopIdx := uint8(bits.TrailingZeros64(oppBishops))
		oppBishops &= oppBishops - 1
		bishopTargets := CalculateBishopMoveBitboard(currBishopIdx, allPieces) & (^(oppPieces.All))
		pinnedPiece := bishopTargets & kingDiagTargets & ourPieces.All
		if pinnedPiece == 0 { // there is no pin
			continue
		}
		pinnedPieceIdx := uint8(bits.TrailingZeros64(pinnedPiece))
		bishopToPinnedSlope := (float32(pinnedPieceIdx)/8 - float32(currBishopIdx)/8) /
			(float32(pinnedPieceIdx%8) - float32(currBishopIdx%8))
		bishopToKingSlope := (float32(ourKingIdx)/8 - float32(currBishopIdx)/8) /
			(float32(ourKingIdx%8) - float32(currBishopIdx%8))
		if bishopToPinnedSlope != bishopToKingSlope { // just an intersection, not a pin
			continue
		}
		allPinnedPieces |= pinnedPiece // store pinned piece
		// if it's a pawn we might be able to capture with it
		// the capture square must also be in allowdest
		if pinnedPiece&ourPieces.Pawns != 0 {
			if (uint64(1)<<currBishopIdx)&allowDest != 0 {
				if (b.Wtomove && (pinnedPieceIdx/8)+1 == currBishopIdx/8) ||
					(!b.Wtomove && pinnedPieceIdx/8 == (currBishopIdx/8)+1) {
					if ((uint64(1) << currBishopIdx) & ourPromotionRank) != 0 { // We get to promote!
						for i := Piece(Knight); i <= Queen; i++ {
							var move Move
							move.Setfrom(Square(pinnedPieceIdx)).Setto(Square(currBishopIdx)).Setpromote(i)
							*moveList = append(*moveList, move)
						}
					} else { // no promotion
						var move Move
						move.Setfrom(Square(pinnedPieceIdx)).Setto(Square(currBishopIdx))
						*moveList = append(*moveList, move)
					}
				}
			}
			continue
		}
		// If it's not a bishop or queen, it can't move
		if pinnedPiece&ourPieces.Bishops == 0 && pinnedPiece&ourPieces.Queens == 0 {
			continue
		}
		// all diag moves, as if it was not pinned
		pinnedPieceAllMoves := CalculateBishopMoveBitboard(pinnedPieceIdx, allPieces) & (^(ourPieces.All))
		// actually available moves
		pinnedTargets := pinnedPieceAllMoves & (bishopTargets | kingDiagTargets | (uint64(1) << currBishopIdx))
		pinnedTargets &= allowDest
		genMovesFromTargets(moveList, Square(pinnedPieceIdx), pinnedTargets)
	}
	return allPinnedPieces
}

// Generate moves involving advancing pawns.
// Only pieces marked nonpinned can be moved. Only squares in allowDest can be moved to.
func (b *Board) pawnPushes(moveList *[]Move, nonpinned uint64, allowDest uint64) {
	targets, doubleTargets := b.pawnPushBitboards(nonpinned)
	targets, doubleTargets = targets&allowDest, doubleTargets&allowDest
	oneRankBack := 8
	if b.Wtomove {
		oneRankBack = -oneRankBack
	}
	// push all pawns by one square
	for targets != 0 {
		target := bits.TrailingZeros64(targets)
		targets &= targets - 1 // unset the lowest active bit
		var canPromote bool
		if b.Wtomove {
			canPromote = target >= 56
		} else {
			canPromote = target <= 7
		}
		var move Move
		move.Setfrom(Square(target + oneRankBack)).Setto(Square(target))
		if canPromote {
			for i := Piece(Knight); i <= Queen; i++ {
				move.Setpromote(i)
				*moveList = append(*moveList, move)
			}
		} else {
			*moveList = append(*moveList, move)
		}
	}
	// push some pawns by two squares
	for doubleTargets != 0 {
		doubleTarget := bits.TrailingZeros64(doubleTargets)
		doubleTargets &= doubleTargets - 1 // unset the lowest active bit
		var move Move
		move.Setfrom(Square(doubleTarget + 2*oneRankBack)).Setto(Square(doubleTarget))
		*moveList = append(*moveList, move)
	}
}

// A helper function that produces bitboards of valid pawn push locations.
func (b *Board) pawnPushBitboards(nonpinned uint64) (targets uint64, doubleTargets uint64) {
	free := (^b.White.All) & (^b.Black.All)
	if b.Wtomove {
		movableWhitePawns := b.White.Pawns & nonpinned
		targets = movableWhitePawns << 8 & free
		doubleTargets = targets << 8 & onlyRank[3] & free
	} else {
		movableBlackPawns := b.Black.Pawns & nonpinned
		targets = movableBlackPawns >> 8 & free
		doubleTargets = targets >> 8 & onlyRank[4] & free
	}
	return
}

// A function that computes available pawn captures.
// Only pieces marked nonpinned can be moved. Only squares in allowDest can be moved to.
func (b *Board) pawnCaptures(moveList *[]Move, nonpinned uint64, allowDest uint64) {
	east, west := b.pawnCaptureBitboards(nonpinned)
	if b.enpassant > 0 { // always allow us to try en-passant captures
		allowDest = allowDest | 1<<b.enpassant
	}
	east, west = east&allowDest, west&allowDest
	dirbitboards := [2]uint64{east, west}
	if !b.Wtomove {
		dirbitboards[0], dirbitboards[1] = dirbitboards[1], dirbitboards[0]
	}
	for dir, board := range dirbitboards { // for east and west
		for board != 0 {
			target := bits.TrailingZeros64(board)
			board &= board - 1
			var move Move
			move.Setto(Square(target))
			canPromote := false
			if b.Wtomove {
				move.Setfrom(Square(target - (9 - (dir * 2))))
				canPromote = target >= 56
			} else {
				move.Setfrom(Square(target + (9 - (dir * 2))))
				canPromote = target <= 7
			}
			if uint8(target) == b.enpassant && b.enpassant != 0 {
				// Apply, check actual legality, then unapply
				// Warning: not thread safe
				var ourPieces, oppPieces *Bitboards
				var enpassantEnemy uint8
				if b.Wtomove {
					enpassantEnemy = uint8(move.To()) - 8
					ourPieces = &(b.White)
					oppPieces = &(b.Black)
				} else {
					enpassantEnemy = uint8(move.To()) + 8
					ourPieces = &(b.Black)
					oppPieces = &(b.White)
				}
				ourPieces.Pawns &= ^(uint64(1) << move.From())
				ourPieces.All &= ^(uint64(1) << move.From())
				ourPieces.Pawns |= (uint64(1) << move.To())
				ourPieces.All |= (uint64(1) << move.To())
				oppPieces.Pawns &= ^(uint64(1) << enpassantEnemy)
				oppPieces.All &= ^(uint64(1) << enpassantEnemy)
				kingInCheck := b.OurKingInCheck()
				ourPieces.Pawns |= (uint64(1) << move.From())
				ourPieces.All |= (uint64(1) << move.From())
				ourPieces.Pawns &= ^(uint64(1) << move.To())
				ourPieces.All &= ^(uint64(1) << move.To())
				oppPieces.Pawns |= (uint64(1) << enpassantEnemy)
				oppPieces.All |= (uint64(1) << enpassantEnemy)
				if kingInCheck {
					continue
				}
			}
			if canPromote {
				for i := Piece(Knight); i <= Queen; i++ {
					move.Setpromote(i)
					*moveList = append(*moveList, move)
				}
				continue
			}
			*moveList = append(*moveList, move)
		}
	}
}

// A helper than generates bitboards for available pawn captures.
func (b *Board) pawnCaptureBitboards(nonpinned uint64) (east uint64, west uint64) {
	notHFile := uint64(0x7F7F7F7F7F7F7F7F)
	notAFile := uint64(0xFEFEFEFEFEFEFEFE)
	var targets uint64
	// TODO(dylhunn): Always try the en passant capture and verify check status, regardless of
	// valid square requirements
	if b.enpassant > 0 { // an en-passant target is active
		targets = (1 << b.enpassant)
	}
	if b.Wtomove {
		targets |= b.Black.All
		ourpawns := b.White.Pawns & nonpinned
		east = ourpawns << 9 & notAFile & targets
		west = ourpawns << 7 & notHFile & targets
	} else {
		targets |= b.White.All
		ourpawns := b.Black.Pawns & nonpinned
		east = ourpawns >> 7 & notAFile & targets
		west = ourpawns >> 9 & notHFile & targets
	}
	return
}

// Generate all knight moves.
// Only pieces marked nonpinned can be moved. Only squares in allowDest can be moved to.
func (b *Board) knightMoves(moveList *[]Move, nonpinned uint64, allowDest uint64) {
	var ourKnights, noFriendlyPieces uint64
	if b.Wtomove {
		ourKnights = b.White.Knights & nonpinned
		noFriendlyPieces = (^b.White.All)
	} else {
		ourKnights = b.Black.Knights & nonpinned
		noFriendlyPieces = (^b.Black.All)
	}
	for ourKnights != 0 {
		currentKnight := bits.TrailingZeros64(ourKnights)
		ourKnights &= ourKnights - 1
		targets := knightMasks[currentKnight] & noFriendlyPieces & allowDest
		genMovesFromTargets(moveList, Square(currentKnight), targets)
	}
}

// Computes king moves without castling.
func (b *Board) kingPushes(moveList *[]Move, ptrToOurBitboards *Bitboards) {
	ourKingLocation := uint8(bits.TrailingZeros64(ptrToOurBitboards.Kings))
	noFriendlyPieces := ^(ptrToOurBitboards.All)

	// TODO(dylhunn): Modifying the board is NOT thread-safe.
	// We only do this to avoid the king danger problem, aka moving away from a
	// checking slider.
	oldKings := ptrToOurBitboards.Kings
	ptrToOurBitboards.Kings = 0
	ptrToOurBitboards.All &= ^(uint64(1) << ourKingLocation)
	targets := kingMasks[ourKingLocation] & noFriendlyPieces
	for targets != 0 {
		target := bits.TrailingZeros64(targets)
		targets &= targets - 1
		if b.UnderDirectAttack(b.Wtomove, uint8(target)) {
			continue
		}
		var move Move
		move.Setfrom(Square(ourKingLocation)).Setto(Square(target))
		*moveList = append(*moveList, move)
	}

	ptrToOurBitboards.Kings = oldKings
	ptrToOurBitboards.All |= (1 << ourKingLocation)
}

// KingMoves BUT PUBLIC
// Generate all available king moves.
// First, if castling is possible, verifies the checking prohibitions on castling.
// Then, outputs castling moves (if any), and king moves.
// Not thread-safe, since the king is removed from the board to compute
// king-danger squares.
func (b *Board) KingMoves(moveList *[]Move) {
	b.kingMoves(moveList)
}

// Generate all available king moves.
// First, if castling is possible, verifies the checking prohibitions on castling.
// Then, outputs castling moves (if any), and king moves.
// Not thread-safe, since the king is removed from the board to compute
// king-danger squares.
func (b *Board) kingMoves(moveList *[]Move) {
	// castling
	var ourKingLocation uint8
	var canCastleQueenside, canCastleKingside bool
	var ptrToOurBitboards *Bitboards
	allPieces := b.White.All | b.Black.All
	if b.Wtomove {
		ourKingLocation = uint8(bits.TrailingZeros64(b.White.Kings))
		ptrToOurBitboards = &(b.White)
		// To castle, we must have rights and a clear path
		kingsideClear := allPieces&((1<<5)|(1<<6)) == 0
		queensideClear := allPieces&((1<<3)|(1<<2)|(1<<1)) == 0
		// skip the king square, since this won't be called while in check
		canCastleQueenside = b.whiteCanCastleQueenside() &&
			queensideClear && !b.anyUnderDirectAttack(true, 2, 3)
		canCastleKingside = b.whiteCanCastleKingside() &&
			kingsideClear && !b.anyUnderDirectAttack(true, 5, 6)
	} else {
		ourKingLocation = uint8(bits.TrailingZeros64(b.Black.Kings))
		ptrToOurBitboards = &(b.Black)
		kingsideClear := allPieces&((1<<61)|(1<<62)) == 0
		queensideClear := allPieces&((1<<57)|(1<<58)|(1<<59)) == 0
		// skip the king square, since this won't be called while in check
		canCastleQueenside = b.blackCanCastleQueenside() &&
			queensideClear && !b.anyUnderDirectAttack(false, 58, 59)
		canCastleKingside = b.blackCanCastleKingside() &&
			kingsideClear && !b.anyUnderDirectAttack(false, 61, 62)
	}
	if canCastleKingside {
		var move Move
		move.Setfrom(Square(ourKingLocation)).Setto(Square(ourKingLocation + 2))
		*moveList = append(*moveList, move)
	}
	if canCastleQueenside {
		var move Move
		move.Setfrom(Square(ourKingLocation)).Setto(Square(ourKingLocation - 2))
		*moveList = append(*moveList, move)
	}

	// non-castling
	b.kingPushes(moveList, ptrToOurBitboards)
}

// Generate all rook moves using magic bitboards.
// Only pieces marked nonpinned can be moved. Only squares in allowDest can be moved to.
func (b *Board) rookMoves(moveList *[]Move, nonpinned uint64, allowDest uint64) {
	var ourRooks, friendlyPieces uint64
	if b.Wtomove {
		ourRooks = b.White.Rooks & nonpinned
		friendlyPieces = b.White.All
	} else {
		ourRooks = b.Black.Rooks & nonpinned
		friendlyPieces = b.Black.All
	}
	allPieces := b.White.All | b.Black.All
	for ourRooks != 0 {
		currRook := uint8(bits.TrailingZeros64(ourRooks))
		ourRooks &= ourRooks - 1
		targets := CalculateRookMoveBitboard(currRook, allPieces) & (^friendlyPieces) & allowDest
		genMovesFromTargets(moveList, Square(currRook), targets)
	}
}

// Generate all bishop moves using magic bitboards.
// Only pieces marked nonpinned can be moved. Only squares in allowDest can be moved to.
func (b *Board) bishopMoves(moveList *[]Move, nonpinned uint64, allowDest uint64) {
	var ourBishops, friendlyPieces uint64
	if b.Wtomove {
		ourBishops = b.White.Bishops & nonpinned
		friendlyPieces = b.White.All
	} else {
		ourBishops = b.Black.Bishops & nonpinned
		friendlyPieces = b.Black.All
	}
	allPieces := b.White.All | b.Black.All
	for ourBishops != 0 {
		currBishop := uint8(bits.TrailingZeros64(ourBishops))
		ourBishops &= ourBishops - 1
		targets := CalculateBishopMoveBitboard(currBishop, allPieces) & (^friendlyPieces) & allowDest
		genMovesFromTargets(moveList, Square(currBishop), targets)
	}
}

// Generate all queen moves using magic bitboards.
// Only pieces marked nonpinned can be moved. Only squares in allowDest can be moved to.
func (b *Board) queenMoves(moveList *[]Move, nonpinned uint64, allowDest uint64) {
	var ourQueens, friendlyPieces uint64
	if b.Wtomove {
		ourQueens = b.White.Queens & nonpinned
		friendlyPieces = b.White.All
	} else {
		ourQueens = b.Black.Queens & nonpinned
		friendlyPieces = b.Black.All
	}
	allPieces := b.White.All | b.Black.All
	for ourQueens != 0 {
		currQueen := uint8(bits.TrailingZeros64(ourQueens))
		ourQueens &= ourQueens - 1
		// bishop motion
		diag_targets := CalculateBishopMoveBitboard(currQueen, allPieces) & (^friendlyPieces) & allowDest
		genMovesFromTargets(moveList, Square(currQueen), diag_targets)
		// rook motion
		ortho_targets := CalculateRookMoveBitboard(currQueen, allPieces) & (^friendlyPieces) & allowDest
		genMovesFromTargets(moveList, Square(currQueen), ortho_targets)
	}
}

// Helper: converts a targets bitboard into moves, and adds them to the moves list.
func genMovesFromTargets(moveList *[]Move, origin Square, targets uint64) {
	for targets != 0 {
		target := bits.TrailingZeros64(targets)
		targets &= targets - 1
		var move Move
		move.Setfrom(origin).Setto(Square(target))
		*moveList = append(*moveList, move)
	}
}

// Variadic function that returns whether any of the specified squares is being attacked
// by the opponent. Potentially expensive.
func (b *Board) anyUnderDirectAttack(byBlack bool, squares ...uint8) bool {
	for _, v := range squares {
		if b.UnderDirectAttack(byBlack, v) {
			return true
		}
	}
	return false
}

func (b *Board) OurKingInCheck() bool {
	byBlack := b.Wtomove
	var origin uint8
	if b.Wtomove {
		origin = uint8(bits.TrailingZeros64(b.White.Kings))
	} else {
		origin = uint8(bits.TrailingZeros64(b.Black.Kings))
	}
	count, _ := b.countAttacks(byBlack, origin, 1)
	return count >= 1
}

// Determine if a square is under attack. Potentially expensive.
func (b *Board) UnderDirectAttack(byBlack bool, origin uint8) bool {
	count, _ := b.countAttacks(byBlack, origin, 1)
	return count >= 1
}

// GetAttackersForSquare
func (b *Board) GetAttackersForSquare(byBlack bool, origin uint8) Bitboards {
	allPieces := b.White.All | b.Black.All
	var opponentPieces *Bitboards
	if byBlack {
		opponentPieces = &(b.Black)
	} else {
		opponentPieces = &(b.White)
	}
	knight_attackers := knightMasks[origin] & opponentPieces.Knights

	diag_candidates := magicBishopBlockerMasks[origin] & allPieces
	diag_dbindex := (diag_candidates * magicNumberBishop[origin]) >> magicBishopShifts[origin]
	origin_diag_rays := magicMovesBishop[origin][diag_dbindex]

	bishop_attackers := origin_diag_rays & opponentPieces.Bishops
	queen_attackers := origin_diag_rays & opponentPieces.Queens

	ortho_candidates := magicRookBlockerMasks[origin] & allPieces
	ortho_dbindex := (ortho_candidates * magicNumberRook[origin]) >> magicRookShifts[origin]
	origin_ortho_rays := magicMovesRook[origin][ortho_dbindex]

	rook_attackers := origin_ortho_rays & opponentPieces.Rooks
	queen_attackers |= (origin_ortho_rays & opponentPieces.Queens)

	king_attackers := kingMasks[origin] & opponentPieces.Kings

	var pawn_attackers_mask uint64 = 0
	if byBlack {
		pawn_attackers_mask = (1 << (origin + 7)) & ^(onlyFile[7])
		pawn_attackers_mask |= (1 << (origin + 9)) & ^(onlyFile[0])
	} else {
		if origin-7 >= 0 {
			pawn_attackers_mask = (1 << (origin - 7)) & ^(onlyFile[0])
		}
		if origin-9 >= 0 {
			pawn_attackers_mask |= (1 << (origin - 9)) & ^(onlyFile[7])
		}
	}
	pawn_attackers_mask &= opponentPieces.Pawns
	return Bitboards{Pawns: pawn_attackers_mask,
		Bishops: bishop_attackers,
		Knights: knight_attackers,
		Rooks:   rook_attackers,
		Queens:  queen_attackers,
		Kings:   king_attackers}
}

// Compute whether an individual square is under direct attack. Potentially expensive.
// Can be asked to abort early, when a certain number of attacks are found.
// The found number might exceed the abortion threshold, since attacks are grouped.
// Also returns the mask of attackers.
func (b *Board) countAttacks(byBlack bool, origin uint8, abortEarly int) (int, uint64) {
	numAttacks := 0
	var blockerDestinations uint64 = 0
	allPieces := b.White.All | b.Black.All
	var opponentPieces *Bitboards
	if byBlack {
		opponentPieces = &(b.Black)
	} else {
		opponentPieces = &(b.White)
	}
	// find attacking knights
	knight_attackers := knightMasks[origin] & opponentPieces.Knights
	numAttacks += bits.OnesCount64(knight_attackers)
	blockerDestinations |= knight_attackers
	if numAttacks >= abortEarly {
		return numAttacks, blockerDestinations
	}
	// find attacking bishops and queens
	diag_candidates := magicBishopBlockerMasks[origin] & allPieces
	diag_dbindex := (diag_candidates * magicNumberBishop[origin]) >> magicBishopShifts[origin]
	origin_diag_rays := magicMovesBishop[origin][diag_dbindex]
	diag_attackers := origin_diag_rays & (opponentPieces.Bishops | opponentPieces.Queens)
	numAttacks += bits.OnesCount64(diag_attackers)
	blockerDestinations |= diag_attackers
	if numAttacks >= abortEarly {
		return numAttacks, blockerDestinations
	}
	// If we found diagonal attackers, add interposed squares to the blocker mask.
	for diag_attackers != 0 {
		curr_attacker := uint8(bits.TrailingZeros64(diag_attackers))
		diag_attackers &= diag_attackers - 1
		diag_attacks := CalculateBishopMoveBitboard(curr_attacker, allPieces)
		attackRay := diag_attacks & origin_diag_rays
		blockerDestinations |= attackRay
	}

	// find attacking rooks and queens
	ortho_candidates := magicRookBlockerMasks[origin] & allPieces
	ortho_dbindex := (ortho_candidates * magicNumberRook[origin]) >> magicRookShifts[origin]
	origin_ortho_rays := magicMovesRook[origin][ortho_dbindex]
	ortho_attackers := origin_ortho_rays & (opponentPieces.Rooks | opponentPieces.Queens)
	numAttacks += bits.OnesCount64(ortho_attackers)
	blockerDestinations |= ortho_attackers
	if numAttacks >= abortEarly {
		return numAttacks, blockerDestinations
	}
	// If we found orthogonal attackers, add interposed squares to the blocker mask.
	for ortho_attackers != 0 {
		curr_attacker := uint8(bits.TrailingZeros64(ortho_attackers))
		ortho_attackers &= ortho_attackers - 1
		ortho_attacks := CalculateRookMoveBitboard(curr_attacker, allPieces)
		attackRay := ortho_attacks & origin_ortho_rays
		blockerDestinations |= attackRay
	}
	// find attacking kings
	// TODO(dylhunn): What if the opponent king can't actually move to the origin square?
	king_attackers := kingMasks[origin] & opponentPieces.Kings
	numAttacks += bits.OnesCount64(king_attackers)
	blockerDestinations |= king_attackers
	if numAttacks >= abortEarly {
		return numAttacks, blockerDestinations
	}
	// find attacking pawns
	var pawn_attackers_mask uint64 = 0
	if byBlack {
		pawn_attackers_mask = (1 << (origin + 7)) & ^(onlyFile[7])
		pawn_attackers_mask |= (1 << (origin + 9)) & ^(onlyFile[0])
	} else {
		if origin-7 >= 0 {
			pawn_attackers_mask = (1 << (origin - 7)) & ^(onlyFile[0])
		}
		if origin-9 >= 0 {
			pawn_attackers_mask |= (1 << (origin - 9)) & ^(onlyFile[7])
		}
	}
	pawn_attackers_mask &= opponentPieces.Pawns
	numAttacks += bits.OnesCount64(pawn_attackers_mask)
	blockerDestinations |= pawn_attackers_mask
	if numAttacks >= abortEarly {
		return numAttacks, blockerDestinations
	}
	return numAttacks, blockerDestinations
}

// Calculates the attack bitboard for a rook. This might include targeted squares
// that are actually friendly pieces, so the proper usage is:
// rookTargets := CalculateRookMoveBitboard(myRookLoc, allPieces) & (^myPieces)
// Externally useful for evaluation functions.
func CalculateRookMoveBitboard(currRook uint8, allPieces uint64) uint64 {
	blockers := magicRookBlockerMasks[currRook] & allPieces
	dbindex := (blockers * magicNumberRook[currRook]) >> magicRookShifts[currRook]
	targets := magicMovesRook[currRook][dbindex]
	return targets
}

// Calculates the attack bitboard for a bishop. This might include targeted squares
// that are actually friendly pieces, so the proper usage is:
// bishopTargets := CalculateBishopMoveBitboard(myBishopLoc, allPieces) & (^myPieces)
// Externally useful for evaluation functions.
func CalculateBishopMoveBitboard(currBishop uint8, allPieces uint64) uint64 {
	blockers := magicBishopBlockerMasks[currBishop] & allPieces
	dbindex := (blockers * magicNumberBishop[currBishop]) >> magicBishopShifts[currBishop]
	targets := magicMovesBishop[currBishop][dbindex]
	return targets
}
