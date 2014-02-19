package generator

import (
	"bufio"
	"strings"
	"testing"
)

// Test that the metadata is generated properly following the Pelican standards
// for it
func TestMetadataGeneration(t *testing.T) {
	content := `Title: My super title
Date: 2010-12-03 10:20
Tags: thats, awesome
Slug: my-super-post

This is the content of my super blog post.
`

	pf := ParsedFile{}
	pf.scanner = bufio.NewScanner(strings.NewReader(content))

	pf.parseMetadata()
	if pf.Title != "My super title" {
		t.Errorf("Title is not the expected: '%s'", pf.Title)
	}
	if pf.Date != "2010-12-03 10:20" {
		t.Errorf("Date is not the expected: '%s'", pf.Date)
	}
	if pf.tags != "thats, awesome" {
		t.Errorf("Tags is not expected: '%s'", pf.tags)
	}
	if pf.Slug != "my-super-post" {
		t.Errorf("Slug is not expected: '%s'", pf.Slug)
	}

	// Check that doesn't read too much
	pf.scanner.Scan()
	if pf.scanner.Text() != "This is the content of my super blog post." {
		t.Errorf("This line should be here!")
	}
}
