package file

import (
	"bufio"
	"html/template"
	"os"
	"strings"
	"time"
)

// ParsedFile holds the struct of a file after we parsed the metadata and its
// content.
type ParsedFile struct {
	Author  string
	Title   string
	Slug    string
	Content template.HTML
	Summary template.HTML

	IsPage   bool
	Category string
	Tags     []string
	Date     time.Time

	// Not to be used by the template
	rawContent string
	summary    string
	status     string // To keep track of the drafts

	file    *os.File
	scanner *bufio.Scanner
}

// New return a new ParsedFile after load it from disk.
func New(path string) (*ParsedFile, error) {
	pf := ParsedFile{
		IsPage:   IsPage(path),
		Category: CategoryFromPath(path),
	}
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	pf.file = file
	pf.scanner = bufio.NewScanner(file) // We need this to seek after parsing metadata

	if err := pf.parse(); err != nil {
		return nil, err
	}

	pf.Summary = HTML(pf.summaryOrFirstParagraph())
	pf.Content = HTML(pf.rawContent)
	return &pf, nil
}

// summaryOrFirstParagraph will use the summary from the markdown or generate a new one from the 1st paragraph.
func (f ParsedFile) summaryOrFirstParagraph() string {
	summary := f.summary

	if summary == "" {
		// Get the first paragraph
		for _, summary := range strings.Split(f.rawContent, "\n\n") {
			if summary != "" {
				return summary
			}
		}
	}

	return summary
}

// IsPublished will return true if the status is not draft
// TODO: I wonder if `draft: true` could be a better way of doing this.
func (f ParsedFile) IsPublished() bool {
	return f.status == "draft"
}
