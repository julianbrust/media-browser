package main

import (
	"github.com/julianbrust/media-browser/browser/shows"
	"github.com/julianbrust/media-browser/cli"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/logger"
	"github.com/julianbrust/media-browser/tmdb"
	"github.com/sirupsen/logrus"
	"log"
)

var (
	customLog *logrus.Logger
	conf      config.Config
)

func init() {
	conf = conf.Get()

	cli.GetArgs(&conf)

	customLog = logger.Init(&conf.Logger.Level)
	conf.PrintConfig(customLog)
}

func main() {
	err := tmdb.VerifyConfig(conf)
	if err != nil {
		log.Fatal(err)
	}
	customLog.Infoln("config verified")

	b := shows.Browser{
		Config: &conf,
		Log:    customLog,
	}
	b.Browse()
}
