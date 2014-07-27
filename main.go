package main

import (
	"flag"
	"log"

	"github.com/agonzalezro/polo/generator"
)

var (
	inputPath  string
	outputPath string
	configFile string
)

func init() {
	flag.StringVar(&inputPath, "input", ".", "path to your articles source files.")
	flag.StringVar(&outputPath, "output", ".", "path where you want to creat the html files.")
	flag.StringVar(&configFile, "config", "config.json", "the settings file to create your site.")
}

func main() {
	flag.Parse()

	db, err := generator.GetDB()
	if err != nil {
		log.Panic(err)
	}
	db.Fill(inputPath)

	config, err := generator.ParseConfigFile(configFile)
	if err != nil {
		log.Panic(err)
	}

	site := generator.NewSite(*db, *config, outputPath)
	if err := site.Write(); err != nil {
		log.Panic(err)
	}
}
