package shows

import (
	"encoding/json"
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
)

func getSeason(conf config.Config, id int32, season int32) (tmdb.ShowSeason, error) {
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

	text := []string{
		"This is the top layer of the app",
		"Browse Seasons:",
	}
	cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)

	b.drawSeasonResults()

	for {
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				err := b.browseShows()
				if err != nil {
					return err
				}
			}
			if ev.Key() == tcell.KeyEnter {
				season, err := getSeason(b.Config, b.Show.Details.ID, b.Show.Season.Index)
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
				s.Clear()
				text := []string{
					"This is the top layer of the app",
					"Browse Shows:",
				}
				cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)
				b.Show.Season.Index = updateSelectionIndex(b.Show.Season.Index, true)
				b.drawSeasonResults()
			}
			if ev.Key() == tcell.KeyUp {
				s.Clear()
				text := []string{
					"This is the top layer of the app",
					"Browse Shows:",
				}
				cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)
				b.Show.Season.Index = updateSelectionIndex(b.Show.Season.Index, false)
				b.drawSeasonResults()
			}
		}
	}
}

func (b Browser) drawSeasonResults() {
	var text []string
	for i, season := range b.Show.Details.Seasons {
		if int32(i) == b.Show.Season.Index {
			text = append(text, "> "+season.Name)
		} else {
			text = append(text, "  "+season.Name)
		}
	}
	cli.DrawText(b.CLI.Screen, 0, 2, 100, 100, b.CLI.Style, text)
}
