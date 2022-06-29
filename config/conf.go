package config

import (
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"path"
)

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

func ReadConf() (Config, error) {
	conf := Config{}

	dirname, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	config := "config.yaml"

	file, err := os.Open(path.Join(dirname, config))
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	d := yaml.NewDecoder(file)

	if err := d.Decode(&conf); err != nil {
		return conf, err
	}

	conf = setDefaults(conf)

	return conf, nil
}

func setDefaults(conf Config) Config {
	if conf.Library.Settings.Language == "" {
		conf.Library.Settings.Language = "en-US"
	}

	return conf
}
