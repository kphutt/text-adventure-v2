# Phase 2 Execution Plan: Rectangles That Fight

**Author**: Karsten Huttelmaier — co-authored with Claude

## Context

Phase 2 is where the combat initiative lives or dies. The backlog says: "If moving a rectangle around a platform arena, jumping over attacks, and slashing another rectangle feels *good* — then everything after this is polish. If it doesn't feel good, no pixel art or screen shake will save it."

Two halves, one phase:
1. **Pure Engine** (`combat/engine/`) — deterministic game logic, zero UI dependencies, fully testable
2. **Combat Testbed** (`cmd/combat-proto/main.go`) — permanent standalone tool wiring engine to pixelbuf + Bubble Tea

The testbed is standalone — it does NOT modify `main.go` or the dungeon game. Phase 3 handles integration via a separate throwaway integration spike that tests how combat mode fits into the existing Bubble Tea main app (mode dispatcher, transitions, game loop coexistence). The testbed survives Phase 3 as a tuning/debugging/regression tool.

## Coordinate Convention

**Y increases downward.** Gravity is positive (accelerates downward). JumpForce is applied as `VelY = -JumpForce` (negative = upward). MaxFallSpeed clamps VelY to a positive maximum. The floor is at the bottom of the arena (high Y value).

## Dependency Rules

```
combat/engine/     imports ONLY stdlib (not world/, not pixelbuf/, not oto)
cmd/combat-proto/  imports combat/engine/, pixelbuf/, bubbletea, oto
pixelbuf/          imports ONLY stdlib (already built)
main.go            NOT TOUCHED in Phase 2
```

Note: This supersedes the Phase 1 plan's dependency diagram, which showed `combat/engine/` importing `world/`. The engine defines its own types; Phase 3 bridges them.

Note: The backlog scopes `combat/input.go` for reuse. Input handling lives in the testbed for now. Phase 3 integration spike will determine whether to extract it to `combat/input.go` or inline it into the main app's mode dispatcher.

---

## Package Structure

```
combat/
  engine/
    consts.go              # Physics constants, arena dimensions, enums
    entity.go              # Rect, Player, Enemy, Platform, InputState
    physics.go             # Gravity, moveAndResolve collision
    player.go              # Player state machine, input handling
    attack.go              # Attack hitbox, damage, knockback
    engine.go              # Engine struct, Tick, NewEngine, Reset, arena builder
    engine_test.go         # Integration tests (multi-tick scenarios)
    entity_test.go         # Rect/AABB unit tests
    physics_test.go        # Gravity + collision resolution tests
    player_test.go         # State machine transition tests
    attack_test.go         # Hitbox + damage tests
    test_helpers_test.go   # testEngine(), emptyInput(), assertNear()

cmd/
  combat-proto/
    main.go                # Permanent standalone combat testbed (Bubble Tea)
```

---

## Type Definitions

### Enums (`consts.go`)

```go
type PlayerState int
const (
    StateIdle PlayerState = iota
    StateRun
    StateJump
    StateFall
    StateAttack
    StateHurt
)

type Dir int
const (
    DirRight Dir = 1
    DirLeft  Dir = -1
)

type Result int
const (
    ResultNone Result = iota
    ResultPlayerWin
    ResultPlayerDead
)
```

### Physics Constants (`consts.go`)

```go
const TickRate = 30
const DT = 1.0 / float64(TickRate) // ~0.0333s

// Coordinate system: Y increases downward. Positive VelY = falling.

// Physics (pixels/sec or pixels/sec^2)
Gravity      = 800.0   // positive = accelerates downward
MaxFallSpeed = 400.0   // VelY clamped to this positive value
JumpForce    = 300.0   // applied as VelY = -JumpForce (negative = upward)
RunSpeed     = 200.0

// Timing windows (seconds)
// Revised from backlog's "4 frames" to 3 frames at 30fps.
CoyoteTime     = 0.1   // ~3 frames — can jump briefly after leaving edge
JumpBufferTime = 0.1   // ~3 frames — press jump slightly before landing
JumpCutMultiplier = 0.5 // release jump early = cut upward velocity (applied ONCE)

// Attack
AttackDuration  = 0.2   // hitbox active ~6 frames
AttackCooldown  = 0.35  // between swings
AttackWidth     = 20.0
AttackHeight    = 14.0
AttackOffsetX   = 14.0  // from player center to hitbox edge
AttackDamage    = 1

// Entities
PlayerWidth, PlayerHeight = 12.0, 20.0
PlayerHP = 5
EnemyWidth, EnemyHeight = 14.0, 20.0
EnemyHP = 3

// Hurt
HurtDuration = 0.5
InvincTime   = 1.0
KnockbackVel = 150.0   // velocity impulse on hit (px/s), NOT position teleport

// Arena
ArenaWidth  = 160       // pixels (= terminal columns)
ArenaHeight = 120       // pixels (= 60 terminal rows at half-block)
```

