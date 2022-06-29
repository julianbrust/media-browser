package cli

import (
	"github.com/gdamore/tcell/v2"
	"log"
)

func SetupScreen() (tcell.Screen, tcell.Style) {
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}

	// Set default text style
	defStyle := tcell.StyleDefault.Background(tcell.ColorReset).Foreground(tcell.ColorReset)
	s.SetStyle(defStyle)

	// Clear screen
	s.Clear()

	return s, defStyle
}

func DrawText(screen tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text []string) {
	row := y1
	col := x1
	for _, line := range text {
		for _, r := range []rune(line) {
			screen.SetContent(col, row, r, nil, style)
			col++
			if col >= x2 {
				row++
				col = x1
			}
			if row > y2 {
				break
			}
		}
		col = x1
		row++
	}
}
