package main

import (
	"flag"

	"github.com/agonzalezro/polo/polo"
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

	db := polo.GetDB()
	db.Fill(inputPath)

	site := polo.NewSite(*db, polo.GetConfig(configFile), outputPath)
	site.Write()
}
