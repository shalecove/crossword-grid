package xgrid

import (
	"strings"
	"testing"
)

func TestParseStrictValid(t *testing.T) {
	rows := []string{
		"A.#",
		"...",
		"#.B",
	}
	g, err := Parse(rows, Options{})
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	if g.Rows != 3 || g.Cols != 3 {
		t.Fatalf("got %dx%d grid, want 3x3", g.Rows, g.Cols)
	}
	if g.At(0, 0) != 'A' {
		t.Errorf("At(0,0) = %q, want 'A'", g.At(0, 0))
	}
	if !g.IsBlack(0, 1) {
		t.Errorf("IsBlack(0,1) = false, want true")
	}
	if g.IsBlack(1, 1) {
		t.Errorf("IsBlack(1,1) = true, want false")
	}
	if got := g.String(); got != strings.Join(rows, "\n") {
		t.Errorf("String() = %q, want %q", got, strings.Join(rows, "\n"))
	}
}

func TestParseNoRows(t *testing.T) {
	if _, err := Parse(nil, Options{}); err == nil {
		t.Fatal("Parse(nil, ...) = nil error, want error")
	}
}

func TestParseZeroWidth(t *testing.T) {
	if _, err := Parse([]string{"", ""}, Options{}); err == nil {
		t.Fatal("Parse with empty rows = nil error, want error")
	}
}

func TestParseStrictRaggedRejected(t *testing.T) {
	rows := []string{"###", "#.", "###"}
	if _, err := Parse(rows, Options{}); err == nil {
		t.Fatal("Parse strict with ragged rows = nil error, want error")
	}
}

func TestParseStrictUnrecognizedChar(t *testing.T) {
	rows := []string{"a.#"}
	if _, err := Parse(rows, Options{}); err == nil {
		t.Fatal("Parse strict with lowercase letter = nil error, want error")
	}
}

func TestParseLenientPadsRaggedRows(t *testing.T) {
	rows := []string{"###", "#.", "###"}
	g, err := Parse(rows, Options{Lenient: true})
	if err != nil {
		t.Fatalf("Parse lenient: %v", err)
	}
	if g.Cols != 3 {
		t.Fatalf("Cols = %d, want 3", g.Cols)
	}
	if !g.IsBlack(1, 2) {
		t.Errorf("padded cell (1,2) should be black")
	}
}

func TestParseLenientCharacterSet(t *testing.T) {
	rows := []string{"a_ ?"}
	g, err := Parse(rows, Options{Lenient: true})
	if err != nil {
		t.Fatalf("Parse lenient: %v", err)
	}
	if g.At(0, 0) != 'A' {
		t.Errorf("lowercase letter not upshifted: got %q", g.At(0, 0))
	}
	if !g.IsBlack(0, 1) || !g.IsBlack(0, 2) {
		t.Errorf("'_' and ' ' should both be treated as black squares")
	}
	if g.IsBlack(0, 3) || g.At(0, 3) != Empty {
		t.Errorf("'?' should be treated as an empty square")
	}
}

func TestParseLenientStillRejectsUnknownChars(t *testing.T) {
	if _, err := Parse([]string{"1"}, Options{Lenient: true}); err == nil {
		t.Fatal("Parse lenient with digit = nil error, want error")
	}
}
