package site

import (
	"fmt"
	"os"
	"path"
	"text/template"

	log "github.com/Sirupsen/logrus"
	assets "github.com/agonzalezro/polo/templates"
)

// TODO: probably to be override from the cmd
const TemplatesRelativePath = "templates"

var templates map[string]*template.Template

// TODO: this could be our own type based on template.Template
func parseFileOrAsset(t *template.Template, p string) (*template.Template, error) {
	if _, err := os.Stat(p); err == nil {
		log.Debug("Loading template from disk: ", p)
		return t.ParseFiles(p)
	}

	b, err := assets.Asset(p)
	if err != nil {
		return nil, err
	}
	log.Debug("Loading template from asset: ", p)
	return t.Parse(string(b))
}

// mustParseCommonTemplates will parse common templates and fatal if it can.
// Common templates are those templates shared.
func mustParseCommonTemplates() *template.Template {
	// We can't just walk a dir to find them because GOPATH will not be available
	// on a binary installation.
	// TODO: it would be cool to use go-bindata list but _bindata is private.
	commonTemplatePaths := []string{
		"templates/base.tmpl",
		"templates/body/analytics.tmpl",
		"templates/body/footer.tmpl",
		"templates/body/footer_scripts.tmpl",
		"templates/body/navbar.tmpl",
		"templates/head/header.tmpl",
		"templates/head/header_scripts.tmpl",
		"templates/head/share_this.tmpl",
	}

	tpl := template.New("common")
	var err error
	for _, p := range commonTemplatePaths {
		tpl, err = parseFileOrAsset(tpl, p)
		if err != nil {
			log.Fatal(err)
		}
	}
	return tpl
}

func mustParseContentTemplates(commonTpl *template.Template, templatesPath string) map[string]*template.Template {
	// Content templates are those templates that are willing to change depending
	// on what we are rending at that moment.
	contentTemplatePaths := make(map[string][]string)
	for templateName, paths := range map[string][]string{
		articleTemplate:  []string{"article/article.tmpl", "article/disqus.tmpl", "article/share_icons.tmpl"},
		archiveTemplate:  []string{"archive.tmpl"},
		categoryTemplate: []string{"category.tmpl"},
		indexTemplate:    []string{"index.tmpl"},
		pageTemplate:     []string{"page.tmpl"},
		tagTemplate:      []string{"tag.tmpl"},
	} {
		contentTemplatePaths[templateName] = paths
	}

	templates = make(map[string]*template.Template) // WARNING package level var
	for name, paths := range contentTemplatePaths {
		tpl, err := commonTpl.Clone()
		if err != nil {
			log.Fatal(err)
		}

		for _, p := range paths {
			p = path.Join(templatesPath, "body", "content", p)
			tpl, err = parseFileOrAsset(tpl, p)
			if err != nil {
				log.Fatal(err)
			}
		}
		templates[name] = tpl
	}

	// Atom template doesn't inherit from any shared template
	tpl := template.New("atom")
	tpl, err := parseFileOrAsset(tpl, "templates/atom.tmpl")
	if err != nil {
		log.Fatal(err)
	}
	templates[atomTemplate] = tpl

	return templates
}

func (s *Site) getTemplates() map[string]*template.Template {
	s.mux.Lock()
	defer s.mux.Unlock()

	if templates != nil {
		return templates
	}

	templatesPath := path.Join(s.templatesBasePath, TemplatesRelativePath)
	log.Debug("Templates path: ", templatesPath)

	commonTpl := mustParseCommonTemplates()
	return mustParseContentTemplates(commonTpl, templatesPath)
}

func (s *Site) getTemplate(name string) (*template.Template, error) {
	if v, ok := s.getTemplates()[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("Template '%s' not found!", name)
}
