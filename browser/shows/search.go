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
		"This is the top layer of the app",
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
			s.Clear()
			dim = cli.GetDimensions(s.Size())
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
				s.Clear()

				b.Query += string(ev.Rune())
				header[2] = "> " + b.Query

				text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, []cli.Content{}, false)
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
				s.Clear()

				inputTrim := len(b.Query)
				if inputTrim > 0 {
					b.Query = b.Query[:inputTrim-1]
				}
				header[2] = "> " + b.Query

				text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, []cli.Content{}, false)
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
		}
	}
}
