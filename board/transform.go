package board

// Mirror mirrors b in-place. If columns is true it mirrors by column,
// otherwise it mirrors by row.
func (b *Board) Mirror(columns bool) {
	b.unsafeMirror(columns)

	b.updateID()
}

func (b *Board) unsafeMirror(columns bool) {
	for i := range b.Markers {
		if columns {
			b.Markers[i].X = b.Size - 1 - b.Markers[i].X
		} else {
			b.Markers[i].Y = b.Size - 1 - b.Markers[i].Y
		}
	}
}

// Rotate rotates b 90 degrees in-place
func (b *Board) Rotate() {
	for i := range b.Markers {
		y := b.Markers[i].Y
		b.Markers[i].Y = b.Markers[i].X
		b.Markers[i].X = b.Size - 1 - y
	}

	b.updateID()
}

// Normalize brings b to a canonical orientation.
//
// If any two boards can be rotated/mirrored to the same form, normalizing
// both will turn them into the same board.
func (b *Board) Normalize() {
	b.updateID()
	minID := b.ID

	for i := 0; i < 8; i++ {
		if i == 4 {
			b.unsafeMirror(false)
		}
		b.Rotate()

		if b.ID < minID {
			minID = b.ID
		}
	}

	b.Markers = markersFromID(b.Size, minID)
	b.updateID()
}

func (b *Board) updateID() {
	b.ID = 0
	b.sortMarkers()
	for i := len(b.Markers) - 1; i >= 0; i-- {
		b.ID *= int64(b.Size)
		b.ID += int64(b.Markers[i].Y)
		b.ID *= int64(b.Size)
		b.ID += int64(b.Markers[i].X)
	}
}
