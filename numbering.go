package xgrid

// StartsAcross reports whether (row, col) is the first square of an across
// entry: it's white, the square to its left is black or off the grid, and
// the square to its right is white and on the grid.
func StartsAcross(g *Grid, row, col int) bool {
	if g.IsBlack(row, col) {
		return false
	}
	if col > 0 && !g.IsBlack(row, col-1) {
		return false
	}
	return col+1 < g.Cols && !g.IsBlack(row, col+1)
}

// StartsDown reports whether (row, col) is the first square of a down
// entry: it's white, the square above it is black or off the grid, and the
// square below it is white and on the grid.
func StartsDown(g *Grid, row, col int) bool {
	if g.IsBlack(row, col) {
		return false
	}
	if row > 0 && !g.IsBlack(row-1, col) {
		return false
	}
	return row+1 < g.Rows && !g.IsBlack(row+1, col)
}

// Number assigns standard crossword clue numbers to a grid. Squares are
// visited in reading order (left to right, top to bottom); a square gets
// the next number if it starts an across entry, a down entry, or both. The
// result is a Rows x Cols matrix where 0 means the square isn't numbered.
func Number(g *Grid) [][]int {
	numbers := make([][]int, g.Rows)
	next := 1
	for r := 0; r < g.Rows; r++ {
		numbers[r] = make([]int, g.Cols)
		for c := 0; c < g.Cols; c++ {
			if StartsAcross(g, r, c) || StartsDown(g, r, c) {
				numbers[r][c] = next
				next++
			}
		}
	}
	return numbers
}
