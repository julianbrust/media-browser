package shows

import (
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
	"github.com/sirupsen/logrus"
	"net/http"
)

// Browser represents the object that holds all the data for the current search through shows,
// their seasons and episodes.
type Browser struct {
	Config  *config.Config // the current app config
	Log     *logrus.Logger // logger instance
	cli.CLI                // the CLI object for the tcell screens used
	Search  []tmdb.Show    // collection of show search results fron the tmdb API
	Query   string         // the current query
	Show                   // object containing all the information for the currently selected show
}

// Show represents all the information for a show.
type Show struct {
	cli.Page                 // the current results' page to display
	Index    int             // the current index of the selected show
	Details  tmdb.ShowDetail // details about the show provided by the tmdb API
	Error    bool            // indicating if there is an error currently
	Season                   // object containing all the information for the currently selected season
}

// Season represents all the information for a season.
type Season struct {
	cli.Page                 // the current results' page to display
	Index    int             // the current index of the selected season
	Details  tmdb.ShowSeason // details about the season provided by the tmdb API
	Episode                  // object containing all the information for the currently selected episode
}

// Episode represents all the information for an episode.
type Episode struct {
	cli.Page                  // the current results' page to display
	Index    int              // the current index of the selected episode
	Details  tmdb.ShowEpisode // details about the episode provided by the tmdb API
	ID       int              // ID of the current episode
}

var (
	server tmdb.Server
)

func init() {
	server = tmdb.Server{
		Client: &http.Client{},
	}
}

// Browse initiates drawing the first screen for browsing shows.
func (b *Browser) Browse() {
	b.Log.Infoln("starting browser")
	b.showSearch()
}
