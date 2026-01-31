# Project Roadmap

This document outlines the future direction for the Go Text Adventure Game, including potential new features and architectural improvements. Items are generally ordered by a combination of impact and effort.

---

## New Game Features

This list focuses on adding new gameplay mechanics and content to make the game more interactive and engaging.

1.  **Simple Item Usage / Puzzles**
    *   **Size**: Extra Small
    *   **Goal**: Allow items to be `use`d for more than just unlocking doors (e.g., `use torch` in a dark room). This adds a new layer of problem-solving.

2.  **Simple, Static NPCs**
    *   **Size**: Small
    *   **Goal**: Add characters to the world who can be `talk`ed to. This makes the world feel more alive and allows for delivering clues and story.

3.  **Saving and Loading**
    *   **Size**: Medium to High
    *   **Goal**: Implement `save` and `load` commands. This is an essential feature that enables longer, more complex dungeons and improves player convenience.

4.  **A Simple Combat System**
    *   **Size**: Medium
    *   **Goal**: Introduce basic enemies and an `attack` command. This adds a new dimension of challenge, risk, and reward to exploration.

5.  **A Scrolling Message Log**
    *   **Size**: Small
    *   **Goal**: Display a short history of the last few messages on screen to prevent players from missing important information.

6.  **Multiple Dungeon "Themes"**
    *   **Size**: Large
    *   **Goal**: Allow the generator to create different kinds of worlds (e.g., "Crypt," "Cavern," "Forest") with unique descriptions and items, massively increasing replayability.

7.  **A Branching Dialogue System**
    *   **Size**: Large to Extra Large
    *   **Goal**: Transform the game into a richer RPG experience by allowing for interactive, choice-based conversations with NPCs.

---

## Architectural Improvements

This list focuses on internal refactoring projects to make the codebase more modular, maintainable, and professional.

1.  **Refactor to the Command Pattern**
    *   **Size**: Small
    *   **Goal**: Replace the large `switch` statement in the game engine with dedicated "command objects." This will make adding new game verbs (like `use`, `talk`, `attack`) much cleaner.

2.  **Implement the New Map Visuals**
    *   **Size**: Medium
    *   **Goal**: Upgrade the `renderer` to use the advanced "Tiled Grid" ASCII art style we designed. This is a major visual and user experience enhancement.

3.  **Implement the Decorator Pattern**
    *   **Size**: Medium
    *   **Goal**: Introduce a system for dynamically adding temporary properties to game objects (e.g., a "darkness" effect on a room, a "glowing" effect on a sword).

4.  **Implement the State Pattern**
    *   **Size**: Large
    *   **Goal**: A major architectural refactoring to manage the overall game flow (e.g., `MainMenu`, `Playing`, `GameOver`). This is crucial for expanding the game beyond a single mode.

---

## Brainstorming: Fun & High-Impact Feature Ideas
*(Ranked by "Fun Factor" for getting new coders excited)*

1.  **A Cheat Code System**
    *   **Effort**: Extra Small
    *   **Fun Factor**: ★★★★★
    *   **The Pitch**: Kids love secrets and feeling powerful. Adding a hidden command like `xyzzy` that grants a special power (like giving all items, or revealing the map) is a classic, fun game development tradition and a very satisfying feature to code.

2.  **Magic, Joke, and Puzzle Items**
    *   **Effort**: Extra Small
    *   **Fun Factor**: ★★★★★
    *   **The Pitch**: This offers maximum creativity for minimal effort. Brainstorming and implementing silly or surprising items (`use 'bouncy ball'`) provides immediate, often humorous, feedback. It's the perfect way for a new coder to add their own personal touch to the game.

3.  **Secret Rooms and Passages**
    *   **Effort**: Low
    *   **Fun Factor**: ★★★★☆
    *   **The Pitch**: The thrill of discovery is a huge motivator. Adding a `search` command that can reveal hidden doors makes players feel clever and encourages deep exploration. It feels like building a real secret into the world.

4.  **Inspectable Scenery and Clues**
    *   **Effort**: Medium
    *   **Fun Factor**: ★★★☆☆
    *   **The Pitch**: This adds a new layer of depth and mystery to the world. Players can `inspect <object>` to get more detailed descriptions or find hidden clues. It's a great feature for creative writing and making the world feel more interactive.