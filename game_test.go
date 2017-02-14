package blokus

import (
	"fmt"
	"strings"
	"testing"
)

func newGameOrDie(t *testing.T) *Game {
	g, err := NewGame(10, []*Piece{newPieceOrDie(t, []Coord{{0, 0}})})
	if err != nil {
		t.Fatalf("NewGame(): got %v, want no error", err)
	}
	return g
}

func newPieceOrDie(t *testing.T, blocks []Coord) *Piece {
	p, err := NewPiece(blocks)
	if err != nil {
		t.Fatalf("NewPiece(): got %v, want no error", err)
	}
	return p
}

func TestNewGame(t *testing.T) {
	p := newPieceOrDie(t, []Coord{Coord{3, 4}})
	g, err := NewGame(22, []*Piece{p})
	if err != nil {
		t.Fatalf("NewGame(): got %v, want no error", err)
	}

	if g.Board == nil {
		t.Fatal("Game board: got nil, want not nil")
	}
	if got, want := len(g.Board.Grid), 22; got != want {
		t.Errorf("Game board size: got %v, want %v", got, want)
	}
	if got, want := len(g.Pieces), 1; got != want {
		t.Fatalf("Game pieces len: got %v, want %v", got, want)
	}
	if got, want := len(g.Pieces[0].Blocks), 1; got != want {
		t.Fatalf("Game pieces[0] blocks len: got %v, want %v", got, want)
	}
	if got, want := g.Pieces[0].Blocks[0], (Coord{3, 4}); got != want {
		t.Errorf("Game pieces[0] blocks[0]: got %v, want %v", got, want)
	}
}

func TestAddPlayer(t *testing.T) {
	ps := []*Piece{
		newPieceOrDie(t, []Coord{{3, 4}}),
		newPieceOrDie(t, []Coord{{5, 6}}),
	}
	g, err := NewGame(22, ps)
	if err != nil {
		t.Fatalf("NewGame(): got %v, want no error", err)
	}

	if err := g.AddPlayer("foo", Red, Coord{0, 0}); err != nil {
		t.Fatalf("AddPlayer(): got %v, want no error", err)
	}

	if got, want := len(g.Players), 1; got != want {
		t.Fatalf("Game players len: got %v, want %v", got, want)
	}
	p := g.Players[0]
	if got, want := p.Name, "foo"; got != want {
		t.Errorf("Player name: got %v, want %v", got, want)
	}
	if got, want := p.Color, Red; got != want {
		t.Errorf("Player color: got %v, want %v", got, want)
	}
	if got, want := p.StartPos, (Coord{0, 0}); got != want {
		t.Errorf("Player start position: got %v, want %v", got, want)
	}
	if got, want := len(p.PlacedPieces), len(ps); got != want {
		t.Errorf("Player placed pieces length: got %v, want %v", got, want)
	}
	for i, placed := range p.PlacedPieces {
		if placed {
			t.Errorf("Player placed pieces slice at index %v: got true, want false", i)
		}
	}
}

func TestAddPlayerInvalidColor(t *testing.T) {
	g := newGameOrDie(t)
	if err, msg := g.AddPlayer("foo", 99, Coord{9, 9}), "color"; err == nil || !strings.Contains(err.Error(), msg) {
		t.Errorf("AddPlayer() with invalid color: got %v, want error to contain %v", err, msg)
	}
}

func TestAddPlayerInvalidStartPosition(t *testing.T) {
	g := newGameOrDie(t)
	if err := g.AddPlayer("foo", Blue, Coord{10, 10}); err == nil || !strings.Contains(err.Error(), "position is out of bounds") {
		t.Errorf("AddPlayer() with invalid start position: got %v, want starting position out of bounds error", err)
	}
}

func TestAddPlayerDupeName(t *testing.T) {
	g := newGameOrDie(t)
	name := "foo"
	if err := g.AddPlayer(name, Blue, Coord{9, 9}); err != nil {
		t.Fatalf("Add first player: got error %v, want no error", err)
	}

	if err := g.AddPlayer(name, Yellow, Coord{0, 0}); err == nil || !strings.Contains(err.Error(), name) {
		t.Errorf("AddPlayer() with duplicate name: got %v, want duplicate name error containing %v", err, name)
	}
}

