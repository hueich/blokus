package blokus

type Rotation uint8

const (
	Rot0 Rotation = iota
	Rot90
	Rot180
	Rot270

	rotEnd
)

func (r Rotation) Normalize() Rotation {
	return r % rotEnd
}

type Orientation struct {
	Rot  Rotation
	Flip bool
}

func (o Orientation) Transform(c Coord) Coord {
	// TODO: Use rotation matrix instead of looping
	for i := 0; i < int(o.Rot); i++ {
		c = rotateCoord(c)
	}
	if o.Flip {
		c = flipCoord(c)
	}
	return c
}

func rotateCoord(c Coord) Coord {
	return Coord{c.Y, -c.X}
}

func flipCoord(c Coord) Coord {
	return Coord{c.X, -c.Y}
}
