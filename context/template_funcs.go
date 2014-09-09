package context

import (
	"fmt"
	"time"

	"github.com/agonzalezro/polo/file"
)

// ArticlesByPage returns the articles of the site paginated.
func (c Context) ArticlesByPage(page int) []*file.ParsedFile {
	var start, end int

	paginationSize := c.Site.Config.PaginationSize

	start = (page - 1) * paginationSize

	if start+paginationSize <= len(c.Site.Articles) {
		end = start + paginationSize
	} else {
		end = len(c.Site.Articles)
	}

	return c.Site.Articles[start:end]
}

// ArticlesByTag returns all the articles belonging to the given tag.
func (c Context) ArticlesByTag(tag string) []*file.ParsedFile {
	var files []*file.ParsedFile

LOOP:
	for _, f := range c.Site.Articles {
		for _, t := range f.Tags {
			if tag == t {
				files = append(files, f)
				continue LOOP
			}
		}
	}

	return files
}

// ArticlesByCategory returns all the articles belonging to the give category.
func (c Context) ArticlesByCategory(category string) []*file.ParsedFile {
	var files []*file.ParsedFile

	for _, f := range c.Site.Articles {
		if f.Category == category {
			files = append(files, f)
		}
	}

	return files
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
	if page == c.Site.NumberOfPages {
		return "#"
	}

	return fmt.Sprintf("/index%d.html", page+1)
}

// ArrayOfPages is a dirty hack because we can not (or I don't know how) do a
// range from X to Y on the template
func (c Context) ArrayOfPages() (pages []int) {
	for i := 1; i < c.Site.NumberOfPages+1; i++ {
		pages = append(pages, i)
	}
	return pages
}

// HumanizeDatetime returns a date or datetime depending of the datetime
// received.
// For example, if the datetime received doesn't have any hour/minutes, the
// hours/minutes part doesn't need to be shown.
func (c Context) HumanizeDatetime(datetime time.Time) string {
	if datetime.Hour()+datetime.Minute() == 0 {
		return datetime.Format("2006-01-02")
	}
	return datetime.Format("2006-01-02 15:04")
}
