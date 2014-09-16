package config

import (
	"encoding/json"
	"os"
)

// Config stores the configurations readed from the JSON file.
type Config struct {
	Author string
	Title  string

	URL     string
	Favicon string

	ShowArchive    bool
	ShowCategories bool
	ShowTags       bool

	PaginationSize int

	DisqusSitename     string
	GoogleAnalyticsID  string
	SharethisPublisher string
}

// ErrorOpeningConfigFile will be raised when the file doesn't exist.
type ErrorOpeningConfigFile error

// ErrorParsingConfigFile will be raised when the JSON config is malformed.
type ErrorParsingConfigFile error

// New returns a New configuration after parse the file received as input.
func New(configFile string) (*Config, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return nil, ErrorOpeningConfigFile(err)
	}

	decoder := json.NewDecoder(file)
	config := &Config{}

	err = decoder.Decode(&config)
	if err != nil {
		return nil, ErrorParsingConfigFile(err)
	}
	return config, nil
}
