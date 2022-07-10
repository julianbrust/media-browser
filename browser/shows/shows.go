package shows

import (
	"encoding/json"
	"errors"
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/tmdb"
	"math"
	"os"
)

// getShowResult retrieves a parsed object for searching shows with a query.
// It requires config parameters, a query string and the requested results page.
func (b *Browser) getShowResult(page int) (tmdb.Show, error) {
	b.Log.Traceln("start getShowResult")

	queries := tmdb.Queries{
		ApiKey:       b.Config.Library.Auth.APIKey,
		Language:     b.Config.Library.Settings.Language,
		AdultContent: b.Config.Library.Settings.AdultContent,
		Query:        b.Query,
		Page:         page,
	}
	searchRes, err := server.SearchTV(queries)
	if err != nil {
		return tmdb.Show{}, err
	}
	searchObj := tmdb.Show{}

	searchBody := json.NewDecoder(searchRes.Body)
	err = searchBody.Decode(&searchObj)
	if err != nil {
		return tmdb.Show{}, err
	}

	if len(searchObj.Results) == 0 {
		return searchObj, errors.New("no results")
	}

	return searchObj, nil
}

// getShow retrieves a parsed object with details for a specific show.
// It requires config parameters and the ID of the show.
func (b *Browser) getShow(id int) (tmdb.ShowDetail, error) {
	b.Log.Traceln("start getShow")

	queries := tmdb.Queries{
		ApiKey:   b.Config.Library.Auth.APIKey,
		Language: b.Config.Library.Settings.Language,
	}
	searchObj := tmdb.ShowDetail{}
	searchRes, err := server.GetTVShow(id, queries)
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

// browseShows starts and handles the CLI screen for browsing shows.
func (b *Browser) browseShows() error {
	b.Log.Traceln("starting browseShows")

	s, defStyle := cli.SetupScreen()
	b.CLI.Screen = s
	b.CLI.Style = defStyle

	b.Show.Page.Current = 1
	b.Show.Page.Results = 10

	header := []string{
		"[↓→↑←: Navigate | ENTER: Confirm | ESC: Back | CTRL+C: Quit]",
		"Browse Shows:",
	}

	text := cli.BuildScreen(b.Show.Page, b.Show.Index, header, b.Show.Page.Content, true)
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
				b.Log.Info("exiting app from shows screen with CTRL+C")

				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEscape {
				b.Log.Traceln("shows: escape")
				s.Fini()
				b.showSearch()
			}
			if ev.Key() == tcell.KeyEnter {
				if b.Show.Error {
					break
				}
				b.Log.Traceln("select show")

				show, err := b.getShow(b.Show.Page.Content[b.Show.Index].ID)
				if err != nil {
					s.Fini()
					b.showSearch()
				}

				b.Show.Details = show

				s.Fini()
				b.Show.Season.Index = 0

				err = b.browseSeasons()
				if err != nil {
					b.Log.Error(err)
					err := b.browseShows()
					if err != nil {
						b.Log.Fatal(err)
					}
				}
				os.Exit(1)
			}
			if ev.Key() == tcell.KeyDown {
				b.Log.Traceln("shows: key down")

				if b.Show.Index < len(b.Show.Page.Content)-1 {
					b.Show.Index++
				}

				var err error
				err = b.getSearchResults(b.Show.Page.Current, b.Show.Page.Results)
				if err != nil {
					b.Show.Error = true
					text = cli.BuildErrorScreen(header, err.Error())
				} else {
					b.Show.Error = false
					text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, b.Show.Page.Content, true)
				}

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyUp {
				b.Log.Traceln("shows: key up")

				if b.Show.Index > 0 {
					b.Show.Index--
				}

				var err error

				err = b.getSearchResults(b.Show.Page.Current, b.Show.Page.Results)
				if err != nil {
					b.Show.Error = true
					text = cli.BuildErrorScreen(header, err.Error())
				} else {
					b.Show.Error = false
					text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, b.Show.Page.Content, true)
				}

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyRight {
				b.Log.Traceln("shows: key right")

				var err error

				err = b.getSearchResults(b.Show.Page.Current+1, b.Show.Page.Results)
				if err != nil {
					b.Show.Error = true
					text = cli.BuildErrorScreen(header, err.Error())
				} else {
					b.Show.Error = false

					if b.Show.Index > len(b.Show.Page.Content)-1 {
						b.Show.Index = len(b.Show.Page.Content) - 1
					}

					text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, b.Show.Page.Content, true)
				}

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyLeft {
				b.Log.Traceln("shows: key left")

				var err error

				err = b.getSearchResults(b.Show.Page.Current-1, b.Show.Page.Results)
				if err != nil {
					b.Show.Error = true
					text = cli.BuildErrorScreen(header, err.Error())
				} else {
					b.Show.Error = false
					text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, b.Show.Page.Content, true)
				}

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
		}
	}
}

// getSearchResults retrieves all necessary data objects for a shows search and creates a cli.Page
// with the results for a specific page and amount of results.
func (b *Browser) getSearchResults(page int, results int) error {
	b.Log.Traceln("starting getSearchResults")

	startIndex := results * (page - 1)
	if startIndex < 0 {
		return nil
	}

	endIndex := startIndex + results - 1

	var err error
	b.Search, err = b.getMissingSearchData(endIndex)
	if err != nil {
		return err
	}

	if startIndex >= b.Search[0].TotalResults {
		return nil
	}
	if endIndex > b.Search[0].TotalResults {
		endIndex = b.Search[0].TotalResults
	}

	b.filterSelectedData(page, results, startIndex, endIndex)

	b.Show.Page.Total = int(math.Ceil(float64(b.Search[0].TotalResults) / float64(results)))

	return nil
}

// getMissingSearchData retrieves additional data for a shows search based on the required endIndex
// of the objects to display.
func (b *Browser) getMissingSearchData(endIndex int) ([]tmdb.Show, error) {
	b.Log.Traceln("starting getMissingSearchData")

	if len(b.Search) > 0 && endIndex > b.Search[0].TotalResults {
		return b.Search, nil
	}

	currentEndIndex := 0

	for currentEndIndex < endIndex {
		currentEndIndex = 0

		for _, entry := range b.Search {
			currentEndIndex += len(entry.Results) - 1
		}

		if currentEndIndex < endIndex {
			reqPage := len(b.Search) + 1
			newContent, err := b.getShowResult(reqPage)
			if err != nil {
				b.Log.Error(err)
				return nil, err
			}
			b.Search = append(b.Search, newContent)

			if newContent.TotalResults < endIndex {
				break
			}
		}

		for _, entry := range b.Search {
			currentEndIndex += len(entry.Results) - 1
		}
	}

	return b.Search, nil
}

// filterSelectedData creates a new cli.Page based on the provided shows search data.
// It defines the data for the Page based on the requested page and amount of results to display,
// and the startIndex and endIndex of the available data.
func (b *Browser) filterSelectedData(page int, results int, startIndex int, endIndex int) {
	b.Log.Traceln("starting filterSelectedData")

	var data []tmdb.ShowResult
	var selectedData []cli.Content

	for _, search := range b.Search {
		data = append(data, search.Results...)
	}

	for i, search := range data {
		if i >= startIndex && i <= endIndex {
			selectedData = append(selectedData, cli.Content{
				Display: search.Name,
				ID:      search.ID,
			})
		}
	}

	maxTabs := math.Ceil(float64(len(data)) / float64(results))

	b.Show.Page = cli.Page{
		Current: page,
		Total:   int(maxTabs),
		Results: results,
		Content: selectedData,
	}
}
