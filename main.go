package main

import (
	"github.com/julianbrust/media-browser/browser/shows"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/logger"
	"github.com/julianbrust/media-browser/tmdb"
	"github.com/sirupsen/logrus"
)

var (
	log  *logrus.Logger
	conf config.Config
)

func init() {
	var err error
	conf, err = config.Get()
	if err != nil {
		log.Fatal(err)
	}
	log = logger.Init(&conf.Logger.Level)
}

func main() {
	err := tmdb.VerifyConfig(conf)
	if err != nil {
		log.Fatal(err)
	}

	b := shows.Browser{
		Config: &conf,
		Log:    log,
	}
	b.Browse()
}
