package pixelbuf

import "testing"

// --- blend tests (table-driven) ---

func TestBlend(t *testing.T) {
	tests := []struct {
		name     string
		src, dst Color
		expected Color
	}{
		{
			name:     "opaque_over_anything",
			src:      Color{255, 0, 0, 255},
			dst:      Color{0, 0, 255, 255},
			expected: Color{255, 0, 0, 255},
		},
		{
			name:     "transparent_over_anything",
			src:      Color{255, 0, 0, 0},
			dst:      Color{0, 0, 255, 255},
			expected: Color{0, 0, 255, 255},
		},
		{
			name:     "half_red_over_blue",
			src:      Color{255, 0, 0, 128},
			dst:      Color{0, 0, 255, 255},
			expected: Color{128, 0, 127, 255},
		},
		{
			name:     "half_red_over_transparent",
			src:      Color{255, 0, 0, 128},
			dst:      Color{0, 0, 0, 0},
			expected: Color{128, 0, 0, 128},
		},
		{
			name:     "self_blend_opaque",
			src:      Color{100, 100, 100, 255},
			dst:      Color{100, 100, 100, 255},
			expected: Color{100, 100, 100, 255},
		},
		{
			name:     "half_white_over_black",
			src:      Color{255, 255, 255, 128},
			dst:      Color{0, 0, 0, 255},
			expected: Color{128, 128, 128, 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := blend(tt.src, tt.dst)
			// Allow ±1 per channel for the >>8 approximation.
			if !colorClose(got, tt.expected, 1) {
				t.Errorf("blend(%v, %v) = %v, want %v", tt.src, tt.dst, got, tt.expected)
			}
		})
	}
}

func colorClose(a, b Color, tolerance uint8) bool {
	return absDiff(a.R, b.R) <= tolerance &&
		absDiff(a.G, b.G) <= tolerance &&
		absDiff(a.B, b.B) <= tolerance &&
		absDiff(a.A, b.A) <= tolerance
}

func absDiff(a, b uint8) uint8 {
	if a > b {
		return a - b
	}
	return b - a
}

// --- Buffer tests ---

func TestNewBuffer(t *testing.T) {
	buf := NewBuffer(3, 2)
	if buf.Width != 3 || buf.Height != 2 {
		t.Errorf("NewBuffer dimensions: got %dx%d, want 3x2", buf.Width, buf.Height)
	}
	for y := 0; y < 2; y++ {
		for x := 0; x < 3; x++ {
			c := buf.At(x, y)
			if c != transparent {
				t.Errorf("NewBuffer pixel (%d,%d) = %v, want transparent", x, y, c)
			}
		}
	}
}

func TestBufferSetAt(t *testing.T) {
	buf := NewBuffer(4, 4)
	positions := [][2]int{{0, 0}, {2, 2}, {3, 3}}
	for _, pos := range positions {
		c := Color{uint8(pos[0] * 50), uint8(pos[1] * 50), 100, 255}
		buf.Set(pos[0], pos[1], c)
		got := buf.At(pos[0], pos[1])
		if got != c {
			t.Errorf("Set/At at (%d,%d): got %v, want %v", pos[0], pos[1], got, c)
		}
	}
}

func TestBufferOutOfBounds(t *testing.T) {
	buf := NewBuffer(2, 2)
	buf.Set(0, 0, red)

	// Out-of-bounds Set: no panic, no effect on valid pixels.
	buf.Set(-1, 0, blue)
	buf.Set(2, 0, blue)
	buf.Set(0, -1, blue)
	buf.Set(0, 2, blue)

	if buf.At(0, 0) != red {
		t.Error("OOB Set corrupted valid pixel")
	}

	// Out-of-bounds At: returns transparent, no panic.
	if buf.At(-1, 0) != transparent {
		t.Error("OOB At(-1,0) should return transparent")
	}
	if buf.At(2, 0) != transparent {
		t.Error("OOB At(2,0) should return transparent")
	}
}

func TestBufferClear(t *testing.T) {
	buf := NewBuffer(3, 3)
	buf.Set(1, 1, red)
	buf.Clear(blue)
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			if buf.At(x, y) != blue {
				t.Errorf("Clear: pixel (%d,%d) = %v, want blue", x, y, buf.At(x, y))
			}
		}
	}
}

func TestBufferInBounds(t *testing.T) {
	buf := NewBuffer(3, 2)
	tests := []struct {
		x, y     int
		expected bool
	}{
		{0, 0, true}, {2, 1, true}, {1, 0, true},
		{-1, 0, false}, {3, 0, false}, {0, -1, false}, {0, 2, false},
	}
	for _, tt := range tests {
		if buf.InBounds(tt.x, tt.y) != tt.expected {
			t.Errorf("InBounds(%d,%d) = %v, want %v", tt.x, tt.y, !tt.expected, tt.expected)
		}
	}
}

// --- FillRect tests ---

func TestFillRect(t *testing.T) {
	buf := NewBuffer(4, 4)
	FillRect(buf, 1, 1, 2, 2, red)

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			inRect := x >= 1 && x < 3 && y >= 1 && y < 3
			got := buf.At(x, y)
			if inRect && got != red {
				t.Errorf("FillRect: pixel (%d,%d) inside rect = %v, want red", x, y, got)
			}
			if !inRect && got != transparent {
				t.Errorf("FillRect: pixel (%d,%d) outside rect = %v, want transparent", x, y, got)
			}
		}
	}
}

