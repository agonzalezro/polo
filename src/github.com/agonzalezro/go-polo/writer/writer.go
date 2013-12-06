package writer

import (
	"fmt"
	"io/ioutil"

	"github.com/agonzalezro/go-polo/parser"
)

func WriteToHtml(document parser.ParsedFile, outputPath string) {
	filePath := fmt.Sprintf("%s/%s.html", outputPath, document.Metadata["slug"])
	ioutil.WriteFile(filePath, document.Content, 0644)
}
