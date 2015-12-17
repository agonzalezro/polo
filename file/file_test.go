package file

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
	"time"
)

func fail(t *testing.T, name string, got interface{}, expected interface{}) {
	t.Error(name, "was not the expected.\n\tgot:", got, "\n\texpected:", expected)
}

// Test that the metadata is parsed properly following the Pelican standards
// for it, that means. Check the examples/ folder for more info about them.
func TestPelicanMetadataParsing(t *testing.T) {
	timeLayout := "2006-01-02 15:04"

	expectedTitle := "My super title"
	expectedDate, _ := time.Parse(timeLayout, "2010-12-03 10:20")
	expectedSlug := "/my-super-post.html"
	expectedText := "This is the content of my super blog post."
	expectedTags := []string{"thats", "awesome"}

	content := fmt.Sprintf(
		"Title: %s\nDate: %s\nTags: thats, awesome\nSlug: %s\n\n%s",
		expectedTitle, expectedDate.Format(timeLayout), expectedSlug, expectedText)

	pf := ParsedFile{}
	pf.scanner = bufio.NewScanner(strings.NewReader(content))

	pf.parseMetadata()
	if pf.Title != expectedTitle {
		fail(t, "Title", pf.Title, expectedTitle)
	}

	if pf.Date != expectedDate {
		fail(t, "Date", pf.Date, expectedDate)
	}
	if pf.Tags[0] != expectedTags[0] || pf.Tags[1] != expectedTags[1] {
		fail(t, "Tags", pf.Tags, expectedTags)
	}
	if pf.Slug != expectedSlug {
		fail(t, "Slug", pf.Slug, expectedSlug)
	}

	// Check that doesn't read too much
	pf.scanner.Scan()
	if pf.scanner.Text() != expectedText {
		t.Errorf("This line should be here!")
	}
}

// Test that the metadata is parsed properly when the standard used for it is
// the Jekyll standard (enclosed between '---' lines).
func TestJekyllMetadataParsing(t *testing.T) {
	expectedTitle := "jekyll test"

	content := fmt.Sprintf(
		"---\nTitle:%s\n---This is the content", expectedTitle)

	pf := ParsedFile{
		scanner: bufio.NewScanner(strings.NewReader(content)),
	}
	pf.parseMetadata()

	if pf.Title != expectedTitle {
		fail(t, "Title", pf.Title, expectedTitle)
	}
}
