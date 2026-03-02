package engine

import "testing"

func TestEngine_PlayerFallsToGround(t *testing.T) {
	e := NewEngine()
	for i := 0; i < 60; i++ {
		e.Tick(emptyInput())
		if e.Player.Grounded {
			return
		}
	}
	t.Error("player did not reach ground within 60 ticks")
}

func TestEngine_EnemyFallsToGround(t *testing.T) {
	e := NewEngine()
	for i := 0; i < 60; i++ {
		e.Tick(emptyInput())
		if e.Enemy.Grounded {
			return
		}
	}
	t.Error("enemy did not reach ground within 60 ticks")
}

func TestEngine_PlayerRunsRight(t *testing.T) {
	e := testEngine()
	startX := e.Player.Pos.X
	for i := 0; i < 5; i++ {
		e.Tick(InputState{Right: true})
	}
	if e.Player.Pos.X <= startX {
		t.Errorf("player X = %v, expected > %v after running right", e.Player.Pos.X, startX)
	}
}

func TestEngine_PlayerJumps(t *testing.T) {
	e := testEngine()
	groundY := e.Player.Pos.Y
	e.Tick(InputState{JumpPress: true, JumpHeld: true})
	if e.Player.VelY >= 0 {
		t.Errorf("VelY = %v, expected < 0 after jump", e.Player.VelY)
	}

	// Tick a few frames — player should rise.
	for i := 0; i < 5; i++ {
		e.Tick(InputState{JumpHeld: true})
	}
	if e.Player.Pos.Y >= groundY {
		t.Errorf("player Y = %v, expected < %v (should have risen)", e.Player.Pos.Y, groundY)
	}

	// Eventually falls back.
	for i := 0; i < 60; i++ {
		e.Tick(emptyInput())
		if e.Player.Grounded {
			return
		}
	}
	t.Error("player did not land after jumping")
}

func TestEngine_FullCombat_PlayerWins(t *testing.T) {
	e := testEngine()

	// Place player directly left of enemy so the attack hitbox will overlap.
	// Both on the floor (Y = FloorY - height).
	floorY := 112.0
	e.Player.Pos.X = e.Enemy.Pos.X - PlayerWidth - AttackOffsetX + AttackWidth/2
	e.Player.Pos.Y = floorY - PlayerHeight
	e.Player.Grounded = true
	e.Player.Facing = DirRight

	e.Enemy.Pos.Y = floorY - EnemyHeight
	e.Enemy.Grounded = true

	// Verify hitbox would overlap enemy from this position.
	hb := AttackHitbox(&Player{
		Pos:         e.Player.Pos,
		Facing:      DirRight,
		State:       StateAttack,
		AttackTimer: AttackDuration,
	})
	if !hb.Overlaps(e.Enemy.Pos) {
		t.Fatalf("setup error: hitbox %+v doesn't overlap enemy %+v", hb, e.Enemy.Pos)
	}

	// Attack EnemyHP times, waiting for invincibility between attacks.
	// Reposition player before each attack since knockback moves the enemy.
	for ticks := 0; ticks < 300; ticks++ {
		if e.Result != ResultNone {
			break
		}

		canAttack := e.Player.State != StateAttack &&
			e.Player.AttackCooldownTimer <= 0 &&
			e.Enemy.InvincTimer <= 0

		if canAttack && e.Enemy.Alive {
			// Reposition player next to enemy (knockback moves enemy each hit).
			e.Player.Pos.X = e.Enemy.Pos.X - PlayerWidth - AttackOffsetX + AttackWidth/2
			e.Tick(InputState{Attack: true})
		} else {
			e.Tick(emptyInput())
		}
	}

	if e.Result != ResultPlayerWin {
		t.Errorf("Result = %v, want ResultPlayerWin (enemy HP=%d, alive=%v)",
			e.Result, e.Enemy.HP, e.Enemy.Alive)
	}
}

func TestEngine_Reset(t *testing.T) {
	e := testEngine()
	// Modify state.
	e.Tick(InputState{Right: true})
	e.Tick(InputState{JumpPress: true, JumpHeld: true})
	e.TickCount = 100

	e.Reset()

	if e.TickCount != 0 {
		t.Errorf("TickCount = %d after reset, want 0", e.TickCount)
	}
	if e.Player.HP != PlayerHP {
		t.Errorf("Player HP = %d, want %d", e.Player.HP, PlayerHP)
	}
	if e.Enemy.HP != EnemyHP {
		t.Errorf("Enemy HP = %d, want %d", e.Enemy.HP, EnemyHP)
	}
	if !e.Enemy.Alive {
		t.Error("Enemy should be alive after reset")
	}
	if e.Result != ResultNone {
		t.Errorf("Result = %v, want ResultNone", e.Result)
	}
	assertNear(t, "Player.Pos.X", e.Player.Pos.X, 20, 0.1)
	assertNear(t, "Enemy.Pos.X", e.Enemy.Pos.X, 120, 0.1)
}

func TestEngine_NoUpdateAfterResult(t *testing.T) {
	e := testEngine()
	e.Result = ResultPlayerWin
	tick := e.TickCount
	e.Tick(InputState{Right: true})
	if e.TickCount != tick {
		t.Error("Tick should be a no-op after result is set")
	}
}
