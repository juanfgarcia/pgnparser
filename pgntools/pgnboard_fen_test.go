package pgntools

import (
	"testing"
)

// Test for generating the correct initial position
func TestInitialPosition(t *testing.T) {
	board := InitPgnBoard()
	got := board.GetFen()
	want := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq"
	assert(t, got, want)
}

// A game where both king moves for testing the castling
func TestMoveKing(t *testing.T){
	board := InitPgnBoard()

	var moveTable = []struct {
		move PgnMove
		fen string
	}{
		{ PgnMove{1, 1, "e3",-1,""}, "rnbqkbnr/pppppppp/8/8/8/4P3/PPPP1PPP/RNBQKBNR b KQkq"},
		{ PgnMove{1, -1, "e6",-1,""}, "rnbqkbnr/pppp1ppp/4p3/8/8/4P3/PPPP1PPP/RNBQKBNR w KQkq"},
		{ PgnMove{1, 1, "Ke2",-1,""}, "rnbqkbnr/pppp1ppp/4p3/8/8/4P3/PPPPKPPP/RNBQ1BNR b kq"},
		{ PgnMove{1, -1, "Ke7",-1,""},"rnbq1bnr/ppppkppp/4p3/8/8/4P3/PPPPKPPP/RNBQ1BNR w -"},
	}

	for _, tt := range moveTable {
		t.Run(tt.move.String(), func(t *testing.T){
			board.UpdateBoard(tt.move, false)
			got := board.GetFen()
			want := tt.fen
			assert(t, got, want)
		})
	}

}

func assert(t *testing.T, got, want string){
	t.Helper()

	if got != want {
        t.Errorf("got '%s' want '%s'", got, want)
    }
}

