package xgrid

import (
	"fmt"
	"strings"
)

// minEntryLen is the shortest an across or down entry is allowed to be in
// strict mode. American-style crosswords disallow two-letter answers.
const minEntryLen = 3

// ValidationError collects every rule a grid broke, so a caller can report
// them all at once instead of fixing and re-running one at a time.
type ValidationError struct {
	Issues []string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("xgrid: grid failed validation (%d issue(s)): %s", len(e.Issues), strings.Join(e.Issues, "; "))
}

// Validate checks a grid against standard crossword conventions.
//
// In strict mode it requires: at least one white square, 180-degree
// rotational symmetry (the pattern of black squares looks the same
// upside down), no isolated white squares, and no across or down entry
// shorter than three letters.
//
// In lenient mode only the check that would make the grid meaningless
// (having no white squares at all) is enforced; the rest are skipped so
// that cryptics, asymmetric grids, and other non-standard styles pass.
func Validate(g *Grid, opts Options) error {
	var issues []string

	if !hasAnyWhiteSquare(g) {
		issues = append(issues, "grid has no white squares")
	}

	if !opts.Lenient {
		if !isSymmetric(g) {
			issues = append(issues, "grid is not 180-degree rotationally symmetric")
		}
		issues = append(issues, checkEntryLengths(g)...)
		issues = append(issues, checkNoIsolatedSquares(g)...)
	}

	if len(issues) == 0 {
		return nil
	}
	return &ValidationError{Issues: issues}
}

func hasAnyWhiteSquare(g *Grid) bool {
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			if !g.IsBlack(r, c) {
				return true
			}
		}
	}
	return false
}

func isSymmetric(g *Grid) bool {
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			mr, mc := g.Rows-1-r, g.Cols-1-c
			if g.IsBlack(r, c) != g.IsBlack(mr, mc) {
				return false
			}
		}
	}
	return true
}

// checkEntryLengths walks every maximal run of white squares in each row
// and each column and flags any run shorter than minEntryLen.
func checkEntryLengths(g *Grid) []string {
	var issues []string

	for r := 0; r < g.Rows; r++ {
		start := -1
		for c := 0; c <= g.Cols; c++ {
			white := c < g.Cols && !g.IsBlack(r, c)
			if white && start == -1 {
				start = c
			} else if !white && start != -1 {
				if length := c - start; length == 1 || (length > 0 && length < minEntryLen) {
					issues = append(issues, fmt.Sprintf("row %d has a %d-letter run starting at column %d (minimum is %d)", r, length, start, minEntryLen))
				}
				start = -1
			}
		}
	}

	for c := 0; c < g.Cols; c++ {
		start := -1
		for r := 0; r <= g.Rows; r++ {
			white := r < g.Rows && !g.IsBlack(r, c)
			if white && start == -1 {
				start = r
			} else if !white && start != -1 {
				if length := r - start; length == 1 || (length > 0 && length < minEntryLen) {
					issues = append(issues, fmt.Sprintf("column %d has a %d-letter run starting at row %d (minimum is %d)", c, length, start, minEntryLen))
				}
				start = -1
			}
		}
	}

	return issues
}

// checkNoIsolatedSquares flags any white square that is not part of a
// multi-letter entry in either direction. A single white square boxed in
// by black squares on all four sides can't be answered by any word.
func checkNoIsolatedSquares(g *Grid) []string {
	var issues []string
	for r := 0; r < g.Rows; r++ {
		for c := 0; c < g.Cols; c++ {
			if g.IsBlack(r, c) {
				continue
			}
			inAcross := (c > 0 && !g.IsBlack(r, c-1)) || (c < g.Cols-1 && !g.IsBlack(r, c+1))
			inDown := (r > 0 && !g.IsBlack(r-1, c)) || (r < g.Rows-1 && !g.IsBlack(r+1, c))
			if !inAcross && !inDown {
				issues = append(issues, fmt.Sprintf("square at row %d, column %d is isolated (no across or down entry passes through it)", r, c))
			}
		}
	}
	return issues
}
