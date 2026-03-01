# Phase 1 Execution Plan: The Pixel Canvas (`pixelbuf/`)

## Context

This is the foundational rendering package for the combat initiative. Every visual element in combat — player, enemy, arena, particles, health bars — ultimately becomes pixels in a buffer rendered to the terminal via half-block characters (`▀`/`▄`) with ANSI 24-bit true color. This package is Phase 1 because nothing visual can exist without it.

`pixelbuf/` is a leaf package: it imports only the Go standard library, knows nothing about games or combat, and sits alongside `renderer/` as a rendering primitive. The combat engine (`combat/engine/`, Phase 2) will **never** import `pixelbuf/` — only the wiring layer (`main.go` / `cmd/combat-proto/`) bridges domain and presentation.

```
Dependency graph (through Phase 3):

world/              (shared domain types)
game/               (dungeon domain, imports world/)
combat/engine/      (combat domain, imports world/, NOT pixelbuf)
generator/          (imports world/)
pixelbuf/           (rendering primitive, imports ONLY stdlib)
renderer/           (dungeon rendering, imports world/)
main.go             (wiring, imports everything)
```

Import path: `text-adventure-v2/pixelbuf` — a package within the existing module, not a separate `go.mod`.

## Design Decisions

### Three reviewers shaped this plan:

**Rob Pike (Go idiom)** — Killed the `Renderer` interface (premature; this codebase has zero interfaces). Killed `color.RGBA` from stdlib (carries premultiplied alpha baggage we never use). Collapsed 5 source files to 2. Killed `Sprite.Name`. Eliminated string-returning format helpers from the hot path.

**Casey Muratori (performance)** — Confirmed `[]Color` array-of-structs layout is correct for our access patterns. Added persistent `strings.Builder` reuse (zero allocations after first frame). Added `[256]string` lookup table for uint8→string conversion. Fixed Blit to pre-clip and use direct slice indexing. Confirmed state tracking ROI is 2500:1.

**Martin Fowler (architecture)** — Drew the boundary: `Sprite` is a rendering-layer type, combat engine must never import it. Flagged the need for pipeline scenario tests (Blit→Render). Made `blend` and `pixels` unexported. Added dependency diagram. Noted Phase 3 must give combat full screen ownership (no lipgloss wrapping of pixelbuf output).

### Patterns We're Using

| Pattern | Why |
|---------|-----|
| **Custom `Color` type** (`struct { R, G, B, A uint8 }`) | Avoids `image/color` import and premultiplied alpha confusion. We control the type fully. Zero external deps. |
| **Package-level `Render` function** | Matches existing `renderer.RenderMap(MapView) string`. No interface — this codebase has none. Consumer defines interface if needed later. |
| **Persistent `strings.Builder` reuse** | Package-level `var` reset each frame. Zero allocations after warmup. Single-goroutine game loop means no concurrency concern. |
| **`[256]string` lookup table** | Pre-computed uint8→decimal strings at `init()`. Eliminates all `strconv.Itoa` allocation in the render hot path. |
| **Pre-clipped Blit with direct slice indexing** | Compute visible rect once, then loop without `At()`/`Set()`/`InBounds()` calls. |
| **`(x + 128) >> 8` alpha blend** | Avoids integer division. Max error: 1/255 per channel — invisible at 8-bit precision. |
| **Unexported `blend`, `pixels` field** | `blend` is an implementation detail of `Blit`. `pixels` forces use of `Set`/`At` API. Exported only when a consumer actually needs them. |
| **Table-driven tests + pipeline scenario tests** | Tables for exact ANSI verification. Scenarios for Blit→Render composition. |
| **Two source files** | `pixelbuf.go` (types + operations) and `render.go` (ANSI output). Matches codebase convention of single-file packages. |

### Patterns We're Skipping

