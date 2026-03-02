package engine

// Rect is an axis-aligned bounding box.
type Rect struct {
	X, Y, W, H float64
}

// Overlaps reports whether a and b have a non-zero area intersection.
// Adjacent rects (sharing an edge) do NOT overlap.
// Zero-size rects never overlap anything.
func (a Rect) Overlaps(b Rect) bool {
	if a.W <= 0 || a.H <= 0 || b.W <= 0 || b.H <= 0 {
		return false
	}
	return a.X < b.X+b.W && a.X+a.W > b.X &&
		a.Y < b.Y+b.H && a.Y+a.H > b.Y
}

// Center returns the center point of the rect.
func (a Rect) Center() (float64, float64) {
	return a.X + a.W/2, a.Y + a.H/2
}

// InputState is a per-frame input snapshot. The engine has no UI knowledge.
type InputState struct {
	Left      bool // movement held
	Right     bool // movement held
	JumpPress bool // jump pressed THIS frame (edge-detected)
	JumpHeld  bool // jump button currently held
	Attack    bool // attack pressed THIS frame
}

// Player is the player entity.
type Player struct {
	Pos    Rect
	VelX   float64
	VelY   float64
	Facing Dir
	State  PlayerState
	HP     int
	MaxHP  int

	Grounded bool

	CoyoteTimer         float64 // counts UP from 0; time since last grounded
	JumpBufferTimer     float64 // counts UP from 0; time since last JumpPress
	JumpCut             bool    // true once variable-height cut applied this jump
	AttackTimer         float64 // counts down during attack
	AttackCooldownTimer float64 // counts down between attacks
	HurtTimer           float64
	InvincTimer         float64
	AttackHit           bool // prevents multi-hit per swing
}

// Enemy is a minimal Phase 2 enemy: stands there, takes hits, dies.
// Has physics for gravity (initial fall) and knockback.
type Enemy struct {
	Pos         Rect
	VelX        float64
	VelY        float64
	Grounded    bool
	HP          int
	MaxHP       int
	Facing      Dir
	HurtTimer   float64
	InvincTimer float64
	Alive       bool
}

// Platform is a solid collidable surface.
type Platform struct {
	Rect Rect
}
