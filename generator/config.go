package generator

import (
	"encoding/json"
	"log"
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

func GetConfig(configFile string) Config {
	file, err := os.Open(configFile)
	if err != nil {
		log.Panicf("Error opening the configuration file: %s\n%v", configFile, err)
	}

	decoder := json.NewDecoder(file)
	config := &Config{}

	err = decoder.Decode(&config)
	if err != nil {
		log.Panicf("Error reading the JSON file: %s\n%v", configFile, err)
	}
	return *config
}
