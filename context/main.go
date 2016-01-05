package context

import (
	"fmt"
	"math"
	"sync"
	"time"

	"github.com/agonzalezro/polo/config"
	"github.com/agonzalezro/polo/file"
)

type Context struct {
	Pages    []file.ParsedFile
	Articles []file.ParsedFile

	Tags, Categories []string

	Config  config.Config
	Updated string

	// Temporal stuff for template rendering
	Article     file.ParsedFile
	Page        file.ParsedFile
	Tag         string
	Category    string
	CurrentPage int

	tagUniquenessMux, categoryUniquenessMux, numberOfPagesMux *sync.Mutex
}

func New(config config.Config) *Context {
	return &Context{
		Config:                config,
		Updated:               time.Now().Format(time.RFC3339),
		tagUniquenessMux:      &sync.Mutex{},
		categoryUniquenessMux: &sync.Mutex{},
		numberOfPagesMux:      &sync.Mutex{},
	}
}

func (c *Context) Copy() *Context {
	return &Context{
		Config:     c.Config,
		Pages:      c.Pages,
		Articles:   c.Articles,
		Tags:       c.Tags,
		Categories: c.Categories,
	}
}

// TODO: fix this; we need a number of pages that doesn't fluctuate with the number of articles in the current context
var numberOfPages int

func (c Context) NumberOfPages() int {
	if numberOfPages != 0 {
		return numberOfPages
	}

	// TODO: we can't use a pointer to context so I doubt that this mutex does anything at all
	c.numberOfPagesMux.Lock()
	defer c.numberOfPagesMux.Unlock()

	numberOfPages = int(
		math.Ceil(
			float64(len(c.Articles)) / float64(c.Config.PaginationSize)))
	return numberOfPages
}

// PreviousSlug "calculates" the previous index slug given the page number.
func (c Context) PreviousSlug(page int) (slug string) {
	switch page {
	case 1:
		slug = "#"
	case 2:
		slug = "/index.html"
	default:
		slug = fmt.Sprintf("/index%d.html", page-1)

	}
	return slug
}

// NextSlug "calculates" the next index slug given the page number.
func (c Context) NextSlug(page int) string {
	if page == c.NumberOfPages() {
		return "#"
	}

	return fmt.Sprintf("/index%d.html", page+1)
}

// AppendUniqueTags will append the tag only if it's not already on the context.
func (c *Context) AppendUniqueTags(newTags []string) {
	c.tagUniquenessMux.Lock()
	defer c.tagUniquenessMux.Unlock()

LOOP:
	for _, newTag := range newTags {
		for _, oldTag := range c.Tags {
			if newTag == oldTag {
				continue LOOP
			}
		}
		c.Tags = append(c.Tags, newTag)
	}
}

// AppendUniqueCategory will append the category just in case that doesn't
// belong to the Context yet.
func (c *Context) AppendUniqueCategory(newCategory string) {
	c.categoryUniquenessMux.Lock()
	defer c.categoryUniquenessMux.Unlock()

	for _, category := range c.Categories {
		if category == newCategory {
			return
		}
	}
	c.Categories = append(c.Categories, newCategory)
}

// Len is needed to implement the sorting interface.
func (c Context) Len() int {
	return len(c.Articles)
}

// Less is a comparator to help us to sort the context Articles by date DESC.
func (c Context) Less(i, j int) bool {
	return c.Articles[i].Date.After(c.Articles[j].Date)
}

// Swap is needed to implement the sorting interface.
func (c Context) Swap(i, j int) {
	c.Articles[i], c.Articles[j] = c.Articles[j], c.Articles[i]
}
