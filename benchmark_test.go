package main

import (
	"testing"

	"github.com/dylhunn/dragontoothmg"
)

func BenchmarkSearchDepth7(b *testing.B) {
	b.ReportAllocs()
	board := dragontoothmg.ParseFen("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		search(&board, 8)
	}
}
