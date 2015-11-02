package site

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"config"
	"file"
)

// Site is the struct that we will during the execution of the program. It will
// be our storage, config...
type Site struct {
	Config     config.Config
	outputPath string

	Articles, Pages []*file.ParsedFile
	Categories      []string
	Tags            []string

	mux sync.Mutex
}

// New returns a new Site object.
func New(config config.Config, outputPath string) *Site {
	site := &Site{Config: config, outputPath: outputPath}
	return site
}

func newUniqueSlugError(slug string) error {
	errorMessage := fmt.Sprintf("The slug '%s' already exists!", slug)
	return errors.New(errorMessage)
}

// AppendUniqueTags will append the tag only if it's not already on the site
// tags.
func (site *Site) AppendUniqueTags(newTags []string) {
LOOP:
	for _, newTag := range newTags {
		for _, oldTag := range site.Tags {
			if newTag == oldTag {
				continue LOOP
			}
		}
		site.Tags = append(site.Tags, newTag)
	}
}

// AppendUniqueCategory will append the category just in case that doesn't
// belong to the Site yet.
func (site *Site) AppendUniqueCategory(newCategory string) {
	for _, category := range site.Categories {
		if category == newCategory {
			return
		}
	}
	site.Categories = append(site.Categories, newCategory)
}

// Populate will read the file system searching for markdown and store those on
// the Site object.
func (site *Site) Populate(root string) error {
	if err := filepath.Walk(root, site.parseAndStore); err != nil {
		return err
	}
	// Sort the articles after we got them
	sort.Sort(site)
	return nil
}

// Append the paths to an array in case that they are markdown files.
func (site *Site) parseAndStore(path string, fileInfo os.FileInfo, inputErr error) (err error) {
	if inputErr != nil {
		return inputErr
	}

	slugsPresence := make(map[string]bool)

	if !fileInfo.Mode().IsDir() && strings.HasSuffix(path, ".md") {
		file, err := file.New(path)
		if err != nil {
			return err
		}

		if _, present := slugsPresence[file.Slug]; present {
			return newUniqueSlugError(file.Slug)
		}
		slugsPresence[file.Slug] = true

		// Add the pages or articles to the proper array on the site
		if file.IsPage {
			site.Pages = append(site.Pages, file)
		} else {
			site.Articles = append(site.Articles, file)

			// Just supported on Articles
			if len(file.Tags) > 0 {
				site.AppendUniqueTags(file.Tags)
			}
			site.AppendUniqueCategory(file.Category)
		}
	}
	return
}
