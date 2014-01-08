package parser

import (
	"html/template"
	"strings"

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

func (file ParsedFile) Html(content string) template.HTML {
	html := blackfriday.MarkdownCommon([]byte(content))
	return template.HTML(html)
}

func (file ParsedFile) Tags() (tags []string) {
	for _, tag := range strings.Split(file.Metadata["tags"], ",") {
		tags = append(tags, tag)
	}
	return tags
}
