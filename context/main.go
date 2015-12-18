package context

import (
	"sync"

	"github.com/agonzalezro/polo/file"
)

type Context struct {
	Pages    []file.ParsedFile
	Articles []file.ParsedFile

	Tags, Categories []string

	tagUniqueness, categoryUniqueness *sync.Mutex
}

func New() *Context {
	return &Context{
		tagUniqueness:      &sync.Mutex{},
		categoryUniqueness: &sync.Mutex{},
	}
}

// AppendUniqueTags will append the tag only if it's not already on the context.
func (c *Context) AppendUniqueTags(newTags []string) {
	c.tagUniqueness.Lock()
	defer c.tagUniqueness.Unlock()

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
	c.categoryUniqueness.Lock()
	defer c.categoryUniqueness.Unlock()

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
