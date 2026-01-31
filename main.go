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
		y := 0

		helpMessage, _ := g.HandleCommand("help")
		for _, line := range strings.Split(helpMessage, "\n") {
			drawText(screen, 0, y, tcell.StyleDefault, line)
			y++
		}

		y++
		for _, line := range strings.Split(g.GetMapString(), "\n") {
			drawText(screen, 0, y, tcell.StyleDefault, line)
			y++
		}
		
		y++
		for _, line := range strings.Split(g.Look(), "\n") {
			drawText(screen, 0, y, tcell.StyleDefault, line)
			y++
		}
		
		drawText(screen, 0, y, tcell.StyleDefault, message)
		y++
		drawText(screen, 0, y, tcell.StyleDefault, "> "+inputStr)
		screen.Show()

		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			var shouldExit bool
			switch ev.Key() {
			case tcell.KeyEscape:
				return
			case tcell.KeyEnter:
				message, shouldExit = g.HandleCommand(inputStr)
				if shouldExit {
					screen.Clear()
					drawText(screen, 0, 0, tcell.StyleDefault, message)
					drawText(screen, 0, 1, tcell.StyleDefault, "Press any key to exit.")
					screen.Show()
					// wait for any key press
					for {
						ev := screen.PollEvent()
						if _, ok := ev.(*tcell.EventKey); ok {
							return
						}
					}
				}
				inputStr = ""
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(inputStr) > 0 {
					inputStr = inputStr[:len(inputStr)-1]
				}
			case tcell.KeyRune:
				var shouldExit bool
				switch ev.Rune() {
				case 'w', 'a', 's', 'd', 'e', 'i', 'u', 'h', 'q':
					message, shouldExit = g.HandleCommand(string(ev.Rune()))
					if shouldExit {
						screen.Clear()
						drawText(screen, 0, 0, tcell.StyleDefault, message)
						drawText(screen, 0, 1, tcell.StyleDefault, "Press any key to exit.")
						screen.Show()
						// wait for any key press
						for {
							ev := screen.PollEvent()
							if _, ok := ev.(*tcell.EventKey); ok {
								return
							}
						}
					}
				default:
					inputStr += string(ev.Rune())
				}
			}
		}
	}
}
