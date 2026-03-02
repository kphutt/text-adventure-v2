package engine

import "testing"

func TestRect_Overlaps_Overlapping(t *testing.T) {
	a := Rect{0, 0, 10, 10}
	b := Rect{5, 5, 10, 10}
	if !a.Overlaps(b) {
		t.Error("expected overlapping rects to overlap")
	}
}

func TestRect_Overlaps_Adjacent(t *testing.T) {
	a := Rect{0, 0, 10, 10}
	b := Rect{10, 0, 10, 10} // shares right edge
	if a.Overlaps(b) {
		t.Error("adjacent rects should not overlap")
	}
}

func TestRect_Overlaps_Separated(t *testing.T) {
	a := Rect{0, 0, 10, 10}
	b := Rect{20, 20, 10, 10}
	if a.Overlaps(b) {
		t.Error("separated rects should not overlap")
	}
}

func TestRect_Overlaps_SameRect(t *testing.T) {
	a := Rect{5, 5, 10, 10}
	if !a.Overlaps(a) {
		t.Error("a rect should overlap itself")
	}
}

func TestRect_Overlaps_ZeroSize(t *testing.T) {
	a := Rect{5, 5, 0, 0}
	b := Rect{0, 0, 10, 10}
	if a.Overlaps(b) {
		t.Error("zero-size rect should not overlap anything")
	}
}

func TestRect_Center(t *testing.T) {
	r := Rect{10, 20, 30, 40}
	cx, cy := r.Center()
	if cx != 25 || cy != 40 {
		t.Errorf("Center() = (%v, %v), want (25, 40)", cx, cy)
	}
}
