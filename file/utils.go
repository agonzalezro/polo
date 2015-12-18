package file

import (
	"html/template"
	"strings"

	"github.com/russross/blackfriday"
)

func IsMarkdown(path string) bool {
	return strings.HasSuffix(path, ".md") || strings.HasSuffix(path, ".markdown")
}

func IsPage(path string) bool {
	return strings.HasPrefix(path, "pages/") || strings.Index(path, "/pages/") > 0
}

func HTML(content string) template.HTML {
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

func CategoryFromPath(path string) string {
	splittedPath := strings.Split(path, "/")
	length := len(splittedPath)
	if length > 1 {
		return splittedPath[length-2]
	}
	return ""
}
