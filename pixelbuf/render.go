package pixelbuf

import (
	"strconv"
	"strings"
)

// Pre-computed uint8 → decimal string lookup table.
// Eliminates all integer-to-string conversion in the render hot path.
var itoa [256]string

func init() {
	for i := 0; i < 256; i++ {
		itoa[i] = strconv.Itoa(i)
	}
}

// renderBuf is reused across Render calls. Grow pre-allocates once per call.
// Single-goroutine game loop — no concurrency concern.
var renderBuf strings.Builder

// Render converts the buffer to an ANSI string using half-block characters.
// Each terminal row represents two pixel rows: the upper pixel as the
// foreground color on '▀' and the lower pixel as the background color.
// Odd-height buffers pair the last row with transparent black.
// A zero-size buffer returns an empty string.
func Render(buf *Buffer) string {
	if buf.Width == 0 || buf.Height == 0 {
		return ""
	}

	renderBuf.Reset()
	renderBuf.Grow(buf.Width * (buf.Height/2 + 1) * 40)

	var lastFG, lastBG Color
	fgSet, bgSet := false, false

	rows := buf.Height / 2
	if buf.Height%2 != 0 {
		rows++
	}

	for row := 0; row < rows; row++ {
		y := row * 2
		for x := 0; x < buf.Width; x++ {
			top := buf.pixels[y*buf.Width+x]

			var bottom Color
			if y+1 < buf.Height {
				bottom = buf.pixels[(y+1)*buf.Width+x]
			}

			if top == bottom {
				writeBGCode(&renderBuf, top, &lastBG, &bgSet)
				renderBuf.WriteByte(' ')
			} else {
				writeFGCode(&renderBuf, top, &lastFG, &fgSet)
				writeBGCode(&renderBuf, bottom, &lastBG, &bgSet)
				renderBuf.WriteString("▀")
			}
		}
		renderBuf.WriteString("\x1b[0m")
		if row < rows-1 {
			renderBuf.WriteByte('\n')
		}
		lastFG, lastBG = Color{}, Color{}
		fgSet, bgSet = false, false
	}

	return renderBuf.String()
}

// writeFGCode writes a foreground ANSI escape only if the color changed.
func writeFGCode(sb *strings.Builder, c Color, last *Color, set *bool) {
	if *set && c == *last {
		return
	}
	sb.WriteString("\x1b[38;2;")
	sb.WriteString(itoa[c.R])
	sb.WriteByte(';')
	sb.WriteString(itoa[c.G])
	sb.WriteByte(';')
	sb.WriteString(itoa[c.B])
	sb.WriteByte('m')
	*last = c
	*set = true
}

// writeBGCode writes a background ANSI escape only if the color changed.
func writeBGCode(sb *strings.Builder, c Color, last *Color, set *bool) {
	if *set && c == *last {
		return
	}
	sb.WriteString("\x1b[48;2;")
	sb.WriteString(itoa[c.R])
	sb.WriteByte(';')
	sb.WriteString(itoa[c.G])
	sb.WriteByte(';')
	sb.WriteString(itoa[c.B])
	sb.WriteByte('m')
	*last = c
	*set = true
}
