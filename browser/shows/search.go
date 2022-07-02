package shows

import (
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/tmdb"
	"os"
)

func (b Browser) showSearch() {
	s, defStyle := cli.SetupScreen()
	b.CLI.Screen = s
	b.CLI.Style = defStyle

	header := []string{
		"[Type in your Search | ENTER: Confirm | ESC, CTRL+C: Quit]",
		"Search TV shows:",
		"> " + b.Query,
	}

	text := cli.BuildScreen(cli.Page{}, b.Show.Index, header, []cli.Content{}, false)

	dim := cli.GetDimensions(s.Size())
	cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)

	for {
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			dim = cli.GetDimensions(s.Size())

			s.Clear()
			cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEnter {
				b.Search = []tmdb.Show{}
				b.Search, b.Show.Page = getSearchResults(b, 1, 10)

				s.Fini()
				err := b.browseShows()
				if err != nil {
					s.Fini()
					b.showSearch()
				}
			}
			if ev.Key() == tcell.KeyRune {
				b.Query += string(ev.Rune())
				header[2] = "> " + b.Query

				text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, []cli.Content{}, false)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
				inputTrim := len(b.Query)
				if inputTrim > 0 {
					b.Query = b.Query[:inputTrim-1]
				}
				header[2] = "> " + b.Query

				text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, []cli.Content{}, false)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
		}
	}
}
