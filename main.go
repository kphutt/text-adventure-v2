package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"text-adventure-v2/game"
)

func drawText(s tcell.Screen, x, y int, style tcell.Style, text string) {
	for i, r := range text {
		s.SetContent(x+i, y, r, nil, style)
	}
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	g := game.NewGame()
	inputStr := ""
	message := ""

	for {
		screen.Clear()

		if g.IsWon {
			drawText(screen, 0, 0, tcell.StyleDefault, "You win! Press any key to exit.")
			screen.Show()
			screen.PollEvent()
			return
		}

		y := 0
		for _, line := range strings.Split(g.Look(), "\n") {
			drawText(screen, 0, y, tcell.StyleDefault, line)
			y++
		}
		y = 10
		for _, line := range strings.Split(g.GetMapString(), "\n") {
			drawText(screen, 0, y, tcell.StyleDefault, line)
			y++
		}
		y++
		helpMessage, _ := g.HandleCommand("help")
		drawText(screen, 0, y, tcell.StyleDefault, helpMessage)
		drawText(screen, 0, 20, tcell.StyleDefault, "> "+inputStr)
		drawText(screen, 0, 15, tcell.StyleDefault, message)
		screen.Show()

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			var shouldExit bool
			var command string
			switch ev.Key() {
			case tcell.KeyEscape:
				return
			case tcell.KeyEnter:
				command = inputStr
				inputStr = ""
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(inputStr) > 0 {
					inputStr = inputStr[:len(inputStr)-1]
				}
			case tcell.KeyRune:
				command = string(ev.Rune())
			}

			if command != "" {
				message, shouldExit = g.HandleCommand(command)
				if shouldExit {
					return
				}
			}
		}
	}
}