| Pattern | Why |
|---------|-----|
| **`Renderer` interface** | Zero implementations besides half-block. Consumer defines the interface when they need polymorphism. |
| **`image/color.RGBA`** | Premultiplied alpha footgun. Unnecessary import. Our `Color` is 4 bytes, same layout, no baggage. |
| **Separate `Sprite` type (for now)** | `Blit` takes `*Buffer` as source. A read-only "sprite" is just a `*Buffer` you don't call `Set()` on. Avoids duplicating the pixel grid type and limiting Blit to only sprite→buffer compositing. If Phase 4 needs a distinct Sprite with animation frames, it's a backward-compatible addition. |
| **`fmt.Stringer` on Buffer** | Rendering is explicit via `Render()`, not implicit via `String()`. |
| **Functional options** | Not enough config knobs. Constructor + struct literal. |
| **Exported palette vars** | Mutable exported vars are a footgun. Palette constants are unexported, used in tests. Consumers write `pixelbuf.Color{255, 0, 0, 255}` — it's 4 fields. |

---

## Package Structure

```
pixelbuf/
├── pixelbuf.go              # Color, Buffer, Blit, FillRect, blend, palette
├── render.go                # Render function, half-block algorithm, lookup table
├── pixelbuf_test.go         # Color, Buffer, Blit, FillRect tests
├── render_test.go           # ANSI output tests (table-driven) + pipeline scenarios
└── test_helpers_test.go     # solidBuffer, checkerBuffer, etc.
```

---

## File 1: `pixelbuf/pixelbuf.go`

### Color Type

```go
// Color represents an RGBA color with 8 bits per channel.
// Alpha is straight (not premultiplied).
type Color struct {
    R, G, B, A uint8
}
```

No methods. Value type. 4 bytes. Compared with `==`.

### Palette (unexported)

```go
var (
    transparent = Color{0, 0, 0, 0}
    black       = Color{0, 0, 0, 255}
    white       = Color{255, 255, 255, 255}
    // ... red, green, blue, yellow, cyan, magenta
)
```

Used internally and in tests. Not exported — consumers construct their own colors.

### blend (unexported)

```go
// blend composites src over dst using source-over alpha blending.
func blend(src, dst Color) Color
```

Implementation uses `(uint16(src.R)*uint16(src.A) + uint16(dst.R)*uint16(255-src.A) + 128) >> 8` per channel. Fast paths: `src.A == 255` → return src. `src.A == 0` → return dst.

### Buffer Type

```go
// Buffer represents a 2D grid of RGBA pixels in row-major order.
type Buffer struct {
    Width  int
    Height int
    pixels []Color // unexported — use Set/At
}
```

Methods:

```go
func NewBuffer(w, h int) *Buffer           // all pixels = transparent black
func (b *Buffer) Set(x, y int, c Color)    // silent no-op if out of bounds
func (b *Buffer) At(x, y int) Color        // returns transparent black if out of bounds
func (b *Buffer) Clear(c Color)            // fill every pixel
func (b *Buffer) InBounds(x, y int) bool   // bounds check
```

`pixels` is unexported. `Width` and `Height` are exported (read-only intent; Go has no `readonly` modifier). `NewBuffer` is the only way to get a consistent Buffer.

### Blit

```go
// Blit draws src onto dst at position (dx, dy).
// Transparent pixels (A=0) are skipped. Semi-transparent pixels are alpha-blended.
// Pixels outside dst bounds are clipped.
func Blit(dst, src *Buffer, dx, dy int)
```

Implementation — pre-clipped, direct slice indexing:

```go
func Blit(dst, src *Buffer, dx, dy int) {
    // Pre-compute clipped source rectangle
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
                continue // transparent, skip
            }
            if pixel.A == 255 {
                dst.pixels[dstRow+sx] = pixel // opaque, overwrite
            } else {
                dst.pixels[dstRow+sx] = blend(pixel, dst.pixels[dstRow+sx])
            }
        }
    }
}
```

Note: both `dst` and `src` are `*Buffer`. This means you can composite sub-buffers (health bar buffer onto main frame), not just sprites. If Phase 4 needs a distinct `Sprite` type with animation frames, `Blit` can be overloaded or a `BlitSprite` variant added — backward-compatible.

### FillRect

```go
// FillRect fills a rectangular region with a solid color, clipped to buffer bounds.
func FillRect(dst *Buffer, x, y, w, h int, c Color)
```

Implementation: clamp to bounds, nested loop, direct `dst.pixels[]` write.

---

## File 2: `pixelbuf/render.go`

