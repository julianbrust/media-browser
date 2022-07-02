package shows

import (
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"os"
)

// showSelection starts and handles the CLI screen for displaying the current selection.
func (b Browser) showSelection() error {
	b.Log.Traceln("starting showSelection")

	s, defStyle := cli.SetupScreen()
	b.CLI.Screen = s
	b.CLI.Style = defStyle

	header := []string{
		"[ENTER, CTRL+C: Quit | ESC: Back]",
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

// getSelection combines all the information about the current selection in lines to display.
func getSelection(show Show, header []string) []string {
	showName := show.Details.Name
	seasonName := show.Season.Details.Name
	epDetails := show.Season.Episode.Details

	text := header

	text = append(text, showName+": "+seasonName+": "+epDetails.Name)
	text = append(text, "Description: "+epDetails.Overview)

	return text
}
