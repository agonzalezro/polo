package utils

import (
	"regexp"
	"strings"
)

/*
 * Transform a string in a proper slug.
 *
 * Note: this will possibly not be the best function if you want to clean non
 * ascii chars as ñ or similar (if you are going to use the slug in a url this
 * is important).
 */
func Slugify(title string) (slug string) {
	re, _ := regexp.Compile(`[^\w\s-]`)
	slug = re.ReplaceAllLiteralString(title, "")

	re, _ = regexp.Compile(`[-\s]+`)
	slug = re.ReplaceAllLiteralString(slug, "-")

	slug = strings.ToLower(slug)
	return
}