### Lookup Table

```go
var itoa [256]string

func init() {
    for i := 0; i < 256; i++ {
        itoa[i] = strconv.Itoa(i)
    }
}
```

256 pre-computed strings. ~1KB total. Lives in L1 cache permanently. Eliminates all integer-to-string conversion in the hot path.

### Persistent Builder

```go
var renderBuf strings.Builder
```

Package-level. `Reset()` at start of each `Render()` call — keeps backing `[]byte` allocated. Zero allocations after the first frame. Safe because the game is single-goroutine.

### Render Function

```go
// Render converts the buffer to an ANSI string using half-block characters.
// Each terminal row represents two pixel rows: upper pixel as foreground
// color on '▀', lower pixel as background color.
// Odd-height buffers pair the last row with black.
func Render(buf *Buffer) string
```

**The algorithm:**

```go
func Render(buf *Buffer) string {
    renderBuf.Reset()
    renderBuf.Grow(buf.Width * (buf.Height / 2 + 1) * 40)

    var lastFG, lastBG Color
    fgSet, bgSet := false, false

    for y := 0; y < buf.Height; y += 2 {
        for x := 0; x < buf.Width; x++ {
            top := buf.pixels[y*buf.Width+x]

            var bottom Color
            if y+1 < buf.Height {
                bottom = buf.pixels[(y+1)*buf.Width+x]
            }
            // else bottom stays zero-value (transparent black)

            if top == bottom {
                // Same color: background + space (saves fg escape code)
                writeBGCode(&renderBuf, top, &lastBG, &bgSet)
                renderBuf.WriteByte(' ')
            } else {
                // Different: fg=top, bg=bottom, emit ▀
                writeFGCode(&renderBuf, top, &lastFG, &fgSet)
                writeBGCode(&renderBuf, bottom, &lastBG, &bgSet)
                renderBuf.WriteString("▀")
            }
        }
        renderBuf.WriteString("\x1b[0m\n")
        lastFG, lastBG = Color{}, Color{}
        fgSet, bgSet = false, false
    }

    return strings.TrimRight(renderBuf.String(), "\n")
}
```

### Unexported Helpers

```go
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
    // identical to writeFGCode but with "\x1b[48;2;" prefix
}
```

Zero allocations. Direct builder writes. Lookup table for digits. State tracking with 2500:1 ROI on terminal bandwidth savings.

---

## Implementation Order

| Step | What | Test First |
|------|------|-----------|
| 1 | `Color` type + `blend` + palette | `blend` math: opaque/transparent fast paths, 50% blend, self-blend |
| 2 | `Buffer` type: `New`, `Set`, `At`, `Clear`, `InBounds` | Roundtrip, bounds safety, clear |
| 3 | `FillRect` | Within bounds, clipped, zero-size |
| 4 | `Blit` | Opaque, transparent, semi-transparent, clipped, off-screen |
| 5 | `Render` + lookup table + builder reuse | Exact ANSI verification for known small buffers |
| 6 | Pipeline scenario tests | Blit sprite onto buffer → Render → verify ANSI |

`go test ./pixelbuf/...` after each step. Green before proceeding.

---

## Test Plan

### `pixelbuf_test.go`

**blend tests (table-driven):**
| Case | src | dst | Expected |
|------|-----|-----|----------|
| Opaque over anything | `{255,0,0,255}` | `{0,0,255,255}` | `{255,0,0,255}` |
| Transparent over anything | `{255,0,0,0}` | `{0,0,255,255}` | `{0,0,255,255}` |
| 50% red over blue | `{255,0,0,128}` | `{0,0,255,255}` | `{128,0,127,...}` |
| 50% red over transparent | `{255,0,0,128}` | `{0,0,0,0}` | `{128,0,0,128}` |
| Self-blend | `{100,100,100,255}` | `{100,100,100,255}` | `{100,100,100,255}` |

**Buffer tests:**
- `TestNewBuffer` — all pixels transparent black
- `TestBufferSetAt` — roundtrip at (0,0), center, (w-1,h-1)
- `TestBufferOutOfBounds` — Set and At on negative, overflow coords: no panic, correct defaults
- `TestBufferClear` — every pixel matches clear color
- `TestBufferInBounds` — true for valid, false for edge cases

