package file

import (
	"bufio"
	"strings"
	"time"
)

// ParsedFile holds the struct of a file after we parsed the metadata and its
// content.
type ParsedFile struct {
	Metadata map[string]string

	Author  string
	Title   string
	Slug    string
	Content string
	status  string // To keep track of the drafts
	summary string

	IsPage bool

	Category string
	Tags     []string
	Date     time.Time

	scanner *bufio.Scanner
}

// New return a new ParsedFile after load it from disk.
func New(path string) (*ParsedFile, error) {
	pf := ParsedFile{}
	pf.IsPage = strings.HasPrefix(path, "pages/") || strings.Index(path, "/pages/") > 0
	if err := pf.load(path); err != nil {
		return nil, err
	}
	return &pf, nil
}
