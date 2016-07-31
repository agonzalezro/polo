package file

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsPublished(t *testing.T) {
	assert := assert.New(t)

	assert.True(ParsedFile{}.IsPublished())

	assert.False(ParsedFile{status: "Draft"}.IsPublished())
	assert.False(ParsedFile{status: "draft"}.IsPublished())
}
