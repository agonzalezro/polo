package file

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/agonzalezro/polo/utils"
)

func parseData(value string) (t time.Time, err error) {
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

// parseMetadata sets the metadata on the ParsedFile.
// If no metadata is found no error is going to be raised.
func (pf *ParsedFile) parseMetadata() (hasMetadata bool, err error) {
	for pf.scanner.Scan() {
		line := pf.scanner.Text()

		metadataSplited := strings.Split(line, ":")
		key := strings.ToLower(metadataSplited[0])
		value := strings.Trim(strings.Join(metadataSplited[1:], ":"), " ")

		switch key {
		case "---":
			// If the metadata is enclosed between lines like this: '---'
			// (Jekyll style) we need to return after process it.
			if hasMetadata == true {
				return true, nil
			}
			hasMetadata = true
		case "tags":
			// Remove all the spaces between comma and tag and
			// add one comma at the beginning and other at the end, this will
			// make the querying much simpler
			for _, tag := range strings.Split(value, ",") {
				pf.Tags = append(pf.Tags, strings.Replace(tag, " ", "", -1))
			}
			hasMetadata = true
		case "date":
			pf.Date, err = parseData(value)
			if err != nil {
				return true, err
			}
			hasMetadata = true
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
			hasMetadata = true
		case "title":
			pf.Title = value
			hasMetadata = true
		case "status":
			pf.status = value
			hasMetadata = true
		case "summary":
			pf.summary = value
			hasMetadata = true
		case "author":
			pf.Author = value
			hasMetadata = true
		default:
			return hasMetadata, nil
		}
	}

	return hasMetadata, nil
}

// Loads the content of the file object from the given filename.
func (pf *ParsedFile) load(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	pf.scanner = bufio.NewScanner(file)
	hasMetadata, err := pf.parseMetadata()
	if err != nil {
		return err
	}
	if !hasMetadata {
		// Rewind the file and reset the scanner
		file.Seek(0, 0)
		pf.scanner = bufio.NewScanner(file)
	}

	isFirstLine := true
	for pf.scanner.Scan() {
		line := pf.scanner.Text()
		if isFirstLine == true {
			if line == "" {
				// Do not read empty lines at the beginning
				continue
			}

			if pf.Title == "" {
				// Needed to remove markdown syntax before storing values on the structs
				re := regexp.MustCompile("^#+\\s*")
				pf.Title = re.ReplaceAllString(line, "")
			}
			if pf.Slug == "" {
				prefix := ""
				if pf.IsPage {
					prefix = "/pages"
				}
				pf.Slug = fmt.Sprintf("%s/%s.html", prefix, utils.Slugify(pf.Title))
			}
			pf.scanner.Scan() // We don't want the title underlining

			isFirstLine = false
		} else {
			pf.Content += line + "\n"
		}
	}

	// Set the category from the filePath
	splittedPath := strings.Split(filePath, "/")
	length := len(splittedPath)
	if length > 1 {
		pf.Category = splittedPath[len(splittedPath)-2]
	}

	return nil
}
