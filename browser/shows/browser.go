package shows

import (
	"github.com/gdamore/tcell/v2"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
)

type Browser struct {
	Config config.Config
	CLI
	Search []tmdb.ShowResult
	Query  string
	Show
}

type Show struct {
	Index   int32
	Details tmdb.ShowDetail
	Season
}

type Season struct {
	Index   int32
	Details tmdb.ShowSeason
	Episode
}

type Episode struct {
	Index int32
}

type CLI struct {
	Screen tcell.Screen
	Style  tcell.Style
}

func Browse(conf config.Config) {
	b := Browser{
		Config: conf,
	}
	b.showSearch()
}
