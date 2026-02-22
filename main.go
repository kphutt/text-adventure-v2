package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"text-adventure-v2/game"
	"text-adventure-v2/renderer"
)

var debugMode = flag.Bool("debug", false, "enable debug logging to debug.log")

var (
	hudStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("39"))                                           // blue
	mapStyle  = lipgloss.NewStyle().Border(lipgloss.RoundedBorder()).Padding(0, 1).BorderForeground(lipgloss.Color("63")) // purple border
	helpStyle = lipgloss.NewStyle().Faint(true)
	lookStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))                                                                                                // soft white
	msgStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))                                                                                                // pink
	itemStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("220"))                                                                                     // gold
	winStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("10")).Border(lipgloss.DoubleBorder()).Padding(1, 3).BorderForeground(lipgloss.Color("11")) // green + gold border
)

const helpText = "w,a,s,d: move | e: take | u: unlock | i: inventory | h: help | q: quit"
const maxLogLines = 5

type model struct {
	game      *game.Game
	textInput textinput.Model
	messages  []string
	won       bool
	debugLog  *log.Logger
}

func initialModel() model {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14")) // cyan prompt
	ti.Focus()

	g := game.NewGame()

	m := model{
		game:      g,
		textInput: ti,
	}

	if *debugMode {
		f, err := os.Create("debug.log")
		if err == nil {
			m.debugLog = log.New(f, "", log.LstdFlags)
			logStartupState(m.debugLog, g)
		}
	}

	return m
}

func logStartupState(dl *log.Logger, g *game.Game) {
	// Log all rooms and their connections
	names := make([]string, 0, len(g.AllRooms))
	for name := range g.AllRooms {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		room := g.AllRooms[name]
		var conns []string
		for dir, exit := range room.Exits {
			conns = append(conns, fmt.Sprintf("%s: %s", dir, exit.Room.Name))
		}
		sort.Strings(conns)
		dl.Printf("[MAP] %s -> %s", name, strings.Join(conns, ", "))

		// Log items in this room
		for _, item := range room.Items {
			dl.Printf("[ITEM] %s in %s", item.Name, name)
		}

		// Log locked exits
		for dir, exit := range room.Exits {
			if exit.Locked {
				dl.Printf("[LOCK] %s -> %s (locked)", name, dir)
			}
		}
	}

	dl.Printf("[START] %s", g.Player.Location.Name)
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
}

func (m *model) handleCommand(command string) {
	scoreBefore := m.game.Score()
	roomBefore := m.game.Player.Location.Name

	response, shouldExit := m.game.HandleCommand(command)

	if m.debugLog != nil {
		m.debugLog.Printf("[CMD] t=%d room=%q cmd=%q", m.game.Turns, roomBefore, command)
		if response != "" {
			m.debugLog.Printf("[RSP] %s", response)
		} else {
			m.debugLog.Printf("[RSP] (empty — moved to %q)", m.game.Player.Location.Name)
		}
		if newScore := m.game.Score(); newScore != scoreBefore {
			m.debugLog.Printf("[SCORE] %d -> %d", scoreBefore, newScore)
		}
	}

	if response != "" {
		m.messages = append(m.messages, response)
	}
	if shouldExit {
		m.won = true
		if m.debugLog != nil {
			m.debugLog.Printf("[WIN] Player won in %d turns with score %d", m.game.Turns, m.game.Score())
		}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if m.won {
			return m, tea.Quit
		}

		switch msg.Type {
		case tea.KeyCtrlC, tea.KeyEsc:
			return m, tea.Quit

		case tea.KeyEnter:
			command := m.textInput.Value()
			if command != "" {
				m.handleCommand(command)
				m.textInput.Reset()
			}
			return m, nil

		default:
			// Instant commands when input is empty
			if m.textInput.Value() == "" {
				key := msg.String()
				switch key {
				case "w", "a", "s", "d", "e", "i", "u", "h", "q":
					m.handleCommand(key)
					return m, nil
				}
			}
		}
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m model) highlightItems(s string) string {
	for _, item := range m.game.Player.Location.Items {
		s = strings.ReplaceAll(s, item.Name, itemStyle.Render(item.Name))
	}
	for _, item := range m.game.Player.Inventory {
		s = strings.ReplaceAll(s, item.Name, itemStyle.Render(item.Name))
	}
	return s
}

func (m model) View() string {
	if m.won {
		lastMsg := ""
		if len(m.messages) > 0 {
			lastMsg = m.messages[len(m.messages)-1]
		}
		return winStyle.Render(fmt.Sprintf(
			"%s\n\nTotal Turns: %d\nFinal Score: %d\n\nPress any key to exit.",
			lastMsg, m.game.Turns, m.game.Score(),
		))
	}

	mapView := renderer.MapView{
		AllRooms:            m.game.AllRooms,
		PlayerLocation:      m.game.Player.Location,
		CurrentLocationName: m.game.Player.Location.Name,
		TurnsTaken:          m.game.Turns,
		Score:               m.game.Score(),
	}

	hudStr := hudStyle.Render(renderer.RenderHUD(mapView))
	mapStr := mapStyle.Render(renderer.RenderMap(mapView))
	lookStr := lookStyle.Render(m.highlightItems(m.game.Look()))
	// Show last N messages as a scrolling log
	start := 0
	if len(m.messages) > maxLogLines {
		start = len(m.messages) - maxLogLines
	}
	logLines := make([]string, 0, maxLogLines)
	for _, msg := range m.messages[start:] {
		logLines = append(logLines, m.highlightItems(msg))
	}
	msgStr := msgStyle.Render(strings.Join(logLines, "\n"))
	inputStr := m.textInput.View()

	helpStr := helpStyle.Render(helpText)

	return lipgloss.JoinVertical(lipgloss.Left, hudStr, helpStr, mapStr, lookStr, msgStr, inputStr)
}

func main() {
	flag.Parse()
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
