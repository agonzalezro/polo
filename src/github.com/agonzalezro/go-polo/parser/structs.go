package parser

import (
	"html/template"
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
	Content  []byte
}

func (file ParsedFile) Html() template.HTML {
	return template.HTML(string(file.Content))
}
