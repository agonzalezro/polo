package site

import (
	"time"

	"context"
	"file"
)

func (s Site) newSiteContext() context.SiteContext {
	return context.SiteContext{
		s.Articles, s.Pages,
		s.Tags,
		s.Categories,
		s.getNumberOfPages(),
		s.Config,
	}
}

// NewContext will create a common Context object to be used by other contexts.
func (s Site) NewContext() *context.Context {
	contextSite := s.newSiteContext()
	return &context.Context{
		Site: contextSite,
	}
}

// NewArticleContext return the context needed to render an Article document.
func (s Site) NewArticleContext(f file.ParsedFile) *context.Context {
	c := s.NewContext()
	c.Article = f
	return c
}

// NewPageContext returns the context needed to render a Page document.
func (s Site) NewPageContext(f file.ParsedFile) *context.Context {
	c := s.NewContext()
	c.Page = f
	return c
}

// NewAtomFeedContext returns the context needed to render the Atom feed.
func (s Site) NewAtomFeedContext(limit int) *context.Context {
	c := s.NewContext()
	c.Updated = time.Now().Format(time.RFC3339)
	c.Articles = s.Articles[:limit]
	return c
}

// NewTagContext returns the context needed to render a tag page.
func (s Site) NewTagContext(tag string) *context.Context {
	c := s.NewContext()
	c.Tag = tag
	c.Articles = c.ArticlesByTag(tag)
	return c
}

// NewCategoryContext returns the context needed to render a Category page.
func (s Site) NewCategoryContext(category string) *context.Context {
	c := s.NewContext()
	c.Category = category
	c.Articles = c.ArticlesByCategory(category)
	return c
}

// NewPaginatedContext returns the context needed to render an Index page.
func (s Site) NewPaginatedContext(page int) *context.Context {
	c := s.NewContext()
	c.PageNumber = page
	c.Articles = c.ArticlesByPage(page)
	return c
}
