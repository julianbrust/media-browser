package shows

import (
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"os"
	"strconv"
	"time"
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

	text := b.getSelection(header)

	dim := cli.GetDimensions(s.Size())
	cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)

	for {
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				b.Log.Info("exiting app from selection screen with CTRL+C")

				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEscape {
				b.Log.Traceln("selection: escape")

				s.Fini()
				err := b.browseEpisodes()
				if err != nil {
					b.Log.Error(err)
					err := b.showSelection()
					if err != nil {
						b.Log.Fatal(err)
					}
				}
			}
			if ev.Key() == tcell.KeyEnter {
				b.Log.Info("exiting app from selection screen with ENTER")

				s.Fini()
				os.Exit(0)
			}
		}
	}
}

// getSelection combines all the information about the current selection in lines to display.
func (b Browser) getSelection(header []string) []string {
	b.Log.Traceln("start getSelection")

	showName := b.Show.Details.Name
	seasonName := b.Show.Season.Details.Name
	epDetails := b.Show.Season.Episode.Details

	text := header

	text = append(text, "  "+showName+": "+seasonName+": "+epDetails.Name)
	if epDetails.Overview != "" {
		text = append(text, "  "+"Description: "+epDetails.Overview)
	}
	if epDetails.AirDate != "" {
		t, err := time.Parse(YYYYMMDD, epDetails.AirDate)
		if err == nil {
			text = append(text, "  "+"Air Date: "+t.Format(TextDate))
		}
	}
	if epDetails.Runtime != 0 {
		text = append(text, "  "+"Runtime: "+strconv.FormatInt(int64(epDetails.Runtime), 10)+" min")
	}

	return text
}

const (
	YYYYMMDD = "2006-01-02"
	TextDate = "January 2, 2006"
)
