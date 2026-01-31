# Dev Log: The Evolution of the ASCII Map

This document summarizes the collaborative design process that led to the final "Tiled Grid" visual style for the game's map renderer.

## Initial Problem

The original map renderer was functional but visually ambiguous. It displayed rooms as disconnected boxes (`[ ]`), making it impossible for the player to know where valid exits were without trial and error.

**Initial State:**
```
[@]   [ ]
   [ ]
```

## Design Journey & Key Iterations

The goal was to create a sharp, professional, and unambiguous map. The design evolved through several key insights and feedback loops.

### Iteration 1: "Interstitial Glyphs"

The first idea was to keep rooms as separate boxes but draw explicit connector symbols (`|`, `-`) in the space *between* them.

**Problem**: This looked disjointed and "weird." The connecting lines appeared to float in space and didn't feel anchored to the rooms.

### Iteration 2: The "Shared Wall" Insight

The key breakthrough came from the feedback that rooms should be **"touching and share a single wall."** This led to the "Tiled Grid" concept, where the map is a continuous grid and its lines form the shared walls.

This concept was then refined through a final series of feedback loops:
1.  **Problem**: Initial mockups with larger room boxes felt cramped when displaying text.
2.  **Solution**: It was decided to limit the content of each room to a **single character** (e.g., `@` or a letter). This simplified the visuals immensely.
3.  **Problem**: My manual ASCII drawings were creating alignment issues and implying "phantom rooms" in empty space.
4.  **Solution**: The final rule was established: the map must be a full grid where any cell not containing a generated room is filled with a placeholder pattern. This ensures perfect alignment.

## The Final Blueprint

This process resulted in a final, approved visual language defined by these rules:

*   **Style**: A continuous "Tiled Grid" with shared walls.
*   **Room Size**: Each room's interior is a `3x1` character space (`| A |`).
*   **Content**: Each room displays only a **single character**.
*   **Connections**: Openings are rendered as gaps in the shared walls, while locked doors are rendered as special symbols (`║` or `═`) in place of a wall segment.
*   **Empty Space**: All cells in the grid that do not contain a room are filled with a placeholder pattern to guarantee alignment.

**Final Approved Example:**
```
+---+---+
| @   B |
+---+   +
| C ║ D |
+---+---+
```

This iterative process, driven by precise user feedback, was crucial in arriving at a professional and highly functional final design.
