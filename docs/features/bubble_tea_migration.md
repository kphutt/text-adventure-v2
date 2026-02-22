# Feature Design: Bubble Tea TUI Migration

## 1. Project Goal

Migrate the game's terminal UI from raw `tcell` to Charm's [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework. This replaces manual event loops, character-by-character drawing, and fragile input handling with a clean Elm-style architecture (Model → Update → View). Styling moves from raw terminal codes to [Lip Gloss](https://github.com/charmbracelet/lipgloss) declarative styles.

### Motivation

- **Input handling** — tcell requires manual key dispatch and string building. The `bubbles/textinput` component handles cursor movement, editing, and focus out of the box, eliminating an entire class of keyboard interception bugs.
- **Rendering** — tcell draws character-by-character at (x, y) coordinates. Bubble Tea's `View()` returns a single string, and Lip Gloss composes styled sections declaratively.
- **Maintainability** — the Elm architecture enforces a strict `state → view` data flow, making the UI easier to reason about and extend.

---

## 2. Architecture

The game's existing separation of concerns makes this a clean swap — only `main.go` changes.

```
UNCHANGED                          REWRITTEN
┌─────────────┐                    ┌──────────────────────┐
│ game/       │                    │ main.go              │
│   game.go   │──HandleCommand()──→│   model (state)      │
│   structs.go│  Score(), Look()   │   Update (input)     │
│   parser.go │                    │   View (render)      │
├─────────────┤                    │                      │
│ renderer/   │──RenderHUD()──────→│   Uses Lip Gloss to  │
│   renderer.go RenderMap()        │   style existing     │
├─────────────┤                    │   renderer output    │
│ world/      │                    └──────────────────────┘
│ generator/  │                    tcell → bubbletea + lipgloss
└─────────────┘
```

### Why the current design enables a clean swap

1. **Decoupled game logic** — `game.HandleCommand()` is a pure function of `(command string) → (message string, shouldExit bool)`. It has no knowledge of the terminal.
2. **Pure renderer functions** — `renderer.RenderHUD()` and `renderer.RenderMap()` accept a `MapView` struct and return plain strings. No screen references.
3. **MapView as view model** — the renderer already consumes a view-model struct, so the new `View()` method just passes the same data through.

---

## 3. Requirements

### Functional Requirements (FR)

- **FR-1 (Elm Architecture)**: The UI shall be structured as a Bubble Tea `Model` with `Init`, `Update`, and `View` methods.
- **FR-2 (Text Input)**: Player input shall use the `bubbles/textinput` component with a `> ` prompt.
- **FR-3 (Instant Commands)**: Single-key commands (w/a/s/d/e/i/u/h/q) shall fire immediately when the text input is empty.
- **FR-4 (Typed Commands)**: Multi-word commands (`go north`, `take key`) shall submit on Enter.
- **FR-5 (Win Screen)**: Upon winning, the game shall display a styled victory message with final turns, final score, and "Press any key to exit."
- **FR-6 (Quit)**: Ctrl+C and Esc shall cleanly exit the game.
- **FR-7 (Styling)**: HUD, map, messages, and win screen shall use Lip Gloss styles (bold, color, borders).

### Non-Functional Requirements (NFR)

- **NFR-1 (Zero game logic changes)**: No modifications to `game/`, `renderer/`, `world/`, or `generator/` packages.
- **NFR-2 (Alt Screen)**: The game shall run in alternate screen mode (full-screen terminal).
- **NFR-3 (Test Stability)**: All existing tests shall continue to pass unchanged.
- **NFR-4 (Fewer LOC)**: The new `main.go` should be ~80 lines, down from ~120.

---

## 4. Visual Mockup

```
╭─────────────────────────────────────────────╮
│  Location: Dank Cellar                      │
│  Turns: 5    Score: 25                      │
│  ──────────────────────────────────────────  │
│                                             │
│  ╭─────────────────╮                        │
│  │    [ ]   [ ]    │                        │
│  │ [@]   [ ]       │                        │
│  │    [ ]          │                        │
│  ╰─────────────────╯                        │
│                                             │
│  You are in a small, damp room.             │
│  Exits:                                     │
│  - north                                    │
│                                             │
│  You took the key.                          │
│                                             │
│  > _                                        │
╰─────────────────────────────────────────────╯
```

---

## 5. Dependencies

| Action | Package |
|--------|---------|
| **Add** | `github.com/charmbracelet/bubbletea` |
| **Add** | `github.com/charmbracelet/bubbles` |
| **Add** | `github.com/charmbracelet/lipgloss` |
| **Remove** | `github.com/gdamore/tcell/v2` |

---

## 6. Implementation Plan

### Step 1: Install dependencies

```sh
go get github.com/charmbracelet/bubbletea github.com/charmbracelet/bubbles github.com/charmbracelet/lipgloss
```

### Step 2: Rewrite main.go

- Define `model` struct holding `*game.Game`, `textinput.Model`, `message string`, `won bool`.
- Implement `Init()` — create game and text input, return `nil` cmd.
- Implement `Update()` — dispatch on `tea.KeyMsg`: quit keys, enter (submit), instant commands when input empty, otherwise delegate to text input.
- Implement `View()` — build `MapView`, call `RenderHUD`/`RenderMap`/`Look()`, join with Lip Gloss styling.
- `main()` — `tea.NewProgram(initialModel(), tea.WithAltScreen()).Run()`.

### Step 3: Add Lip Gloss styles

```go
hudStyle  = lipgloss.NewStyle().Bold(true)
mapStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1)
msgStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
winStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10"))
```

### Step 4: Clean up

```sh
go mod tidy   # removes tcell, adds charm deps
```

---

## 7. Files Changed

| File | Change |
|------|--------|
| `main.go` | Full rewrite (~80 lines replacing ~120 lines) |
| `go.mod` | +bubbletea, +bubbles, +lipgloss, -tcell |
| `go.sum` | Auto-updated |

**Zero changes to:** `game/`, `renderer/`, `world/`, `generator/`

---

## 8. Testing Plan

- `go build ./...` — compiles successfully
- `go test ./...` — all existing tests pass (game, generator, renderer)
- Manual play test:
  - Game launches in alt screen
  - HUD shows Location, Turns, Score with bold styling
  - ASCII map renders inside a bordered box
  - Instant commands (w/a/s/d) work when input is empty
  - Typing `score` + Enter works without key interception
  - `take sword`, `go north` work normally
  - Win condition shows styled victory screen with final score
  - Ctrl+C and Esc quit cleanly

---

## 9. Rollback

```sh
git checkout main.go && go mod tidy
```
