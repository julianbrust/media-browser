package config

import (
	"gopkg.in/yaml.v3"
	"os"
	"path"
)

// Config represent the configuration object for the app
type Config struct {
	Library struct {
		Auth struct {
			APIKey string `yaml:"apiKey"`
		} `yaml:"auth"`
		Settings struct {
			AdultContent bool   `yaml:"adultContent"`
			Language     string `yaml:"language"`
		} `yaml:"settings"`
	} `yaml:"library"`
}

// ReadConf parses the YAML configuration from the config.yaml file
// located in the same directory as the app.
func ReadConf() (Config, error) {
	conf := Config{}

	dirname, err := os.Getwd()
	if err != nil {
		return Config{}, err
	}
	config := "config.yaml"

	file, err := os.Open(path.Join(dirname, config))
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err = d.Decode(&conf); err != nil {
		return Config{}, err
	}

	conf = setDefaults(conf)

	return conf, nil
}

// setDefaults provides some default parameters if nothing is defined.
func setDefaults(conf Config) Config {
	if conf.Library.Settings.Language == "" {
		conf.Library.Settings.Language = "en-US"
	}

	return conf
}
