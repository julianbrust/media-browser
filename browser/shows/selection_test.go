package shows

import (
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestGetSelection(t *testing.T) {
	header := []string{
		"Header 1",
		"Header 2",
	}

	b := Browser{
		Config: &config.Config{},
		Query:  "breaking bad",
		Log:    logrus.New(),
	}
	b.Show.Details.Name = "ShowName"
	b.Show.Season.Details.Name = "SeasonName"
	b.Show.Season.Episode.Details = tmdb.ShowEpisode{
		Name:     "EpisodeName",
		Overview: "Overview",
	}

	text := b.getSelection(header)

	expected := []string{
		"Header 1",
		"Header 2",
		"  ShowName: SeasonName: EpisodeName",
		"  Description: Overview",
	}

	for i, line := range expected {
		if text[i] != line {
			t.Errorf("Expected line string '%v', got %v", line, text[i])
		}
	}
}
