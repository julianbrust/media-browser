package shows

import (
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"math"
	"os"
)

// browseEpisodes starts and handles the CLI screen for browsing episodes.
func (b Browser) browseEpisodes() error {
	b.Log.Traceln("starting browseEpisodes")

	s, defStyle := cli.SetupScreen()
	b.CLI.Screen = s
	b.CLI.Style = defStyle

	b.Show.Season.Episode.Page.Current = 1
	b.Show.Season.Episode.Page.Results = 10

	header := []string{
		"[↓→↑←: Navigate | ENTER: Confirm | ESC: Back | CTRL+C: Quit]",
		"Browse Episodes:",
	}
	b.Show.Season.Episode.Page = b.getEpisodeResults(b.Show.Season.Episode.Page.Current, b.Show.Season.Episode.Page.Results)

	text := cli.BuildScreen(
		b.Show.Season.Episode.Page, b.Show.Season.Episode.Index, header, b.Show.Season.Episode.Page.Content, true)

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
			if ev.Key() == tcell.KeyCtrlC {
				b.Log.Info("exiting app from episodes screen with CTRL+C")

				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEscape {
				b.Log.Traceln("episodes: escape")

				s.Fini()
				err := b.browseSeasons()
				if err != nil {
					b.Log.Error(err)
					err := b.browseEpisodes()
					if err != nil {
						b.Log.Fatal(err)
					}
				}
			}
			if ev.Key() == tcell.KeyEnter {
				b.Log.Traceln("select episode")

				b.Show.Season.Episode = b.getCurrentEpisode()

				s.Fini()
				err := b.showSelection()
				if err != nil {
					b.Log.Error(err)
					err := b.browseEpisodes()
					if err != nil {
						b.Log.Fatal(err)
					}
				}
			}
			if ev.Key() == tcell.KeyDown {
				b.Log.Traceln("episodes: key down")

				if b.Show.Season.Episode.Index < len(b.Show.Season.Episode.Page.Content)-1 {
					b.Show.Season.Episode.Index++
				}

				b.Show.Season.Episode.Page = b.getEpisodeResults(
					b.Show.Season.Episode.Page.Current, b.Show.Season.Episode.Page.Results)

				text = cli.BuildScreen(
					b.Show.Season.Episode.Page, b.Show.Season.Episode.Index, header, b.Show.Season.Episode.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyUp {
				b.Log.Traceln("episodes: key up")

				if b.Show.Season.Episode.Index > 0 {
					b.Show.Season.Episode.Index--
				}

				b.Show.Season.Episode.Page = b.getEpisodeResults(
					b.Show.Season.Episode.Page.Current, b.Show.Season.Episode.Page.Results)

				text = cli.BuildScreen(
					b.Show.Season.Episode.Page, b.Show.Season.Episode.Index, header, b.Show.Season.Episode.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyRight {
				b.Log.Traceln("episodes: key right")

				b.Show.Season.Episode.Page = b.getEpisodeResults(
					b.Show.Season.Episode.Page.Current+1, b.Show.Season.Episode.Page.Results)

				if b.Show.Season.Episode.Index > len(b.Show.Season.Episode.Page.Content)-1 {
					b.Show.Season.Episode.Index = len(b.Show.Season.Episode.Page.Content) - 1
				}

				text = cli.BuildScreen(
					b.Show.Season.Episode.Page, b.Show.Season.Episode.Index, header, b.Show.Season.Episode.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyLeft {
				b.Log.Traceln("episodes: key left")

				b.Show.Season.Episode.Page = b.getEpisodeResults(
					b.Show.Season.Episode.Page.Current-1, b.Show.Season.Episode.Page.Results)

				text = cli.BuildScreen(
					b.Show.Season.Episode.Page, b.Show.Season.Episode.Index, header, b.Show.Season.Episode.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
		}
	}
}

// getEpisodeResults provides the requested selection of results in a cli.Page. It updates the
// requested data of the previous epPage with data from all available episodes based on the requested
// page and number of results per page.
func (b Browser) getEpisodeResults(page int, results int) cli.Page {
	b.Log.Traceln("start getEpisodeResults")

	startIndex := results * (page - 1)
	if startIndex < 0 || startIndex > len(b.Show.Season.Details.Episodes) {
		return b.Show.Season.Episode.Page
	}
	endIndex := startIndex + results - 1
	if endIndex > len(b.Show.Season.Details.Episodes) {
		endIndex = len(b.Show.Season.Details.Episodes)
	}

	var content []cli.Content

	for i, season := range b.Show.Season.Details.Episodes {
		if i >= startIndex && i <= endIndex {
			content = append(content, cli.Content{
				Display: season.Name,
				ID:      season.ID,
			})
		}
	}

	maxTabs := math.Ceil(float64(len(b.Show.Season.Details.Episodes)) / float64(results))

	return cli.Page{
		Current: page,
		Total:   int(maxTabs),
		Results: results,
		Content: content,
	}
}

// getCurrentEpisode updates the current Episode details from the requested episode in the list of
// available episodes.
func (b Browser) getCurrentEpisode() Episode {
	b.Log.Traceln("start getCurrentEpisode")

	for _, ep := range b.Show.Season.Details.Episodes {
		if ep.ID == b.Show.Season.Episode.Page.Content[b.Show.Season.Episode.Index].ID {
			b.Show.Season.Episode.Details = ep
		}
	}

	return b.Show.Season.Episode
}