func TestAddPlayerDupeColor(t *testing.T) {
	g := newGameOrDie(t)
	color := Blue
	if err := g.AddPlayer("foo", color, Coord{9, 9}); err != nil {
		t.Fatalf("Add first player: got error %v, want no error", err)
	}

	if err := g.AddPlayer("bar", color, Coord{0, 0}); err == nil || !strings.Contains(err.Error(), color.String()) {
		t.Errorf("AddPlayer() with duplicate color: got %v, want duplicate color error containing %v", err, color)
	}
}

func TestAddPlayerDupeStartPosition(t *testing.T) {
	g := newGameOrDie(t)
	if err := g.AddPlayer("foo", Blue, Coord{9, 9}); err != nil {
		t.Fatalf("Add first player: got error %v, want no error", err)
	}

	if err := g.AddPlayer("bar", Yellow, Coord{9, 9}); err == nil || !strings.Contains(err.Error(), "position already occupied") {
		t.Errorf("AddPlayer() with duplicate start position: got %v, want position already occupied error", err)
	}
}

func TestAddPlayerAutoAssignColor(t *testing.T) {
	g := newGameOrDie(t)
	if err := g.AddPlayer("foo", colorEmpty, Coord{0, 0}); err != nil {
		t.Errorf("Add player with no color: got error %v, want no error", err)
	}
	if got, want := len(g.Players), 1; got != want {
		t.Errorf("Num players after AddPlayer(): got %v, want %v", got, want)
	}
	if got, want := g.Players[0].Color, Color(1); got != want {
		t.Errorf("Player color after AddPlayer(): got %v, want %v", got, want)
	}
}

func TestAddPlayerAutoAssignColorNoMoreColors(t *testing.T) {
	g := newGameOrDie(t)
	for i := 1; i < int(colorEnd); i++ {
		if err := g.AddPlayer(fmt.Sprintf("foo_%d", i), colorEmpty, Coord{0, i}); err != nil {
			t.Fatalf("Add player %d with no color: got error %v, want no error", i, err)
		}
	}
	colors := make(map[Color]bool)
	for _, p := range g.Players {
		colors[p.Color] = true
	}
	if got, want := len(colors), int(colorEnd)-1; got != want {
		t.Errorf("Num colors after adding max players: got %v, want %v", got, want)
	}
	if err := g.AddPlayer("bar", colorEmpty, Coord{1, 0}); err == nil || !strings.Contains(err.Error(), "color") {
		t.Errorf("Add extra player with no color: got %v, want error about color", err)
	}
}

func newGameWithTwoPlayersAndTwoPieces(t *testing.T, size int) *Game {
	tps := []*Piece{
		newPieceOrDie(t, []Coord{{0, 0}, {1, 0}, {2, 0}}),
		newPieceOrDie(t, []Coord{{0, 0}, {1, 0}, {1, 1}, {2, 1}}),
	}
	g, err := NewGame(size, tps)
	if err != nil {
		t.Fatalf("NewGame(ot %v, want no error", err)
	}
	if err := g.AddPlayer("foo", Blue, Coord{0, 0}); err != nil {
		t.Fatalf("AddPlayer(foo): got %v, want no error", err)
	}
	if err := g.AddPlayer("bar", Yellow, Coord{size - 1, size - 1}); err != nil {
		t.Fatalf("AddPlayer(bar): got %v, want no error", err)
	}
	return g
}

func TestPassTurn(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.PassTurn(g.Players[0]); err != nil {
		t.Errorf("PassTurn(player0): got %v, want no error", err)
	}
	if got, want := len(g.Moves), 1; got != want {
		t.Fatalf("Num moves after PassTurn(): got %v, want %v", got, want)
	}
	if got := g.Moves[0]; !got.IsPass() {
		t.Errorf("Moves[0]: got %v, want IsPass()=true", got)
	}
}

func TestPassTurnNilPlayer(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.PassTurn(nil); err == nil || !strings.Contains(err.Error(), "Invalid player") {
		t.Errorf("PassTurn(nil): got %v, want error about invalid player", err)
	}
}

func TestPassTurnWrongPlayerTurn(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.PassTurn(g.Players[1]); err == nil || !strings.Contains(err.Error(), "Turn") {
		t.Errorf("PassTurn(player1): got %v, want error about wrong turn", err)
	}
}

