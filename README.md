# media-browser
Simple CLI to browse through an index of movies and TV shows.

The app uses the [TMDB](https://developers.themoviedb.org/) database to browse TV shows.
A personal API key is required.

The CLI offers a search input to find TV shows and browse through results for shows, seasons and episodes.
It displays detailed information about the selected episode.

## How to run it

### Config file

The CLI can be set up with a config file for the API key and other optional parameters.

Create a file `config.yaml` in the same directory as the app. `config-example.yaml` can be used as a reference.

```yaml
library:
  auth:
    # your TMDB API key
    apiKey: <your-key>
  settings:
    # toggle to allow adult content as part of query results
    adultContent: false
    # the display language of the results (as ISO 639-1)
    language: en-US
logger:
  # desired log level
  # trace | debug | info | warn | error | fatal
  level: info
```

### Command-line Arguments

The following arguments can be passed to the command:

| Key        | Value                                                                | Description                  |
|------------|----------------------------------------------------------------------|------------------------------|
| --key      | string                                                               | API key                      |
| --adult    | true &#124; false                                                    | enable/disable adult content |
| --language | en-US                                                                | language to be used          |
| --log      | trace &#124; debug &#124; info &#124; warn &#124; error &#124; fatal | log level                    |


Example `go run . --key <your-key> --adult false --language en-US --log info`

## Logging

Logs will be written to `/tmp/media-browser.log`
