package engine

// applyGravity accelerates velY downward and clamps to MaxFallSpeed.
func applyGravity(velY *float64) {
	*velY += Gravity * DT
	if *velY > MaxFallSpeed {
		*velY = MaxFallSpeed
	}
}

// moveAndResolve moves an entity by its velocity and resolves collisions
// against all platforms. Uses split-axis resolution: move X then resolve,
// move Y then resolve.
func moveAndResolve(pos *Rect, velX, velY *float64, grounded *bool, platforms []Platform) {
	// --- X axis ---
	pos.X += *velX * DT
	for _, p := range platforms {
		if !pos.Overlaps(p.Rect) {
			continue
		}
		if *velX > 0 {
			// Moving right — push left
			pos.X = p.Rect.X - pos.W
		} else if *velX < 0 {
			// Moving left — push right
			pos.X = p.Rect.X + p.Rect.W
		}
		*velX = 0
	}

	// --- Y axis ---
	pos.Y += *velY * DT
	*grounded = false
	for _, p := range platforms {
		if !pos.Overlaps(p.Rect) {
			continue
		}
		if *velY >= 0 {
			// Falling or stationary — land on top
			pos.Y = p.Rect.Y - pos.H
			*grounded = true
		} else {
			// Rising — bonk ceiling
			pos.Y = p.Rect.Y + p.Rect.H
		}
		*velY = 0
	}
}
