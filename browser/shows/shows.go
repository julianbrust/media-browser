package shows

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
	"math"
	"os"
)

// getShowResult retrieves a parsed object for searching shows with a query.
// It requires config parameters, a query string and the requested results page.
func getShowResult(conf config.Config, query string, page int) (tmdb.Show, error) {
	queries := tmdb.Queries{
		ApiKey:       conf.Library.Auth.APIKey,
		Language:     conf.Library.Settings.Language,
		AdultContent: conf.Library.Settings.AdultContent,
		Query:        query,
		Page:         page,
	}
	searchRes, err := tmdb.SearchTV(queries)
	if err != nil {
		return tmdb.Show{}, err
	}
	searchObj := tmdb.Show{}

	searchBody := json.NewDecoder(searchRes.Body)
	err = searchBody.Decode(&searchObj)
	if err != nil {
		return tmdb.Show{}, err
	}

	return searchObj, nil
}

// getShow retrieves a parsed object with details for a specific show.
// It requires config parameters and the ID of the show.
func getShow(conf config.Config, id int) (tmdb.ShowDetail, error) {
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

// browseShows starts and handles the CLI screen for browsing shows.
func (b Browser) browseShows() error {
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
				s.Fini()
				os.Exit(0)
			}
			if ev.Key() == tcell.KeyEscape {
				s.Fini()
				b.showSearch()
			}
			if ev.Key() == tcell.KeyEnter {
				show, err := getShow(b.Config, b.Show.Page.Content[b.Show.Index].ID)
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
				os.Exit(1)
			}
			if ev.Key() == tcell.KeyDown {
				if b.Show.Index < len(b.Show.Page.Content)-1 {
					b.Show.Index++
				}

				b.Search, b.Show.Page = getSearchResults(b, b.Show.Page.Current, b.Show.Page.Results)

				text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, b.Show.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyUp {
				if b.Show.Index > 0 {
					b.Show.Index--
				}

				b.Search, b.Show.Page = getSearchResults(b, b.Show.Page.Current, b.Show.Page.Results)

				text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, b.Show.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyRight {
				b.Search, b.Show.Page = getSearchResults(b, b.Show.Page.Current+1, b.Show.Page.Results)

				if b.Show.Index > len(b.Show.Page.Content)-1 {
					b.Show.Index = len(b.Show.Page.Content) - 1
				}

				text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, b.Show.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
			if ev.Key() == tcell.KeyLeft {
				b.Search, b.Show.Page = getSearchResults(b, b.Show.Page.Current-1, b.Show.Page.Results)

				text = cli.BuildScreen(b.Show.Page, b.Show.Index, header, b.Show.Page.Content, true)

				s.Clear()
				cli.DrawScreen(b.CLI.Screen, b.CLI.Style, dim, text)
			}
		}
	}
}

// getSearchResults retrieves all necessary data objects for a shows search and creates a cli.Page
// with the results for a specific page and amount of results.
func getSearchResults(b Browser, page int, results int) ([]tmdb.Show, cli.Page) {
	startIndex := results * (page - 1)
	if startIndex < 0 {
		return b.Search, b.Show.Page
	}

	endIndex := startIndex + results - 1

	b.Search = getMissingSearchData(b, endIndex)

	if startIndex > b.Search[0].TotalResults {
		return b.Search, b.Show.Page
	}
	if endIndex > b.Search[0].TotalResults {
		endIndex = b.Search[0].TotalResults
	}

	resPage := filterSelectedData(b.Search, page, results, startIndex, endIndex)

	resPage.Total = int(math.Ceil(float64(b.Search[0].TotalResults) / float64(results)))

	return b.Search, resPage
}

// getMissingSearchData retrieves additional data for a shows search based on the required endIndex
// of the objects to display.
func getMissingSearchData(b Browser, endIndex int) []tmdb.Show {
	if len(b.Search) > 0 && endIndex > b.Search[0].TotalResults {
		return b.Search
	}

	currentEndIndex := 0

	for currentEndIndex < endIndex {
		currentEndIndex = 0

		for _, entry := range b.Search {
			currentEndIndex += len(entry.Results) - 1
		}

		if currentEndIndex < endIndex {
			reqPage := len(b.Search) + 1
			newContent, err := getShowResult(b.Config, b.Query, reqPage)
			if err != nil {
				fmt.Println(err)
			}
			b.Search = append(b.Search, newContent)
		}

		for _, entry := range b.Search {
			currentEndIndex += len(entry.Results) - 1
		}
	}

	return b.Search
}

// filterSelectedData creates a new cli.Page based on the provided shows search data.
// It defines the data for the Page based on the requested page and amount of results to display,
// and the startIndex and endIndex of the available data.
func filterSelectedData(showSearch []tmdb.Show, page int, results int, startIndex int, endIndex int) cli.Page {
	var data []tmdb.ShowResult
	var selectedData []cli.Content

	for _, search := range showSearch {
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

	maxTabs := math.Round(float64(len(data)) / float64(results))

	resPage := cli.Page{
		Current: page,
		Total:   int(maxTabs),
		Results: results,
		Content: selectedData,
	}

	return resPage
}
