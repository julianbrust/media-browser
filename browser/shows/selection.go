package shows

import (
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"os"
)

func (b Browser) showSelection() error {
	s, defStyle := cli.SetupScreen()
	b.CLI.Screen = s
	b.CLI.Style = defStyle

	text := []string{
		"This is the top layer of the app",
		"You selected:",
	}
	cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)

	b.drawSelection()

	for {
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				err := b.browseEpisodes()
				if err != nil {
					return err
				}
			}
			if ev.Key() == tcell.KeyEnter {
				s.Fini()
				os.Exit(0)
			}
		}
	}
}

func (b Browser) drawSelection() {
	showName := b.Show.Details.Name
	seasonName := b.Show.Season.Details.Name
	epIndex := b.Show.Season.Episode.Index
	epDetails := b.Show.Season.Details.Episodes[epIndex]

	text := []string{
		showName + ": " + seasonName + ": " + epDetails.Name,
		"Description: " + epDetails.Overview,
	}

	cli.DrawText(b.CLI.Screen, 0, 2, 100, 100, b.CLI.Style, text)
}
