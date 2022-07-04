package config

import (
	"github.com/sirupsen/logrus"
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
	Logger struct {
		Level string `yaml:"level"`
	} `yaml:"logger"`
}

// Get parses the YAML configuration from the config.yaml file
// located in the same directory as the app.
func (conf Config) Get() Config {
	dirname, _ := os.Getwd()

	config := "config.yaml"

	file, _ := os.Open(path.Join(dirname, config))

	defer file.Close()

	d := yaml.NewDecoder(file)

	_ = d.Decode(&conf)

	newConf := conf.setDefaults()

	return newConf
}

// setDefaults provides some default parameters if nothing is defined.
func (conf Config) setDefaults() Config {
	if conf.Library.Settings.Language == "" {
		conf.Library.Settings.Language = "en-US"
	}
	if conf.Logger.Level == "" {
		conf.Logger.Level = "info"
	}

	return conf
}

func (conf Config) PrintConfig(log *logrus.Logger) {
	log.Infof("Settings: adultContent: %v", conf.Library.Settings.AdultContent)
	log.Infof("Settings: language: %v", conf.Library.Settings.Language)
	log.Infof("Logger: level: %v", conf.Logger.Level)
}