### Core Types (`entity.go`)

**Rect** — AABB primitive:
```go
type Rect struct { X, Y, W, H float64 }
func (a Rect) Overlaps(b Rect) bool   // strict inequality: non-zero area intersection
func (a Rect) Center() (float64, float64)
```

**InputState** — per-frame input snapshot (engine has no UI knowledge):
```go
type InputState struct {
    Left, Right bool    // movement held
    JumpPress   bool    // jump pressed THIS frame (edge-detected)
    JumpHeld    bool    // jump button currently held
    Attack      bool    // attack pressed THIS frame
}
```

**Player:**
```go
type Player struct {
    Pos       Rect
    VelX, VelY float64
    Facing     Dir
    State      PlayerState
    HP, MaxHP  int
    Grounded   bool

    CoyoteTimer         float64 // counts UP from 0; time since last grounded
    JumpBufferTimer     float64 // counts UP from 0; time since last JumpPress
    JumpCut             bool    // true once variable-height cut applied this jump
    AttackTimer         float64 // counts down during attack
    AttackCooldownTimer float64 // counts down between attacks
    HurtTimer           float64
    InvincTimer         float64
    AttackHit           bool    // prevents multi-hit per swing
}
```

**Enemy** (Phase 2 minimal — stands there, takes hits, dies; has physics for gravity + knockback):
```go
type Enemy struct {
    Pos         Rect
    VelX, VelY  float64   // needed for gravity (initial fall) and knockback
    Grounded    bool
    HP, MaxHP   int
    Facing      Dir
    HurtTimer   float64
    InvincTimer float64
    Alive       bool
}
```

**Platform:**
```go
type Platform struct { Rect Rect }
```

### Engine (`engine.go`)

```go
type Engine struct {
    Player    Player
    Enemy     Enemy
    Platforms []Platform
    Result    Result
    TickCount int
}

func NewEngine() *Engine
func (e *Engine) Reset()
func (e *Engine) Tick(input InputState)
```

`NewEngine()` initializes:
- Player at (20, 0) with HP=PlayerHP, MaxHP=PlayerHP, Facing=DirRight
- Enemy at (120, 0) with HP=EnemyHP, MaxHP=EnemyHP, Facing=DirLeft, Alive=true
- JumpBufferTimer = JumpBufferTime (prevents false first-frame trigger from Go zero value)
- All 7 platforms (floor, walls, platforms 1-4)

Note: `Dir` has no zero-value variant (DirRight=1, DirLeft=-1). Facing fields MUST be set explicitly. `Enemy.Alive` defaults to false in Go — MUST be set to true.

`Tick()` orchestration order:
1. Handle player input + state machine (uses CoyoteTimer/JumpBufferTimer from *previous* frame — this one-frame lag is intentional and standard for coyote time)
2. Apply gravity to player and enemy
3. Move and resolve collisions for player, then enemy
4. Update coyote timer: if player grounded, reset to 0; else increment by DT
5. Process attack hitbox / damage
6. Decrement enemy timers (clamp to 0)
7. Decrement player timers (clamp to 0)
8. Check win/lose conditions (enemy death checked first; if both die same frame, player wins)

**Entity collision:** Player and enemy bodies do not collide — entities pass through each other. This is intentional (matches Hollow Knight-style design where body collision creates frustrating movement traps). The attack hitbox is the only player-enemy interaction.

---

