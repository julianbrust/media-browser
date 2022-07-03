package shows

import (
	"encoding/json"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	Config config.Config
)

func TestGetSearchResults(t *testing.T) {
	bytes, err := ioutil.ReadFile("../../testing/tmdb/get-tv.json")
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
		Config: &Config,
		Query:  "breaking bad",
		Log:    logrus.New(),
	}
	b.Config.Library.Auth.APIKey = "string"
	b.Config.Library.Settings.AdultContent = false
	b.Config.Library.Settings.Language = "en-US"

	shows, page, err := b.getSearchResults(1, 5)
	if err != nil {
		t.Error(err)
	}

	if len(shows) > 0 && len(shows[0].Results) > 0 {
		if shows[0].Results[0].ID != 1396 {
			t.Errorf("Expected show ID '1396', got %v", shows[0].Results[0].ID)
		}
	} else {
		t.Error("Expected show result, got none")
	}
	if page.Current != 1 {
		t.Errorf("Expected current page '1', got %v", page.Current)
	}
	if page.Total != 1 {
		t.Errorf("Expected total pages '1', got %v", page.Total)
	}
}

func TestGetMissingSearchData(t *testing.T) {
	bytes, err := ioutil.ReadFile("../../testing/tmdb/get-tv.json")
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
		Config: &Config,
		Query:  "breaking bad",
		Log:    logrus.New(),
	}
	b.Config.Library.Auth.APIKey = "string"
	b.Config.Library.Settings.AdultContent = false
	b.Config.Library.Settings.Language = "en-US"

	shows, page, err := b.getSearchResults(1, 5)
	if err != nil {
		t.Error(err)
	}

	if len(shows) > 0 && len(shows[0].Results) > 0 {
		if shows[0].Results[0].ID != 1396 {
			t.Errorf("Expected show ID '1396', got %v", shows[0].Results[0].ID)
		}
	} else {
		t.Error("Expected show result, got none")
	}
	if page.Current != 1 {
		t.Errorf("Expected current page '1', got %v", page.Current)
	}
	if page.Total != 1 {
		t.Errorf("Expected total pages '1', got %v", page.Total)
	}
}

func TestFilterSelectedData(t *testing.T) {
	bytes, err := ioutil.ReadFile("../../testing/tmdb/get-tv.json")
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
		Config: &Config,
		Query:  "breaking bad",
		Log:    logrus.New(),
	}
	b.Config.Library.Auth.APIKey = "string"
	b.Config.Library.Settings.AdultContent = false
	b.Config.Library.Settings.Language = "en-US"

	var result tmdb.Show

	json.Unmarshal(bytes, &result)

	b.Search = append(b.Search, result)

	page := b.filterSelectedData(1, 5, 0, 2)

	if page.Current != 1 {
		t.Errorf("Expected current page '1', got %v", page.Current)
	}
	if page.Total != 1 {
		t.Errorf("Expected total pages '1', got %v", page.Total)
	}
}
