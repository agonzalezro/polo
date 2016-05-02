package context

import "github.com/agonzalezro/polo/file"

// FilterByPage returns the articles of the site paginated.
func (c Context) FilterByPage(page int) []file.ParsedFile {
	var start, end int

	paginationSize := c.Config.PaginationSize
	start = (page - 1) * paginationSize

	if start > len(c.Articles) || start < 0 {
		return []file.ParsedFile{}
	}

	if start+paginationSize <= len(c.Articles) {
		end = start + paginationSize
	} else {
		end = len(c.Articles)
	}

	return c.Articles[start:end]
}

// FilterByTag returns all the articles belonging to the given tag.
func (c Context) FilterByTag(tag string) []file.ParsedFile {
	var files []file.ParsedFile

LOOP:
	for _, f := range c.Articles {
		for _, t := range f.Tags {
			if tag == t {
				files = append(files, f)
				continue LOOP
			}
		}
	}

	return files
}

// FilterByCategory returns all the articles belonging to the give category.
func (c Context) FilterByCategory(category string) []file.ParsedFile {
	var files []file.ParsedFile

	for _, f := range c.Articles {
		if f.Category == category {
			files = append(files, f)
		}
	}

	return files
}
