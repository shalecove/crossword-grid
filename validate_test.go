package xgrid

import (
	"strings"
	"testing"
)

func mustParse(t *testing.T, lines []string, opts Options) *Grid {
	t.Helper()
	g, err := Parse(lines, opts)
	if err != nil {
		t.Fatalf("Parse: %v", err)
	}
	return g
}

func TestValidateAllWhiteGridPasses(t *testing.T) {
	g := mustParse(t, []string{"...", "...", "..."}, Options{})
	if err := Validate(g, Options{}); err != nil {
		t.Fatalf("Validate: %v", err)
	}
}

func TestValidateNoWhiteSquares(t *testing.T) {
	g := mustParse(t, []string{"###", "###", "###"}, Options{})
	err := Validate(g, Options{})
	if err == nil {
		t.Fatal("Validate = nil error, want error")
	}
	if !strings.Contains(err.Error(), "no white squares") {
		t.Errorf("Validate error = %q, want it to mention no white squares", err)
	}
}

func TestValidateStrictRejectsAsymmetricGrid(t *testing.T) {
	rows := []string{
		"#....",
		".....",
		".....",
		".....",
		".....",
	}
	g := mustParse(t, rows, Options{})

	err := Validate(g, Options{})
	if err == nil {
		t.Fatal("Validate = nil error, want error")
	}
	ve, ok := err.(*ValidationError)
	if !ok {
		t.Fatalf("Validate error type = %T, want *ValidationError", err)
	}
	if len(ve.Issues) != 1 || !strings.Contains(ve.Issues[0], "rotationally symmetric") {
		t.Errorf("Issues = %v, want a single symmetry issue", ve.Issues)
	}

	if err := Validate(g, Options{Lenient: true}); err != nil {
		t.Errorf("Validate lenient: %v, want nil (only the no-white-squares check runs)", err)
	}
}

func TestValidateFlagsShortEntriesAndIsolatedSquares(t *testing.T) {
	g := mustParse(t, []string{"###", "#.#", "###"}, Options{})
	err := Validate(g, Options{})
	if err == nil {
		t.Fatal("Validate = nil error, want error")
	}
	if !strings.Contains(err.Error(), "isolated") {
		t.Errorf("Validate error = %q, want it to mention an isolated square", err)
	}
}

func TestIsSymmetric(t *testing.T) {
	sym := mustParse(t, []string{"#.#", "...", "#.#"}, Options{})
	if !isSymmetric(sym) {
		t.Error("isSymmetric = false, want true")
	}

	asym := mustParse(t, []string{"#..", "...", "..."}, Options{})
	if isSymmetric(asym) {
		t.Error("isSymmetric = true, want false")
	}
}

func TestHasAnyWhiteSquare(t *testing.T) {
	white := mustParse(t, []string{"#.#"}, Options{})
	if !hasAnyWhiteSquare(white) {
		t.Error("hasAnyWhiteSquare = false, want true")
	}

	black := mustParse(t, []string{"###"}, Options{})
	if hasAnyWhiteSquare(black) {
		t.Error("hasAnyWhiteSquare = true, want false")
	}
}

func TestCheckEntryLengths(t *testing.T) {
	g := mustParse(t, []string{"..#", "..#", "###"}, Options{})
	issues := checkEntryLengths(g)

	want := []string{
		"row 0 has a 2-letter run starting at column 0 (minimum is 3)",
		"row 1 has a 2-letter run starting at column 0 (minimum is 3)",
		"column 0 has a 2-letter run starting at row 0 (minimum is 3)",
		"column 1 has a 2-letter run starting at row 0 (minimum is 3)",
	}
	if len(issues) != len(want) {
		t.Fatalf("checkEntryLengths = %v, want %v", issues, want)
	}
	for i, w := range want {
		if issues[i] != w {
			t.Errorf("issues[%d] = %q, want %q", i, issues[i], w)
		}
	}
}

func TestCheckNoIsolatedSquares(t *testing.T) {
	g := mustParse(t, []string{"###", "#.#", "###"}, Options{})
	issues := checkNoIsolatedSquares(g)
	want := "square at row 1, column 1 is isolated (no across or down entry passes through it)"
	if len(issues) != 1 || issues[0] != want {
		t.Errorf("checkNoIsolatedSquares = %v, want [%q]", issues, want)
	}
}

func TestValidationErrorMessage(t *testing.T) {
	err := &ValidationError{Issues: []string{"a", "b"}}
	want := "xgrid: grid failed validation (2 issue(s)): a; b"
	if err.Error() != want {
		t.Errorf("Error() = %q, want %q", err.Error(), want)
	}
}
