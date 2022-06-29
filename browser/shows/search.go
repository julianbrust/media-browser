package shows

import (
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"os"
)

func (b Browser) showSearch() {
	s, defStyle := cli.SetupScreen()
	b.CLI.Screen = s
	b.CLI.Style = defStyle

	text := []string{
		"This is the top layer of the app",
		"Search TV shows:",
		"> " + b.Query,
	}
	cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)

	for {
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEnter {
				results, err := getShowResults(b.Config, b.Query)
				if err != nil {
					s.Fini()
					b.showSearch()
				}
				b.Search = results

				s.Fini()
				err = b.browseShows()
				if err != nil {
					s.Fini()
					b.showSearch()
				}
			}
			if ev.Key() == tcell.KeyRune {
				b.Query += string(ev.Rune())
				text[2] = "> " + b.Query
				cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)
			}
			if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
				s.Clear()

				inputTrim := len(b.Query)
				if inputTrim > 0 {
					b.Query = b.Query[:inputTrim-1]
				}
				text[2] = "> " + b.Query
				cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)
			}
		}
	}
}
