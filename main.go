package main

import (
	"flag"

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

	db := generator.GetDB()
	db.Fill(inputPath)

	site := generator.NewSite(*db, generator.GetConfig(configFile), outputPath)
	site.Write()
}
