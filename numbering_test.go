package xgrid

import (
	"reflect"
	"testing"
)

func TestNumberSimpleGrid(t *testing.T) {
	// 1 2 3
	// 4 . .
	// . . .
	rows := []string{
		"...",
		"...",
		".#.",
	}
	g := mustParse(t, rows, Options{})

	got := Number(g)
	want := [][]int{
		{1, 2, 3},
		{4, 0, 0},
		{0, 0, 0},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Number() = %v, want %v", got, want)
	}
}

func TestNumberAllBlack(t *testing.T) {
	g := mustParse(t, []string{"###", "###"}, Options{})
	got := Number(g)
	want := [][]int{{0, 0, 0}, {0, 0, 0}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Number() = %v, want %v", got, want)
	}
}

func TestNumberIsolatedSquareUnnumbered(t *testing.T) {
	g := mustParse(t, []string{"###", "#.#", "###"}, Options{})
	got := Number(g)
	if got[1][1] != 0 {
		t.Errorf("Number()[1][1] = %d, want 0 (isolated square starts no entry)", got[1][1])
	}
}

func TestNumberOnlyDownStart(t *testing.T) {
	// (0,1) has a black square to its right, so it can't start an across
	// entry, but a white square below it, so it starts a down entry.
	rows := []string{
		"#.",
		"..",
	}
	g := mustParse(t, rows, Options{})
	got := Number(g)
	if got[0][1] == 0 {
		t.Errorf("Number()[0][1] = 0, want nonzero (starts a down entry)")
	}
	if got[1][0] == 0 {
		t.Errorf("Number()[1][0] = 0, want nonzero (starts an across entry)")
	}
	if got[1][1] != 0 {
		t.Errorf("Number()[1][1] = %d, want 0 (continuation of both entries)", got[1][1])
	}
}

func TestStartsAcrossAndStartsDown(t *testing.T) {
	g := mustParse(t, []string{"#.#", "...", "#.#"}, Options{})

	if StartsAcross(g, 0, 1) {
		t.Error("StartsAcross(0,1) = true, want false (single square, no entry)")
	}
	if !StartsAcross(g, 1, 0) {
		t.Error("StartsAcross(1,0) = false, want true")
	}
	if StartsAcross(g, 1, 1) {
		t.Error("StartsAcross(1,1) = true, want false (not the first square of the run)")
	}
	if !StartsDown(g, 0, 1) {
		t.Error("StartsDown(0,1) = false, want true")
	}
	if StartsDown(g, 1, 1) {
		t.Error("StartsDown(1,1) = true, want false (not the first square of the run)")
	}
}
