// Combat prototype testbed — standalone Bubble Tea app for tuning the
// combat engine. Wires combat/engine to pixelbuf for half-block rendering.
//
// Keys: A/D = move, Space = jump, F = attack, R = restart, Q/Esc = quit
package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"text-adventure-v2/combat/engine"
	"text-adventure-v2/pixelbuf"

	tea "charm.land/bubbletea/v2"
	"github.com/ebitengine/oto/v3"
)

// Colors.
var (
	bgColor    = pixelbuf.Color{R: 20, G: 20, B: 30, A: 255}
	platColor  = pixelbuf.Color{R: 80, G: 80, B: 100, A: 255}
	playerCol  = pixelbuf.Color{R: 100, G: 200, B: 255, A: 255}
	enemyCol   = pixelbuf.Color{R: 255, G: 80, B: 80, A: 255}
	attackCol  = pixelbuf.Color{R: 255, G: 255, B: 100, A: 255}
	hpFullCol  = pixelbuf.Color{R: 80, G: 220, B: 80, A: 255}
	hpEmptyCol = pixelbuf.Color{R: 80, G: 20, B: 20, A: 255}
	whiteCol   = pixelbuf.Color{R: 255, G: 255, B: 255, A: 255}
)

const (
	tickDuration = time.Second / time.Duration(engine.TickRate)
	hudRows      = 3 // rows reserved below the frame for HUD text

	// 100ms ≈ 3 ticks at 30fps. Compromise between coast on release
	// (reduced from 30px to 20px) and stutter on hold start (~150ms gap
	// before first auto-repeat).
	fallbackTimeout = 100 * time.Millisecond
)

// Key names as returned by KeyPressMsg.String() in bubbletea v2.
// "space" is the v2 representation — v1 used " ".
const (
	keyLeft  = "a"
	keyRight = "d"
	keyJump  = "space"
	keyAtk   = "f"
	keyReset = "r"
)

type tickMsg time.Time

