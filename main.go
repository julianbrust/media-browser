package main

import (
	"github.com/julianbrust/media-browser/browser/shows"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
	"github.com/sirupsen/logrus"
	"os"
)

var log *logrus.Logger

func init() {
	file := "/tmp/media-browser.log"
	log = logrus.New()

	err := os.Remove(file)
	if err != nil {
		log.Fatal(err)
	}

	logFile, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}

	log.SetOutput(logFile)
}

func main() {
	conf, err := config.ReadConf()
	if err != nil {
		log.Fatal(err)
	}

	err = tmdb.VerifyConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	shows.Browse(conf)
}
