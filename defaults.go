package blokus

const (
	// Default game board size.
	DefaultBoardSize = 20
)

// Default set of pieces.
func DefaultPieces() []*Piece {
	pieces := make([]*Piece, 0)

	// 1 block pieces.
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
	}))

	// 2 block pieces.
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
	}))

	// 3 block pieces.
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{1, 1},
	}))

	// 4 block pieces.
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{2, 1},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{1, 1},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{1, 1},
		{0, 1},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{1, 1},
		{2, 1},
	}))

	// 5 block pieces.
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
		{4, 0},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
		{3, 1},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{2, 1},
		{3, 1},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{1, 1},
		{2, 1},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{0, 1},
		{2, 1},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{3, 0},
		{2, 1},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{2, 1},
		{2, -1},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{2, 1},
		{2, 2},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{1, 1},
		{2, 1},
		{2, 2},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{1, 1},
		{1, 2},
		{2, 2},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{1, 1},
		{1, 2},
		{2, 1},
	}))
	pieces = append(pieces, NewPieceOrNil([]Coord{
		{0, 0},
		{1, 0},
		{2, 0},
		{1, 1},
		{1, -1},
	}))

	return pieces
}
