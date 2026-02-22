# Box-drawing map with fog of war

## Status
Accepted

## Context
The original map used flat `[ ]`/`[@]` symbols with no walls, corridors, or visited state. Players couldn't tell where passages were or which rooms they'd been to.

## Decision
Buffer-based renderer draws 5x3 room boxes connected by corridors. Interior markers show `@` (player), `.` (visited), or space (unvisited). Locked doors use double-line box-drawing characters.

## Consequences
Map is now the visual star of the game. Room topology is clear at a glance. Fog of war adds exploration incentive.
