package shows

import (
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/tmdb"
	"os"
)

// showSearch starts and handles the CLI screen for typing a search.
func (b *Browser) showSearch() {
	b.Log.Traceln("starting showSearch")

	s, defStyle := cli.SetupScreen()
	b.CLI.Screen = s
	b.CLI.Style = defStyle

	header := []string{
		"[Type in your Search | ENTER: Confirm | ESC, CTRL+C: Quit]",
		"Search TV shows:",
		"  > " + b.Query,
	}

	text := cli.BuildScreen(cli.Page{}, b.Show.Index, header, []cli.Content{}, false)

	dim := cli.GetDimensions(s.Size())
	cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)

	s.ShowCursor(len(header[len(header)-1]), len(header)-1)

	b.Log.Traceln("finished drawing initial search screen")

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
				b.Log.Info("exiting app from search screen")

				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEnter {
				b.Log.Tracef("start search with: %v", b.Query)

				b.Search = []tmdb.Show{}
				var err error

				err = b.getSearchResults(1, 10)
				if err != nil {
					b.Log.Trace(err)
					text = cli.BuildErrorScreen(header, err.Error())

					s.Clear()
					cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
				} else {
					s.Fini()
					b.Show.Index = 0

					err = b.browseShows()
					if err != nil {
						s.Fini()
						b.showSearch()
					}
				}
			}
			if ev.Key() == tcell.KeyRune {
				if len(b.Query) < 50 {
					b.Query += string(ev.Rune())
					b.Log.Tracef("updated search text: %v", b.Query)

					header[2] = "  > " + b.Query

					text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, []cli.Content{}, false)

					s.Clear()
					cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)

					s.ShowCursor(len(header[len(header)-1]), len(header)-1)
				} else {
					b.Log.Traceln("query length limit reached")
				}
			}
			if ev.Key() == tcell.KeyBackspace || ev.Key() == tcell.KeyBackspace2 {
				inputTrim := len(b.Query)
				if inputTrim > 0 {
					b.Query = b.Query[:inputTrim-1]
				}
				b.Log.Tracef("updated search text: %v", b.Query)

				header[2] = "  > " + b.Query

				text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, []cli.Content{}, false)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)

				s.ShowCursor(len(header[len(header)-1]), len(header)-1)
			}
		}
	}
}
