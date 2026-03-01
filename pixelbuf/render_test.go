package pixelbuf

import (
	"strings"
	"testing"
)

// fg returns an ANSI foreground escape for a color.
func fg(c Color) string {
	return "\x1b[38;2;" + itoa[c.R] + ";" + itoa[c.G] + ";" + itoa[c.B] + "m"
}

// bg returns an ANSI background escape for a color.
func bg(c Color) string {
	return "\x1b[48;2;" + itoa[c.R] + ";" + itoa[c.G] + ";" + itoa[c.B] + "m"
}

const reset = "\x1b[0m"

func TestRender(t *testing.T) {
	tests := []struct {
		name     string
		buf      *Buffer
		expected string
	}{
		{
			name:     "empty_buffer",
			buf:      NewBuffer(0, 0),
			expected: "",
		},
		{
			name: "1x2_solid_red",
			buf:  solidBuffer(1, 2, red),
			// top == bottom → bg + space (same-color optimization).
			expected: bg(red) + " " + reset,
		},
		{
			name: "1x2_red_over_blue",
			buf: func() *Buffer {
				b := NewBuffer(1, 2)
				b.Set(0, 0, red)
				b.Set(0, 1, blue)
				return b
			}(),
			// Different colors → fg(top) + bg(bottom) + ▀
			expected: fg(red) + bg(blue) + "▀" + reset,
		},
		{
			name: "1x1_odd_height",
			buf: func() *Buffer {
				b := NewBuffer(1, 1)
				b.Set(0, 0, red)
				return b
			}(),
			// Odd height: top=red, bottom=transparent black (zero value).
			expected: fg(red) + bg(transparent) + "▀" + reset,
		},
		{
			name: "2x2_solid_black",
			buf:  solidBuffer(2, 2, black),
			// All same color → bg + space, twice. State tracking: bg emitted once.
			expected: bg(black) + "  " + reset,
		},
		{
			name: "state_tracking_same_row",
			buf:  solidBuffer(4, 2, green),
			// 4 columns, all same color. BG emitted once, then 4 spaces.
			expected: bg(green) + "    " + reset,
		},
		{
			name: "2x2_checkerboard",
			buf: func() *Buffer {
				b := NewBuffer(2, 2)
				b.Set(0, 0, red)
				b.Set(1, 0, blue)
				b.Set(0, 1, blue)
				b.Set(1, 1, red)
				return b
			}(),
			// Col 0: top=red, bottom=blue → fg(red)+bg(blue)+▀
			// Col 1: top=blue, bottom=red → fg(blue)+bg(red)+▀
			expected: fg(red) + bg(blue) + "▀" + fg(blue) + bg(red) + "▀" + reset,
		},
		{
			name: "line_reset_two_rows",
			buf: func() *Buffer {
				b := NewBuffer(1, 4)
				b.Set(0, 0, red)
				b.Set(0, 1, red)
				b.Set(0, 2, blue)
				b.Set(0, 3, blue)
				return b
			}(),
			// Row pair 0: red/red → bg(red) + space + reset
			// Row pair 1: blue/blue → bg(blue) + space + reset
			// Colors re-emitted after line reset.
			expected: bg(red) + " " + reset + "\n" + bg(blue) + " " + reset,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Render(tt.buf)
			if got != tt.expected {
				t.Errorf("Render:\n  got:  %q\n  want: %q", got, tt.expected)
			}
		})
	}
}

// --- Pipeline scenario tests ---

func TestPipelineFillRectThenRender(t *testing.T) {
	// 8x4 black buffer → FillRect 4x2 red at (2,1) → Render → verify.
	buf := solidBuffer(8, 4, black)
	FillRect(buf, 2, 1, 4, 2, red)

	result := Render(buf)
	lines := strings.Split(result, "\n")

	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines, got %d", len(lines))
	}

	// Line 0 pairs rows 0+1. Row 0 is all black. Row 1 has red at cols 2-5.
	// Cols 0-1: top=black, bottom=black → same → bg(black)+space
	// Cols 2-5: top=black, bottom=red → different → fg(black)+bg(red)+▀
	// Cols 6-7: top=black, bottom=black → same → bg(black)+space
	line0 := lines[0]
	if !strings.Contains(line0, bg(red)) {
		t.Error("Line 0 should contain red background for the filled region")
	}

	// Line 1 pairs rows 2+3. Row 2 has red at cols 2-5. Row 3 is all black.
	// Cols 2-5: top=red, bottom=black → different → fg(red)+bg(black)+▀
	line1 := strings.TrimSuffix(lines[1], reset)
	if !strings.Contains(line1, fg(red)) {
		t.Error("Line 1 should contain red foreground for the filled region")
	}
}

func TestPipelineBlitThenRender(t *testing.T) {
	// 6x4 black buffer → Blit 2x2 blue at (1,1) → Render → verify.
	dst := solidBuffer(6, 4, black)
	src := solidBuffer(2, 2, blue)
	Blit(dst, src, 1, 1)

	result := Render(dst)
	lines := strings.Split(result, "\n")

	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines, got %d", len(lines))
	}

	// Row pair 0 (rows 0+1): col 1-2 at row 1 are blue.
	if !strings.Contains(lines[0], bg(blue)) {
		t.Error("Line 0 should contain blue from the blitted buffer")
	}

	// Row pair 1 (rows 2+3): col 1-2 at row 2 are blue.
	if !strings.Contains(lines[1], fg(blue)) || !strings.Contains(lines[1], bg(black)) {
		t.Error("Line 1 should contain blue foreground from the blitted buffer")
	}
}

func TestPipelineBlitClippedThenRender(t *testing.T) {
	// 4x4 buffer → Blit 3x3 red at (2,2) → only top-left 2x2 of src visible.
	dst := solidBuffer(4, 4, black)
	src := solidBuffer(3, 3, red)
	Blit(dst, src, 2, 2)

	result := Render(dst)

	// The red pixels should appear in the output.
	if !strings.Contains(result, bg(red)) && !strings.Contains(result, fg(red)) {
		t.Error("Clipped blit should still produce visible red pixels in render output")
	}

	// Verify it didn't crash and produced 2 lines.
	lines := strings.Split(result, "\n")
	if len(lines) != 2 {
		t.Fatalf("Expected 2 lines, got %d", len(lines))
	}
}
