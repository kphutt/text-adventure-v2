package engine

import "testing"

func attackingPlayer(facing Dir) Player {
	return Player{
		Pos:         Rect{50, 92, PlayerWidth, PlayerHeight},
		Facing:      facing,
		State:       StateAttack,
		AttackTimer: AttackDuration,
		HP:          PlayerHP,
		MaxHP:       PlayerHP,
	}
}

func aliveEnemy(x float64) Enemy {
	return Enemy{
		Pos:    Rect{x, 92, EnemyWidth, EnemyHeight},
		HP:     EnemyHP,
		MaxHP:  EnemyHP,
		Facing: DirLeft,
		Alive:  true,
	}
}

func TestAttackHitbox_FacingRight(t *testing.T) {
	p := attackingPlayer(DirRight)
	hb := AttackHitbox(&p)
	if hb.W == 0 {
		t.Fatal("expected non-zero hitbox")
	}
	cx, _ := p.Pos.Center()
	// Hitbox should be to the right of player center.
	if hb.X < cx {
		t.Errorf("hitbox X=%v should be >= player center X=%v when facing right", hb.X, cx)
	}
}

func TestAttackHitbox_FacingLeft(t *testing.T) {
	p := attackingPlayer(DirLeft)
	hb := AttackHitbox(&p)
	if hb.W == 0 {
		t.Fatal("expected non-zero hitbox")
	}
	cx, _ := p.Pos.Center()
	// Hitbox right edge should be to the left of player center.
	if hb.X+hb.W > cx {
		t.Errorf("hitbox right edge=%v should be <= player center X=%v when facing left", hb.X+hb.W, cx)
	}
}

func TestAttackHitbox_NotAttacking(t *testing.T) {
	p := attackingPlayer(DirRight)
	p.State = StateIdle
	hb := AttackHitbox(&p)
	if hb.W != 0 || hb.H != 0 {
		t.Errorf("expected zero hitbox when not attacking, got %+v", hb)
	}
}

func TestAttackHitbox_TimerExpired(t *testing.T) {
	p := attackingPlayer(DirRight)
	p.AttackTimer = 0
	hb := AttackHitbox(&p)
	if hb.W != 0 {
		t.Error("expected zero hitbox when timer expired")
	}
}

func TestProcessAttack_HitsEnemy(t *testing.T) {
	p := attackingPlayer(DirRight)
	hb := AttackHitbox(&p)
	// Place enemy overlapping with hitbox.
	e := aliveEnemy(hb.X)
	processAttack(&p, &e)
	if e.HP != EnemyHP-AttackDamage {
		t.Errorf("enemy HP = %d, want %d", e.HP, EnemyHP-AttackDamage)
	}
	assertNear(t, "HurtTimer", e.HurtTimer, HurtDuration, 0.01)
	assertNear(t, "InvincTimer", e.InvincTimer, InvincTime, 0.01)
	if !p.AttackHit {
		t.Error("AttackHit should be true")
	}
}

func TestProcessAttack_EnemyInvincible(t *testing.T) {
	p := attackingPlayer(DirRight)
	hb := AttackHitbox(&p)
	e := aliveEnemy(hb.X)
	e.InvincTimer = 0.5
	processAttack(&p, &e)
	if e.HP != EnemyHP {
		t.Errorf("enemy HP = %d, should be unchanged at %d", e.HP, EnemyHP)
	}
}

func TestProcessAttack_MultiHitPrevention(t *testing.T) {
	p := attackingPlayer(DirRight)
	hb := AttackHitbox(&p)
	e := aliveEnemy(hb.X)

	processAttack(&p, &e)
	hpAfterFirst := e.HP
	e.InvincTimer = 0 // clear invincibility to isolate AttackHit check
	processAttack(&p, &e)
	if e.HP != hpAfterFirst {
		t.Errorf("enemy HP = %d, expected no change due to AttackHit = %d", e.HP, hpAfterFirst)
	}
}

func TestProcessAttack_EnemyDies(t *testing.T) {
	p := attackingPlayer(DirRight)
	hb := AttackHitbox(&p)
	e := aliveEnemy(hb.X)
	e.HP = 1
	processAttack(&p, &e)
	if e.HP != 0 {
		t.Errorf("enemy HP = %d, want 0", e.HP)
	}
	if e.Alive {
		t.Error("enemy should be dead")
	}
}

func TestProcessAttack_KnockbackDirection_Right(t *testing.T) {
	p := attackingPlayer(DirRight)
	hb := AttackHitbox(&p)
	e := aliveEnemy(hb.X) // enemy to the right of player
	processAttack(&p, &e)
	if e.VelX <= 0 {
		t.Errorf("enemy VelX = %v, expected positive (pushed right)", e.VelX)
	}
}

func TestProcessAttack_KnockbackDirection_Left(t *testing.T) {
	p := attackingPlayer(DirLeft)
	hb := AttackHitbox(&p)
	// Place enemy to the left of player.
	e := aliveEnemy(hb.X)
	processAttack(&p, &e)
	if e.VelX >= 0 {
		t.Errorf("enemy VelX = %v, expected negative (pushed left)", e.VelX)
	}
}

func TestProcessAttack_KnockbackWall(t *testing.T) {
	// Enemy near right wall. After knockback + moveAndResolve, enemy
	// should stop at the wall, not pass through.
	p := attackingPlayer(DirRight)
	p.Pos.X = 130

	rightWall := Platform{Rect{156, 0, 4, 120}}
	floor := Platform{Rect{0, 112, 160, 8}}

	e := Enemy{
		Pos:    Rect{148, 92, EnemyWidth, EnemyHeight},
		HP:     EnemyHP,
		MaxHP:  EnemyHP,
		Facing: DirLeft,
		Alive:  true,
	}

	// Manually set attack state so hitbox overlaps enemy.
	p.State = StateAttack
	p.AttackTimer = AttackDuration
	p.Facing = DirRight

	// Verify hitbox actually overlaps.
	hb := AttackHitbox(&p)
	if !hb.Overlaps(e.Pos) {
		t.Skipf("hitbox %+v doesn't overlap enemy %+v — adjust positions", hb, e.Pos)
	}

	processAttack(&p, &e)
	if e.VelX <= 0 {
		t.Fatalf("expected positive knockback, got VelX=%v", e.VelX)
	}

	// Simulate collision resolution.
	platforms := []Platform{rightWall, floor}
	grounded := e.Grounded
	moveAndResolve(&e.Pos, &e.VelX, &e.VelY, &grounded, platforms)

	rightEdge := e.Pos.X + e.Pos.W
	if rightEdge > 156+0.1 {
		t.Errorf("enemy right edge = %v, should not exceed wall at 156", rightEdge)
	}
}

func TestProcessAttack_DeadEnemy_NoHit(t *testing.T) {
	p := attackingPlayer(DirRight)
	hb := AttackHitbox(&p)
	e := aliveEnemy(hb.X)
	e.Alive = false
	e.HP = 0
	processAttack(&p, &e)
	if p.AttackHit {
		t.Error("should not hit dead enemy")
	}
}
