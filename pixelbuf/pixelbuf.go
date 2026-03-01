// Package pixelbuf provides a pixel buffer and half-block terminal renderer
// for graphical output. It is independent of the dungeon map renderer in
// the renderer/ package.
package pixelbuf

// Color represents an RGBA color with 8 bits per channel.
// Alpha is straight (not premultiplied).
type Color struct {
	R, G, B, A uint8
}

// Palette colors (unexported — consumers construct their own).
var (
	transparent = Color{0, 0, 0, 0}
	black       = Color{0, 0, 0, 255}
	white       = Color{255, 255, 255, 255}
	red         = Color{255, 0, 0, 255}
	green       = Color{0, 255, 0, 255}
	blue        = Color{0, 0, 255, 255}
	yellow      = Color{255, 255, 0, 255}
	cyan        = Color{0, 255, 255, 255}
	magenta     = Color{255, 0, 255, 255}
)

// blend composites src over dst using source-over alpha blending.
// Uses (x + 128) >> 8 approximation to avoid integer division.
func blend(src, dst Color) Color {
	if src.A == 255 {
		return src
	}
	if src.A == 0 {
		return dst
	}
	sa := uint16(src.A)
	da := uint16(255 - src.A)
	return Color{
		R: uint8((uint16(src.R)*sa + uint16(dst.R)*da + 128) >> 8),
		G: uint8((uint16(src.G)*sa + uint16(dst.G)*da + 128) >> 8),
		B: uint8((uint16(src.B)*sa + uint16(dst.B)*da + 128) >> 8),
		A: uint8((sa + uint16(dst.A)*da/255)),
	}
}

// Buffer represents a 2D grid of RGBA pixels in row-major order.
type Buffer struct {
	Width  int
	Height int
	pixels []Color
}

// NewBuffer creates a buffer with all pixels set to transparent black.
func NewBuffer(w, h int) *Buffer {
	return &Buffer{
		Width:  w,
		Height: h,
		pixels: make([]Color, w*h),
	}
}

// Set writes a color at (x, y). Out-of-bounds writes are silently ignored.
func (b *Buffer) Set(x, y int, c Color) {
	if !b.InBounds(x, y) {
		return
	}
	b.pixels[y*b.Width+x] = c
}

// At reads the color at (x, y). Out-of-bounds reads return transparent black.
func (b *Buffer) At(x, y int) Color {
	if !b.InBounds(x, y) {
		return Color{}
	}
	return b.pixels[y*b.Width+x]
}

// Clear fills every pixel with the given color.
func (b *Buffer) Clear(c Color) {
	for i := range b.pixels {
		b.pixels[i] = c
	}
}

// InBounds reports whether (x, y) is within the buffer dimensions.
func (b *Buffer) InBounds(x, y int) bool {
	return x >= 0 && x < b.Width && y >= 0 && y < b.Height
}

// FillRect fills a rectangular region with a solid color, clipped to buffer bounds.
func FillRect(dst *Buffer, x, y, w, h int, c Color) {
	// Clamp to buffer bounds.
	x0 := max(0, x)
	y0 := max(0, y)
	x1 := min(dst.Width, x+w)
	y1 := min(dst.Height, y+h)

	for py := y0; py < y1; py++ {
		row := py * dst.Width
		for px := x0; px < x1; px++ {
			dst.pixels[row+px] = c
		}
	}
}

// Blit draws src onto dst at position (dx, dy).
// Transparent pixels (A=0) are skipped. Semi-transparent pixels are alpha-blended.
// Pixels outside dst bounds are clipped.
func Blit(dst, src *Buffer, dx, dy int) {
	// Pre-compute clipped source rectangle.
	srcX0 := max(0, -dx)
	srcY0 := max(0, -dy)
	srcX1 := min(src.Width, dst.Width-dx)
	srcY1 := min(src.Height, dst.Height-dy)

	for sy := srcY0; sy < srcY1; sy++ {
		srcRow := sy * src.Width
		dstRow := (dy+sy)*dst.Width + dx
		for sx := srcX0; sx < srcX1; sx++ {
			pixel := src.pixels[srcRow+sx]
			if pixel.A == 0 {
				continue
			}
			if pixel.A == 255 {
				dst.pixels[dstRow+sx] = pixel
			} else {
				dst.pixels[dstRow+sx] = blend(pixel, dst.pixels[dstRow+sx])
			}
		}
	}
}
