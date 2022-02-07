[![License: GPL v3](https://img.shields.io/badge/License-GPL%20v3-blue.svg)](http://www.gnu.org/licenses/gpl-3.0)
[![Build Status](https://travis-ci.org/dylhunn/dragontoothmg.svg?branch=master)](https://travis-ci.org/dylhunn/dragontoothmg)
[![Documentation](https://img.shields.io/badge/Documentation-GoDoc-green.svg)](https://godoc.org/github.com/dylhunn/dragontoothmg)


Dragontooth Movegen | Dylan D. Hunn
==================================

Dragontooth Movegen is a fast, no-compromises chess move generator written entirely in Go. It provides a simple API for `GenerateLegalMoves()`. It also provides `Board` and `Move` types, `Apply()` and `Unapply()` functionality, and easy-to-use Zobrist-backed `hash`ing of board positions. FEN parsing/serializing and Move parsing/serializing are supported out of the box.

`Dragontoothmg` is based on *magic bitboards* for maximum performance, and generates legal moves only using *pinned piece tables*.

**This project is currently stable and fully functional.** Optimizations are underway, to improve on the benchmarks listed below.

Repo summary
============

Here is a summary of the important files in the repo:

| **File**         | **Description**                                                                                                                                         |
|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| movegen.go   | This is the "primary" source file. Functions are located here if, and only if, they are performance critical and executed to generate moves in-game. |
| types.go     | This file contains the Board and Moves types, along with some supporting helper functions and types.                                                 |
| constants.go | All constants for move generation are hard-coded here, along with functions to compute the magic bitboard lookup tables when the file loads.         |
| util.go      | This file contains supporting library functions, for FEN reading and conversions.                                                                    |
| apply.go     | This provides functions to apply and unapply moves to the board. (Useful for Perft as well.)                                                         |
| perft.go     | The actual Perft implementation is contained in this file.                                                                                           |

API
===

Here are significant API calls that this library provides. For invocation details, see the [docs](https://godoc.org/github.com/dylhunn/dragontoothmg).

| **Function**         | **Description**                                                                                                                                         |
|--------------|------------------------------------------------------------------------------------------------------------------------------------------------------|
| GenerateLegalMoves   | A fast way to generate all moves in the current position. |
| Board.Apply     | Apply a move to the board. Returns a function that allows it to be unapplied.                                                         |                                                      |
| Perft     | Standard "performance test," which recursively counts all of the moves from a position to a given depth.                                                         |
| ParseFen     | Construct a Board from a standard chess FEN string.                                               |
| Board.ToFen | Convert a Board to a standard FEN string.         |
| Board.Hash     | Generate a hash value for a Board, using the Zobrist method.                                                                                           |
| ParseMove     | Parse a long-algbraic notation move from a string.                                                                                           |
| Move.String     | Convert a Move to a string, in normal long-algebraic notation.                                                                                           |

Installing and building the library
===================================

This project requires Go 1.9. As of the time of writing, 1.9 is still a pre-release version. You can get it by cloning the official [Go Repo](https://github.com/golang/go), and building it yourself. There are [instructions for this](https://golang.org/doc/install/source). (Once Go 1.9 comes out in a few weeks, this will become unnecessary.)

To use this package in your own code, make sure your `GO_PATH` environment variable is correctly set, and install it using `go get`:

    go get github.com/dylhunn/dragontoothmg

Then, you can include it in your project:

	import "github.com/dylhunn/dragontoothmg"

Alternatively, you can clone it yourself:

    git clone https://github.com/dylhunn/dragontoothmg.git

Testing and benchmarking
========================

To run all tests, `cd` into the directory and use:

	go test -v

The `-v` shows verbose progress output, since some of the Perft tests can take some time.

To run benchmarks:

	go run bench/runbench.go

Current benchmark results are around 60 million NPS (nodes per second) on a modern Intel i5. This [significantly outperforms](http://i68.tinypic.com/r8rwow.png) the best current Go chess engines, and is about 40% of the performance of the Stockfish move generator. (Not bad for a garbage-collected language!) Improvements are continually underway, and results will vary on your machine.

![Sample Benchmark Results](/benchmarks.png?raw=true "Sample Benchmark Results")

Documentation and examples
==========================

You can find the documentation [here](https://godoc.org/github.com/dylhunn/dragontoothmg).

Here is a simple example invocation:

    // Read a position from a FEN string
    board := dragontoothmg.ParseFen("1Q2rk2/2p2p2/1n4b1/N7/2B1Pp1q/2B4P/1QPP4/4K2R b K e3 4 30")
    // Generate all legal moves
    moveList := board.GenerateLegalMoves()
    // For every legal move
    for _, currMove := range moveList {
        // Apply it to the board
        unapplyFunc := board.Apply(currMove)
        // Print the move, the new position, and the hash of the new position
        fmt.Println("Moved to:", &currMove) // Reference converts Move to string automatically
        fmt.Println("New position is:", b.ToFen())
        fmt.Println("This new position has Zobrist hash:", board.Hash())
        // Unapply the move
        unapplyFunc()
    }