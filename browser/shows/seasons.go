package shows

import (
	"encoding/json"
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
	"math"
	"os"
)

func getSeason(conf config.Config, id int, season int) (tmdb.ShowSeason, error) {
	queries := tmdb.Queries{
		ApiKey:   conf.Library.Auth.APIKey,
		Language: conf.Library.Settings.Language,
	}
	searchObj := tmdb.ShowSeason{}

	searchRes, err := tmdb.GetTVShowSeason(id, season, queries)

	if err != nil {
		return searchObj, err
	}

	searchBody := json.NewDecoder(searchRes.Body)
	err = searchBody.Decode(&searchObj)
	if err != nil {
		return searchObj, err
	}

	return searchObj, nil
}

func (b Browser) browseSeasons() error {
	s, defStyle := cli.SetupScreen()
	b.CLI.Screen = s
	b.CLI.Style = defStyle

	b.Show.Season.Page.Current = 1
	b.Show.Season.Page.Results = 10

	header := []string{
		"[↓→↑←: Navigate | ENTER: Confirm | ESC: Back | CTRL+C: Quit]",
		"Browse Seasons:",
	}

	b.Show.Season.Page = getSeasonResults(
		b.Show.Season.Page, b.Show.Details.Seasons, b.Show.Season.Page.Current, b.Show.Season.Page.Results)

	text := cli.BuildScreen(b.Show.Season.Page, b.Show.Season.Index, header, b.Show.Season.Page.Content, true)

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
				err := b.browseShows()
				if err != nil {
					return err
				}
			}
			if ev.Key() == tcell.KeyEnter {
				currentSeasonIndex := getCurrentSeasonIndex(b.Show.Season, b.Show.Details.Seasons)

				season, err := getSeason(b.Config, b.Show.Details.ID, currentSeasonIndex)
				if err != nil {
					s.Fini()
					err := b.browseSeasons()
					if err != nil {
						return err
					}
				}

				b.Show.Season.Details = season

				s.Fini()
				err = b.browseEpisodes()
				if err != nil {
					s.Fini()
					err := b.browseSeasons()
					if err != nil {
						return err
					}
				}
			}
			if ev.Key() == tcell.KeyDown {
				if b.Show.Season.Index < len(b.Show.Season.Page.Content)-1 {
					b.Show.Season.Index++
				}

				b.Show.Season.Page = getSeasonResults(
					b.Show.Season.Page, b.Show.Details.Seasons, b.Show.Season.Page.Current, b.Show.Season.Page.Results)

				text = cli.BuildScreen(b.Show.Season.Page, b.Show.Season.Index, header, b.Show.Season.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyUp {
				if b.Show.Season.Index > 0 {
					b.Show.Season.Index--
				}

				b.Show.Season.Page = getSeasonResults(
					b.Show.Season.Page, b.Show.Details.Seasons, b.Show.Season.Page.Current, b.Show.Season.Page.Results)

				text = cli.BuildScreen(b.Show.Season.Page, b.Show.Season.Index, header, b.Show.Season.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyRight {
				b.Show.Season.Page = getSeasonResults(
					b.Show.Season.Page, b.Show.Details.Seasons, b.Show.Season.Page.Current+1, b.Show.Season.Page.Results)

				text = cli.BuildScreen(b.Show.Season.Page, b.Show.Season.Index, header, b.Show.Season.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyLeft {
				b.Show.Season.Page = getSeasonResults(
					b.Show.Season.Page, b.Show.Details.Seasons, b.Show.Season.Page.Current-1, b.Show.Season.Page.Results)

				text = cli.BuildScreen(b.Show.Season.Page, b.Show.Season.Index, header, b.Show.Season.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
		}
	}
}

func getSeasonResults(seasonPage cli.Page, seasons []tmdb.ShowSeasonDetail, page int, results int) cli.Page {
	startIndex := results * (page - 1)
	if startIndex < 0 || startIndex > len(seasons) {
		return seasonPage
	}
	endIndex := startIndex + results - 1
	if endIndex > len(seasons) {
		endIndex = len(seasons)
	}

	var content []cli.Content

	for i, season := range seasons {
		if i >= startIndex && i <= endIndex {
			content = append(content, cli.Content{
				Display: season.Name,
				ID:      season.ID,
			})
		}
	}

	maxTabs := math.Ceil(float64(len(seasons)) / float64(results))

	return cli.Page{
		Current: page,
		Total:   int(maxTabs),
		Results: results,
		Content: content,
	}
}

func getCurrentSeasonIndex(season Season, seasons []tmdb.ShowSeasonDetail) int {
	var index int

	for i, s := range seasons {
		if s.ID == season.Page.Content[season.Index].ID {
			index = i
		}
	}

	return index
}
