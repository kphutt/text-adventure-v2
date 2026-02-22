package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"text-adventure-v2/game"
	"text-adventure-v2/renderer"
)

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

type model struct {
	game      *game.Game
	textInput textinput.Model
	message   string
	won       bool
}

func initialModel() model {
	ti := textinput.New()
	ti.Prompt = "> "
	ti.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("14")) // cyan prompt
	ti.Focus()

	return model{
		game:      game.NewGame(),
		textInput: ti,
	}
}

func (m model) Init() tea.Cmd {
	return textinput.Blink
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
				response, shouldExit := m.game.HandleCommand(command)
				m.message = response
				m.textInput.Reset()
				if shouldExit {
					m.won = true
				}
			}
			return m, nil

		default:
			// Instant commands when input is empty
			if m.textInput.Value() == "" {
				key := msg.String()
				switch key {
				case "w", "a", "s", "d", "e", "i", "u", "h", "q":
					response, shouldExit := m.game.HandleCommand(key)
					m.message = response
					if shouldExit {
						m.won = true
					}
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
		return winStyle.Render(fmt.Sprintf(
			"%s\n\nTotal Turns: %d\nFinal Score: %d\n\nPress any key to exit.",
			m.message, m.game.Turns, m.game.Score(),
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
	msgStr := msgStyle.Render(m.highlightItems(m.message))
	inputStr := m.textInput.View()

	helpStr := helpStyle.Render(helpText)

	return lipgloss.JoinVertical(lipgloss.Left, hudStr, helpStr, mapStr, lookStr, msgStr, inputStr)
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
