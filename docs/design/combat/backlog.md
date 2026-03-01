# Combat Initiative: Development Backlog

## Context

The combat system is the most ambitious feature on the roadmap (item #4). It's not one project — it's an initiative spanning multiple projects, each building on the last. This document breaks the initiative into **bite-sized phases** that can each go through their own requirements → design → implement cycle without overwhelming context.

The brainstorm and design philosophy live in companion docs:
- `docs/design/combat/brainstorm.md` — Meier philosophy, deep dives, wild ideas
- `docs/design/combat/principles.md` — The 10 design principles (the laws)

## Decomposition Philosophy

**Why this order matters.** Sid Meier's Principle 8 says "Prototype the fun first." That means every phase should get us closer to PLAYABLE, not just "more infrastructure built." But you can't see anything without a renderer, and you can't fight without physics. So Phase 1 is unavoidably foundational. After that, each phase produces something you can play, see, or feel.

**What "bite-sized" means.** Each phase should be:
- Self-contained enough that a fresh session can pick it up without holding the entire codebase in context
- Has clear boundaries — which packages to create or modify, and what NOT to touch
- Has a "done" condition you can verify by running tests or playing the game
- Can go through requirements → design → implement in 1-2 focused sessions

**Decision gates.** After certain phases, we stop and ask: "Should we keep going, or pivot?" These gates are more important than the phases themselves. They're where bad projects die early and good projects gain conviction.

---

# THE BITE-SIZED PHASES

## Phase 1: The Pixel Canvas

*Build the paintbrush before painting.*

**The goal.** Create `pixelbuf/`, a standalone Go package that converts a 2D grid of RGBA colors into beautiful half-block terminal output. This is the rendering engine for everything visual in combat. It knows nothing about games, combat, or dungeons — it just paints pixels.

**Why this is first.** Every visual element in combat — the player, the enemy, the arena, particles, health bars — ultimately becomes pixels in a buffer that gets rendered to the terminal. Without this, we can't see anything. With it, we can see everything.

**Why this is the right size.** The package is perfectly isolated. Zero game dependencies. A developer only needs to think about one thing: "2D color grid → ANSI string." The API surface is small: Buffer (create, set, get, clear), Color (RGBA, transparency), Sprite (named pixel arrays), Draw (blit, fill), Render (buffer → half-block string). That's ~5 files, ~300-400 lines, plus tests.

**The exciting part.** When this is done, you can write a 20-line program that draws a colored rectangle on screen. It sounds trivial, but it's the moment where "text adventure game" becomes "thing that can render pixel art in a terminal." The same rendering technique that runs Game Boy Color emulators at 60fps in a terminal. The technical foundation for the jaw-drop.

**What's in scope:**
- `pixelbuf/buffer.go` — Buffer struct (width, height, 2D color array), New, Set, Get, Clear
- `pixelbuf/color.go` — Color type (RGBA), basic palette constants, alpha blending
- `pixelbuf/sprite.go` — Sprite type (named 2D color array with transparency)
- `pixelbuf/draw.go` — Blit (draw sprite onto buffer), FillRect
- `pixelbuf/render.go` — Buffer → half-block ANSI string (the core algorithm: pair rows, use `▀`/`▄`/`█` with fg/bg colors)
- Full test coverage — exact ANSI output verification for known small buffers

**What's NOT in scope:** Game logic, physics, input handling, Bubble Tea integration.

**Done when:** `go test ./pixelbuf/...` passes. A test can create a 4x4 buffer, set some colors, call Render(), and get back the exact expected ANSI escape string.

---

## Phase 2: Rectangles That Fight

*The Principle 8 milestone: "I can't stop playtesting the rectangle version."*

**The goal.** Build the combat engine (`combat/engine/`) and a standalone playable prototype that renders colored rectangles via pixelbuf. A player rectangle that runs, jumps, and swings a nail at an enemy rectangle on a platform arena. No sprites. No particles. No polish. Just rectangles — and the question: **is this fun?**

**Why this is THE critical phase.** This is where the entire combat initiative lives or dies. If moving a rectangle around a platform arena, jumping over attacks, and slashing another rectangle feels *good* — if the weight of the jump is right, the attack timing is snappy, the collision feels fair — then everything after this is polish. If it doesn't feel good, no pixel art or screen shake will save it. We find out HERE, before investing in anything else.

**The two halves (one phase, natural seam):**

*First half — Pure Engine.* `combat/engine/` is pure logic with zero UI dependencies. Fully deterministic. Fully testable with scripted inputs. This is where the game design lives:
- Physics: gravity, velocity, AABB collision resolution (move X, resolve; move Y, resolve)
- Player state machine: idle, run, jump, fall, attack, hurt. Variable-height jump, coyote time (4 frames), jump buffering (4 frames)
- Enemy: starts simple — stands on a platform, takes hits, blinks when hurt, dies when HP hits 0. Real AI comes in Phase 4.
- Arena: single screen, no scrolling, ~5 platforms at varying heights, boundary walls
- Health: player HP, enemy HP, damage on hit, invincibility frames, death check
- Engine.Tick(dt, actions) → pure function, advances simulation one frame

*Second half — Playable Prototype.* Wire the engine to pixelbuf and Bubble Tea in a standalone program (`cmd/combat-proto/main.go`). Colored rectangles on a black background. 30fps game loop via `tea.Tick`. Manual key state tracking (held/pressed detection — Bubble Tea delivers KeyMsg events, we maintain a map with timeout-based release). Press 'R' to restart instantly — fast iteration.

**The playtesting question.** When this is done, sit down and play it for 5 minutes. If you want to keep playing, Phase 2 succeeded. If you want to stop, something is wrong with the fundamentals — and we fix it HERE, before building on top of it.

**What's in scope:**
- `combat/engine/` — entities, physics, player, enemy (basic), arena, health, hitbox, engine tick loop
- `cmd/combat-proto/main.go` — standalone Bubble Tea program, renders engine state via pixelbuf
- `combat/input.go` — key state tracking (reusable later in Phase 3)
- Tests for every engine subsystem (physics, collision, player state transitions, damage)

**What's NOT in scope:** Dungeon integration, mode switching, sprites, particles, screen shake, real enemy AI. The enemy just stands there and takes hits. That's enough to test attack timing and feel.

**Done when:** `go run cmd/combat-proto/` launches a full-screen terminal with a player rectangle that can run, jump, and slash an enemy rectangle. The enemy takes damage and eventually dies. You can restart and play again.

### >> GATE: "Is the fun there?"

*Stop here and play the rectangle prototype. Multiple sessions. Show it to someone else. Ask: does the movement feel right? Is the attack timing satisfying? Is jumping fun? If YES → proceed. If NO → tune the physics and player controller until it IS fun. Do NOT proceed to Phase 3 until rectangles are fun. This gate protects the entire initiative from building on a weak foundation.*

---

## Phase 3: The Mode Switch

*The "holy shit" moment — even with rectangles.*

**The goal.** Wire combat into the dungeon game. Player walks into a room with an enemy → screen transforms → full-screen combat → fight → screen returns to dungeon. The enemy is gone, loot is on the ground. This is Principle 5: "The Contrast IS the Spectacle." Even with rectangles, the contrast between a quiet box-drawing map and a full-screen colored arena should make someone say "wait, WHAT?"

**Why this is a separate phase from Phase 2.** Phase 2 is about getting the combat FEEL right in isolation. Phase 3 is about the INTEGRATION — how two different games become one. These are fundamentally different problems. Phase 2 is game design (physics, timing, controls). Phase 3 is architecture (mode switching, state management, data flow).

**The architectural change.** The root Bubble Tea model in `main.go` becomes a thin dispatcher. Currently it's a single `model` struct. After Phase 3:

```
main.go (root model)
  mode: "explore" | "combat"
  ├── exploreModel (current game, barely changed)
  └── combat.Model (full-screen platformer from Phase 2)
```

The root model's `Update()` dispatches to whichever sub-model is active. The root model's `View()` calls whichever sub-model's `View()`. Mode transitions happen via `tea.Cmd` messages.

**The data flow.** This is where HP becomes the universal currency (Principle 2):
1. `game.Move()` enters a room with an enemy → sets `Game.PendingCombat`
2. Explore model sees PendingCombat → emits `StartCombatMsg`
3. Root model receives msg → builds `combat.Model` with player HP, enemy stats → switches to combat mode
4. Player fights. HP changes. Enemy takes damage.
5. Combat ends (win or lose) → `CombatResultMsg` with updated HP, enemy defeated flag, loot
6. Root model applies results to game state → switches back to explore mode
7. If enemy was defeated, it's removed from the room. Loot items appear.

**Minimal world/game changes:**
- `world/structs.go`: Add `HP`, `MaxHP` to Player. Add `EnemyType` struct. Add `Enemy *EnemyType` field to Room.
- `game/structs.go`: Add `PendingCombat` field to Game.
- `game/game.go`: In `Move()`, after entering a room, check for undefeated enemy → set PendingCombat. In `NewGame()`, initialize Player HP.
- Generator: hardcode ONE enemy in a non-start, non-treasure room. Sophisticated placement comes in Phase 5.

**The transition for now:** Simple — instant cut to combat, instant cut back. No wipe, no fade, no animation. Save transition polish for Phase 4. The contrast between modes is dramatic enough without effects.

**What's in scope:**
- `main.go` refactor: root model dispatcher, explore sub-model extraction, combat mode delegation
- `world/structs.go` + `game/structs.go`: HP, enemy, PendingCombat fields
- `game/game.go`: Move() enemy check, NewGame() HP init
- `generator/`: hardcode one enemy placement (minimal, not the full pipeline)
- Message types: `StartCombatMsg`, `CombatResultMsg`
- All existing tests must still pass

**What's NOT in scope:** Transition animations, sprites, juice, strategic enemy placement, loot drops, combat items.

**Done when:** `go run .` — play the dungeon normally, walk into the enemy room, screen switches to full-screen combat (rectangles), fight and win, screen returns to dungeon, enemy is gone from the room. HP carries between modes.

### >> GATE: "Is the mode switch magical?"

*Show the game to someone who hasn't seen it. Don't tell them about combat. Watch their face when they walk into the enemy room. If there's a reaction — surprise, delight, "whoa" — the contrast works and we proceed. If they shrug, we need to make the gap bigger (simpler dungeon? bigger combat canvas? faster transition?).*

---

## Phase 4: Game Feel

*Where "working prototype" becomes "I can't believe this is a terminal."*

**The goal.** Replace rectangles with pixel art. Add hit-stop, screen shake, knockback, particles. Give the enemy real AI. Add health bars, invincibility flashing, death animations, transition effects. This phase is 100% polish — no new systems, no new game mechanics. Just making everything that already works FEEL incredible.

**Why this is one phase, not three.** You might think "sprites" and "screen shake" and "enemy AI" are separate projects. They're not. They're all the same thing: game feel. They all answer the same question: "does hitting this enemy feel GOOD?" You can't evaluate screen shake without sprites to shake. You can't evaluate enemy AI without hit feedback to show when your attacks land. These elements are synergistic — they multiply each other's impact. Ship them together.

**The juice ingredients (priority order):**

1. **Hit-stop** (highest impact, simplest to implement). When the player's attack hits the enemy, freeze the entire simulation for 3 frames (100ms). Both characters stay in their pose. This is the single most impactful piece of game feel in Hollow Knight. It makes hits feel HEAVY. Without it, attacks feel like swinging at air. With it, every hit feels like it matters.

2. **Screen shake.** On hit, offset the entire render by ±2 pixels for 4-5 frames. Oscillating. Bigger shake for enemy hits on the player. Combined with hit-stop, attacks feel like they have WEIGHT.

3. **Knockback.** Struck entities receive a velocity impulse away from the attacker. The enemy slides back when hit. The player is launched away when hurt. Movement = consequence.

4. **Sprites.** Replace colored rectangles with pixel art. Player: ~12x16 pixel warrior/knight silhouette. Enemy (Husk Guard): ~12x16 menacing figure. Hardcoded as Go arrays — no image loading needed. Flip horizontally based on facing direction.

5. **Hit flash.** Struck entity renders all-white for 1-2 frames. Visual confirmation of a hit landing.

6. **Invincibility frames.** After being hurt, the entity blinks (alternating visible/invisible) for ~1 second. Player gets generous i-frames (30 frames). Prevents stun-lock death.

7. **Particles.** Small colored pixels that spawn at the hit point, fly outward with random velocity, and fade over 10-15 frames. Sparks on attack hit. Dust puff on landing. Shower of particles on enemy death.

8. **Enemy AI.** The Husk Guard gets a real state machine: Patrol (walk back and forth) → Detect (turn to face player when in range) → Windup (telegraph attack, 0.5s) → Lunge (fast dash with hitbox, 0.2s) → Recover (stand still, vulnerable, 0.67s) → repeat. Simple but readable. The player learns the pattern.

9. **Health bars.** Player HP (hearts) and enemy HP bar rendered at screen edges via pixelbuf. Shows the stakes at a glance.

10. **Transition effects.** Replace the instant cut with: entering combat = brief flash-to-black then arena reveal. Exiting combat = enemy death particles + brief pause + flash back to dungeon. Not the full cinematic sequence from the brainstorm — save that for a later polish pass. Just enough to make the mode switch feel intentional rather than abrupt.

**What's in scope:** All of the above. All in `combat/engine/` (game logic) and `combat/` (rendering). No changes to dungeon game, world structs, or generator.

**What's NOT in scope:** New game mechanics, items, loot, score changes, additional enemy types, sound.

**Done when:** `go run .` — walk into the enemy room, screen transitions smoothly, pixel art player fights a pixel art Husk Guard with telegraphed attacks, hits feel meaty (stop + shake + flash + knockback + sparks), enemy dies dramatically, screen returns to dungeon. Someone watching over your shoulder says "wait, is this a terminal?"

---

## Phase 5: The Economy

*Where Meier's philosophy comes alive. HP as universal currency. Items that bridge both worlds. Every room is a decision.*

**The goal.** Close the dungeon-combat feedback loop. HP persists and matters. Health potions exist and force decisions. Defeated enemies drop loot. Combat items (sword, shield) affect both modes. The generator places 2-3 enemies strategically. Score incorporates combat performance. Inventory cap forces tradeoffs.

**Why this is last among the bite-sized phases.** You can't tune the economy until combat feels good (Phase 4). You can't evaluate HP persistence until the mode switch works (Phase 3). You can't test item decisions until you can fight enemies and take damage. This phase needs EVERYTHING before it to exist. But once it does, this is where the game transcends "tech demo" and becomes something you'd actually want to play through multiple times.

**The systems this phase introduces:**

*HP Persistence.* Already wired in Phase 3, but now it MATTERS. With multiple enemies, every hit costs something beyond the current fight. A player who finishes Fight #1 with 3/5 HP faces Fight #2 already weakened. The brainstorm's HP budget: 5 hearts starting, weak enemies deal 1 heart, strong enemies deal 2.

*Health Potions.* A new item type the player can find in rooms. Using one restores 2 hearts. The decision: drink now (safety) or save for later (what if there's a harder fight ahead)? Can also be used mid-combat. Found in ~30% of non-combat rooms.

*Combat Items.* Sword: doubles attack damage (fights end faster, fewer chances to take hits). Shield: blocks 1 damage per hit (survival, but fights take longer). These are inventory items — carrying one means you can't carry something else. With a 3-4 item inventory cap and the key taking one slot, every item is a decision.

*Loot Drops.* Defeated enemies have a chance to drop items: health (60%), random item (30%), nothing (10%). This makes fighting WORTH IT — combat is both the thing that hurts you (HP loss) and the thing that heals you (health drops). Beautiful tension.

*Strategic Enemy Placement.* The generator gets a `placeEnemies()` step. 1-2 enemies on the main path (mandatory fights). 0-1 enemies on side paths (optional, guarding better loot). Never in the start room or treasure room. This creates the risk/reward topology from the brainstorm.

*Score Integration.* The existing score formula (`inventory*10 + visitedRooms*5`) expands: `+15 per enemy defeated`, `+HP remaining * 5` at game end. Different playstyles produce different scores — the explorer, the fighter, the speedrunner.

**What's in scope:**
- `world/structs.go`: Item types (weapon, consumable), enemy loot tables
- `game/game.go`: Use command for potions, inventory cap enforcement, combat item effects
- `game/game.go`: Score formula update
- `generator/enemies.go` (new): `placeEnemies()` function, enemy templates, BFS-aware placement
- `generator/generator.go`: Add placeEnemies step to pipeline
- Combat engine: item effects (damage multiplier, damage reduction), health drop spawning
- Test helpers for combat+economy scenarios

**What's NOT in scope:** Multiple enemy types (just the Husk Guard with stat variations), sound, boss fights, dungeon themes, style scoring.

**Done when:** Full playthrough where you explore the dungeon, find items, make inventory decisions, fight 2-3 enemies with HP carrying between fights, use a health potion at a critical moment, find the key, unlock the door, reach the treasure. Score reflects combat performance. You feel the "one more room" pull. `go test ./...` passes.

### >> GATE: "Is it one good game?"

*The Covert Action check. Ask: "Do the dungeon and combat feed each other, or fight each other?" Does HP persistence make dungeon rooms feel like decisions? Does combat feel like an exclamation point in the exploration story, or a detour from it? Do fights stay under 90 seconds? Does the player want to explore one more room? If YES → the foundation is solid, proceed to the epic rocks. If NO → identify which principle is being violated (probably #1 or #9) and fix it before adding more.*

---

# THE EPIC ROCKS

After Phase 5, the combat system is complete and playable. Everything beyond this is expansion — making an already-fun game richer, more varied, and more surprising. These rocks are listed in rough priority order but each is independent. They can be reordered based on what's most exciting or impactful at the time.

---

## Rock: Enemy Variety

*Different enemies. Different patterns. Eventually, different genres.*

The Husk Guard is one enemy type with one behavior pattern. This rock adds 2-3 more enemy types, each with distinct AI, different visual designs, and different combat strategies required to beat them. A fast enemy that dashes and retreats. A ranged enemy that throws projectiles. A heavy enemy that hits hard but moves slow. Each requires the player to adapt.

The long-term vision (from the brainstorm) is the **Pirates! model** — different enemy types trigger entirely different minigame genres. The Goblin Horde triggers Space Invaders. The Mimic triggers Pong. This is the "curiosity" hook that makes players WANT to find enemies. But it's a massive expansion. Start with 2-3 platformer enemy types. Graduate to genre-switching when the architecture proves it can support it.

Generator integration: enemy type selection based on room position (tougher enemies further from start), enemy type variety per run (never the same type twice in a row).

---

## Rock: Sound & Music

*The mode switch you can HEAR.*

This is NOT polish. Audio is half the spectacle (see brainstorm deep dive). The moment combat triggers, a driving chiptune track should PUNCH in, synchronized with the visual transition. Dungeon exploration: silence or faint ambient. Combat: 8-bit energy. The audio contrast hits you the same moment the visual contrast does.

Technical approach: `gopxl/beep` v2 for audio playback. An `audio/` package with a manager goroutine that receives commands via `tea.Cmd`. OGG Vorbis for music tracks, WAV for SFX. Sound palette: nail slash (metallic SHING), hit impact (meaty THUD), enemy death (satisfying CRASH + sparkle), player hurt (sharp + bass rumble), item pickup (bright chime).

Free music sources: Soundimage.org, Ozzed.net, itch.io asset packs. Or compose original tracks with FamiStudio.

---

## Rock: Boss Fights

*The climax of every run.*

The Dragon Boss (or equivalent) guards the treasure room. This is the final fight — the one everything has been building toward. Longer arena (maybe scrolling?). Multi-phase behavior (enrages at 50% HP). Bigger, more dramatic pixel art. The stakes: you've been managing HP across 3-4 fights for 10 minutes. You have the key. One fight stands between you and the treasure. Everything is on the line.

Mechanically: the boss uses all the patterns the player has learned from lesser enemies, plus unique attacks. It's a test of mastery, not a gear check. 90-120 second fights (stretching Principle 1 to its 2-minute boss limit).

The death animation should be the most dramatic thing in the game. Extended slow-mo, massive particle shower, screen flash. The player just beat a boss in a TERMINAL. Make them feel it.

---

## Rock: Advanced Combat Mechanics

*Skill ceiling without skill floor.*

Dash (horizontal burst, one air dash, resets on landing). Directional attacks (slash up/down based on held keys). Pogo bouncing (downward attack on enemy = bounce upward — the signature Hollow Knight move). Dash-canceling for advanced players.

These don't change the core loop — basic attack spam still beats basic enemies. But expert players can chain dash-attacks, pogo-bounce across the arena, and combo for style points. This is Principle 4: "The Player Is the Star." The player discovers these moves organically and feels clever.

Optional: Style Scoring (Devil May Cry style meter: D → C → B → A → S based on hit chains without taking damage). Higher style = more score/better loot. Rewards mastery without punishing beginners.

---

## Rock: Dungeon Themes

*The arena should feel like the room you entered, expanded to fill the screen.*

Different visual themes for combat arenas: Crypt (dark stone, chains, torchlight), Cavern (dripping water, mushrooms, crystal formations), Library (bookshelves as platforms, floating pages). The arena theme matches the room description the player just read, connecting the two modes narratively.

This also opens the door to dungeon themes in exploration mode — different room descriptions, item sets, and atmospheres for each run. "Crypt" vs "Cavern" vs "Forest" dungeons with matching combat arenas.

Generator integration: dungeon theme selected at generation time, determines room descriptions, item pool, enemy pool, and arena visuals.

---

## Rock: Transition Cinematics

*The full spectacle from the brainstorm.*

Replace the simple flash-to-black transition with the full cinematic sequence: Warning text ("Something stirs in the darkness...") → Freeze → Column-by-column wipe to black → Arena fade-in (background → platforms → enemy materializes) → Player appears → Health bars slide in → FIGHT. Exit: death animation → particles → score popup → wipe back to dungeon.

This is the thing people screenshot and share. The transition should feel like a movie cut, not a loading screen. Time the audio crossfade to the visual transition for maximum impact.

---

## Rock: The Last Stand

*When HP hits 0, you get one final shot.*

From the brainstorm's wild ideas. When the player's HP reaches 0 in combat, instead of instant death: time slows dramatically, the screen desaturates, and the player gets ONE final attack. If it kills the enemy, the player survives with 1 HP. If it misses, game over.

This single mechanic transforms every near-death moment from frustration into a clutch opportunity. It makes every defeat feel like "I ALMOST had it" instead of "that was unfair." Pure Meier: the player feels like the star even in defeat.

---

# DOCUMENT ORGANIZATION

```
docs/design/combat/
  brainstorm.md      # WHY  — Meier philosophy, deep dives, wild ideas
  principles.md      # LAWS — the 10 design principles
  backlog.md         # WHAT & WHEN — this document (phased development plan)
```

When it's time to build a specific phase, it can get its own requirements/design doc if needed. The brainstorm and principles are always there as reference — pull up `principles.md` during any implementation decision.
