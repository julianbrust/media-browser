package main

import (
	"github.com/julianbrust/media-browser/browser/shows"
	"github.com/julianbrust/media-browser/config"
	"github.com/julianbrust/media-browser/tmdb"
	"log"
	"os"
)

func main() {
	//TODO: clean up logging
	file := "/tmp/media-browser.log"

	err := os.Remove(file)
	if err != nil {
		panic(err)
	}

	logFile, err := os.OpenFile(file, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Panic(err)
	}
	defer logFile.Close()

	log.SetOutput(logFile)

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
