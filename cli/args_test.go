package cli

import (
	"github.com/julianbrust/media-browser/config"
	"os"
	"testing"
)

func TestGetArgs(t *testing.T) {
	conf := config.Config{}
	conf.Library.Auth.APIKey = "secret"
	conf.Library.Settings.AdultContent = false
	conf.Library.Settings.Language = "en-US"
	conf.Logger.Level = "trace"

	os.Args = append(os.Args, "--key", "update")
	os.Args = append(os.Args, "--adult", "true")
	os.Args = append(os.Args, "--language", "de")
	os.Args = append(os.Args, "--log", "info")

	res := GetArgs(&conf)

	if res.Library.Auth.APIKey != "update" {
		t.Errorf("Expected key 'update', got %v", res.Library.Auth.APIKey)
	}
	if conf.Library.Settings.AdultContent != true {
		t.Errorf("Expected adult setting 'true', got %v", res.Library.Settings.AdultContent)
	}
	if conf.Library.Settings.Language != "de" {
		t.Errorf("Expected language 'de', got %v", res.Library.Settings.Language)
	}
	if conf.Logger.Level != "info" {
		t.Errorf("Expected log level 'info', got %v", res.Logger.Level)
	}
}
