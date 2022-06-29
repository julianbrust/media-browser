package shows

import (
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
)

func (b Browser) browseEpisodes() error {
	s, defStyle := cli.SetupScreen()
	b.CLI.Screen = s
	b.CLI.Style = defStyle

	text := []string{
		"This is the top layer of the app",
		"Browse Episodes:",
	}
	cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)

	b.drawEpisodeResults()

	for {
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				err := b.browseSeasons()
				if err != nil {
					return err
				}
			}
			if ev.Key() == tcell.KeyEnter {
				s.Fini()
				b.showSelection()
			}
			if ev.Key() == tcell.KeyDown {
				s.Clear()
				text := []string{
					"This is the top layer of the app",
					"Browse Episodes:",
				}
				cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)
				b.Show.Season.Episode.Index = updateSelectionIndex(b.Show.Season.Episode.Index, true)
				b.drawEpisodeResults()
			}
			if ev.Key() == tcell.KeyUp {
				s.Clear()
				text := []string{
					"This is the top layer of the app",
					"Browse Episodes:",
				}
				cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)
				b.Show.Season.Episode.Index = updateSelectionIndex(b.Show.Season.Episode.Index, false)
				b.drawEpisodeResults()
			}
		}
	}
}

func (b Browser) drawEpisodeResults() {
	var text []string

	for i, episode := range b.Show.Season.Details.Episodes {
		if int32(i) == b.Show.Season.Episode.Index {
			text = append(text, "> "+episode.Name)
		} else {
			text = append(text, "  "+episode.Name)
		}
	}
	cli.DrawText(b.CLI.Screen, 0, 2, 100, 100, b.CLI.Style, text)
}
