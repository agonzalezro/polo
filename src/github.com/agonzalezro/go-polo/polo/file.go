package polo

import (
	"bufio"
	"errors"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"

	"github.com/agonzalezro/go-polo/polo/utils"
	"github.com/russross/blackfriday"
)

type ParsedFile struct {
	Metadata map[string]string

	Title   string
	Slug    string
	Content string
	isPage  bool

	tags string
	Date string

	scanner *bufio.Scanner
}

/*
 * Set the know metadata to the current object.
 *
 * The supported metadata for now is:
 *   - tags: going to be transformed to a []string
 *   - date: a string in ISO8601 format
 */
func (pf *ParsedFile) parseMetadata() error {
	for pf.scanner.Scan() {
		line := pf.scanner.Text()

		if strings.HasPrefix(line, "---") {
			pf.scanner.Scan() // Read the last ---
			return nil
		}

		metadataSplited := strings.Split(line, ":")
		key := strings.ToLower(metadataSplited[0])
		value := strings.Trim(strings.Join(metadataSplited[1:], ":"), " ")
		switch {
		case key == "tags":
			pf.tags = value
		case key == "date":
			pf.Date = value
		}
		// It's possible that : is on the value of the metadata too (example: a date)
	}

	return errors.New("The metadata section was not properly closed!")
}

// Loads the content of the file object from the given filename.
func (pf *ParsedFile) load(filePath string) {
	isFirstLine := true

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	pf.scanner = bufio.NewScanner(file)
	for pf.scanner.Scan() {
		line := pf.scanner.Text()
		if line == "---" {
			if err := pf.parseMetadata(); err != nil {
				log.Fatal(err)
			}
			continue
		}

		if isFirstLine {
			pf.Title = line
			// TODO: check if there is a slug key on the metadata and don't assign it in that case
			pf.Slug = utils.Slugify(line)
			pf.scanner.Scan() // We don't want the title underlining

			isFirstLine = false
		} else {
			pf.Content += line + "\n"
		}

	}
}

// Split the tags into a list.
func (pf *ParsedFile) Tags() []string {
	return strings.Split(pf.tags, ",")
}

// Function to be called from the templates. It render safe HTML code.
func (file ParsedFile) Html(content string) template.HTML {
	html := blackfriday.MarkdownCommon([]byte(content))
	return template.HTML(html)
}

// Store the file in a "permanent" storage.
func (file ParsedFile) save(db *DB) error {
	query := `
    INSERT INTO files (title, slug, content, tags, date, is_page)
    VALUES ("%s", "%s", "%s", "%s", "%s", %b)
    `

	// SQLite doesn't support booleans :(
	isPageInt := 0
	if file.isPage {
		isPageInt = 1
	}

	filledQuery := fmt.Sprintf(query, file.Title, file.Slug, file.Content, file.tags, file.Date, isPageInt)
	if _, err := db.connection.Exec(filledQuery); err != nil {
		return err
	}
	return nil
}
