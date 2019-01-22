package main

import (
	"errors"

	dt "github.com/dylhunn/dragontoothmg"
)

type transpositionFlag uint16

type transpositionEntry struct {
	value int16
	depth int16
	move  dt.Move
	flag  transpositionFlag
	key   uint64
}

type transpositionMapping []transpositionEntry

// transposition table for bookkeeping already evaluated positions
var transpositionTable transpositionMapping

// table for keeping PV
var hashMoveTable []dt.Move

// tables for keeping beta-cutoff moves(but not a capture or promotion, we search them anyway early on)
var killerOneTable []dt.Move
var killerTwoTable []dt.Move

var errNoTranspositionEntry = errors.New("No entry")

const (
	// EXACT value from search
	EXACT transpositionFlag = iota
	// LOWERBOUND alpha from search
	LOWERBOUND
	// UPPERBOUND beta from search
	UPPERBOUND
)

func (t transpositionMapping) put(board *dt.Board, trEntry transpositionEntry) {
	h := board.Hash()
	idx := (h & transpositionTableSize)
	entry := t[idx]

	if entry.depth <= trEntry.depth {
		trEntry.key = h ^ (uint64(trEntry.value)<<48 | uint64(trEntry.depth)<<32 | uint64(trEntry.move)<<16 | uint64(trEntry.flag))
		t[idx] = trEntry
	}
}

func (t transpositionMapping) get(board *dt.Board) (transpositionEntry, error) {
	h := board.Hash()
	idx := (h & transpositionTableSize)
	entry := t[idx]

	if entry.key != (h ^ (uint64(entry.value)<<48 | uint64(entry.depth)<<32 | uint64(entry.move)<<16 | uint64(entry.flag))) {
		return transpositionEntry{key: h}, errNoTranspositionEntry
	}
	return entry, nil
}

func addKiller(move dt.Move, depth int) {
	if killerOneTable[depth] == 0 {
		killerOneTable[depth] = move
	} else if move != killerOneTable[depth] {
		killerTwoTable[depth] = move
	}
}