## Core Algorithms

### Physics — "Move X, Resolve; Move Y, Resolve"

```go
func applyGravity(velY *float64)
func moveAndResolve(pos *Rect, velX, velY *float64,
    grounded *bool, platforms []Platform)
```

`applyGravity`: `*velY += Gravity * DT`, then clamp: `*velY = min(*velY, MaxFallSpeed)`.

`moveAndResolve`:
1. Move X by `*velX * DT`
2. For each overlapping platform: push entity out horizontally, zero `*velX`
3. Move Y by `*velY * DT`
4. Set `*grounded = false`
5. For each overlapping platform:
   - If falling (`*velY >= 0`): land on top, set `*grounded = true`
   - If rising (`*velY < 0`): bonk ceiling, push down
   - Zero `*velY`

All platforms are fully solid from every direction. No one-way platforms in Phase 2.

Called once per tick for player AND once for enemy (enemy has velocity for gravity + knockback).

**Known limitation:** X-then-Y resolution creates a directional bias at platform corners — a player falling diagonally into a platform edge may be pushed sideways instead of landing on top. Add a test for this case. If it feels wrong during Step 7, the fix is to resolve the axis with the smaller penetration first.

### Player State Machine

States: Idle, Run, Jump, Fall, Attack, Hurt

Key mechanics:
- **Coyote time:** CoyoteTimer resets to 0 when grounded (step 4), increments when airborne. Jump allowed in step 1 if `CoyoteTimer < CoyoteTime`. Because step 1 runs before step 4, the timer reflects LAST frame's grounded state — this is correct and standard.
- **Jump buffering:** JumpBufferTimer resets to 0 on JumpPress, increments by DT otherwise. Jump executes if `JumpBufferTimer < JumpBufferTime AND (grounded OR CoyoteTimer < CoyoteTime)`.
- **Variable-height jump:** If `!JumpHeld` and `VelY < 0` (ascending) and `!JumpCut`, multiply VelY by JumpCutMultiplier and set `JumpCut = true`. The `JumpCut` flag resets when a new jump is initiated. This prevents repeated halving every frame.
- **Attack:** Initiates when Attack=true and cooldown expired. Player can't move during attack. After AttackDuration, transitions back to Idle/Fall. Cooldown starts on attack end.
- **Hurt:** Blocks all input. Lasts HurtDuration, then transitions to Idle/Fall. Phase 2 enemy doesn't attack — this state is structural only. Phase 4 will define player knockback vector when enemy attacks are added.

### Attack System

```go
func AttackHitbox(p *Player) Rect   // exported for prototype rendering
func processAttack(p *Player, e *Enemy)
```

- **AttackHitbox:** Returns `Rect{}` (zero value) if `p.State != StateAttack` or `p.AttackTimer <= 0`. Otherwise returns hitbox positioned relative to player center, offset by Facing direction.
- **processAttack** guards: enemy alive, AttackHitbox overlaps enemy, not already hit this swing (`AttackHit`), enemy not invincible (`InvincTimer <= 0`)
- **On hit:** Decrement enemy HP (clamp to 0). Set enemy HurtTimer = HurtDuration, InvincTimer = InvincTime. Apply knockback: set `enemy.VelX = KnockbackVel` (or `-KnockbackVel`) pushing enemy away from player. Set `AttackHit = true`.
- **Enemy dies** when HP <= 0, set `Alive = false`.

---

## Arena Layout

160 x 120 pixels. At half-block rendering: 160 columns x 60 terminal rows (+ 2-3 rows for HUD text below).

**Terminal size requirement:** Prototype checks terminal size at startup via `tea.WindowSizeMsg`. If < 160 columns or < 63 rows, display an error message suggesting the user resize or maximize their terminal.

```
       0    20        62        104       156
  0    +wall+                             +wall+
       |    |                             |    |
  44   |    |      +------+               |    |
       |    |      |Plat 4|               |    |
       |    |      |(70,44|               |    |
  64   |    |   +----------+              |    |
       |    |   | Plat 2   |              |    |
       |    |   |(62,64,36)|              |    |
  84   |    |+--------+  +--------+       |    |
       |    ||Plat 1   |  |Plat 3  |      |    |
       |    ||(20,84)  |  |(104,84)|      |    |
 112   +----+==========================+--+----+
       |              Floor                     |
 120   +----------------------------------------+
```

