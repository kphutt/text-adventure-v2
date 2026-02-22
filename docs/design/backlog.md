# Design backlog

Ideas not yet tied to an initiative. See ROADMAP.md for the full prioritized list with sizing.

## Gameplay

### Simple item usage / puzzles

Allow items to be `use`d beyond unlocking doors (e.g., `use torch` in a dark room). Requires a generic item-effect system — could pair naturally with the command pattern refactor.

### Simple, static NPCs

Add characters to rooms who respond to `talk`. Delivers clues and story without needing a full dialogue tree. Prerequisite for the branching dialogue system later.

### A simple combat system

Basic enemies and an `attack` command. Needs HP tracking, turn-based resolution, and death/respawn handling. Medium scope — touches game logic, generator (enemy placement), and renderer (health display).

### Saving and loading

`save` and `load` commands. Requires serializing full game state (player, rooms, inventory, visited, locks). Medium-to-high scope due to pointer-heavy world graph.

### Multiple dungeon themes

Generator produces themed worlds ("Crypt," "Cavern," "Forest") with unique descriptions and items. Large scope — needs theme data, description templates, and generator config.

### A branching dialogue system

Choice-based NPC conversations. Depends on static NPCs shipping first. Large-to-extra-large scope.

## Architecture

### Refactor to the command pattern

Replace the `HandleCommand` switch statement with command objects. Unblocks clean addition of new verbs (`use`, `talk`, `attack`). Small scope, high leverage.

### Implement the decorator pattern

Dynamically attach temporary properties to game objects (darkness on a room, glowing on a sword). Medium scope.

### Implement the state pattern

Manage game flow states (`MainMenu`, `Playing`, `GameOver`). Required before adding menus or multiple game modes. Large scope.

## Fun / quick wins

### Cheat code system

Hidden command (`xyzzy`) that grants powers — reveal map, give all items. Extra small scope, high fun factor.

### Magic, joke, and puzzle items

Creative items with surprising effects (`use bouncy ball`). Extra small scope, great for personalization.

### Secret rooms and passages

`search` command reveals hidden doors. Low scope, strong exploration payoff.

### Inspectable scenery and clues

`inspect <object>` for detailed descriptions and hidden clues. Medium scope, adds world depth.

## Done

### ~~Scrolling message log~~ ✅

Last N messages shown on screen. Implemented with Bubble Tea.

### ~~New map visuals~~ ✅

Box-drawing room layout with corridors, locked door indicators, and fog of war (visited `.`, player `@`, unvisited empty).
