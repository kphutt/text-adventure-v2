# Check both sides of exit for locks

## Status
Accepted

## Context
The generator only sets `Locked: true` on one direction of an exit (doorRoom -> treasureRoom). The renderer draws corridors from only the east/south side to avoid double-drawing. If the lock happened to be on the north/west side, the corridor rendered as unlocked.

## Decision
Corridor drawing checks both the drawn exit and the reverse exit for locked status using `isExitLocked()`.

## Consequences
Locked doors display correctly regardless of which direction the generator locked. No changes needed in the generator.
