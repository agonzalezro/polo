package file

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/agonzalezro/polo/utils"
)

func parseDate(value string) (t time.Time, err error) {
	acceptedFormats := []string{
		"2006-01-02 15:04",
		"2006-1-2 15:04",
		"2006-01-02",
		"2006-1-2",
	}
	for _, format := range acceptedFormats {
		if t, err := time.Parse(format, value); err == nil {
			return t, nil
		}
	}
	return t, fmt.Errorf("Accepted date/time formats are: %v", acceptedFormats)
}

var NoMetadataFound = errors.New("No metadata found!")

// parseMetadata sets the metadata on the ParsedFile.
// If no metadata is found no error is going to be raised.
func (pf *ParsedFile) parseMetadata() (err error) {
	var count int

	for pf.scanner.Scan() {
		count++
		line := pf.scanner.Text()

		// In case that the metadata starts like :date:
		if strings.HasPrefix(line, ":") {
			line = line[1:]
		}
		metadataSplited := strings.Split(line, ":")
		key := strings.ToLower(metadataSplited[0])
		value := strings.Trim(strings.Join(metadataSplited[1:], ":"), " ")

		switch key {
		case "---":
			// If the metadata is enclosed between lines like this: '---'
			// (Jekyll style) we need to return after process it.
			if count > 1 {
				goto END
			}
		case "tags":
			// Remove all the spaces between comma and tag and
			// add one comma at the beginning and other at the end, this will
			// make the querying much simpler
			for _, tag := range strings.Split(value, ",") {
				pf.Tags = append(pf.Tags, strings.Replace(tag, " ", "", -1))
			}
		case "date":
			pf.Date, err = parseDate(value)
			if err != nil {
				return err
			}
		case "slug":
			prefix := "/"
			if strings.HasPrefix(value, "/") {
				prefix = "" // Just to be sure that we don't duplicate the /
			}
			suffix := ".html"
			if strings.HasSuffix(value, ".html") || strings.HasSuffix(value, ".html") {
				suffix = "" // And don't duplicate the html either
			}
			pf.Slug = fmt.Sprintf("%s%s%s", prefix, value, suffix)
		case "status":
			pf.status = value
		case "summary":
			pf.summary = value
		case "author":
			pf.Author = value
		case "title":
			pf.Title = value
		default:
			goto END
		}
	}

END:
	// TODO: not the best way to check this. Find a cleaner way.
	allUnset := func() bool {
		return (pf.Tags == nil && pf.Date.IsZero() && pf.Slug == "" && pf.status == "" && pf.Summary == "" && pf.Author == "" && pf.Title == "")
	}
	if count <= 2 && allUnset() {
		return NoMetadataFound
	}
	return nil
}

// parse parses the content and metadata storing the on private fields of the struct
func (pf *ParsedFile) parse() error {
	var err error
	if err = pf.parseMetadata(); err != nil {
		switch err {
		case NoMetadataFound:
			// We have already read part of the file but we didn't found metadata.
			// Rewind the file and reset the scanner
			pf.file.Seek(0, 0)
			pf.scanner = bufio.NewScanner(pf.file)
		default:
			return err
		}
	}

	isFirstLine := true
	for pf.scanner.Scan() {
		line := pf.scanner.Text()
		if isFirstLine == true {
			if line == "" {
				// Ignore empty lines at the beginning of the file
				continue
			}

			if pf.Title == "" {
				// This is needed to remove markdown syntax before storing values on the structs
				re := regexp.MustCompile("^#+\\s*")
				pf.Title = re.ReplaceAllString(line, "")
			}

			if pf.Slug == "" {
				pf.Slug = fmt.Sprintf("/%s.html", utils.Slugify(pf.Title))
			}

			pf.scanner.Scan() // We don't want the title underlining

			isFirstLine = false
			continue
		}

		pf.rawContent += line + "\n"
	}

	return nil
}
