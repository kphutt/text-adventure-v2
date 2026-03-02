package engine

import "testing"

func TestApplyGravity_IncreasesVelocity(t *testing.T) {
	vel := 0.0
	applyGravity(&vel)
	want := Gravity * DT
	assertNear(t, "velY after one tick", vel, want, 0.01)
}

func TestApplyGravity_ClampsToMaxFallSpeed(t *testing.T) {
	vel := MaxFallSpeed - 1.0
	applyGravity(&vel)
	if vel > MaxFallSpeed {
		t.Errorf("velY = %v, should be clamped to MaxFallSpeed %v", vel, MaxFallSpeed)
	}
}

func TestApplyGravity_AlreadyAtMax(t *testing.T) {
	vel := MaxFallSpeed
	applyGravity(&vel)
	assertNear(t, "velY at max", vel, MaxFallSpeed, 0.01)
}

func TestMoveAndResolve_FallOntoFloor(t *testing.T) {
	floor := []Platform{{Rect{0, 100, 200, 10}}}
	// Start close enough that one tick of movement reaches the floor.
	// pos bottom edge = 79+20 = 99. After move: 99 + 200*DT = 105.67 -> overlaps floor.
	pos := Rect{50, 79, 12, 20}
	velX, velY := 0.0, 200.0
	grounded := false

	moveAndResolve(&pos, &velX, &velY, &grounded, floor)

	if !grounded {
		t.Error("expected grounded after landing")
	}
	assertNear(t, "pos.Y", pos.Y, 80.0, 0.1)
	assertNear(t, "velY", velY, 0.0, 0.01)
}

func TestMoveAndResolve_WallCollision(t *testing.T) {
	wall := []Platform{{Rect{100, 0, 10, 200}}}
	pos := Rect{85, 50, 12, 20}
	velX, velY := 200.0, 0.0
	grounded := false

	moveAndResolve(&pos, &velX, &velY, &grounded, wall)

	if pos.X+pos.W > 100.0+0.01 {
		t.Errorf("entity should be pushed out of wall: pos.X=%v, right edge=%v", pos.X, pos.X+pos.W)
	}
	assertNear(t, "velX", velX, 0.0, 0.01)
}

func TestMoveAndResolve_CeilingBonk(t *testing.T) {
	ceiling := []Platform{{Rect{0, 0, 200, 10}}}
	pos := Rect{50, 15, 12, 20}
	velX, velY := 0.0, -200.0
	grounded := false

	moveAndResolve(&pos, &velX, &velY, &grounded, ceiling)

	if pos.Y < 10.0-0.01 {
		t.Errorf("entity should be pushed below ceiling: pos.Y=%v", pos.Y)
	}
	assertNear(t, "velY", velY, 0.0, 0.01)
	if grounded {
		t.Error("should not be grounded after ceiling bonk")
	}
}

func TestMoveAndResolve_FreeMovement(t *testing.T) {
	platforms := []Platform{{Rect{0, 200, 200, 10}}} // far below
	pos := Rect{50, 50, 12, 20}
	velX, velY := 100.0, 50.0
	grounded := false

	startX, startY := pos.X, pos.Y
	moveAndResolve(&pos, &velX, &velY, &grounded, platforms)

	assertNear(t, "pos.X", pos.X, startX+100.0*DT, 0.1)
	assertNear(t, "pos.Y", pos.Y, startY+50.0*DT, 0.1)
	if grounded {
		t.Error("should not be grounded in free fall")
	}
}

func TestMoveAndResolve_DiagonalCorner(t *testing.T) {
	// Entity falling diagonally onto a platform corner.
	// Should land on top (grounded), not be pushed sideways.
	plat := []Platform{{Rect{50, 100, 40, 10}}}
	// Start just above and slightly to the left of the platform
	pos := Rect{46, 78, 12, 20}
	velX, velY := 30.0, 100.0
	grounded := false

	moveAndResolve(&pos, &velX, &velY, &grounded, plat)

	if !grounded {
		t.Error("expected entity to land on platform, not be pushed sideways")
	}
	assertNear(t, "pos.Y", pos.Y, 80.0, 0.1)
}
