package board

// Mirror mirrors b in-place. If columns is true it mirrors by column,
// otherwise it mirrors by row.
func (b *Board) Mirror(columns bool) {
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
}
