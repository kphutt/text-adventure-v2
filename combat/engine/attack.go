package engine

// AttackHitbox returns the active attack hitbox for the player, or a zero
// Rect if the player is not attacking. Exported for testbed rendering.
func AttackHitbox(p *Player) Rect {
	if p.State != StateAttack || p.AttackTimer <= 0 {
		return Rect{}
	}
	cx, cy := p.Pos.Center()
	hb := Rect{
		W: AttackWidth,
		H: AttackHeight,
		Y: cy - AttackHeight/2,
	}
	if p.Facing == DirRight {
		hb.X = cx + AttackOffsetX - AttackWidth/2
	} else {
		hb.X = cx - AttackOffsetX - AttackWidth/2
	}
	return hb
}

// processAttack checks the attack hitbox against the enemy and applies
// damage + knockback on hit.
func processAttack(p *Player, e *Enemy) {
	if !e.Alive {
		return
	}
	hb := AttackHitbox(p)
	if hb.W == 0 {
		return
	}
	if p.AttackHit {
		return
	}
	if e.InvincTimer > 0 {
		return
	}
	if !hb.Overlaps(e.Pos) {
		return
	}

	// Hit confirmed.
	e.HP -= AttackDamage
	if e.HP < 0 {
		e.HP = 0
	}
	e.HurtTimer = HurtDuration
	e.InvincTimer = InvincTime

	// Knockback: push enemy away from player.
	pcx, _ := p.Pos.Center()
	ecx, _ := e.Pos.Center()
	if ecx >= pcx {
		e.VelX = KnockbackVel
	} else {
		e.VelX = -KnockbackVel
	}

	p.AttackHit = true

	if e.HP <= 0 {
		e.Alive = false
	}
}
