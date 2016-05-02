package file

import (
	"bufio"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func newTestParsedFile(s string) *ParsedFile {
	return &ParsedFile{
		scanner: bufio.NewScanner(strings.NewReader(s)),
	}
}

// TestNoMetadata will test that an error NoMetadataFound is raised when
// parseMetadata is called without metadata.
func TestNoMetadata(t *testing.T) {
	assert := assert.New(t)

	for _, input := range []string{
		"",
		"---\n---\n",
		"This is my awesome title\n===\nAnd this is some content",
	} {
		pf := newTestParsedFile(input)
		assert.Equal(NoMetadataFound, pf.parseMetadata())
	}
}

// Test that the metadata is parsed properly following the Pelican standards
// for it, that means. Check the examples/ folder for more info about them.
func TestPelicanMetadataParsing(t *testing.T) {
	assert := assert.New(t)

	timeLayout := "2006-01-02 15:04"

	expectedTitle := "My super title"
	expectedDate, _ := time.Parse(timeLayout, "2010-12-03 10:20")
	expectedSlug := "/my-super-post.html"
	expectedText := "This is the content of my super blog post."
	expectedTags := []string{"thats", "awesome"}

	content := fmt.Sprintf(
		"Title: %s\nDate: %s\nTags: thats, awesome\nSlug: %s\n\n%s",
		expectedTitle, expectedDate.Format(timeLayout), expectedSlug, expectedText)

	pf := newTestParsedFile(content)
	err := pf.parse()
	assert.NoError(err)

	assert.Equal(expectedTitle, pf.Title)
	assert.Equal(expectedDate, pf.Date)
	assert.Equal(expectedTags, pf.Tags)
	assert.Equal(expectedSlug, pf.Slug)
}

// Test that the metadata is parsed properly when the standard used for it is
// the Jekyll standard (enclosed between '---' lines).
func TestJekyllMetadataParsing(t *testing.T) {
	assert := assert.New(t)
	expectedTitle := "jekyll test"

	content := fmt.Sprintf("---\nTitle:%s\n---\nThis is the content", expectedTitle)
	pf := newTestParsedFile(content)

	err := pf.parseMetadata()
	assert.NoError(err)

	assert.Equal(expectedTitle, pf.Title)
}

// Test that defining only one piece of metadata (ex date) the title & slug are
// generated properly.
func TestOnlyOneMetadata(t *testing.T) {
	assert := assert.New(t)
	expectedTitle := "Generated title"
	expectedDate := "1984-04-04 16:16"

	content := fmt.Sprintf("Date: %s\n\n%s\n---\nAnd some content", expectedDate, expectedTitle)
	pf := newTestParsedFile(content)

	err := pf.parseMetadata()
	assert.NoError(err)

	date, err := parseDate(expectedDate)
	assert.NoError(err)
	assert.Equal(date, pf.Date)
}
