package tmdb

import (
	"errors"
	"github.com/julianbrust/media-browser/config"
	"golang.org/x/text/language"
)

// VerifyConfig runs a test request to tmdb to verify the required config parameters.
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

// verifyLanguage tests if the provided language parameter is compliant to BCP 47 and ISO 639-1.
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

// verifyKey runs a test request to verify the provided API Key.
func verifyKey(key string) error {
	if key == "" {
		return errors.New("missing API key")
	}

	queries := Queries{ApiKey: key}

	res, err := server.GetTVLatest(queries)
	if err != nil {
		return err
	}
	if res.StatusCode >= 400 {
		return errors.New("invalid API key")
	}

	return nil
}
