package parser

import (
	"html/template"

	"github.com/russross/blackfriday"
)

type Config struct {
	Author string
	Title  string

	PaginationSize int

	DisqusSitename     string
	GoogleAnalyticsId  string
	SharethisPublisher string
}

type ParsedFile struct {
	Metadata map[string]string

	Title   string
	Slug    string
	Content string
}

func (file ParsedFile) Html() template.HTML {
	html := blackfriday.MarkdownCommon([]byte(file.Content))
	return template.HTML(html)
}
