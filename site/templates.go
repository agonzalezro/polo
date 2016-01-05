package site

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"

	log "github.com/Sirupsen/logrus"
)

const TemplatesRelativePath = "templates"

var (
	templates map[string]*template.Template

	contentTemplatePaths map[string][]string // As for example an article, that needs disqus, analytics...
)

func init() {
	contentTemplatePaths = make(map[string][]string)
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
}

// TODO: refactor this monster
// TODO: use Asset
// Also not sure if it's better to return templates or just change the global var
func (s *Site) getTemplates() (map[string]*template.Template, error) {
	s.mux.Lock()
	defer s.mux.Unlock()

	if templates != nil {
		return templates, nil
	}

	templates = make(map[string]*template.Template)
	templatesPath := path.Join(s.templatesBasePath, TemplatesRelativePath)

	commonTpl := template.New("common")

	walkFn := func(p string, info os.FileInfo, err error) error {
		isNoContentTemplate := !info.IsDir() && !strings.Contains(p, "content/") && strings.HasSuffix(p, ".tmpl")
		if isNoContentTemplate {
			log.Debugf("Loading template: %s", p)
			commonTpl, err = commonTpl.ParseFiles(p) // TODO: Use Asset here as well
			if err != nil {
				log.Warning(err)
			}
		}
		return err
	}

	filepath.Walk(templatesPath, walkFn)

	for name, paths := range contentTemplatePaths {
		tpl, err := commonTpl.Clone()
		if err != nil {
			return nil, err
		}
		// tpl, err := tpl.AddParseTree("base", commonTpl.Tree)
		// if err != nil {
		// 	return nil, err
		// }
		for _, p := range paths {
			p = path.Join(s.templatesBasePath, TemplatesRelativePath, "body", "content", p)
			log.Debugf("Loading template: %s", p)
			tpl, err = tpl.ParseFiles(p) // TODO: use Assets
			if err != nil {
				return nil, err
			}
		}
		templates[name] = tpl
	}

	tpl := template.New("atom")
	tpl, err := tpl.ParseFiles("templates/atom.tmpl") // TODO: use Asset
	if err != nil {
		return nil, err
	}
	templates[atomTemplate] = tpl

	return templates, nil
}

func (s *Site) getTemplate(name string) (*template.Template, error) {
	templates, err := s.getTemplates()
	if err != nil {
		return nil, err
	}

	if v, ok := templates[name]; ok {
		return v, nil
	}
	return nil, fmt.Errorf("Template '%s' not found!", name)
}
