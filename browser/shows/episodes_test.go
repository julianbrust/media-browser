package shows

import (
	"encoding/json"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"testing"
)

func TestGetEpisodeResults(t *testing.T) {
	bytes, err := ioutil.ReadFile("../../testing/tmdb/get-show-season.json")
	if err != nil {
		log.Fatal(err)
	}

	var season tmdb.ShowSeason

	err = json.Unmarshal(bytes, &season)
	if err != nil {
		log.Fatal(err)
	}

	b := Browser{
		Config: &config.Config{},
		Log:    logrus.New(),
	}
	b.Show.Season.Details.Episodes = season.Episodes

	page := b.getEpisodeResults(1, 3)

	if page.Current != 1 {
		t.Errorf("Expected episode result index '1', got %v", page.Current)
	}
	if len(page.Content) != 3 {
		t.Errorf("Expected content amount '3', got %v", len(page.Content))
	}
}

func TestGetCurrentEpisode(t *testing.T) {
	bytes, err := ioutil.ReadFile("../../testing/tmdb/get-show-season.json")
	if err != nil {
		log.Fatal(err)
	}

	var season tmdb.ShowSeason

	err = json.Unmarshal(bytes, &season)
	if err != nil {
		log.Fatal(err)
	}

	b := Browser{
		Config: &config.Config{},
		Log:    logrus.New(),
	}
	b.Show.Season.Details.Episodes = season.Episodes
	b.Show.Season.Episode.Index = 1
	b.Show.Season.Episode.Page = cli.Page{
		Content: []cli.Content{
			{
				Display: "Pilot",
				ID:      62085,
			},
			{
				Display: "Cat's in the Bag...",
				ID:      62086,
			},
		},
	}

	episode := b.getCurrentEpisode()

	if episode.Details.ID != 62086 {
		t.Errorf("Expected episode id '62086', got %v", episode.ID)
	}
}