func TestPlacePieceAlreadyPlaced(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)

	player := g.Players[0]
	if err := g.PlacePiece(player, 0, Orientation{Rot0, false}, Coord{0, 0}); err != nil {
		t.Fatalf("PlacePiece 1st time: got %v, want no error", err)
	}

	if err := g.PlacePiece(player, 0, Orientation{Rot0, false}, Coord{2, 2}); err == nil || !strings.Contains(err.Error(), "already placed") {
		t.Errorf("PlacePiece() 2nd time: got %v, want piece already placed error", err)
	}
}

func TestPlacePieceOutOfTurn(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)

	player := g.Players[1]
	if err := g.PlacePiece(player, 0, Orientation{Rot0, false}, Coord{0, 0}); err == nil || !strings.Contains(strings.ToLower(err.Error()), "turn") {
		t.Errorf("PlacePiece(): got %v, want out of turn error", err)
	}
}

func TestPlacePieceValid(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	numMovesBefore := len(g.Moves)

	player := g.Players[0]
	if err := g.PlacePiece(player, 0, Orientation{Rot270, false}, Coord{0, 0}); err != nil {
		t.Fatalf("PlacePiece(): got %v, want no error", err)
	}

	spaces := []Coord{
		{0, 0},
		{0, 1},
		{0, 2},
	}
	for _, s := range spaces {
		if got := g.Board.Grid[s.X][s.Y]; got != player.Color {
			t.Errorf("Board space [%v]: got %v, want occupied by piece colored %v", s, got, player.Color)
		}
	}
	if !player.PlacedPieces[0] {
		t.Errorf("Player placed piece: got false, want true")
	}

	// Move was recorded in the game.
	if got, want := len(g.Moves), numMovesBefore+1; got != want {
		t.Errorf("Number of moves: got %v, want %v", got, want)
	}
	gotMove := g.Moves[len(g.Moves)-1]
	wantMove := &Move{
		Player:     player,
		PieceIndex: 0,
		Orient:     Orientation{Rot270, false},
		Loc:        Coord{0, 0},
	}
	if *gotMove != *wantMove {
		t.Errorf("Move info: got %v, want %v", *gotMove, *wantMove)
	}
	if got := gotMove.IsPass(); got {
		t.Errorf("Move.IsPass(): got %v, want false", got)
	}

	// Should not advance turn
	if got, want := g.CurPlayerIndex, 0; got != want {
		t.Errorf("Player turn after placing piece: got %v, want %v", got, want)
	}
}

func TestCheckPiecePlacementAtStartingPosition(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.checkPiecePlacement(g.Players[0], g.Pieces[1], Coord{0, 0}); err != nil {
		t.Fatalf("checkPiecePlacement(): got %v, want no error", err)
	}
}

func TestCheckPiecePlacementWithValidTouchingCorner(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.PlacePiece(g.Players[0], 0, Orientation{Rot0, false}, Coord{0, 0}); err != nil {
		t.Fatalf("PlacePiece(blue 1st piece): got %v, want no error", err)
	}

	if err := g.checkPiecePlacement(g.Players[0], g.Pieces[1], Coord{3, 1}); err != nil {
		t.Fatalf("checkPiecePlacement(blue 2nd piece): got %v, want no error", err)
	}
}

func TestCheckPiecePlacementNoValidCorner(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.checkPiecePlacement(g.Players[0], g.Pieces[0], Coord{2, 2}); err == nil || !strings.Contains(err.Error(), "no corner touching") {
		t.Errorf("checkPiecePlacement(): got %v, want no corner touching error", err)
	}
}

func TestCheckPiecePlacementOutOfBounds(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.checkPiecePlacement(g.Players[0], g.Pieces[0], Coord{19, 19}); err == nil || !strings.Contains(err.Error(), "out of bounds") {
		t.Errorf("checkPiecePlacement(): got %v, want out of bounds error", err)
	}
}

