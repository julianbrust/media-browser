package shows

import (
	"encoding/json"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSeason(t *testing.T) {
	bytes, err := ioutil.ReadFile("../../testing/tmdb/get-show-season.json")
	if err != nil {
		log.Fatal(err)
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(bytes)
	}))

	defer ts.Close()

	server.BaseURL = ts.URL

	b := Browser{
		Config: &config.Config{},
		Query:  "breaking bad",
		Log:    logrus.New(),
	}
	b.Config.Library.Auth.APIKey = "string"
	b.Config.Library.Settings.Language = "en-US"

	season, err := b.getSeason(1396, 1)
	if err != nil {
		t.Error(err)
	}

	if season.ID != 3572 {
		t.Errorf("Expected season ID '3572', got %v", season.ID)
	}
}

func TestGetSeasonResults(t *testing.T) {
	bytes, err := ioutil.ReadFile("../../testing/tmdb/get-show-details.json")
	if err != nil {
		log.Fatal(err)
	}

	var season tmdb.ShowDetail

	err = json.Unmarshal(bytes, &season)
	if err != nil {
		log.Fatal(err)
	}

	b := Browser{
		Config: &config.Config{},
		Log:    logrus.New(),
	}
	b.Show.Details.Seasons = season.Seasons

	page := b.getSeasonResults(1, 3)

	if page.Current != 1 {
		t.Errorf("Expected season result index '1', got %v", page.Current)
	}
	if len(page.Content) != 3 {
		t.Errorf("Expected content amount '3', got %v", len(page.Content))
	}
}

func TestGetCurrentSeasonNumber(t *testing.T) {
	bytes, err := ioutil.ReadFile("../../testing/tmdb/get-show-details.json")
	if err != nil {
		log.Fatal(err)
	}

	var season tmdb.ShowDetail

	err = json.Unmarshal(bytes, &season)
	if err != nil {
		log.Fatal(err)
	}

	b := Browser{
		Config: &config.Config{},
		Log:    logrus.New(),
	}
	b.Show.Details.Seasons = season.Seasons
	b.Show.Season.Index = 1
	b.Show.Season.Page = cli.Page{
		Content: []cli.Content{
			{
				Display: "Season 1",
				ID:      3572,
			},
			{
				Display: "Season 2",
				ID:      3573,
			},
		},
	}

	num := b.getCurrentSeasonNumber()

	if num != 2 {
		t.Errorf("Expected season number '1', got %v", num)
	}
}
