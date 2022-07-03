# media-browser
Simple CLI to browse through an index of movies and TV shows.

The app uses the [TMDB](https://developers.themoviedb.org/) database to browse TV shows.
A personal API key is required.

The CLI offers a search input to find TV shows and browse through results for shows, seasons and episodes.
It displays detailed information about the selected episode.

## How to run it

The CLI requires a config file for the API key and other optional parameters.

Create a file `config.yaml` in the same directory as the app. `config-example.yaml` can be used as a reference.

```yaml
library:
  auth:
    # your TMDB API key
    apiKey: <your-key>
  settings:
    # toggle to allow adult content as part of query results
    adultContent: false
    # toggle for the display language of the results (as ISO 639-1)
    language: en-US
logger:
  # desired log level
  # trace | debug | info | warn | error | fatal
  level: info
```

Run the app with `go run .`

## Logging

Logs will be written to `/tmp/media-browser.log`
