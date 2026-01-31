# Go Text Adventure Game

[![Go Version](https://img.shields.io/badge/go-1.25-blue.svg)](https://golang.org/)
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

## Design Details and Architecture

For an in-depth look at the requirements, game design decisions, and the software architecture behind the procedural map generation, please refer to the [DESIGN.md](DESIGN.md) document.

## Acknowledgments

This project was developed with the assistance of Google's Gemini CLI, which served as an AI pair-programming partner for design, implementation, and refactoring.