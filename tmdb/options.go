package tmdb

import (
	"errors"
	"github.com/julianbrust/media-browser/config"
	"golang.org/x/text/language"
)

func VerifyConfig(conf config.Config) error {
	err := verifyKey(conf.Library.Auth.APIKey)
	if err != nil {
		return err
	}

	err = verifyLanguage(conf.Library.Settings.Language)
	if err != nil {
		return err
	}

	return nil
}

func verifyLanguage(lang string) error {
	_, err := language.Parse(lang)
	if err != nil {
		if _, ok := err.(language.ValueError); ok {
			return errors.New("unsupported language provided")
		} else {
			return nil
		}
	}

	return nil
}

func verifyKey(key string) error {
	queries := Queries{ApiKey: key}

	res, err := GetTVLatest(queries)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return errors.New("invalid API key")
	}

	return nil
}