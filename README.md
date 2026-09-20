# crossword-grid

A Go library for parsing and validating crossword grids. No executable,
no dependencies — import the package and call it from your own tool.

## Why

Every crossword construction tool ends up needing the same low-level
piece: something that turns a block of text into a grid, checks whether
it follows the conventions of a standard puzzle, and gives useful errors
when it doesn't. This library is that piece, kept separate from any
particular editor, file format, or clue database.

A grid follows a few conventions once it's built:

- 180-degree rotational symmetry (the black-square pattern looks the
  same rotated upside down)
- no isolated white squares — every white square sits in at least one
  across or down entry
- no two-letter entries
- at least one white square, obviously

Not every grid needs to follow these. Cryptic grids are often
asymmetric; some layouts use short entries on purpose. Rather than pick
one behavior, the library defaults to strict and gives you an explicit
way to opt out.

## Usage

```go
package main

import (
	"fmt"

	"github.com/shalecove/crossword-grid"
)

func main() {
	rows := []string{
		"###.....###",
		"##.......##",
		"#.........#",
		".....#.....",
		"...........",
		"....###....",
		"...........",
		".....#.....",
		"#.........#",
		"##.......##",
		"###.....###",
	}

	grid, err := xgrid.Parse(rows, xgrid.Options{})
	if err != nil {
		fmt.Println("parse error:", err)
		return
	}

	if err := xgrid.Validate(grid, xgrid.Options{}); err != nil {
		fmt.Println("invalid grid:", err)
		return
	}

	fmt.Println("grid is valid:")
	fmt.Println(grid)
}
```

Strict mode rejects ragged input outright:

```go
rows := []string{
	"###",
	"#.",   // shorter than the first row
	"###",
}

_, err := xgrid.Parse(rows, xgrid.Options{})
// err: xgrid: row 1 has 2 columns, want 3 (grid is ragged; use Options.Lenient to pad it)

grid, err := xgrid.Parse(rows, xgrid.Options{Lenient: true})
// short row is padded with black squares instead of failing
```

The same escape hatch applies to `Validate`: pass `Options{Lenient: true}`
to accept an asymmetric grid, one with two-letter entries, or one built
from an unusual character set (lowercase letters, `_` for black squares,
`?` for empty squares are all accepted in lenient parsing).

## Status

Early. The grid model, parsing, structural validation, and clue
numbering (`Number`) are in place. Entry extraction, fill/pattern
matching, grid generation, and .puz/.ipuz import/export are not written
yet — see the roadmap in the issue tracker.

## License

MIT, see LICENSE.