| Platform | X | Y | W | H |
|----------|---|---|---|---|
| Floor | 0 | 112 | 160 | 8 |
| Left wall | 0 | 0 | 4 | 120 |
| Right wall | 156 | 0 | 4 | 120 |
| Platform 1 | 20 | 84 | 36 | 6 |
| Platform 2 | 62 | 64 | 36 | 6 |
| Platform 3 | 104 | 84 | 36 | 6 |
| Platform 4 | 70 | 44 | 20 | 6 |

Player starts at (20, 0), enemy at (120, 0). Both fall to floor via gravity (~15 ticks at 30fps).

Max jump height: `v^2 / (2g) = 300^2 / 1600 = 56 px` — enough to reach any platform from the one below.

---

## Combat Testbed (`cmd/combat-proto/main.go`)

### Input Handling — Timeout-Based Key Release

Bubble Tea has no key-up events. Solution: `map[string]time.Time` tracking last-seen time per key. A key is "held" if seen within the timeout window. Edge detection for jump/attack (true only on transition from not-held to held).

**Key release timeout: start at 50ms.** This is a tuning parameter — if movement stutters, increase toward 100ms. If controls feel laggy (especially variable-height jump), decrease toward 35ms. Test on Windows Terminal with default key repeat settings.

```go
func (m *model) buildInput(now time.Time) engine.InputState
```

Keys: A/D = move, Space = jump, F = attack, R = restart, Q/Esc = quit.

### Game Loop

```go
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    case tea.KeyMsg:  // track key timestamp
    case tickMsg:     // buildInput -> engine.Tick(input) -> return tick()
}
```

30fps via `tea.Tick(tickDuration, ...)`. Each tick: build input, call `engine.Tick()`, schedule next tick.

### Terminal Size Check

On `tea.WindowSizeMsg`, store width/height. In `View()`, if terminal is too small (< 160 cols or < 63 rows), render a message: "Terminal too small (need 160x63, have WxH). Please resize." instead of the game frame.

### Rendering Pipeline

Each frame in `View()`:
1. `buf.Clear(bgColor)`
2. `FillRect` each platform (gray) — cast float64 to int via `int()`
3. `FillRect` enemy (red, blink white during hurt)
4. `FillRect` player (cyan, blink during invincibility)
5. `FillRect` attack hitbox when active (yellow)
6. Draw HP pips at top of screen (green/red squares)
7. `pixelbuf.Render(buf)` -> ANSI string
8. Append HUD text below frame (controls, win/lose message)

Colors:
- Background: dark blue-gray `{20, 20, 30, 255}`
- Platforms: medium gray-blue `{80, 80, 100, 255}`
- Player: bright cyan `{100, 200, 255, 255}`
- Enemy: bright red `{255, 80, 80, 255}`
- Attack hitbox: bright yellow `{255, 255, 100, 255}`
- HP full: green `{80, 220, 80, 255}`

### Sound Spike (Step 6b)

Minimal procedural audio via `github.com/ebitengine/oto/v3` (cross-platform: WASAPI on Windows, Core Audio on macOS, ALSA on Linux). New dependency added to `go.mod`.

**Scope:** 2 sounds, procedurally generated (no sound files):
- **Sword hit:** Short noise burst (~50ms) with downward frequency sweep. Satisfying impact feel.
- **Jump:** Quick rising sine tone (~30ms). Subtle feedback.

**Architecture:**
- Sound generation lives entirely in the testbed (not the engine)
- The testbed detects state changes each frame: `player entered StateAttack`, `enemy.HP decreased`, `player entered StateJump`
- On state change, fire-and-forget a sound via oto
- ~50-80 lines of sound code in the testbed

**Incremental path:**
- Phase 2: Procedural spike (2 sounds, validates pipeline)
- Phase 4: Full sound system — `beep` library for WAV/OGG file playback, mixing, volume. Richer designed sounds.
- Future: Composed audio assets. Atmospheric retro style (think Celeste/Hyper Light Drifter — dark, resonant, impactful within pixel-art register).