func TestCheckPiecePlacementOverlapExistingPiece(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 3)
	// Place an existing piece
	if err := g.PlacePiece(g.Players[0], 1, Orientation{Rot0, false}, Coord{0, 0}); err != nil {
		t.Fatalf("Place initial piece: got %v, want no error", err)
	}

	if err := g.checkPiecePlacement(g.Players[1], g.Pieces[1], Coord{0, 1}); err == nil || !strings.Contains(err.Error(), "occupied") {
		t.Errorf("checkPiecePlacement(): got %v, want occupied error", err)
	}
}

func TestCheckPiecePlacementTouchingSameColor(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 10)
	if err := g.PlacePiece(g.Players[0], 0, Orientation{Rot0, false}, Coord{0, 0}); err != nil {
		t.Fatalf("PlacePiece(blue 1st piece): got %v, want no error", err)
	}

	if err := g.checkPiecePlacement(g.Players[0], g.Pieces[1], Coord{2, 1}); err == nil || !strings.Contains(err.Error(), "another blue piece") {
		t.Errorf("checkPiecePlacement(blue 2nd piece): got error %v, want next to another blue piece error", err)
	}
}

func TestCheckPiecePlacementTouchingAnotherColor(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 3)
	if err := g.PlacePiece(g.Players[0], 0, Orientation{Rot0, false}, Coord{0, 0}); err != nil {
		t.Fatalf("PlacePiece(blue 1st piece): got %v, want no error", err)
	}
	if err := g.checkPiecePlacement(g.Players[1], g.Pieces[0], Coord{0, 2}); err != nil {
		t.Fatalf("PlacePiece(yellow 1st piece): got %v, want no error", err)
	}
}

func TestAdvanceTurnNoPlayers(t *testing.T) {
	g := newGameOrDie(t)
	if err := g.AdvanceTurn(); err == nil || !strings.Contains(err.Error(), "no players") {
		t.Errorf("AdvanceTurn() with no players: got %v, want no players error", err)
	}
	if got := g.CurPlayerIndex; got != 0 {
		t.Fatalf("After AdvanceTurn() returned error, curPlayerIndex: got %v, want 0", got)
	}
}

func TestAdvanceTurn(t *testing.T) {
	g := newGameOrDie(t)
	if err := g.AddPlayer("foo", Red, Coord{0, 0}); err != nil {
		t.Fatalf("Error adding player 'foo': %v", err)
	}
	if err := g.AddPlayer("bar", Blue, Coord{0, 9}); err != nil {
		t.Fatalf("Error adding player 'bar': %v", err)
	}

	if got, want := g.Players[g.CurPlayerIndex].Name, "foo"; got != want {
		t.Fatalf("Before AdvanceTurn(): got player %v, want player %v", got, want)
	}
	if err := g.AdvanceTurn(); err != nil {
		t.Fatalf("AdvanceTurn() with players: got %v, want no error", err)
	}
	if got, want := g.Players[g.CurPlayerIndex].Name, "bar"; got != want {
		t.Errorf("After AdvanceTurn(): got player %v, want player %v", got, want)
	}
}

func TestIsGameEndNoPlayers(t *testing.T) {
	g := newGameOrDie(t)
	if got := g.IsGameEnd(); got {
		t.Errorf("IsGameEnd() with no players: got %v, want false", got)
	}
}

func TestIsGameEnd(t *testing.T) {
	g := newGameWithTwoPlayersAndTwoPieces(t, 5)
	if got := g.IsGameEnd(); got {
		t.Errorf("IsGameEnd() with no moves: got %v, want false", got)
	}
	if err := g.PassTurn(g.Players[0]); err != nil {
		t.Fatalf("PassTurn(player0): got %v, want no error", err)
	}
	if err := g.AdvanceTurn(); err != nil {
		t.Fatalf("AdvanceTurn(player0): got %v, want no error", err)
	}
	if got := g.IsGameEnd(); got {
		t.Errorf("IsGameEnd() with 1 pass: got %v, want false", got)
	}
	if err := g.PassTurn(g.Players[1]); err != nil {
		t.Fatalf("PassTurn(player1): got %v, want no error", err)
	}
	if err := g.AdvanceTurn(); err != nil {
		t.Fatalf("AdvanceTurn(player1): got %v, want no error", err)
	}
	if got := g.IsGameEnd(); !got {
		t.Errorf("IsGameEnd() with 2 passes: got %v, want true", got)
	}
}
