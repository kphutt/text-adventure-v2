# CLAUDE.md

## Project
Terminal-based text adventure game in Go with procedurally generated worlds.

## Build & Test
```sh
go build ./...       # compile
go test ./...        # run all tests
go run .             # play the game (interactive TUI, needs a real terminal)
```

## Architecture
- `game/` — core game logic (HandleCommand, Move, Take, Drop, Unlock, Score). Pure logic, no UI.
- `renderer/` — ASCII map and HUD rendering. Pure functions: `RenderMap(MapView)`, `RenderHUD(MapView)` return strings.
- `generator/` — procedural world generation (room layout, key/door placement, BFS validation).
- `world/` — shared data types (Room, Player, Item, Exit).
- `main.go` — Bubble Tea TUI (Model/Update/View). Only file that imports UI libraries.

## Key patterns
- Game logic and rendering are fully decoupled from the UI framework.
- `MapView` struct is the view model bridging game state to renderer.
- Instant commands (w/a/s/d/e/i/u/h/q) fire when text input is empty; otherwise keys go to the text input.
- Items in room/inventory are highlighted with Lip Gloss styling in the View layer.

## Dependencies
- `github.com/charmbracelet/bubbletea` — TUI framework
- `github.com/charmbracelet/bubbles` — text input component
- `github.com/charmbracelet/lipgloss` — styling/layout

## Testing notes
- Tests live in `game/*_test.go`, `generator/generator_test.go`, `renderer/renderer_test.go`.
- Test helpers in `game/test_helpers.go` build specific world layouts (simple, items, locks, win).
- All test helpers must initialize `VisitedRooms` with the start room.
- Game logic tests don't depend on any UI library.