func TestFillRectClipped(t *testing.T) {
	buf := NewBuffer(3, 3)
	FillRect(buf, -1, -1, 3, 3, green)

	// Should only fill the visible portion: (0,0) to (1,1).
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			inRect := x < 2 && y < 2
			got := buf.At(x, y)
			if inRect && got != green {
				t.Errorf("FillRect clipped: (%d,%d) = %v, want green", x, y, got)
			}
			if !inRect && got != transparent {
				t.Errorf("FillRect clipped: (%d,%d) = %v, want transparent", x, y, got)
			}
		}
	}
}

func TestFillRectZeroSize(t *testing.T) {
	buf := NewBuffer(2, 2)
	FillRect(buf, 0, 0, 0, 0, red)
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			if buf.At(x, y) != transparent {
				t.Errorf("FillRect zero-size: (%d,%d) should be transparent", x, y)
			}
		}
	}
}

func TestFillRectFullyOffScreen(t *testing.T) {
	buf := NewBuffer(2, 2)
	FillRect(buf, 5, 5, 2, 2, red)
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			if buf.At(x, y) != transparent {
				t.Errorf("FillRect off-screen: (%d,%d) should be transparent", x, y)
			}
		}
	}
}

// --- Blit tests ---

func TestBlitOpaque(t *testing.T) {
	dst := NewBuffer(4, 4)
	dst.Clear(black)
	src := solidBuffer(2, 2, red)

	Blit(dst, src, 1, 1)

	for y := 0; y < 4; y++ {
		for x := 0; x < 4; x++ {
			inSrc := x >= 1 && x < 3 && y >= 1 && y < 3
			got := dst.At(x, y)
			if inSrc && got != red {
				t.Errorf("Blit opaque: (%d,%d) = %v, want red", x, y, got)
			}
			if !inSrc && got != black {
				t.Errorf("Blit opaque: (%d,%d) = %v, want black", x, y, got)
			}
		}
	}
}

func TestBlitTransparent(t *testing.T) {
	dst := NewBuffer(3, 3)
	dst.Clear(blue)

	src := NewBuffer(2, 2)
	src.Set(0, 0, red)
	// (1,0), (0,1), (1,1) stay transparent.

	Blit(dst, src, 0, 0)

	if dst.At(0, 0) != red {
		t.Errorf("Blit transparent: (0,0) = %v, want red", dst.At(0, 0))
	}
	if dst.At(1, 0) != blue {
		t.Errorf("Blit transparent: (1,0) = %v, want blue (preserved)", dst.At(1, 0))
	}
}

func TestBlitSemiTransparent(t *testing.T) {
	dst := NewBuffer(1, 1)
	dst.Set(0, 0, Color{0, 0, 255, 255})

	src := NewBuffer(1, 1)
	src.Set(0, 0, Color{255, 0, 0, 128})

	Blit(dst, src, 0, 0)

	got := dst.At(0, 0)
	expected := Color{128, 0, 127, 255}
	if !colorClose(got, expected, 1) {
		t.Errorf("Blit semi-transparent: got %v, want ~%v", got, expected)
	}
}

func TestBlitClippedNegative(t *testing.T) {
	dst := NewBuffer(3, 3)
	dst.Clear(black)
	src := solidBuffer(2, 2, red)

	Blit(dst, src, -1, -1)

	// Only bottom-right pixel of src (1,1) should land at dst (0,0).
	if dst.At(0, 0) != red {
		t.Errorf("Blit clipped negative: (0,0) = %v, want red", dst.At(0, 0))
	}
	if dst.At(1, 0) != black {
		t.Errorf("Blit clipped negative: (1,0) = %v, want black", dst.At(1, 0))
	}
}

func TestBlitClippedOverflow(t *testing.T) {
	dst := NewBuffer(3, 3)
	dst.Clear(black)
	src := solidBuffer(2, 2, green)

	Blit(dst, src, 2, 2)

	// Only top-left pixel of src (0,0) should land at dst (2,2).
	if dst.At(2, 2) != green {
		t.Errorf("Blit clipped overflow: (2,2) = %v, want green", dst.At(2, 2))
	}
	if dst.At(1, 2) != black {
		t.Errorf("Blit clipped overflow: (1,2) = %v, want black", dst.At(1, 2))
	}
}

func TestBlitFullyOffScreen(t *testing.T) {
	dst := NewBuffer(2, 2)
	dst.Clear(black)
	src := solidBuffer(2, 2, red)

	Blit(dst, src, 10, 10)

	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			if dst.At(x, y) != black {
				t.Errorf("Blit off-screen: (%d,%d) = %v, want black", x, y, dst.At(x, y))
			}
		}
	}
}

func TestBlitEmptySrc(t *testing.T) {
	dst := NewBuffer(2, 2)
	dst.Clear(blue)
	src := NewBuffer(0, 0)

	Blit(dst, src, 0, 0) // should not panic or change anything

	if dst.At(0, 0) != blue {
		t.Errorf("Blit empty src: (0,0) = %v, want blue", dst.At(0, 0))
	}
}
