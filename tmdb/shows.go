package tmdb

import (
	"net/http"
	"strconv"
)

type Queries struct {
	ApiKey       string
	AdultContent bool
	Language     string
	Page         int
	Query        string
}

//type ShowsSearch struct {
//	Search         Show
//	Query          string
//	SelectionIndex int32
//}

type ShowSelection struct {
	Details      ShowDetail
	SeasonIndex  int
	EpisodeIndex int
}

type Show struct {
	ResBody
	Results []ShowResult `json:"results"`
}

type ShowResult struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type ShowDetail struct {
	ID               int                `json:"id"`
	Name             string             `json:"name"`
	NumberOfEpisodes int                `json:"number_of_episodes"`
	NumberOfSeasons  int                `json:"number_of_seasons"`
	Seasons          []ShowSeasonDetail `json:"seasons"`
}

type ShowSeasonDetail struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Overview     string `json:"overview"`
	EpisodeCount int    `json:"episode_count"`
	SeasonNumber int    `json:"season_number"`
}

type ShowSeason struct {
	ID           int           `json:"id"`
	Name         string        `json:"name"`
	SeasonNumber int           `json:"season_number"`
	Episodes     []ShowEpisode `json:"episodes"`
	Overview     string        `json:"overview"`
}

type ShowEpisode struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Overview string `json:"overview"`
}

type ResBody struct {
	Page         int `json:"page"`
	TotalPages   int `json:"total_pages"`
	TotalResults int `json:"total_results"`
}

func GetTVLatest(queries Queries) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api.themoviedb.org/3/tv/latest", nil)

	q := req.URL.Query()
	q.Add("api_key", queries.ApiKey)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func SearchTV(queries Queries) (*http.Response, error) {
	client := &http.Client{}

	req, _ := http.NewRequest("GET", "https://api.themoviedb.org/3/search/tv", nil)

	q := req.URL.Query()
	q.Add("api_key", queries.ApiKey)
	q.Add("language", queries.Language)
	q.Add("include_adult", strconv.FormatBool(queries.AdultContent))
	q.Add("page", strconv.FormatInt(int64(queries.Page), 10))
	q.Add("query", queries.Query)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetTVShow(id int, queries Queries) (*http.Response, error) {
	client := &http.Client{}

	addr := "https://api.themoviedb.org/3/tv/" + strconv.FormatInt(int64(id), 10)
	req, _ := http.NewRequest("GET", addr, nil)

	q := req.URL.Query()
	q.Add("api_key", queries.ApiKey)
	q.Add("language", queries.Language)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func GetTVShowSeason(id int, season int, queries Queries) (*http.Response, error) {
	client := &http.Client{}

	addr := "https://api.themoviedb.org/3/tv/" + strconv.FormatInt(int64(id), 10) +
		"/season/" + strconv.FormatInt(int64(season), 10)

	req, _ := http.NewRequest("GET", addr, nil)

	q := req.URL.Query()
	q.Add("api_key", queries.ApiKey)
	q.Add("language", queries.Language)
	req.URL.RawQuery = q.Encode()

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
