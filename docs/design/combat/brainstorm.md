# Combat System Brainstorm: Silksong-Inspired Platformer Mode

**Author**: Karsten Huttelmaier — co-authored with Claude

## The Vision

When the player encounters an enemy in the dungeon, the entire screen transforms into a **full-screen, real-time side-scrolling platformer** — inspired by Hollow Knight / Silksong. The dungeon exploration stays small and understated (unchanged), making the mode switch a jaw-dropping contrast.

Long-term vision: different enemy types trigger different minigames (Space Invaders, Breakout, Pong, etc.), but we start with one killer platformer mode.

---

## Applying Sid Meier's Design Philosophy

### The Covert Action Rule — Our Biggest Risk

Meier's most instructive failure: Covert Action combined an action game with a mystery game. Each was fun alone, but *"together, they fought with each other."* This is **exactly our risk** — we're putting a platformer inside a roguelike. Two games in one.

**But Pirates! solved this.** Pirates! had ship combat, sword fighting, dancing, sailing — all different minigames inside one pirate RPG. The difference from Covert Action? Each minigame was **SHORT** (30-90 seconds) and **SERVED** the meta-game rather than competing with it.

**Rule for our combat**: Each fight must be 30-90 seconds. MAX 2 minutes for a boss. If a single combat encounter takes longer than that, it becomes its own game and fights the dungeon exploration for attention. The platformer combat must leave the player wanting MORE and eager to explore the next room — not exhausted and forgetting what they were doing in the dungeon.

### "A Game Is a Series of Interesting Decisions"

Reflexes alone aren't decisions. Our combat needs **choices that matter**:

**During combat:**
- Attack the enemy during their recovery window (big damage, risky timing) or play safe and wait?
- Use your one air dash offensively (close distance) or save it as an escape?
- Go for the downward nail-bounce (pogo attack, stylish) or stay grounded (safe)?
- Use a health potion now (survive but lose the item) or try to win without it?

**Before combat:**
- Which items from the dungeon do you bring into the fight? The sword (more damage) or the shield (more defense)? You found both but can only equip one.
- Do you even fight, or try to flee and find another path?

**After combat:**
- The enemy drops a reward. Take it (fills inventory slot) or leave it (save space for the key)?

**Meta-decisions across the dungeon run:**
- HP is a shared resource across ALL fights. Every hit you take in one fight is HP you don't have for the next. This makes dungeon exploration decisions matter — do you fight the optional enemy for the reward, knowing you'll burn HP?

### "The Player Must Be the Star"

Meier: *"Our job is to impress you with yourself."*

- **Generous timing windows**: Coyote time (jump after walking off edge), jump buffering (press jump before landing), generous attack hitboxes. The player should feel precise and skillful even when the game is being forgiving.
- **Enemies telegraph HARD**: Every attack has a clear windup animation. The player should feel smart for dodging, not frustrated by surprise hits.
- **Overkill moments**: When you land the killing blow, time slows briefly and the screen flashes. The player feels POWERFUL.
- **Skill ceiling without skill floor**: Basic attack spam can beat easy enemies. But expert players can dash-cancel, pogo-bounce, and combo for style points or faster kills. Reward mastery without punishing beginners.

### "Player Psychology Trumps Mathematics"

Meier's GDC 2010 keynote: Players take credit for wins and blame the game for losses. This means:

- **Deaths must feel fair**: Every death should feel like "I could have dodged that." Never kill the player with something they couldn't see coming. Always telegraph.
- **Near-misses feel exciting, not frustrating**: If the player dodges by 1 pixel, they should feel like a genius, not feel like the hitbox was too big.
- **Generous invincibility frames**: After getting hit, give the player a full second of invincibility. This prevents the "stun-lock death" that makes players rage-quit.
- **Health pickups from enemies**: When you hit an enemy, occasionally drop a small health orb (like Hollow Knight's Soul mechanic). This rewards aggression and makes the player feel like they're healing BECAUSE of their skill.

### "Simple Systems Interacting to Create Complexity"

Don't make combat complex. Make the INTERACTION between combat and dungeon create complexity:

- **Dungeon items affect combat**: Sword = more damage. Shield = damage reduction. Torch = lights up dark arenas. Potions = mid-fight healing.
- **Combat rewards affect dungeon**: Defeated enemies drop keys, reveal hidden passages, or leave loot.
- **HP persists across rooms**: Every fight is a COST. Exploration becomes risk management.
- **Score combines both worlds**: Rooms visited + items collected + enemies defeated + style bonuses.

These are all simple systems, but they interact to create the "multiple irons in the fire" feeling that drives Meier's "one more turn" loop — except ours is "one more room."

### "One More Room" — The Compulsion Loop

Applying Meier's "one more turn" architecture:

- **Multiple things in progress**: You're 2 rooms from the treasure. You have the key but low HP. There's an optional side path with a health potion but also an enemy. Do you risk it?
- **Curiosity about the next encounter**: What enemy type is in the next room? What minigame will it be? The MYSTERY drives exploration.
- **Visible progress**: The map fills in as you explore. Enemies crossed off. Score climbing. You can SEE your progress.
- **Near-term payoffs feeding long-term goals**: Defeating an enemy feels great NOW. But the loot from it helps you reach the treasure LATER.

### "Make It Epic" — The Mode Switch as Spectacle

The moment the screen transforms from a quiet little box-drawing map into a full-screen pixel-art battle arena should be the **defining moment** of this game. This is the thing people screenshot and share.

**The transition:**
1. Player walks into a room with an enemy
2. The dungeon text says: *"A shadow moves in the darkness..."*
3. Brief pause (0.5s) — tension builds
4. Screen WIPES to black
5. Full-screen arena FADES IN with the enemy revealed
6. Health bars appear. Controls prompt shown briefly.
7. FIGHT.

**The return:**
1. Enemy defeated — dramatic death animation (particles, flash)
2. Victory text: *"Enemy defeated! +15 score"*
3. Brief pause (1s) — savor the win
4. Screen WIPES back to the dungeon map
5. The room now shows the enemy is gone, loot on the ground

### "First 15 Minutes" — The First Combat Must Be Perfect

Meier: *"The most important part of a game is the first and last 15 minutes."*

The first time a player encounters an enemy, everything must click:
- The mode switch is dramatic and surprising
- The controls are immediately intuitive (move, jump, attack — that's it for the first fight)
- The first enemy is EASY — the player wins and feels powerful
- The visual quality is jaw-dropping (full color pixel art in a terminal?!)
- Return to dungeon is smooth — the player immediately understands the loop

Don't throw dash, directional attacks, or complex patterns at the player in fight #1. Save those for later enemies. The first fight teaches: move, jump, attack, win.

### "Double It or Cut It in Half" — Tuning Strategy

When we're tuning combat, follow Meier's rule:
- Enemy feels too hard? **Cut its HP in half.** Don't tweak by 10%.
- Attacks feel wimpy? **Double the damage.** See if it's even the right variable.
- Combat takes too long? **Cut the arena in half.** Or halve the enemy's aggression cooldown.
- Hit feedback isn't satisfying? **Double the screen shake.** Double the hit-stop frames.

We'll know we're done tuning when combat feels like the player is barely surviving but always winning.

### The Pirates! Model — Different Minigames as Different Genres

This is the long-term vision that makes the game special:

| Enemy Type | Minigame Genre | Dungeon Context |
|-----------|---------------|-----------------|
| **Husk Guard** | Side-scrolling platformer (Silksong) | Guards hallways and doors |
| **Goblin Horde** | Space Invaders (wave shooter) | Ambush in large rooms |
| **Skeleton Archer** | Breakout/Arkanoid (deflect shots) | Ranged enemy in open areas |
| **Mimic** | Pong (bizarre, comedic) | Disguised as a treasure chest |
| **Dragon Boss** | Full platformer boss fight (extended) | Guards the treasure room |

Each minigame is SHORT, SIMPLE, and DIFFERENT. The player never knows what genre the next room will throw at them. This is the "curiosity" that drives the "one more room" loop.

**The Covert Action safety valve**: Each minigame is self-contained. 30-90 seconds. Simple controls. The dungeon exploration is always the meta-game that ties everything together.

### What Meier Would Cut

Applying "know when to stop":
- ~~Experience/leveling system~~ — Too RPG. The dungeon items ARE the progression.
- ~~Multiple equipment slots~~ — Weapon + one accessory. That's enough for interesting decisions.
- ~~Elemental damage types~~ — Too complex. Save for v3.
- ~~Enemy respawning~~ — Defeated = defeated. Clearing a room should feel permanent and satisfying.
- ~~Branching combat dialogue~~ — You're fighting, not talking. Keep it pure.

---

## Wild Ideas to Explore

**Style Scoring**: Like Devil May Cry's style meter. Chain attacks without getting hit — style rank goes up (D -> C -> B -> A -> S). Higher style = more score/loot. Makes combat feel expressive.

**The "Last Stand" Mechanic**: When HP hits 0 in combat, you don't die immediately. Time slows. You get ONE final attack. If you kill the enemy with it, you survive with 1 HP. If you miss, game over. Makes every defeat feel like "I ALMOST had it."

**Enemy Learning**: Enemies in later rooms have slightly different patterns based on how you fought earlier ones. If you always dash-attack, later enemies start countering dashes. The game adapts to you. (This is advanced — save for later.)

**The Dungeon as Training**: Each dungeon room before a combat encounter could have environmental hints about the enemy. "Claw marks on the wall" = aggressive enemy. "Scorch marks" = fire-based attack patterns. The roguelike exploration becomes intelligence gathering for the combat.

**Combo System**: Specific sequences of attacks create special moves. Not required — but discoverable. Players who experiment are rewarded. Like finding a hidden chest in a Zelda game but for combat mechanics.

---

## DEEP DIVE: The Dungeon-Combat Feedback Loop

*This is where the game lives or dies. Not in the combat mechanics. Not in the graphics. In how the two halves FEED each other.*

### The Core Loop: Explore -> Prepare -> Fight -> Reward -> Explore

Every great game has a core loop that players repeat but that never feels repetitive because each cycle changes the context. Civilization's is: Explore -> Expand -> Exploit -> Exterminate. Ours is:

```
   +-----------------------------------------------------+
   |                                                     |
   v                                                     |
EXPLORE --> DISCOVER --> DECIDE --> FIGHT --> REWARD -----+
  map          item         risk?      combat     loot
  rooms        enemy        prepare?   (mode      changes
  paths        clue         retreat?    switch)    dungeon
```

Each step feeds the next. The REWARD from combat changes what you EXPLORE next. The DISCOVERY in exploration changes how you FIGHT. This is Meier's "multiple irons in the fire" — the player always has something cooking at every stage.

### HP as Universal Currency

This is the single most important design decision. **HP persists across the entire dungeon run.** This one rule transforms every room in the dungeon from "just another room" into a decision:

- *"I have 6 HP. The next room has an enemy. I could fight (risk losing HP) or take the long way around (costs turns, might miss loot)."*
- *"I just beat an enemy with only 2 HP left. There's a side room with a health potion... but also another enemy between me and it."*
- *"I'm at full HP with the key. Do I rush the locked door, or clear the optional enemies for score?"*

**HP is not a combat stat. It's an EXPLORATION stat.** It determines how much of the dungeon you can afford to explore. This is the interaction between simple systems that Meier talks about.

#### HP Budget for a Full Run

The numbers should create tension without frustration. Rough design:

| Value | Amount | Why |
|-------|--------|-----|
| Player starting HP | 5 hearts | Enough to survive 3-4 encounters with mistakes |
| Weak enemy damage | 1 heart per hit | Forgiving — new players can take 5 hits |
| Strong enemy damage | 2 hearts per hit | Scary — 2-3 mistakes and you're dead |
| Health potion healing | 2 hearts | Found in ~30% of rooms. A real decision to use |
| Player attack damage | Varies by weapon | Sword = 2, bare fists = 1 |
| Typical enemy HP | 5-8 hits to kill | Fights last 30-60 seconds |

With 5 hearts and ~3 enemies on the path to the treasure, the player needs to average less than ~1.5 hits taken per fight. That's tight enough to feel tense but generous enough that recovery is possible.

### Items That Bridge Both Worlds

Currently items are just "things you carry." With combat, every item becomes a **strategic decision**:

| Item | Dungeon Use | Combat Effect | Decision |
|------|-------------|---------------|----------|
| **Key** | Unlocks doors | No combat use | Must carry — can't drop it |
| **Sword** | None | 2x attack damage | Carry or leave for inventory space? |
| **Shield** | None | Block 1 damage per hit | Trade offense for defense? |
| **Health Potion** | Can drink anytime | Can drink mid-combat | Use now or save for the boss? |
| **Torch** | Lights up dark room descriptions | Reveals enemy patterns faster | Quality of life vs. combat slot |
| **Lucky Charm** | +bonus score per room | Enemies drop better loot | Risk/reward for score chasers |

The player can only carry ~3-4 items (keep inventory small — forces choices). Key takes one slot permanently. That leaves 2-3 slots for real decisions. This is Meier: "a game is a series of interesting decisions."

**The killer scenario:** You're carrying the key, a sword, and a health potion. You find a shield. Do you drop the health potion (lose safety net) or the sword (weaker attacks, longer fights)? There is no right answer. That's a perfect Meier decision.

### Loot and Combat Rewards

Defeating enemies must feel WORTH IT, or players will avoid combat:

| Reward Type | Effect | Frequency |
|------------|--------|-----------|
| **Health drop** (small) | Restore 1 heart | Common (~60% of fights) |
| **Item drop** | Random item (potion, weapon, charm) | Uncommon (~30%) |
| **Room clear bonus** | +15 score | Always |
| **Style bonus** | +5 score per style rank | If style system is added |
| **Path reveal** | Enemy drops a "map fragment" showing adjacent rooms | Rare (~15%) |

The health drop is critical — it means fighting is how you RECOVER HP, not just how you lose it. This creates a beautiful tension: combat is both the thing that hurts you AND the thing that heals you. Players who fight well are rewarded with sustain. Players who fight poorly drain their HP pool. The game rewards skill without punishing beginners too harshly.

### The Risk/Reward Topology

The dungeon layout itself becomes a strategic map:

```
         [Health Potion]
              |
[Start] -- [Enemy] -- [Key Room]
              |
          [Optional]     [Locked Door]
          [Enemy]            |
              |          [TREASURE]
          [Rare Loot]
```

**The main path** has 1-2 mandatory enemies between start and treasure. These are unavoidable — the player MUST engage.

**Side paths** have optional enemies guarding good loot (health potions, weapons, score bonuses). The player CHOOSES whether the reward is worth the HP cost.

This means the generator doesn't just place enemies randomly — it creates a **risk/reward topology** where every branch is a meaningful choice.

### Score as Narrative

Score isn't just a number. It tells the story of HOW you played:

```
Final Score Breakdown:
  Rooms explored:    8 x 5  =  40
  Items collected:   3 x 10 =  30
  Enemies defeated:  4 x 15 =  60
  Style bonuses:           =  20
  HP remaining:      3 x 5  =  15  (new: reward for finishing healthy)
  ────────────────────────────────
  TOTAL:                      165
```

Two players can win with vastly different scores:
- **The Explorer**: High room count, low combat score. Found every corner.
- **The Fighter**: Low room count, high combat score. Cleared every enemy.
- **The Speedrunner**: Low everything except HP bonus. Beelined for the treasure.

All three feel valid. Meier: "Allow for as many choices and play styles as possible."

---

## DEEP DIVE: The Visual Spectacle

*The mode switch is not just a UI change. It's the emotional peak of the game.*

### The Contrast IS the Feature

The dungeon is deliberately understated:
- Small box-drawing map (5x3 rooms, thin corridors)
- Monochrome palette with subtle color accents
- Quiet, contemplative, turn-based
- Fits in a corner of the screen

The combat is deliberately MAXIMAL:
- Full-screen pixel art filling every terminal cell
- Rich color palette (dark atmosphere with bright action accents)
- Fast, fluid, 30fps real-time action
- Particles, screen shake, hit-stop, flash effects

The gap between these two extremes IS the spectacle. The less flashy the dungeon is, the more jaw-dropping the combat entrance becomes. Don't upgrade the dungeon visuals. Keep them understated. The contrast does the work.

### The Transition Sequence

This is the player's first impression of combat. It has to be PERFECT.

**Entering combat (1.5 seconds total):**
1. **The Warning** (0.3s): Player enters room — message appears: *"Something stirs in the darkness..."* The map view is still showing. Tension builds.
2. **The Freeze** (0.2s): Screen freezes. Brief pause. The calm before the storm.
3. **The Wipe** (0.3s): Screen rapidly fills with black, column by column (like a curtain closing). The dungeon disappears.
4. **The Reveal** (0.5s): Combat arena fades in from black. Background first, then platforms, then the enemy materializes.
5. **The Standoff** (0.2s): Player character appears. Health bars slide in from the edges. Controls flash briefly at the bottom.

**Exiting combat (1.0 seconds):**
1. **The Kill** (0.3s): Final hit — enemy flashes white — explodes into particles. Brief slow-mo.
2. **The Reward** (0.3s): Loot and score pop up on screen.
3. **The Return** (0.4s): Screen wipes back to dungeon. Enemy removed from room. Loot items appear.

These transitions should feel like a **movie cut**, not a loading screen.

### The Arena: Dark, Atmospheric, Alive

The combat arena isn't just platforms and a flat background. It's a PLACE:

**Background layers (parallax if we get fancy, flat is fine for v1):**
- Deep background: Dark stone walls, faint torchlight glow
- Mid-ground: Pillars, chains, cobwebs — set dressing that doesn't affect gameplay
- Foreground: The platforms and floor the player interacts with

**Lighting mood:**
- Overall dark palette (deep blues, grays, blacks)
- Player character slightly brighter than surroundings (they're the star)
- Enemy eyes or weapon glow with accent color (red, green, purple depending on type)
- Hit effects are BRIGHT (white flashes, yellow sparks) — stand out against the darkness

**The arena should feel like the dungeon room you entered, expanded to fill the screen.** If the room description said "A damp cellar," the arena has dripping water pixels and mossy platforms. If it said "A forgotten library," the platforms are bookshelves. This connects the two modes narratively.

### Frame-by-Frame Hit Feedback

This is what separates "good" from "INCREDIBLE." Each element is simple but they combine:

**When the player's nail hits an enemy:**
1. Frame 0: Attack hitbox overlaps enemy hurtbox — HIT detected
2. Frames 0-2: **HIT-STOP** — entire simulation freezes. Both characters stay in their current pose. This is the "weight" that makes hits feel meaty. 3 frames = 100ms. Crucial.
3. Frame 0: **HIT FLASH** — enemy renders as ALL WHITE for this frame
4. Frame 1-2: Enemy flickers (alternating white/normal)
5. Frame 3: Simulation resumes. **KNOCKBACK** — enemy receives velocity impulse (200px/s) away from player
6. Frames 3-8: **SCREEN SHAKE** — render offset oscillates +/-2 pixels for 5 frames
7. Frame 3: **PARTICLES** — 4-6 small pixel particles spawn at hit point, fly outward with random velocity, fade over 10 frames
8. Frame 3: Enemy enters **HURT** state with invincibility frames (prevents multi-hit from one swing)

**When an enemy hits the player:**
- Same feedback but MORE dramatic (bigger screen shake, longer hit-stop)
- Player flashes red instead of white
- Player knocked back further
- Invincibility lasts longer (30 frames / 1 second — generous, as Meier would want)
- If HP drops to 1, the screen flashes red briefly. Warning without being annoying.

**When the enemy dies:**
- Extended hit-stop (6 frames)
- Enemy flashes white rapidly
- Explodes into a shower of particles (20-30 particles, larger, slower)
- Brief screen flash
- Time resumes at half speed for 10 frames (slow-mo victory moment)
- Loot/reward appears where enemy was standing

This is 100% "juice." None of it affects gameplay. All of it affects how the game FEELS. And feeling is everything.

---

## DEEP DIVE: The "One More Room" Addiction

*Civilization has "one more turn." We need "one more room." Here's how every system conspires to keep the player exploring.*

### The Three Hooks Running Simultaneously

At any moment during the dungeon, the player has at least three reasons to keep playing:

**Hook 1 — The Puzzle**: "Where is the key? Where is the locked door? Can I reach the treasure?" This is the existing game's hook. It works. Don't change it.

**Hook 2 — The Economy**: "I have 3 HP and a sword. The next room might have an enemy. But it also might have a health potion. Do I risk it?" Every room is a gamble. HP and items are the chips.

**Hook 3 — The Spectacle**: "What will the next combat look like? Will I fight something new?" The sheer novelty and visual drama of combat makes the player WANT to find enemies. Most roguelikes make you dread enemies. Ours makes you seek them out because fighting is FUN and looks INCREDIBLE.

When these three hooks overlap — the key is behind an enemy that you don't know if you can afford to fight — that's when the game sings. That's "one more room."

### Meier's "Multiple Irons in the Fire"

At any given moment, the player should be tracking 3-5 things:

In the **dungeon**:
- "I need to find the key" (goal)
- "I have 3 HP, can I afford another fight?" (resource)
- "That side room might have a health potion" (opportunity)
- "I've explored 7/10 rooms" (progress)
- "My score is 120, can I beat my best?" (meta-goal)

In **combat**:
- "The enemy is winding up — dodge NOW" (immediate)
- "I have 2 HP — one more hit kills me" (stakes)
- "The enemy is at low HP — press the attack?" (decision)
- "I have a health potion — use it or try to finish without?" (resource)

Between the two modes, the player's mind is never idle. There's always something to think about, decide, or anticipate.

### The Information Drip

The player should never have perfect information. Mystery drives exploration.

- **Fog of war** (already exists): You can't see rooms you haven't visited. This works.
- **Enemy mystery**: You know a room has an enemy (dungeon text says "You hear growling nearby") but you don't know WHAT kind or how strong until you enter.
- **Loot mystery**: You don't know what an enemy will drop until you beat it.
- **Path mystery**: There might be a health potion on the other side of an optional enemy. Or there might be nothing.

This is Meier's "cliffhanger structure" — each room resolves some questions but opens new ones.

### The Session Arc

A single game run should take 5-15 minutes and feel like a complete story:

1. **The Opening** (rooms 1-3): Safe exploration. Find items. Learn the layout. Maybe one easy fight. Build confidence.
2. **The Midgame** (rooms 4-7): Tension rises. Harder enemies. Resource decisions matter. Find the key. Stakes increase.
3. **The Crisis** (rooms 7-9): Low HP. The locked door found. One or two fights standing between you and victory. Every room is "do I risk it?"
4. **The Climax**: The final fight before the treasure room. All resources on the line.
5. **The Resolution**: Unlock the door. Reach the treasure. Score revealed. "Play again?"

This arc mirrors a movie structure. The player feels like they went on a JOURNEY, not just solved a puzzle.

---

## DEEP DIVE: Sound and Music

*The mode switch you can HEAR.*

### Sound Is Half the Spectacle

The mode-switch impact needs AUDIO as much as visuals. Think FPS Chess — the audio shift is just as jarring and delightful as the visual one.

For our game:
- **Dungeon exploration**: Silence, or the faintest ambient hum. Dripping water. Distant echoes. Barely there. Understated — just like the visuals. This is the quiet before the storm.
- **Combat mode**: Driving chiptune track kicks in IMMEDIATELY when the screen transforms. Fast tempo, heavy bass, 8-bit energy. The audio contrast hits you the same moment the visual contrast does.
- **The crossfade**: 2-3 second transition. Dungeon ambient fades out. Brief silence (the "freeze" moment). Then combat music PUNCHES in with a hard downbeat timed to the arena reveal.

This isn't polish. This IS the feature. The mode switch should assault two senses at once.

### The Sound Palette

| Sound | When | Feel |
|-------|------|------|
| **Nail slash** | Player attacks | Sharp, metallic *SHING* — crispy and fast |
| **Hit impact** | Attack connects | Deep, meaty *THUD* — weight and power |
| **Enemy death** | Final blow | Satisfying *CRASH* + sparkle — like breaking glass |
| **Player hurt** | Taking damage | Sharp pain sound + bass rumble — scary but not annoying |
| **Item pickup** | Collecting loot | Bright chime — *ding!* — instant dopamine |
| **Health restore** | Healing | Warm ascending tone — relief |
| **Door unlock** | Using key | Heavy mechanical *CHUNK* — accomplishment |
| **Combat enter** | Mode switch | Rising tension -> BOOM -> combat music |
| **Combat victory** | Enemy defeated | Music hits a triumphant riff, then fades |

Every sound reinforces the player-as-star principle. Attacks sound POWERFUL. Victories sound TRIUMPHANT. The player feels like a hero because they sound like one.

### Technical Approach

- **Library**: `gopxl/beep` v2 — actively maintained Go audio library. Supports WAV, OGG, MP3, and tracker formats (MOD/XM).
- **No terminal conflicts**: Audio runs in a separate goroutine, completely independent from Bubble Tea's terminal I/O.
- **Architecture**: An `audio/` package with a manager that receives commands via `tea.Cmd`. Play, stop, crossfade — all async, never blocking the game loop.
- **Music format**: OGG Vorbis for pre-made tracks (good compression) or MOD/XM for authentic chiptune feel and tiny file sizes.
- **SFX format**: WAV for sound effects (low latency, no decode overhead).

```
audio/
  manager.go    # Playback goroutine, crossfade logic, volume control
  types.go      # PlayTrack, PlaySFX, FadeOut message types
assets/
  music/
    dungeon.ogg     # Quiet ambient exploration track
    combat.ogg      # Driving chiptune battle track
    victory.ogg     # Short victory fanfare
  sfx/
    slash.wav, hit.wav, death.wav, pickup.wav, hurt.wav, unlock.wav
```

### Free Music Sources

- **Soundimage.org** — Royalty-free chiptune tracks, dungeon/arcade style
- **Ozzed.net** — Creative Commons 8-bit music library
- **itch.io** — Free chiptune asset packs (music + SFX)
- **FamiStudio** — Free NES-style tracker if we want to compose our own

### When to Add Sound

Sound is NOT in the first prototype. Principle 8: "Prototype the fun first." Rectangles before sprites. Silence before music. But sound is a HIGH-PRIORITY second pass because it transforms the mode switch from "cool" to "holy crap."

---

## DEEP DIVE: Renderer Swappability

*Build for today, architect for tomorrow.*

### The Architecture Already Supports It

The `pixelbuf/` package is the abstraction boundary. The combat engine writes to a pixel buffer (abstract 2D color grid). It doesn't know what happens after that. Today, `pixelbuf.Render()` converts it to half-block terminal characters. Tomorrow, it could be:

```
combat/engine/  ->  pixelbuf.Buffer  ->  pixelbuf.Render()        [half-block, today]
                                     ->  sixel.Render()           [Sixel protocol, future]
                                     ->  kitty.Render()           [Kitty protocol, future]
                                     ->  web.Render()             [WebSocket to browser, future]
                                     ->  sdl.Render()             [SDL2 window, if we ever leave terminal]
```

The key insight: the `Buffer` type (a 2D grid of RGBA colors) is renderer-agnostic. Any system that can consume pixels can render the game. No extra abstraction layer needed — the pixel buffer IS the abstraction.

### What This Means Practically

- Start with half-block rendering (universal, proven, good enough).
- If someone wants Sixel support later, they write ONE function: `SixelRender(buf *Buffer) string`. The engine, physics, sprites, and everything else are untouched.
- If the game ever outgrows the terminal, the combat engine could render to an SDL window. The engine code stays identical — only the render backend changes.

### What NOT to Do

Don't build a "renderer interface" with multiple implementations today. That's premature abstraction. The half-block renderer is the only one we need. The clean package boundary means swapping is trivial IF we ever need it. YAGNI until proven otherwise.

---

## Rendering: Best Possible Terminal Graphics

After researching all terminal rendering approaches, **half-block pixel art** is the best option:

| Approach | Color | Resolution | Terminal Support | Verdict |
|----------|-------|-----------|-----------------|---------|
| **Half-block** | Full 24-bit per sub-pixel | 2x vertical | Universal | **Best choice** |
| Braille | Monochrome per cell | 2x4 per cell | Good | No color = no good |
| Sixel/Kitty | Full bitmap | Arbitrary | Fragmented | Unreliable |

This is the same technique used by gambatte-terminal (a full Game Boy Color emulator in terminal at 60fps) and movy (Go terminal graphics engine with alpha blending). It's the gold standard.

- **Resolution**: Dynamically sized via `tea.WindowSizeMsg` to fill the entire screen. Full-screen combat.
- **Full 24-bit color**: Two independent colors per terminal cell via fg/bg ANSI escapes.
- **Sprites**: Hardcoded Go pixel art (~12x16 pixel characters). Flip horizontally based on facing. Flash white on hit.
- **Alpha blending**: Support transparency in sprites for layered rendering (background -> platforms -> entities -> effects).
- **Particle effects**: Dust on landing, slash arcs on attack, sparks on hit — all rendered as small animated sprites.

No new rendering dependencies — just Unicode + ANSI escapes + lipgloss for HUD.
