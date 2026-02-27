# Go Text Adventure Game

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A simple, terminal-based text adventure game written in Go, now featuring procedurally generated and unique worlds for endless replayability!

## Key Features

*   **Procedurally Generated Worlds**: Every new game offers a unique map layout, ensuring a fresh experience with each playthrough.
*   **Dynamic Puzzles**: Key items, locked doors, and treasure rooms are strategically placed to guarantee a solvable and engaging challenge.
*   **Terminal-Based UI**: A classic text adventure experience directly in your terminal.
*   **Clear Map Visualization**: An in-game map helps you navigate the generated world.

## How to Play

1.  **Run the game:**
    ```bash
    go run .
    ```
2.  **Instant Commands (No Enter key needed):**
    *   `w`, `a`, `s`, `d`: Move north, west, south, and east.
    *   `e`: Take the first available item in the room.
    *   `i`: View your inventory.
    *   `u`: Attempt to unlock a door.
    *   `l`: Look around the current room.
    *   `q`: Quit the game.
3.  **Typed Commands (Enter key needed):**
    *   `go [direction]`: Move in a specific direction (e.g., `go north`).
    *   `take [item name]`: Pick up a specific item from the room.
    *   `drop [item name]`: Drop an item from your inventory.
    *   `help`: Display the list of available commands.
    *   `quit`: Quit the game.

## The Goal

The goal of the game is to find the key, unlock the door, and reach the treasure room!

Good luck, adventurer!

## The Engineering Problem

Procedural generation often produces unsolvable states — a locked door with no key, a room with no exits, a treasure behind a wall you can't reach. The generator must produce worlds that are always winnable without constraining them into boring layouts.

## Generation Guarantees

Every generated world is validated before the game starts. The generator runs BFS-based checks (`generator/validator.go`, 12 tests) that enforce:

- **Solvability** — the key is reachable from the start room
- **Key-before-lock ordering** — the key is always placed before the locked door on the critical path (`generator/puzzler.go`)
- **Treasure is locked** — the treasure room is unreachable without first obtaining the key
- **Connected traversal** — all rooms on the critical path are reachable via BFS

If validation fails, the world is regenerated. See [DESIGN.md](DESIGN.md) for the full constraint model.

## Design Details and Architecture

For an in-depth look at the requirements, game design decisions, and the software architecture of the project, please refer to the [DESIGN.md](DESIGN.md) document.

## Project Roadmap

To see the list of planned features, architectural improvements, and other brainstormed ideas for the future of this project, please refer to the [ROADMAP.md](ROADMAP.md) document.

## Documentation Structure

This project utilizes a structured documentation approach to keep design details and development logs organized.

*   **[`docs/features/`](docs/features/)**: Contains detailed design documents for specific features (e.g., HUD and Score System). These documents include requirements, mockups, and implementation plans.
*   **[`docs/dev_log/`](docs/dev_log/)**: Contains developer logs that summarize key design decisions, brainstorming sessions, and the evolution of complex features (e.g., the map rendering visual design).

## Acknowledgments

This project was developed with the assistance of Google's Gemini CLI, which served as an AI pair-programming partner for design, implementation, and refactoring.