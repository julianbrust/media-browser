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

	header := []string{
		"This is the top layer of the app",
		"You selected:",
	}

	text := getSelection(b.Show, header)

	dim := cli.GetDimensions(s.Size())
	cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)

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

func getSelection(show Show, header []string) []string {
	showName := show.Details.Name
	seasonName := show.Season.Details.Name
	epDetails := show.Season.Episode.Details

	text := header

	text = append(text, showName+": "+seasonName+": "+epDetails.Name)
	text = append(text, "Description: "+epDetails.Overview)

	return text
}
