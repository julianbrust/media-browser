package tmdb

import (
	"net/http"
	"strconv"
)

// Queries represents a collection of possible query parameters for tmdb requests.
type Queries struct {
	ApiKey       string // API Key to authenticate with tmdb
	AdultContent bool   // toggle to include adult content
	Language     string // ISO 639-1 language option
	Page         int    // requested page of results
	Query        string // query string
}

// Show represents the response object for tmdb request GET /search/tv.
type Show struct {
	Page         int          `json:"page"`
	TotalPages   int          `json:"total_pages"`
	TotalResults int          `json:"total_results"`
	Results      []ShowResult `json:"results"`
}

// ShowResult represents the object for a show result in tmdb request GET /search/tv.
type ShowResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// ShowDetail represents the response object for tmdb request GET /tv/{tv_id}.
type ShowDetail struct {
	ID               int                `json:"id"`
	Name             string             `json:"name"`
	NumberOfEpisodes int                `json:"number_of_episodes"`
	NumberOfSeasons  int                `json:"number_of_seasons"`
	Seasons          []ShowDetailSeason `json:"seasons"`
}

// ShowDetailSeason represents the object for a Season in tmdb request GET /tv/{tv_id}.
type ShowDetailSeason struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	EpisodeCount int    `json:"episode_count"`
	SeasonNumber int    `json:"season_number"`
}

// ShowSeason represents the response object for tmdb request GET /tv/{tv_id}/season/{season_number}.
type ShowSeason struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	SeasonNumber int           `json:"season_number"`
	Episodes     []ShowEpisode `json:"episodes"`
	Overview     string        `json:"overview"`
}

// ShowEpisode represents the object for an Episode in tmdb request GET /tv/{tv_id}/season/{season_number}.
type ShowEpisode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Overview string `json:"overview"`
}

func GetTVLatest(queries Queries) (*http.Response, error) {
	req, _ := http.NewRequest("GET", "https://api.themoviedb.org/3/tv/latest", nil)

	q := req.URL.Query()
	q.Add("api_key", queries.ApiKey)
	req.URL.RawQuery = q.Encode()

	res, err := Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// SearchTV runs a request to tmdb for GET /search/tv.
func SearchTV(queries Queries) (*http.Response, error) {
	req, _ := http.NewRequest("GET", "https://api.themoviedb.org/3/search/tv", nil)

	q := req.URL.Query()
	q.Add("api_key", queries.ApiKey)
	q.Add("language", queries.Language)
	q.Add("include_adult", strconv.FormatBool(queries.AdultContent))
	q.Add("page", strconv.FormatInt(int64(queries.Page), 10))
	q.Add("query", queries.Query)
	req.URL.RawQuery = q.Encode()

	res, err := Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetTVShow runs a request to tmdb for GET /tv/{tv_id}.
func GetTVShow(id int, queries Queries) (*http.Response, error) {
	addr := "https://api.themoviedb.org/3/tv/" + strconv.FormatInt(int64(id), 10)
	req, _ := http.NewRequest("GET", addr, nil)

	q := req.URL.Query()
	q.Add("api_key", queries.ApiKey)
	q.Add("language", queries.Language)
	req.URL.RawQuery = q.Encode()

	res, err := Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

// GetTVShowSeason runs a request to tmdb for GET /tv/{tv_id}/season/{season_number}.
func GetTVShowSeason(id int, season int, queries Queries) (*http.Response, error) {
	addr := "https://api.themoviedb.org/3/tv/" + strconv.FormatInt(int64(id), 10) +
		"/season/" + strconv.FormatInt(int64(season), 10)

	req, _ := http.NewRequest("GET", addr, nil)

	q := req.URL.Query()
	q.Add("api_key", queries.ApiKey)
	q.Add("language", queries.Language)
	req.URL.RawQuery = q.Encode()

	res, err := Client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
