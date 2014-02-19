package polo

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	Author string
	Title  string

	PaginationSize int

	DisqusSitename     string
	GoogleAnalyticsId  string
	SharethisPublisher string
}

func GetConfig(configFile string) Config {
	file, err := os.Open(configFile)
	if err != nil {
		log.Panic(err)
	}

	decoder := json.NewDecoder(file)
	config := &Config{}

	err = decoder.Decode(&config)
	if err != nil {
		log.Panic(err)
	}
	return *config
}
