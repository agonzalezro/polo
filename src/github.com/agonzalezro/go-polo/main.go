package main

import (
	"flag"
	"fmt"

	"github.com/agonzalezro/go-polo/parser"
	"github.com/agonzalezro/go-polo/reader"
	"github.com/agonzalezro/go-polo/writer"
)

var (
	inputPath    string
	outputPath   string
	settingsFile string
)

func init() {
	flag.StringVar(&inputPath, "input", ".", "path to your articles source files.")
	flag.StringVar(&outputPath, "output", ".", "path where you want to creat the html files.")
	flag.StringVar(&settingsFile, "settings", "settings.yaml", "the settings file to create your site.")
}

func main() {
	flag.Parse()

	pageFilePaths, articleFilePaths := reader.GetPagesAndArticles(inputPath)
	fmt.Printf("PAGES:\n")
	pages := parser.ParseFiles(pageFilePaths)
	for _, document := range pages {
		writer.WriteToHtml(document, outputPath)
	}

	fmt.Printf("\nARTICLES:\n")
	articles := parser.ParseFiles(articleFilePaths)
	for _, document := range articles {
		writer.WriteToHtml(document, outputPath)
	}
}
