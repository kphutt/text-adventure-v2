package engine

import (
	"math"
	"testing"
)

// testEngine creates an engine and ticks until both player and enemy
// are grounded. Panics if it takes more than 60 ticks.
func testEngine() *Engine {
	e := NewEngine()
	for i := 0; i < 60; i++ {
		e.Tick(InputState{})
		if e.Player.Grounded && e.Enemy.Grounded {
			return e
		}
	}
	panic("testEngine: player/enemy not grounded after 60 ticks")
}

// emptyInput returns a zero-valued InputState.
func emptyInput() InputState {
	return InputState{}
}

// assertNear fails if |got - want| > epsilon.
func assertNear(t *testing.T, name string, got, want, epsilon float64) {
	t.Helper()
	if math.Abs(got-want) > epsilon {
		t.Errorf("%s = %v, want ~%v (epsilon %v)", name, got, want, epsilon)
	}
}
