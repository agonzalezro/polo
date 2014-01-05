package main

import (
	"flag"

	"github.com/agonzalezro/go-polo/parser"
	"github.com/agonzalezro/go-polo/reader"
	"github.com/agonzalezro/go-polo/writer"
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

	pageFilePaths, articleFilePaths := reader.GetPagesAndArticles(inputPath)

	config := parser.ParseConfig(configFile)
	pages := parser.ParseFiles(pageFilePaths)
	articles := parser.ParseFiles(articleFilePaths)

	site := writer.Site{Config: config, Pages: pages, Articles: articles, OutputPath: outputPath}
	site.WriteSite()
}
