package main

import (
	"flag"
	"log"

	"github.com/agonzalezro/polo/generator"
)

func main() {
	var (
		inputPath = flag.String("input", ".",
			"path to your articles source files.")
		outputPath = flag.String("output", ".",
			"path where you want to creat the html files.")
		configFile = flag.String("config", "config.json",
			"the settings file to create your site.")
	)
	flag.Parse()

	db, err := generator.NewDB()
	if err != nil {
		log.Panic(err)
	}
	if err := db.Fill(*inputPath); err != nil {
		log.Panic(err)
	}

	config, err := generator.ParseConfigFile(*configFile)
	if err != nil {
		log.Panic(err)
	}

	site := generator.NewSite(*db, *config, *outputPath)
	if err := site.Write(); err != nil {
		log.Panic(err)
	}
}
