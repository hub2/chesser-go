package main

import (
	"errors"

	dt "github.com/dylhunn/dragontoothmg"
)

type transpositionFlag int

type transpositionEntry struct {
	value int
	depth int
	move  dt.Move
	flag  transpositionFlag
}

type transpositionMapping map[uint64]transpositionEntry

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
	entry, ok := t[h]

	if !ok || entry.depth <= trEntry.depth {
		t[h] = trEntry
	}
}

func (t transpositionMapping) get(board *dt.Board) (transpositionEntry, error) {
	h := board.Hash()
	entry, ok := t[h]

	if !ok {
		return transpositionEntry{}, errNoTranspositionEntry
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
