// Package xgrid represents crossword grids and checks whether they follow
// the conventions of a standard crossword puzzle: rectangular shape,
// rotational symmetry, minimum entry length, and so on.
//
// A grid is built from plain text: one string per row, '#' for a black
// square, '.' for an empty white square, and A-Z for a filled letter.
package xgrid

import (
	"errors"
	"fmt"
	"strings"
)

const (
	Black byte = '#'
	Empty byte = '.'
)

// Options controls how permissive parsing and validation are. The zero
// value is strict: ragged rows, unrecognized characters, and grids that
// break standard crossword conventions are all rejected. Set Lenient to
// relax those checks for input that didn't come from a trusted source,
// or for non-standard puzzle styles (cryptics, grids with no symmetry,
// two-letter entries, and so on).
type Options struct {
	Lenient bool
}

// Grid is an immutable rectangular arrangement of cells.
type Grid struct {
	Rows, Cols int
	cells      [][]byte
}

// Parse builds a Grid from one string per row. In strict mode every row
// must have the same length and every character must be '#', '.', or an
// uppercase letter. In lenient mode: rows shorter than the widest row are
// padded with black squares, lowercase letters are upshifted, '_' and ' '
// are treated as black squares, and '?' is treated as an empty square.
func Parse(lines []string, opts Options) (*Grid, error) {
	if len(lines) == 0 {
		return nil, errors.New("xgrid: at least one row is required")
	}

	width := len(lines[0])
	if opts.Lenient {
		for _, line := range lines {
			if len(line) > width {
				width = len(line)
			}
		}
	}
	if width == 0 {
		return nil, errors.New("xgrid: rows must have at least one column")
	}

	cells := make([][]byte, len(lines))
	for i, line := range lines {
		if !opts.Lenient && len(line) != width {
			return nil, fmt.Errorf("xgrid: row %d has %d columns, want %d (grid is ragged; use Options.Lenient to pad it)", i, len(line), width)
		}
		row := make([]byte, width)
		for j := 0; j < width; j++ {
			var raw byte = Black
			if j < len(line) {
				raw = line[j]
			}
			norm, err := normalizeCell(raw, opts)
			if err != nil {
				return nil, fmt.Errorf("xgrid: row %d, col %d: %w", i, j, err)
			}
			row[j] = norm
		}
		cells[i] = row
	}

	return &Grid{Rows: len(lines), Cols: width, cells: cells}, nil
}

func normalizeCell(b byte, opts Options) (byte, error) {
	switch {
	case b == Black:
		return Black, nil
	case b == Empty:
		return Empty, nil
	case b >= 'A' && b <= 'Z':
		return b, nil
	case opts.Lenient && b >= 'a' && b <= 'z':
		return b - ('a' - 'A'), nil
	case opts.Lenient && (b == '_' || b == ' '):
		return Black, nil
	case opts.Lenient && b == '?':
		return Empty, nil
	default:
		return 0, fmt.Errorf("unrecognized cell character %q (use Options.Lenient to relax this)", rune(b))
	}
}

// At returns the raw cell value at (row, col). It panics if the
// coordinates are out of range, same as a slice index would.
func (g *Grid) At(row, col int) byte {
	return g.cells[row][col]
}

// IsBlack reports whether the cell at (row, col) is a black square.
func (g *Grid) IsBlack(row, col int) bool {
	return g.cells[row][col] == Black
}

// String renders the grid back to its line-per-row text form.
func (g *Grid) String() string {
	var b strings.Builder
	for i, row := range g.cells {
		if i > 0 {
			b.WriteByte('\n')
		}
		b.Write(row)
	}
	return b.String()
}
