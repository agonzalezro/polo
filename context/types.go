package context

import (
	"github.com/agonzalezro/polo/config"
	"github.com/agonzalezro/polo/file"
)

// Context stores the temporal context needed to render a page.
type Context struct {
	Article, Page   file.ParsedFile
	Articles, Pages []*file.ParsedFile

	Tags          []string
	Tag, Category string
	PageNumber    int

	Updated string

	Site SiteContext
}

// SiteContext is a subcontext where we store the Site globals and the
// configuration.
type SiteContext struct {
	Articles, Pages []*file.ParsedFile
	Tags            []string
	Categories      []string

	NumberOfPages int

	Config config.Config
}
