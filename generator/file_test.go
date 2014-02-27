package generator

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
)

// Test that the metadata is generated properly following the Pelican standards
// for it
func TestMetadataGeneration(t *testing.T) {
	expectedTitle := "My super title"
	expectedDate := "2010-12-03 10:20"
	expectedSlug := "my-super-post"
	expectedText := "This is the content of my super blog post."

	content := fmt.Sprintf(
		"Title: %s\nDate: %s\nTags: thats, awesome\nSlug: %s\n\n%s",
		expectedTitle, expectedDate, expectedSlug, expectedText)

	pf := ParsedFile{}
	pf.scanner = bufio.NewScanner(strings.NewReader(content))

	pf.parseMetadata()
	if pf.Title != expectedTitle {
		t.Errorf("Title is not the expected: '%s'", pf.Title)
	}
	if pf.Date != expectedDate {
		t.Errorf("Date is not the expected: '%s'", pf.Date)
	}
	if pf.tags != ",thats,awesome," {
		t.Errorf("Tags is not expected: '%s'", pf.tags)
	}
	if pf.Slug != expectedSlug {
		t.Errorf("Slug is not expected: '%s'", pf.Slug)
	}

	// Check that doesn't read too much
	pf.scanner.Scan()
	if pf.scanner.Text() != expectedText {
		t.Errorf("This line should be here!")
	}
}
