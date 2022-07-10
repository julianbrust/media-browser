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

func TestGetShowResult(t *testing.T) {
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
		Config: &config.Config{},
		Query:  "breaking bad",
		Log:    logrus.New(),
	}
	b.Config.Library.Auth.APIKey = "string"
	b.Config.Library.Settings.AdultContent = false
	b.Config.Library.Settings.Language = "en-US"

	show, err := b.getShowResult(1)
	if err != nil {
		t.Error(err)
	}

	if show.Results[0].ID != 1396 {
		t.Errorf("Expected show ID '1396', got %v", show.Results[0].ID)
	}
}

func TestGetShow(t *testing.T) {
	bytes, err := ioutil.ReadFile("../../testing/tmdb/get-show-details.json")
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
		Log:    logrus.New(),
	}
	b.Config.Library.Auth.APIKey = "string"
	b.Config.Library.Settings.Language = "en-US"

	show, err := b.getShow(1396)
	if err != nil {
		t.Error(err)
	}

	if show.ID != 1396 {
		t.Errorf("Expected show ID '1396', got %v", show.ID)
	}
	if len(show.Seasons) != 6 {
		t.Errorf("Expected season count '6', got %v", len(show.Seasons))
	}
}

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
		Config: &config.Config{},
		Query:  "breaking bad",
		Log:    logrus.New(),
	}
	b.Config.Library.Auth.APIKey = "string"
	b.Config.Library.Settings.AdultContent = false
	b.Config.Library.Settings.Language = "en-US"

	err = b.getSearchResults(1, 5)
	if err != nil {
		t.Error(err)
	}

	if len(b.Search) > 0 && len(b.Search[0].Results) > 0 {
		if b.Search[0].Results[0].ID != 1396 {
			t.Errorf("Expected show ID '1396', got %v", b.Search[0].Results[0].ID)
		}
	} else {
		t.Error("Expected show result, got none")
	}
	if b.Show.Page.Current != 1 {
		t.Errorf("Expected current page '1', got %v", b.Show.Page.Current)
	}
	if b.Show.Page.Total != 1 {
		t.Errorf("Expected total pages '1', got %v", b.Show.Page.Total)
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
		Config: &config.Config{},
		Query:  "breaking bad",
		Log:    logrus.New(),
	}
	b.Config.Library.Auth.APIKey = "string"
	b.Config.Library.Settings.AdultContent = false
	b.Config.Library.Settings.Language = "en-US"

	err = b.getSearchResults(1, 5)
	if err != nil {
		t.Error(err)
	}

	if len(b.Search) > 0 && len(b.Search[0].Results) > 0 {
		if b.Search[0].Results[0].ID != 1396 {
			t.Errorf("Expected show ID '1396', got %v", b.Search[0].Results[0].ID)
		}
	} else {
		t.Error("Expected show result, got none")
	}
	if b.Show.Page.Current != 1 {
		t.Errorf("Expected current page '1', got %v", b.Show.Page.Current)
	}
	if b.Show.Page.Total != 1 {
		t.Errorf("Expected total pages '1', got %v", b.Show.Page.Total)
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
		Config: &config.Config{},
		Query:  "breaking bad",
		Log:    logrus.New(),
	}
	b.Config.Library.Auth.APIKey = "string"
	b.Config.Library.Settings.AdultContent = false
	b.Config.Library.Settings.Language = "en-US"

	var result tmdb.Show

	json.Unmarshal(bytes, &result)

	b.Search = append(b.Search, result)

	b.filterSelectedData(1, 5, 0, 2)

	if b.Show.Page.Current != 1 {
		t.Errorf("Expected current page '1', got %v", b.Show.Page.Current)
	}
	if b.Show.Page.Total != 1 {
		t.Errorf("Expected total pages '1', got %v", b.Show.Page.Total)
	}
}