---

## Test Plan

All tests are deterministic — no randomness, no time.Now(). Tests call `Tick()` with scripted InputState values and assert state.

**Float comparison convention:** Use `assertNear(t, got, want, epsilon)` helper (epsilon = 0.1 pixels) instead of exact equality for positions. Exact equality is fine for integers (HP), booleans, and enums.

### entity_test.go
- Rect.Overlaps: overlapping, adjacent (no overlap), separated, same rect, zero-size
- Rect.Center: verify center calculation

### physics_test.go
- applyGravity: increases velocity, clamps to MaxFallSpeed
- moveAndResolve: falls onto floor (grounded), wall collision (pushed out), ceiling bonk, free movement (no collision)
- moveAndResolve: diagonal approach to platform corner — verify player lands on top (not pushed sideways)

### player_test.go
- Idle with no input: stays idle, VelX=0
- Run left/right: VelX set, Facing updated, State=Run
- Jump from ground: VelY=-JumpForce, State=Jump
- Coyote time: jump succeeds shortly after leaving edge, fails after timeout
- Jump buffer: early jump press executes on landing
- Variable-height jump: release early cuts velocity ONCE (verify JumpCut flag prevents repeated halving)
- Attack initiation: State=Attack, timer set
- Attack cooldown: can't attack while cooldown > 0
- Hurt blocks input: no movement during hurt state

### attack_test.go
- AttackHitbox facing right/left: correct position
- AttackHitbox when not attacking: returns zero Rect
- processAttack hits enemy: HP decremented, timers set, VelX set (knockback)
- Enemy invincible: no damage
- AttackHit prevents multi-hit per swing
- Enemy dies at HP=0, Alive=false
- Knockback direction: enemy pushed away from player (VelX sign matches direction)
- Knockback does not push enemy through wall: place enemy near right wall, hit, verify enemy stops at wall

### engine_test.go (integration)
- Player falls to ground: loop Tick with empty input until Grounded, assert < 60 ticks
- Enemy falls to ground: same — both entities have gravity
- Player runs right: grounded + Right input -> X increases
- Player jumps: grounded + JumpPress -> VelY < 0, rises then falls
- Full combat: position player near enemy, attack 3 times -> ResultPlayerWin
- Reset: run ticks, call Reset(), verify initial state restored
- No update after result: Tick() is no-op when Result != None

### test_helpers_test.go
- `testEngine()` — creates engine, ticks until player AND enemy are grounded (max 60 ticks, panic if exceeded)
- `emptyInput()` — returns zero InputState
- `assertNear(t, got, want, epsilon float64)` — float comparison helper

---

## Implementation Order

| Step | What | Verify |
|------|------|--------|
| 0 | Save this plan to `docs/design/combat/phase2-plan.md` | N/A |
| 1 | Types + constants (`consts.go`, `entity.go`, `entity_test.go`) | `go test ./combat/engine/` |
| 2 | Physics + collision (`physics.go`, `physics_test.go`) | tests pass |
| 3 | Player state machine (`player.go`, `player_test.go`) | tests pass |
| 4 | Attack + damage (`attack.go`, `attack_test.go`) | tests pass |
| 5 | Engine integration (`engine.go`, `engine_test.go`, `test_helpers_test.go`) | tests pass, `go test ./...` all pass |
| 6 | Combat testbed (`cmd/combat-proto/main.go`) | `go run ./cmd/combat-proto/` is playable |
| 6b | Sound spike: add 1-2 procedural sounds via `oto` | Sword hit and jump are audible during gameplay |
| 7 | Playtest + tune constants (including sound) | Gate: "Is the fun there?" |

---

## Design Decisions

