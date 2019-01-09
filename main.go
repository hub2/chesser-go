package main

import (
	"github.com/dylhunn/dragontoothmg"
)

func main() {
	board := dragontoothmg.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")

	search(&board, 10)

	// moveList := board.GenerateLegalMoves()
	// // For every legal move
	// for _, currMove := range moveList {
	// 	// Apply it to the board
	// 	unapplyFunc := board.Apply(currMove)
	// 	// Print the move, the new position, and the hash of the new position
	// 	fmt.Println("Moved to:", &currMove) // Reference converts Move to string automatically
	// 	fmt.Println("New position is:", board.ToFen())
	// 	fmt.Println("This new position has Zobrist hash:", board.Hash())
	// 	// Unapply the move
	// 	unapplyFunc()
	// }

}
