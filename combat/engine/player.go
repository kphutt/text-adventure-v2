package engine

// updatePlayer processes input and updates the player state machine.
// Called at step 1 of Tick — uses coyote/jump-buffer timers from the
// PREVIOUS frame (one-frame lag is intentional and standard).
func updatePlayer(p *Player, input InputState) {
	// Hurt blocks all input.
	if p.State == StateHurt {
		p.VelX = 0
		return
	}

	// Attack blocks movement but not state transitions out of attack.
	if p.State == StateAttack {
		p.VelX = 0
		if p.AttackTimer <= 0 {
			// Attack finished — transition out.
			if p.Grounded {
				p.State = StateIdle
			} else {
				p.State = StateFall
			}
			p.AttackCooldownTimer = AttackCooldown
			p.AttackHit = false
		}
		return
	}

	// --- Jump buffering ---
	if input.JumpPress {
		p.JumpBufferTimer = 0
	} else {
		p.JumpBufferTimer += DT
	}

	// --- Movement ---
	p.VelX = 0
	if input.Left {
		p.VelX = -RunSpeed
		p.Facing = DirLeft
	}
	if input.Right {
		p.VelX = RunSpeed
		p.Facing = DirRight
	}

	// --- Jump ---
	canJump := p.Grounded || p.CoyoteTimer < CoyoteTime
	wantsJump := p.JumpBufferTimer < JumpBufferTime
	if canJump && wantsJump {
		p.VelY = -JumpForce
		p.State = StateJump
		p.JumpCut = false
		p.JumpBufferTimer = JumpBufferTime // consume the buffer
		p.CoyoteTimer = CoyoteTime         // consume coyote time
		return
	}

	// --- Variable-height jump ---
	if !input.JumpHeld && p.VelY < 0 && !p.JumpCut {
		p.VelY *= JumpCutMultiplier
		p.JumpCut = true
	}

	// --- Attack ---
	if input.Attack && p.AttackCooldownTimer <= 0 {
		p.State = StateAttack
		p.AttackTimer = AttackDuration
		p.AttackHit = false
		p.VelX = 0
		return
	}

	// --- State from velocity/grounded ---
	if p.Grounded {
		if p.VelX != 0 {
			p.State = StateRun
		} else {
			p.State = StateIdle
		}
	} else {
		if p.VelY < 0 {
			p.State = StateJump
		} else {
			p.State = StateFall
		}
	}
}
