package engine

// Engine is the deterministic combat simulation. It has no UI dependencies
// and imports only stdlib.
type Engine struct {
	Player    Player
	Enemy     Enemy
	Platforms []Platform
	Result    Result
	TickCount int
}

// NewEngine creates an engine with the standard arena, player, and enemy.
func NewEngine() *Engine {
	e := &Engine{}
	e.init()
	return e
}

func (e *Engine) init() {
	e.Player = Player{
		Pos:             Rect{20, 0, PlayerWidth, PlayerHeight},
		Facing:          DirRight,
		State:           StateIdle,
		HP:              PlayerHP,
		MaxHP:           PlayerHP,
		JumpBufferTimer: JumpBufferTime, // prevent false first-frame trigger
	}
	e.Enemy = Enemy{
		Pos:   Rect{120, 0, EnemyWidth, EnemyHeight},
		HP:    EnemyHP,
		MaxHP: EnemyHP,
		Facing: DirLeft,
		Alive: true,
	}
	e.Platforms = []Platform{
		{Rect{0, 112, 160, 8}},    // Floor
		{Rect{0, 0, 4, 120}},      // Left wall
		{Rect{156, 0, 4, 120}},    // Right wall
		{Rect{20, 84, 36, 6}},     // Platform 1
		{Rect{62, 64, 36, 6}},     // Platform 2
		{Rect{104, 84, 36, 6}},    // Platform 3
		{Rect{70, 44, 20, 6}},     // Platform 4
	}
	e.Result = ResultNone
	e.TickCount = 0
}

// Reset restores the engine to its initial state.
func (e *Engine) Reset() {
	e.init()
}

// Tick advances the simulation by one frame.
func (e *Engine) Tick(input InputState) {
	if e.Result != ResultNone {
		return
	}
	e.TickCount++

	// 1. Player input + state machine.
	updatePlayer(&e.Player, input)

	// 2. Apply gravity to player and enemy.
	applyGravity(&e.Player.VelY)
	if e.Enemy.Alive {
		applyGravity(&e.Enemy.VelY)
	}

	// 3. Move and resolve collisions.
	moveAndResolve(&e.Player.Pos, &e.Player.VelX, &e.Player.VelY,
		&e.Player.Grounded, e.Platforms)
	if e.Enemy.Alive {
		moveAndResolve(&e.Enemy.Pos, &e.Enemy.VelX, &e.Enemy.VelY,
			&e.Enemy.Grounded, e.Platforms)
	}

	// 4. Update coyote timer.
	if e.Player.Grounded {
		e.Player.CoyoteTimer = 0
	} else {
		e.Player.CoyoteTimer += DT
	}

	// 5. Process attack hitbox / damage.
	processAttack(&e.Player, &e.Enemy)

	// 6. Decrement enemy timers.
	e.Enemy.HurtTimer = max(0, e.Enemy.HurtTimer-DT)
	e.Enemy.InvincTimer = max(0, e.Enemy.InvincTimer-DT)

	// 7. Decrement player timers.
	e.Player.AttackTimer = max(0, e.Player.AttackTimer-DT)
	e.Player.AttackCooldownTimer = max(0, e.Player.AttackCooldownTimer-DT)
	e.Player.HurtTimer = max(0, e.Player.HurtTimer-DT)
	e.Player.InvincTimer = max(0, e.Player.InvincTimer-DT)

	// Transition out of hurt when timer expires.
	if e.Player.State == StateHurt && e.Player.HurtTimer <= 0 {
		if e.Player.Grounded {
			e.Player.State = StateIdle
		} else {
			e.Player.State = StateFall
		}
	}

	// 8. Check win/lose conditions (enemy death first).
	if !e.Enemy.Alive {
		e.Result = ResultPlayerWin
	} else if e.Player.HP <= 0 {
		e.Result = ResultPlayerDead
	}
}
