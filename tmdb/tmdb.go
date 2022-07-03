package tmdb

import "net/http"

type Server struct {
	Client  *http.Client
	BaseURL string
}

var (
	server Server
)

func init() {
	server = Server{
		Client: &http.Client{},
	}
}