| Decision | Rationale |
|----------|-----------|
| No `world/` import in engine | Supersedes Phase 1 dependency diagram. Engine defines its own types. Phase 3 bridges. |
| No one-way platforms | Adds complexity. All platforms fully solid. Deferrable. |
| Timeout-based key release (50ms, tunable) | Bubble Tea has no key-up events. Start conservative; increase if stuttery, decrease if laggy. |
| Edge detection in prototype, not engine | Engine receives clean InputState. Prototype handles raw keyboard quirks. |
| Player/enemy start in air | Tests physics from frame one. Both entities fall ~15 ticks to floor. |
| Enemy has VelX/VelY/Grounded | Needed for gravity (initial fall) and velocity-based knockback. Prepares for Phase 4 AI movement. |
| Knockback is velocity impulse, not position teleport | Integrates with moveAndResolve so knockback respects walls. Position teleport clips through geometry. |
| Player-enemy bodies don't collide | Intentional. Avoids frustrating movement traps. Attack hitbox is the only interaction. |
| JumpCut bool prevents repeated velocity halving | Without it, releasing jump halves VelY every frame while ascending. |
| JumpBufferTimer initialized above threshold | Prevents false jump trigger on first frame from Go zero value. |
| AttackHit boolean per swing | Prevents multi-hit from sustained overlap. Resets on new attack. |
| Fixed arena (160x120) with size check | Testbed shows error if terminal too small. Phase 3 adds dynamic sizing. |
| float64 physics, epsilon test assertions | Simple. Tests use assertNear (0.1px epsilon) to avoid float noise failures. |
| Flat structs, not ECS | ~5 entities total. ECS is overkill. Flat structs = direct field access. |
| iota + switch state machine | Go-idiomatic. Faster than interface dispatch. Zero allocations. |
| oto added as new dependency | Cross-platform audio (WASAPI/CoreAudio/ALSA). Engine stays stdlib-only; oto is testbed-only. |
| HP clamped to 0 on damage | Prevents negative HP edge cases. |
| Timers clamped to 0 on decrement | Prevents negative timer accumulation over long sessions. |
| combat/input.go deferred to Phase 3 | Input handling lives in testbed for now. Phase 3 integration spike determines final location. |
| Testbed is permanent, not throwaway | Survives Phase 3 as tuning/debugging/regression tool. Clean structure, not over-engineered. |
| tcell ruled out | Previously used and replaced with Bubble Tea. Not revisiting. |
| Phase 3 integration spike is throwaway | Separate experiment testing Bubble Tea mode dispatcher. The testbed and the spike serve different purposes. |

## Forward-Looking Notes

1. **Phase 3 bridges engine to dungeon.** `main.go` becomes a mode dispatcher. A wiring layer maps `world.Player.HP` <-> `engine.Player.HP`. Don't bake world types into the engine now.
2. **Phase 4 adds real enemy AI.** The Enemy struct will grow a State field and AI logic. VelX/VelY already exist for knockback, so AI movement is a natural extension.
3. **Phase 4 adds hit-stop.** Engine.Tick() will need a "freeze frames remaining" counter that skips simulation. Backward-compatible addition to Tick().
4. **Tuning constants is the whole point of Step 7.** Use Principle 7: "Double it or cut it in half." If jump feels weak, double JumpForce. If attack feels slow, halve AttackCooldown. All constants in one file. Key release timeout is also a tuning parameter.
5. **The testbed is permanent.** `cmd/combat-proto/` survives as a standalone tuning/debugging/regression tool. Clean structure, but don't over-engineer. Phase 3 builds a separate throwaway integration spike to test the Bubble Tea mode dispatcher.
6. **X-then-Y collision bias.** If platform-edge landing feels off during Step 7, the fix is to resolve the axis with the smaller penetration first. Start with the simpler X-then-Y approach; upgrade only if playtesting reveals the problem.
7. **Sound roadmap.** Phase 2 spike validates the audio pipeline with procedural synthesis. Phase 4 upgrades to `beep` library for WAV/OGG playback with designed sound assets. Target vibe: atmospheric retro (Celeste/Hyper Light Drifter), not orchestral (Hollow Knight). Dark, resonant, impactful within the pixel-art register.

## Verification

```sh
go test ./combat/engine/...      # all pass
go test ./combat/engine/... -v   # see individual names
go test ./combat/engine/... -cover  # target: >90%
go vet ./combat/engine/...       # no warnings
go build ./...                   # existing code unaffected
go test ./...                    # all existing tests still pass
go run ./cmd/combat-proto/       # playable testbed
```
