package cli

import (
	"github.com/julianbrust/media-browser/config"
	"os"
	"strconv"
	"strings"
)

// GetArgs updates the config object with parameters provided by command-line arguments.
func GetArgs(conf *config.Config) *config.Config {
	args := os.Args[1:]

	for i, element := range args {
		if i == len(args)-1 {
			break
		}
		if strings.HasPrefix(element, "--") {
			if element == "--key" {
				conf.Library.Auth.APIKey = args[i+1]
			}
			if element == "--adult" {
				adult, err := strconv.ParseBool(args[i+1])
				if err == nil {
					conf.Library.Settings.AdultContent = adult
				}
			}
			if element == "--language" {
				conf.Library.Settings.Language = args[i+1]
			}
			if element == "--log" {
				conf.Logger.Level = args[i+1]
			}
		}
	}

	return conf
}