func tick() tea.Cmd {
	return tea.Tick(tickDuration, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// --- Procedural sound ---

const sampleRate = 44100

var otoCtx *oto.Context

func initAudio() {
	op := &oto.NewContextOptions{
		SampleRate:   sampleRate,
		ChannelCount: 1,
		Format:       oto.FormatSignedInt16LE,
	}
	var ready chan struct{}
	var err error
	otoCtx, ready, err = oto.NewContext(op)
	if err != nil {
		return
	}
	<-ready
}

func playSound(samples []int16) {
	if otoCtx == nil {
		return
	}
	buf := &bytes.Buffer{}
	for _, s := range samples {
		binary.Write(buf, binary.LittleEndian, s)
	}
	p := otoCtx.NewPlayer(bytes.NewReader(buf.Bytes()))
	p.Play()
	go func() {
		for p.IsPlaying() {
			time.Sleep(5 * time.Millisecond)
		}
		p.Close()
	}()
}

func genHitSound() []int16 {
	dur := 0.05
	n := int(float64(sampleRate) * dur)
	samples := make([]int16, n)
	for i := range samples {
		t := float64(i) / float64(sampleRate)
		frac := t / dur
		freq := 800.0 - 600.0*frac
		sine := math.Sin(2.0 * math.Pi * freq * t)
		noise := (float64((i*1103515245+12345)%65536) / 32768.0) - 1.0
		mix := sine*0.6 + noise*0.4
		env := 1.0 - frac
		samples[i] = int16(mix * env * 16000)
	}
	return samples
}

func genJumpSound() []int16 {
	dur := 0.03
	n := int(float64(sampleRate) * dur)
	samples := make([]int16, n)
	for i := range samples {
		t := float64(i) / float64(sampleRate)
		frac := t / dur
		freq := 300.0 + 500.0*frac
		sine := math.Sin(2.0 * math.Pi * freq * t)
		env := 1.0 - frac*0.5
		samples[i] = int16(sine * env * 10000)
	}
	return samples
}

var hitSamples, jumpSamples []int16

func initSounds() {
	hitSamples = genHitSound()
	jumpSamples = genJumpSound()
}

// --- Model ---

type model struct {
	eng   *engine.Engine
	buf   *pixelbuf.Buffer
	scale float64 // render scale: engine pixels -> buffer pixels

	// Input mode (set once when KeyboardEnhancementsMsg arrives).
	hasKeyReleases bool

	// Press/release mode (KR) — used when hasKeyReleases is true:
	held    map[string]bool // keys currently down (set on press, cleared on release)
	pressed map[string]bool // keys newly pressed this tick (one-shot accumulator)

	// Fallback mode (FB) — used when hasKeyReleases is false:
	fallbackKeys map[string]time.Time // last-seen time per key
	prevHeld     map[string]bool      // previous tick's held snapshot (edge detection)

	width  int
	height int
	frame  string
}

func newModel() model {
	return model{
		eng:          engine.NewEngine(),
		held:         make(map[string]bool),
		pressed:      make(map[string]bool),
		prevHeld:     make(map[string]bool),
		fallbackKeys: make(map[string]time.Time),
		scale:        1.0,
	}
}

// resizeBuf recomputes the pixel buffer and scale to fit the terminal.
func (m *model) resizeBuf() {
	if m.width < 20 || m.height < 5 {
		return
	}
	// Available pixel space: full width, height minus HUD rows (half-block = 2 px/row).
	maxW := m.width
	maxH := (m.height - hudRows) * 2

	scaleX := float64(maxW) / float64(engine.ArenaWidth)
	scaleY := float64(maxH) / float64(engine.ArenaHeight)
	m.scale = min(scaleX, scaleY)

	bufW := int(float64(engine.ArenaWidth) * m.scale)
	bufH := int(float64(engine.ArenaHeight) * m.scale)
	// Half-block rendering needs even height.
	bufH = bufH &^ 1

	if m.buf == nil || m.buf.Width != bufW || m.buf.Height != bufH {
		m.buf = pixelbuf.NewBuffer(bufW, bufH)
	}
}

func (m model) Init() tea.Cmd {
	return tick()
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyboardEnhancementsMsg:
		m.hasKeyReleases = msg.SupportsEventTypes()
		return m, nil

	case tea.KeyPressMsg:
		key := msg.String()
		switch key {
		case "q", "esc", "ctrl+c":
			return m, tea.Quit
		case keyReset:
			m.eng.Reset()
			return m, nil
		}
		if m.hasKeyReleases {
			m.held[key] = true
			if !msg.IsRepeat {
				m.pressed[key] = true
			}
		} else {
			m.fallbackKeys[key] = time.Now()
		}
		return m, nil

	case tea.KeyReleaseMsg:
		delete(m.held, msg.String())
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.resizeBuf()
		return m, nil

	case tickMsg:
		now := time.Now()
		input := m.buildInput(now)

		prevState := m.eng.Player.State
		prevHP := m.eng.Enemy.HP

		m.eng.Tick(input)

		if m.eng.Player.State == engine.StateJump && prevState != engine.StateJump {
			playSound(jumpSamples)
		}
		if m.eng.Enemy.HP < prevHP {
			playSound(hitSamples)
		}

		m.renderFrame()
		return m, tick()
	}
	return m, nil
}

func (m *model) buildInput(now time.Time) engine.InputState {
	var input engine.InputState

	if m.hasKeyReleases {
		// --- Press/release mode (KR) ---
		// Continuous: read held map directly (no copy needed).
		input.Left = m.held[keyLeft]
		input.Right = m.held[keyRight]
		input.JumpHeld = m.held[keyJump]
		// One-shot: read pressed map, then clear it.
		input.JumpPress = m.pressed[keyJump]
		input.Attack = m.pressed[keyAtk]
		clear(m.pressed)
	} else {
		// --- Fallback mode (FB) ---
		// A key is "held" if seen within the timeout window.
		currentHeld := make(map[string]bool, len(m.fallbackKeys))
		for key, lastSeen := range m.fallbackKeys {
			if now.Sub(lastSeen) <= fallbackTimeout {
				currentHeld[key] = true
			}
		}
		input.Left = currentHeld[keyLeft]
		input.Right = currentHeld[keyRight]
		input.JumpHeld = currentHeld[keyJump]
		// Edge detection: held now but not last tick.
		input.JumpPress = currentHeld[keyJump] && !m.prevHeld[keyJump]
		input.Attack = currentHeld[keyAtk] && !m.prevHeld[keyAtk]
		m.prevHeld = currentHeld
	}

	return input
}

// s scales an engine coordinate to buffer pixels.
func (m *model) s(v float64) int {
	return int(v * m.scale)
}

func (m *model) renderFrame() {
	if m.buf == nil {
		return
	}
	e := m.eng
	buf := m.buf

	buf.Clear(bgColor)

	// Platforms.
	for _, p := range e.Platforms {
		r := p.Rect
		pixelbuf.FillRect(buf, m.s(r.X), m.s(r.Y), m.s(r.W), m.s(r.H), platColor)
	}

	// Enemy.
	if e.Enemy.Alive {
		col := enemyCol
		if e.Enemy.HurtTimer > 0 {
			if int(e.Enemy.HurtTimer*20)%2 == 0 {
				col = whiteCol
			}
		}
		pixelbuf.FillRect(buf, m.s(e.Enemy.Pos.X), m.s(e.Enemy.Pos.Y),
			m.s(e.Enemy.Pos.W), m.s(e.Enemy.Pos.H), col)
	}

	// Player.
	col := playerCol
	if e.Player.InvincTimer > 0 {
		if int(e.Player.InvincTimer*10)%2 == 0 {
			col = pixelbuf.Color{R: 40, G: 80, B: 100, A: 255}
		}
	}
	pixelbuf.FillRect(buf, m.s(e.Player.Pos.X), m.s(e.Player.Pos.Y),
		m.s(e.Player.Pos.W), m.s(e.Player.Pos.H), col)

	// Attack hitbox.
	hb := engine.AttackHitbox(&e.Player)
	if hb.W > 0 {
		pixelbuf.FillRect(buf, m.s(hb.X), m.s(hb.Y), m.s(hb.W), m.s(hb.H), attackCol)
	}

	// HP pips at top of screen.
	pipW := max(2, m.s(3))
	pipH := max(2, m.s(3))
	pipGap := max(1, m.s(1))
	drawHP(buf, m.s(6), m.s(2), pipW, pipH, pipGap, e.Player.HP, e.Player.MaxHP)
	if e.Enemy.Alive {
		enemyHPx := buf.Width - m.s(6) - e.Enemy.MaxHP*(pipW+pipGap)
		drawHP(buf, enemyHPx, m.s(2), pipW, pipH, pipGap, e.Enemy.HP, e.Enemy.MaxHP)
	}

	m.frame = pixelbuf.Render(buf)
}

func drawHP(buf *pixelbuf.Buffer, x, y, pipW, pipH, pipGap, hp, maxHP int) {
	for i := 0; i < maxHP; i++ {
		col := hpEmptyCol
		if i < hp {
			col = hpFullCol
		}
		pixelbuf.FillRect(buf, x+i*(pipW+pipGap), y, pipW, pipH, col)
	}
}

func (m model) View() tea.View {
	var content string
	if m.width < 20 || m.height < 5 {
		content = fmt.Sprintf("\n  Terminal too small (%dx%d). Please resize.\n", m.width, m.height)
	} else if m.buf == nil {
		content = "\n  Initializing...\n"
	} else {
		var sb strings.Builder
		sb.WriteString(m.frame)
		sb.WriteByte('\n')

		switch m.eng.Result {
		case engine.ResultPlayerWin:
			sb.WriteString("  VICTORY!  Press R to restart, Q to quit")
		case engine.ResultPlayerDead:
			sb.WriteString("  DEFEATED  Press R to restart, Q to quit")
		default:
			mode := "FB"
			if m.hasKeyReleases {
				mode = "KR"
			}
			sb.WriteString(fmt.Sprintf(
				"  A/D move  SPACE jump  F attack  R restart  Q quit  [%s]", mode))
		}
		sb.WriteByte('\n')

		content = sb.String()
	}

	v := tea.NewView(content)
	v.AltScreen = true
	v.KeyboardEnhancements.ReportEventTypes = true
	return v
}

func main() {
	initAudio()
	initSounds()

	p := tea.NewProgram(newModel())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
