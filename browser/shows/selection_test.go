package shows

import (
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
	"github.com/sirupsen/logrus"
	"testing"
)

func TestGetSelectionFull(t *testing.T) {
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
		AirDate:  YYYYMMDD,
		Runtime:  25,
	}

	text := b.getSelection(header)

	expected := []string{
		"Header 1",
		"Header 2",
		"  ShowName: SeasonName: EpisodeName",
		"  Description: Overview",
		"  Air Date: January 2, 2006",
		"  Runtime: 25 min",
	}

	for i, line := range expected {
		if text[i] != line {
			t.Errorf("Expected line string '%v', got %v", line, text[i])
		}
	}
}

func TestGetSelectionMinimum(t *testing.T) {
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
		Name: "EpisodeName",
	}

	text := b.getSelection(header)

	expected := []string{
		"Header 1",
		"Header 2",
		"  ShowName: SeasonName: EpisodeName",
	}

	for i, line := range expected {
		if text[i] != line {
			t.Errorf("Expected line string '%v', got %v", line, text[i])
		}
	}
}