**FillRect tests:**
- Within bounds: verify all pixels in rect, none outside
- Partially off-screen: clipped to bounds
- Zero-size: no effect
- Fully off-screen: no pixels changed

**Blit tests:**
- Fully opaque src: dst pixels overwritten
- Transparent pixels: dst preserved where `A == 0`
- Semi-transparent: alpha blending verified
- Partially off-screen (negative dx/dy, exceeds bounds): clipped, no panic
- Fully off-screen: no pixels changed
- Empty src buffer: no effect

### `render_test.go` (table-driven)

| Case | Buffer Setup | Verify |
|------|-------------|--------|
| 1x2 solid red | 1w × 2h, all red | BG red + space (same-color opt) |
| 1x2 red/blue | top=red, bottom=blue | FG red + BG blue + ▀ |
| 2x2 checkerboard | alternating red/blue | Color changes per cell |
| 1x1 odd height | single red pixel | Red top paired with black bottom |
| 4x4 solid black | all black | Minimal escape codes |
| 0x0 empty | zero dimensions | Empty string |
| State tracking | row of same color | Escape emitted once, not per pixel |
| Line reset | 2 rows, different colors | Colors re-emitted after newline |

**Pipeline scenario tests (2-3):**
1. Create 8x4 black buffer → FillRect a 4x2 red rect at (2,1) → Render → verify red block appears at correct position in ANSI output
2. Create 6x4 black buffer → Create 2x2 blue "sprite" buffer → Blit at (1,1) → Render → verify blue pixels at correct terminal cells
3. Create 4x4 buffer → Blit partially off-screen sprite → Render → verify clipped output is correct

### `test_helpers_test.go`

```go
func solidBuffer(w, h int, c Color) *Buffer
func checkerBuffer(w, h int, a, b Color) *Buffer
```

---

## Verification

```sh
go test ./pixelbuf/...           # all pass
go test ./pixelbuf/... -v        # see individual names
go test ./pixelbuf/... -cover    # target: >90% coverage
go vet ./pixelbuf/...            # no warnings
go build ./...                   # existing code unaffected
go test ./...                    # existing tests still pass
```

## Boundary: What This Does NOT Touch

- No changes to `game/`, `renderer/`, `generator/`, `world/`, or `main.go`
- No `go.mod` dependency additions — `pixelbuf/` imports only `strings` and `strconv` from stdlib
- No Bubble Tea integration
- No combat logic, physics, input handling
- No `Sprite` type — deferred to Phase 2/4 when animation frames create a need

## Forward-Looking Notes (For Future Phase Planners)

These are NOT in scope for Phase 1. Documented here to prevent Phase 2+ from making avoidable mistakes:

1. **`combat/engine/` must never import `pixelbuf/`**. Entity dimensions live in `world/` or `combat/types/`. The wiring layer maps engine entities to pixel buffers.
2. **Combat mode owns the full screen**. No lipgloss wrapping of pixelbuf ANSI output. The mode dispatcher in `main.go` is a clean either/or.
3. **Phase 4 will need `BlitTinted`** for hit-flash (render entity all-white). This is a backward-compatible function addition. Don't add it now.
4. **Phase 2 may want a `Sprite` type** with animation frames, flip-horizontal. Adding it is backward-compatible — `Blit` already works with `*Buffer`, a new `BlitSprite(*Buffer, *Sprite, dx, dy)` can coexist.
5. **A `combat/view` or `combat/renderer` package** (analogous to `renderer.MapView`) will be needed to translate engine state into draw calls. Don't bake this logic into `main.go`.

## Performance Budget

Casey Muratori's estimates for a 120x60 pixel buffer (60x30 terminal cells):
- Frame output size: ~37 KB with state tracking (~74 KB without)
- Sustained throughput at 30fps: ~1.1 MB/s
- Terminal headroom: ~2x on modern emulators (Windows Terminal)
- Allocations per frame after warmup: **zero** (persistent builder)
- The render hot path: no `fmt.Sprintf`, no `strconv.Itoa`, no `At()`/`Set()` method calls — direct slice access and lookup table writes
