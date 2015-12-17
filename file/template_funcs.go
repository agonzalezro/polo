package file

import (
	"html/template"
	"strings"

	"github.com/russross/blackfriday"
)

// HTML return a template.HTML with the original content rendered as markdown.
func (file ParsedFile) HTML(content string) template.HTML {
	// set up the HTML renderer
	htmlFlags := 0
	htmlFlags |= blackfriday.HTML_USE_SMARTYPANTS
	renderer := blackfriday.HtmlRenderer(htmlFlags, "", "")

	// set up the parser
	extensions := 0
	extensions |= blackfriday.EXTENSION_NO_INTRA_EMPHASIS
	extensions |= blackfriday.EXTENSION_TABLES
	extensions |= blackfriday.EXTENSION_FENCED_CODE
	extensions |= blackfriday.EXTENSION_AUTOLINK
	extensions |= blackfriday.EXTENSION_STRIKETHROUGH
	extensions |= blackfriday.EXTENSION_SPACE_HEADERS

	html := blackfriday.Markdown([]byte(content), renderer, extensions)
	return template.HTML(html)
}

// Summary returns the document summary in case that it was defined on the
// metadata, or it creates a small summary (1st paragraph) if not.
func (file ParsedFile) Summary() string {
	if file.summary != "" {
		return file.summary
	}
	// Avoid empty lines
	for _, content := range strings.Split(file.Content, "\n\n") {
		if content != "" {
			return content
		}
	}
	return ""
}
