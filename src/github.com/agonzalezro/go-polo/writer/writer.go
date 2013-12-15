package writer

import (
	"fmt"
	"log"
	"os"
	"text/template"

	"github.com/agonzalezro/go-polo/parser"
)

func DumpAll(pages []parser.ParsedFile, articles []parser.ParsedFile, outputPath string) {
	indexFile := fmt.Sprintf("%s/index.html", outputPath)

	template, err := template.ParseFiles("templates/base.html")
	if err != nil {
		log.Panic(err)
	}
	file, err := os.Create(indexFile)
	if err != nil {
		log.Panic(err)
	}
	err = template.Execute(file, parser.Site{pages, articles})
}

func articleToHtml(document parser.ParsedFile, outputPath string) {
	filePath := fmt.Sprintf("%s/%s.html", outputPath, document.Metadata["slug"])

	template, err := template.ParseFiles("templates/base.html")
	if err != nil {
		log.Panic(err)
	}
	file, err := os.Create(filePath)
	err = template.Execute(file, document)
	if err != nil {
		log.Panic(err)
	}
}
