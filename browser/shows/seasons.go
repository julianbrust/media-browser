package shows

import (
	"encoding/json"
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/tmdb"
	"math"
	"os"
)

// getSeason retrieves a parsed object for searching for a specific season.
// It requires config parameters, the id of the show and the number for the required season.
func (b *Browser) getSeason(id int, season int) (tmdb.ShowSeason, error) {
	b.Log.Traceln("starting getSeason")

	queries := tmdb.Queries{
		ApiKey:   b.Config.Library.Auth.APIKey,
		Language: b.Config.Library.Settings.Language,
	}
	searchObj := tmdb.ShowSeason{}

	searchRes, err := server.GetTVShowSeason(id, season, queries)

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

// browseSeasons starts and handles the CLI screen for browsing seasons.
func (b *Browser) browseSeasons() error {
	b.Log.Traceln("starting browseSeasons")

	s, defStyle := cli.SetupScreen()
	b.CLI.Screen = s
	b.CLI.Style = defStyle

	b.Show.Season.Page.Current = 1
	b.Show.Season.Page.Results = 10

	header := []string{
		"[↓→↑←: Navigate | ENTER: Confirm | ESC: Back | CTRL+C: Quit]",
		"Browse Seasons:",
	}

	b.getSeasonResults(b.Show.Season.Page.Current, b.Show.Season.Page.Results)

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
				b.Log.Info("exiting app from season screen with CTRL+C")

				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEscape {
				b.Log.Traceln("seasons: escape")

				s.Fini()
				err := b.browseShows()
				if err != nil {
					b.Log.Error(err)
					err := b.browseSeasons()
					if err != nil {
						b.Log.Fatal(err)
					}
				}
			}
			if ev.Key() == tcell.KeyEnter {
				b.Log.Traceln("select season")

				currentSeasonNumber := b.getCurrentSeasonNumber()

				season, err := b.getSeason(b.Show.Details.ID, currentSeasonNumber)
				if err != nil {
					s.Fini()
					err = b.browseSeasons()
					if err != nil {
						b.Log.Fatal(err)
					}
				}

				b.Show.Season.Details = season

				s.Fini()
				b.Show.Season.Episode.Index = 0

				err = b.browseEpisodes()
				if err != nil {
					b.Log.Error(err)
					err := b.browseSeasons()
					if err != nil {
						b.Log.Fatal(err)
					}
				}
			}
			if ev.Key() == tcell.KeyDown {
				b.Log.Traceln("seasons: key down")

				if b.Show.Season.Index < len(b.Show.Season.Page.Content)-1 {
					b.Show.Season.Index++
				}

				b.getSeasonResults(b.Show.Season.Page.Current, b.Show.Season.Page.Results)

				text = cli.BuildScreen(b.Show.Season.Page, b.Show.Season.Index, header, b.Show.Season.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyUp {
				b.Log.Traceln("seasons: key up")

				if b.Show.Season.Index > 0 {
					b.Show.Season.Index--
				}

				b.getSeasonResults(b.Show.Season.Page.Current, b.Show.Season.Page.Results)

				text = cli.BuildScreen(b.Show.Season.Page, b.Show.Season.Index, header, b.Show.Season.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyRight {
				b.Log.Traceln("seasons: key right")

				b.getSeasonResults(b.Show.Season.Page.Current+1, b.Show.Season.Page.Results)

				if b.Show.Season.Index > len(b.Show.Season.Page.Content)-1 {
					b.Show.Season.Index = len(b.Show.Season.Page.Content) - 1
				}

				text = cli.BuildScreen(b.Show.Season.Page, b.Show.Season.Index, header, b.Show.Season.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyLeft {
				b.Log.Traceln("seasons: key left")

				b.getSeasonResults(b.Show.Season.Page.Current-1, b.Show.Season.Page.Results)

				text = cli.BuildScreen(b.Show.Season.Page, b.Show.Season.Index, header, b.Show.Season.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
		}
	}
}

// getSeasonResults creates a new cli.Page based on the provided season data.
// It defines the data for the Page based on the requested page and amount of results to display.
func (b *Browser) getSeasonResults(page int, results int) {
	startIndex := results * (page - 1)
	if startIndex < 0 || startIndex >= len(b.Show.Details.Seasons) {
		return
	}
	endIndex := startIndex + results - 1
	if endIndex > len(b.Show.Details.Seasons) {
		endIndex = len(b.Show.Details.Seasons)
	}

	var content []cli.Content

	for i, season := range b.Show.Details.Seasons {
		if i >= startIndex && i <= endIndex {
			content = append(content, cli.Content{
				Display: season.Name,
				ID:      season.ID,
			})
		}
	}

	maxTabs := math.Ceil(float64(len(b.Show.Details.Seasons)) / float64(results))

	b.Show.Season.Page = cli.Page{
		Current: page,
		Total:   int(maxTabs),
		Results: results,
		Content: content,
	}
}

// getCurrentSeasonNumber returns the seasonNumber in the retrieved data of seasons based on the
// currently selected season.
func (b *Browser) getCurrentSeasonNumber() int {
	var seasonNumber int

	for _, s := range b.Show.Details.Seasons {
		if s.ID == b.Show.Season.Page.Content[b.Show.Season.Index].ID {
			seasonNumber = s.SeasonNumber
		}
	}

	return seasonNumber
}
