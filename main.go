package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	dt "github.com/dylhunn/dragontoothmg"
	"github.com/golang-collections/collections/stack"
)

var startingFen = "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"

func main() {
	fmt.Println("Chesser-go")

	commandStack := stack.New()
	input := ""

	board := dt.ParseFen(startingFen)
	reader := bufio.NewReader(os.Stdin)
mainloop:
	for {
		if commandStack.Len() > 0 {
			input = commandStack.Pop().(string)
		} else {
			var err error
			input, err = reader.ReadString('\n')
			if err != nil {
				panic(err)
			}
			input = strings.TrimSpace(input)
		}
		fmt.Fprintf(os.Stderr, "[C] %s\n", input)

		switch input {
		case "quit":
			break mainloop
		case "uci":
			fmt.Println("uciok")
		case "isready":
			fmt.Println("readyok")
		case "ucinewgame":
			commandStack.Push("position fen " + startingFen)
		}
		// position fen <fen>
		// position startpos moves ruch ruch ruch
		if strings.HasPrefix(input, "position") {
			board = dt.ParseFen(startingFen)
			params := strings.SplitN(input, " ", 3)
			if params[1] == "fen" {
				fen := params[2]
				board = dt.ParseFen(fen)
			} else if params[1] == "startpos" {
				if len(params) > 2 {
					k := strings.Split(params[2], " ")
					if k[0] == "moves" {
						for _, mv := range k[1:] {
							parsedMove, err := dt.ParseMove(mv)
							//fmt.Fprintf(os.Stderr, "[C] parsed move: %s\n", mv)

							if err != nil {
								panic(err)
							}
							board.Apply(parsedMove)
						}
					}

				}
			}
		}
		if strings.HasPrefix(input, "go") {
			depth := 11
			//movetime := -1
			params := strings.Split(input, " ")
			if len(params) > 1 {
				i := 0
				for i < len(params) {
					param := params[i]
					if param == "depth" {
						i++
						depth, _ = strconv.Atoi(params[i])
					}
					if param == "movetime" {
						i++
						_, _ = strconv.Atoi(params[i])
					}
					i++
				}
			}
			_, move := search(&board, depth)
			fmt.Printf("bestmove %s\n", move.String())
			fmt.Fprintf(os.Stderr, "bestmove %s\n", move.String())
		}
	}

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
