package shows

import (
	"encoding/json"
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
)

func getShowResults(conf config.Config, query string) ([]tmdb.ShowResult, error) {
	queries := tmdb.Queries{
		ApiKey:       conf.Library.Auth.APIKey,
		Language:     conf.Library.Settings.Language,
		AdultContent: conf.Library.Settings.AdultContent,
		Query:        query,
		Page:         1,
	}
	searchRes, err := tmdb.SearchTV(queries)
	if err != nil {
		return nil, err
	}
	searchObj := tmdb.Show{}

	searchBody := json.NewDecoder(searchRes.Body)
	err = searchBody.Decode(&searchObj)
	if err != nil {
		return nil, err
	}

	return searchObj.Results, nil
}

func getShow(conf config.Config, id int32) (tmdb.ShowDetail, error) {
	queries := tmdb.Queries{
		ApiKey:   conf.Library.Auth.APIKey,
		Language: conf.Library.Settings.Language,
	}
	searchObj := tmdb.ShowDetail{}
	searchRes, err := tmdb.GetTVShow(id, queries)
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

func (b Browser) browseShows() error {
	s, defStyle := cli.SetupScreen()
	b.CLI.Screen = s
	b.CLI.Style = defStyle

	text := []string{
		"This is the top layer of the app",
		"Browse Shows:",
	}
	cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)

	b.drawShowResults()

	for {
		s.Show()

		ev := s.PollEvent()

		switch ev := ev.(type) {
		case *tcell.EventResize:
			s.Sync()
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyEscape || ev.Key() == tcell.KeyCtrlC {
				s.Fini()
				b.showSearch()
			}
			if ev.Key() == tcell.KeyEnter {
				show, err := getShow(b.Config, b.Search[b.Show.Index].ID)
				if err != nil {
					s.Fini()
					b.showSearch()
				}

				b.Show.Details = show

				s.Fini()
				err = b.browseSeasons()
				if err != nil {
					s.Fini()
					err := b.browseShows()
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
				b.Show.Index = updateSelectionIndex(b.Show.Index, true)
				b.drawShowResults()
			}
			if ev.Key() == tcell.KeyUp {
				s.Clear()
				text := []string{
					"This is the top layer of the app",
					"Browse Shows:",
				}
				cli.DrawText(b.CLI.Screen, 0, 0, 100, 100, b.CLI.Style, text)
				b.Show.Index = updateSelectionIndex(b.Show.Index, false)
				b.drawShowResults()
			}
		}
	}
}

func updateSelectionIndex(index int32, down bool) int32 {
	i := index
	if down {
		i++
	} else {
		i--
	}
	//TODO: number of results/pagination?
	if i >= 0 && i <= 6 {
		index = i
	}

	return index
}

func (b Browser) drawShowResults() {
	var text []string
	for i, result := range b.Search {
		if int32(i) == b.Show.Index {
			text = append(text, "> "+result.Name)
		} else {
			text = append(text, "  "+result.Name)
		}
	}
	cli.DrawText(b.CLI.Screen, 0, 2, 100, 100, b.CLI.Style, text)
}
