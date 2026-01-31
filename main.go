package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"text-adventure-v2/game"
	"text-adventure-v2/renderer"
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
		
		// Create the view model for the renderer
		mapView := renderer.MapView{
			AllRooms:       g.AllRooms,
			PlayerLocation: g.Player.Location,
		}
		mapString := renderer.RenderMap(mapView)

		for _, line := range strings.Split(mapString, "\n") {
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
			var shouldExit bool // This should be defined once
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
				// Only if not a special command (w,a,s,d,e,i,u,h,q)
				switch ev.Rune() {
				case 'w', 'a', 's', 'd', 'e', 'i', 'u', 'h', 'q':
					command = string(ev.Rune())
				default:
					inputStr += string(ev.Rune())
				}
			}

			if command != "" { // Process command if one was entered
				message, shouldExit = g.HandleCommand(command)
				if shouldExit {
					screen.Clear()
					drawText(screen, 0, 0, tcell.StyleDefault, message) // Use the message from g.HandleCommand
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
			}
		}
	}
}