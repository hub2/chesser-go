package main

import (
	"testing"

	dt "github.com/dylhunn/dragontoothmg"
)

func TestEvalFunc(t *testing.T) {
	board := dt.ParseFen(startingFen)
	val := evalBoard(&board, nil)
	t.Log(val)
}
