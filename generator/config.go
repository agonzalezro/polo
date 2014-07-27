package generator

import (
	"encoding/json"
	"os"
)

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
	GoogleAnalyticsId  string
	SharethisPublisher string
}

type ErrorOpeningConfigFile error
type ErrorParsingConfigFile error

func ParseConfigFile(configFile string) (*Config, error) {
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
