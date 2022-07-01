package shows

import (
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
)

type Browser struct {
	Config config.Config
	cli.CLI
	Search []tmdb.Show
	Query  string
	Show
}

type Show struct {
	cli.Page
	Index   int
	Details tmdb.ShowDetail
	Season
}

type Season struct {
	cli.Page
	Index   int
	Details tmdb.ShowSeason
	Episode
}

type Episode struct {
	cli.Page
	Index   int
	Details tmdb.ShowEpisode
	ID      int
}

func Browse(conf config.Config) {
	b := Browser{
		Config: conf,
	}
	b.showSearch()
}
