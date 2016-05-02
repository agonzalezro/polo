package utils

import (
	"regexp"
	"strings"
)

// Slugify transform a string in a proper slug.
func Slugify(title string) (slug string) {
	// TODO: weird behaviour expected with non-ascci characters.
	re, _ := regexp.Compile(`[^\w\s-]`)
	slug = re.ReplaceAllLiteralString(title, "")

	re, _ = regexp.Compile(`[-\s]+`)
	slug = re.ReplaceAllLiteralString(slug, "-")

	slug = strings.ToLower(slug)
	return
}
