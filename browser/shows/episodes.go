package shows

import (
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/tmdb"
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
	b.Show.Season.Episode.Page = getEpisodeResults(
		b.Show.Season.Episode.Page, b.Show.Season.Details.Episodes, b.Show.Season.Episode.Page.Current, b.Show.Season.Episode.Page.Results)

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
				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				err := b.browseSeasons()
				if err != nil {
					return err
				}
			}
			if ev.Key() == tcell.KeyEnter {
				b.Show.Season.Episode = getCurrentEpisode(b.Show.Season.Episode, b.Show.Season.Details.Episodes)

				s.Fini()
				err := b.showSelection()
				if err != nil {
					return err
				}
			}
			if ev.Key() == tcell.KeyDown {
				if b.Show.Season.Episode.Index < len(b.Show.Season.Episode.Page.Content)-1 {
					b.Show.Season.Episode.Index++
				}

				b.Show.Season.Episode.Page = getEpisodeResults(
					b.Show.Season.Episode.Page, b.Show.Season.Details.Episodes, b.Show.Season.Episode.Page.Current, b.Show.Season.Episode.Page.Results)

				text = cli.BuildScreen(
					b.Show.Season.Episode.Page, b.Show.Season.Episode.Index, header, b.Show.Season.Episode.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyUp {
				if b.Show.Season.Episode.Index > 0 {
					b.Show.Season.Episode.Index--
				}

				b.Show.Season.Episode.Page = getEpisodeResults(
					b.Show.Season.Episode.Page, b.Show.Season.Details.Episodes, b.Show.Season.Episode.Page.Current, b.Show.Season.Episode.Page.Results)

				text = cli.BuildScreen(
					b.Show.Season.Episode.Page, b.Show.Season.Episode.Index, header, b.Show.Season.Episode.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyRight {
				b.Show.Season.Episode.Page = getEpisodeResults(
					b.Show.Season.Episode.Page, b.Show.Season.Details.Episodes, b.Show.Season.Episode.Page.Current+1, b.Show.Season.Episode.Page.Results)

				if b.Show.Season.Episode.Index > len(b.Show.Season.Episode.Page.Content)-1 {
					b.Show.Season.Episode.Index = len(b.Show.Season.Episode.Page.Content) - 1
				}

				text = cli.BuildScreen(
					b.Show.Season.Episode.Page, b.Show.Season.Episode.Index, header, b.Show.Season.Episode.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyLeft {
				b.Show.Season.Episode.Page = getEpisodeResults(
					b.Show.Season.Episode.Page, b.Show.Season.Details.Episodes, b.Show.Season.Episode.Page.Current-1, b.Show.Season.Episode.Page.Results)

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
func getEpisodeResults(epPage cli.Page, episodes []tmdb.ShowEpisode, page int, results int) cli.Page {
	startIndex := results * (page - 1)
	if startIndex < 0 || startIndex > len(episodes) {
		return epPage
	}
	endIndex := startIndex + results - 1
	if endIndex > len(episodes) {
		endIndex = len(episodes)
	}

	var content []cli.Content

	for i, season := range episodes {
		if i >= startIndex && i <= endIndex {
			content = append(content, cli.Content{
				Display: season.Name,
				ID:      season.ID,
			})
		}
	}

	maxTabs := math.Ceil(float64(len(episodes)) / float64(results))

	return cli.Page{
		Current: page,
		Total:   int(maxTabs),
		Results: results,
		Content: content,
	}
}

// getCurrentEpisode updates the current Episode details from the requested episode in the list of
// available episodes.
func getCurrentEpisode(episode Episode, episodes []tmdb.ShowEpisode) Episode {
	for _, ep := range episodes {
		if ep.ID == episode.Page.Content[episode.Index].ID {
			episode.Details = ep
		}
	}

	return episode
}
