package engine

import "testing"

func groundedPlayer() Player {
	return Player{
		Pos:             Rect{50, 92, PlayerWidth, PlayerHeight},
		Facing:          DirRight,
		State:           StateIdle,
		HP:              PlayerHP,
		MaxHP:           PlayerHP,
		Grounded:        true,
		JumpBufferTimer: JumpBufferTime, // prevent false trigger
	}
}

func TestPlayer_IdleNoInput(t *testing.T) {
	p := groundedPlayer()
	updatePlayer(&p, InputState{})
	if p.State != StateIdle {
		t.Errorf("State = %v, want StateIdle", p.State)
	}
	assertNear(t, "VelX", p.VelX, 0, 0.01)
}

func TestPlayer_RunRight(t *testing.T) {
	p := groundedPlayer()
	updatePlayer(&p, InputState{Right: true})
	if p.State != StateRun {
		t.Errorf("State = %v, want StateRun", p.State)
	}
	assertNear(t, "VelX", p.VelX, RunSpeed, 0.01)
	if p.Facing != DirRight {
		t.Errorf("Facing = %v, want DirRight", p.Facing)
	}
}

func TestPlayer_RunLeft(t *testing.T) {
	p := groundedPlayer()
	updatePlayer(&p, InputState{Left: true})
	if p.State != StateRun {
		t.Errorf("State = %v, want StateRun", p.State)
	}
	assertNear(t, "VelX", p.VelX, -RunSpeed, 0.01)
	if p.Facing != DirLeft {
		t.Errorf("Facing = %v, want DirLeft", p.Facing)
	}
}

func TestPlayer_JumpFromGround(t *testing.T) {
	p := groundedPlayer()
	updatePlayer(&p, InputState{JumpPress: true, JumpHeld: true})
	if p.State != StateJump {
		t.Errorf("State = %v, want StateJump", p.State)
	}
	assertNear(t, "VelY", p.VelY, -JumpForce, 0.01)
}

func TestPlayer_CoyoteTime_JumpSucceeds(t *testing.T) {
	p := groundedPlayer()
	p.Grounded = false
	p.CoyoteTimer = CoyoteTime * 0.5 // within coyote window
	updatePlayer(&p, InputState{JumpPress: true, JumpHeld: true})
	if p.State != StateJump {
		t.Errorf("expected jump during coyote time, State = %v", p.State)
	}
	assertNear(t, "VelY", p.VelY, -JumpForce, 0.01)
}

func TestPlayer_CoyoteTime_JumpFails(t *testing.T) {
	p := groundedPlayer()
	p.Grounded = false
	p.CoyoteTimer = CoyoteTime + 0.01 // past coyote window
	updatePlayer(&p, InputState{JumpPress: true, JumpHeld: true})
	if p.State == StateJump {
		t.Error("should not jump after coyote time expired")
	}
}

func TestPlayer_JumpBuffer(t *testing.T) {
	p := groundedPlayer()
	p.Grounded = false
	p.CoyoteTimer = CoyoteTime + 1 // well past coyote

	// Press jump while airborne.
	updatePlayer(&p, InputState{JumpPress: true, JumpHeld: true})
	if p.State == StateJump {
		t.Fatal("should not jump while airborne past coyote")
	}

	// Now land (set grounded). JumpBufferTimer was reset to 0 by the press above.
	// It then incremented by DT, so it should still be < JumpBufferTime.
	p.Grounded = true
	updatePlayer(&p, InputState{JumpHeld: true}) // no new press, but buffer active
	if p.State != StateJump {
		t.Errorf("jump buffer should have triggered, State = %v, JumpBufferTimer = %v", p.State, p.JumpBufferTimer)
	}
}

func TestPlayer_VariableHeightJump_CutsVelocity(t *testing.T) {
	p := groundedPlayer()
	// Initiate jump.
	updatePlayer(&p, InputState{JumpPress: true, JumpHeld: true})
	if p.VelY >= 0 {
		t.Fatal("expected negative VelY after jump")
	}

	// Release jump while ascending.
	origVel := p.VelY
	p.Grounded = false
	updatePlayer(&p, InputState{JumpHeld: false})
	assertNear(t, "VelY after cut", p.VelY, origVel*JumpCutMultiplier, 0.01)
	if !p.JumpCut {
		t.Error("JumpCut should be true")
	}
}

func TestPlayer_VariableHeightJump_CutsOnlyOnce(t *testing.T) {
	p := groundedPlayer()
	updatePlayer(&p, InputState{JumpPress: true, JumpHeld: true})

	// Release to cut.
	p.Grounded = false
	updatePlayer(&p, InputState{JumpHeld: false})
	velAfterCut := p.VelY

	// Release again — should not cut further.
	updatePlayer(&p, InputState{JumpHeld: false})
	assertNear(t, "VelY should not change", p.VelY, velAfterCut, 0.01)
}

func TestPlayer_AttackInitiation(t *testing.T) {
	p := groundedPlayer()
	updatePlayer(&p, InputState{Attack: true})
	if p.State != StateAttack {
		t.Errorf("State = %v, want StateAttack", p.State)
	}
	assertNear(t, "AttackTimer", p.AttackTimer, AttackDuration, 0.01)
	assertNear(t, "VelX", p.VelX, 0, 0.01)
}

func TestPlayer_AttackCooldown(t *testing.T) {
	p := groundedPlayer()
	p.AttackCooldownTimer = 0.1 // cooldown active
	updatePlayer(&p, InputState{Attack: true})
	if p.State == StateAttack {
		t.Error("should not attack while cooldown is active")
	}
}

func TestPlayer_AttackEnds_TransitionsToIdle(t *testing.T) {
	p := groundedPlayer()
	p.State = StateAttack
	p.AttackTimer = 0 // attack just finished
	updatePlayer(&p, InputState{})
	if p.State != StateIdle {
		t.Errorf("State = %v, want StateIdle after attack ends on ground", p.State)
	}
	assertNear(t, "AttackCooldownTimer", p.AttackCooldownTimer, AttackCooldown, 0.01)
}

func TestPlayer_AttackEnds_TransitionsToFall(t *testing.T) {
	p := groundedPlayer()
	p.State = StateAttack
	p.AttackTimer = 0
	p.Grounded = false
	updatePlayer(&p, InputState{})
	if p.State != StateFall {
		t.Errorf("State = %v, want StateFall after attack ends in air", p.State)
	}
}

func TestPlayer_HurtBlocksInput(t *testing.T) {
	p := groundedPlayer()
	p.State = StateHurt
	p.HurtTimer = 0.3
	updatePlayer(&p, InputState{Right: true, JumpPress: true, Attack: true})
	assertNear(t, "VelX", p.VelX, 0, 0.01)
	if p.State != StateHurt {
		t.Errorf("State = %v, want StateHurt", p.State)
	}
}
