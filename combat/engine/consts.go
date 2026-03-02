package engine

// PlayerState represents the player's current action state.
type PlayerState int

const (
	StateIdle   PlayerState = iota
	StateRun
	StateJump
	StateFall
	StateAttack
	StateHurt
)

// Dir represents a horizontal facing direction.
type Dir int

const (
	DirRight Dir = 1
	DirLeft  Dir = -1
)

// Result represents the outcome of a combat encounter.
type Result int

const (
	ResultNone       Result = iota
	ResultPlayerWin
	ResultPlayerDead
)

// Tick rate and delta time.
const (
	TickRate = 30
	DT       = 1.0 / float64(TickRate)
)

// Physics constants (pixels/sec or pixels/sec^2).
// Y increases downward. Positive VelY = falling.
const (
	Gravity      = 800.0
	MaxFallSpeed = 400.0
	JumpForce    = 300.0
	RunSpeed     = 200.0
)

// Timing windows (seconds).
const (
	CoyoteTime        = 0.1
	JumpBufferTime    = 0.1
	JumpCutMultiplier = 0.5
)

// Attack constants.
const (
	AttackDuration = 0.2
	AttackCooldown = 0.35
	AttackWidth    = 20.0
	AttackHeight   = 14.0
	AttackOffsetX  = 14.0
	AttackDamage   = 1
)

// Entity dimensions and HP.
const (
	PlayerWidth  = 12.0
	PlayerHeight = 20.0
	PlayerHP     = 5

	EnemyWidth  = 14.0
	EnemyHeight = 20.0
	EnemyHP     = 3
)

// Hurt / invincibility.
const (
	HurtDuration = 0.5
	InvincTime   = 1.0
	KnockbackVel = 150.0
)

// Arena dimensions (pixels).
const (
	ArenaWidth  = 160
	ArenaHeight = 120
)
