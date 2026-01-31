# Feature Design: HUD and Score System

## 1. Project Goal

To enhance the player experience by providing a persistent on-screen Heads-Up Display (HUD) that shows the player's current location and a running score based on the number of turns taken.

---

## 2. Requirements

### Functional Requirements (FR)

*   **FR-1 (HUD Display)**: The game shall display a persistent status bar (HUD) at the top of the terminal window during active gameplay.
*   **FR-1.1 (Location Display)**: The HUD shall clearly show the player's current room name, formatted as `Location: [Current Room Name]`.
*   **FR-1.2 (Score Display)**: The HUD shall display the total number of "turns taken" by the player, formatted as `Turns: [Number]`.
*   **FR-1.3 (HUD Separator)**: A visible separator line (e.g., a row of dashes) shall be displayed immediately below the HUD to clearly distinguish it from the main game output.
*   **FR-2 (Turn Definition)**: A "turn" shall increment after any successful command that alters the game state (e.g., movement, take, drop, unlock). Commands that only query state (e.g., `look`, `help`, `inventory`) or fail shall not increment a turn.
*   **FR-3 (Win Screen Score)**: Upon winning the game, the final total number of turns taken shall be prominently displayed on the win screen.

### Non-Functional Requirements (NFR)

*   **NFR-1 (Readability)**: The HUD text shall be clear and easily distinguishable from other game elements.
*   **NFR-2 (Performance)**: The rendering of the HUD shall not introduce any noticeable performance degradation.
*   **NFR-3 (Consistency)**: The HUD's position and content shall remain consistent, updating dynamically only when relevant game state changes.
*   **NFR-4 (Integration)**: The feature will be implemented using the existing `renderer` and `game` packages, with `main.go` orchestrating the display.
*   **NFR-5 (No Regressions)**: This feature shall not alter any existing game mechanics or cause test failures.

---

## 3. Visual Mockup

This mockup shows how the game screen will look with the new HUD integrated at the top.

```
Location: Dank Cellar
Turns: 5
-----------------------------------------------

Instant Commands: w,a,s,d (move), e (take)...
Typed Commands: go [dir], take [item]...

   [ ]   [ ]
[@]   [ ]
   [ ]

You are in a small, damp room. A faint...
Exits:
- north

You took the key.
> _
```

---

## 4. Software Design & Implementation Plan

### `game` Package Changes

1.  **`game/structs.go`**:
    *   The `Game` struct will be modified to include a field for tracking turns.
    ```go
    type Game struct {
        Player   *world.Player
        AllRooms map[string]*world.Room
        IsWon    bool
        Turns    int // New field for score
    }
    ```

2.  **`game/game.go`**:
    *   The `HandleCommand` function will be updated to increment `g.Turns`. A simple way to achieve this without complex return values is to check the verb. If the command is a state-changing action and the returned message does not indicate an error, the turn counter is incremented.
    *   The `NewGame` function will initialize `Turns` to 0.

### `renderer` Package Changes

1.  **`renderer/renderer.go`**:
    *   The `MapView` struct will be updated to include the data needed for the HUD.
    ```go
    type MapView struct {
        AllRooms          map[string]*world.Room
        PlayerLocation    *world.Room
        CurrentLocationName string // New field
        TurnsTaken        int    // New field
    }
    ```
    *   A new function, `RenderHUD(view MapView) string`, will be created. It will take the `MapView` and return a formatted string containing the Location, Turns, and the separator line, ready for printing.

### `main.go` (Orchestrator) Changes

1.  **Main Game Loop**:
    *   The main `for` loop will be updated to orchestrate the new rendering flow.
    *   On each iteration, after a command is handled, it will populate the new `CurrentLocationName` and `TurnsTaken` fields in the `MapView` struct it creates.
    *   It will call `hudString := renderer.RenderHUD(mapView)`.
    *   It will call `mapString := renderer.RenderMap(mapView)`.
    *   It will then draw the `hudString` at the top of the screen (starting at `y=0`), and adjust the starting `y` coordinate for drawing the `mapString` and other content to be below the HUD.

2.  **Win Condition**:
    *   The code block that handles the win state (`shouldExit == true`) will be modified to include the final `g.Turns` value in the victory message string.

### Testing Plan

*   A new test, `TestRenderHUD`, will be added to `renderer/renderer_test.go` to ensure the HUD string is formatted correctly.
*   Existing tests in `game/game_test.go` (like `TestMovement` and `TestTakeAndDrop`) will be enhanced to assert that `g.Turns` increments only after successful, state-changing actions.
