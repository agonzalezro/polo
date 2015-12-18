package site

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/agonzalezro/polo/context"
	"github.com/agonzalezro/polo/file"
)

type Siteable interface {
	Load() error
	Write() error
}

type Site struct {
	port           int
	source, output string
	config         Config

	slugs map[string]bool
	mux   *sync.Mutex

	Context *context.Context
}

func New(source, output, configPath string) (*Site, error) {
	config, err := NewConfig(configPath)
	if err != nil {
		return nil, err
	}

	s := Site{
		source:  source,
		output:  output,
		config:  *config,
		Context: context.New(),
		mux:     &sync.Mutex{},
	}

	return &s, s.Load()
}

func (s *Site) AddSlug(slug string) {
	s.mux.Lock()
	defer s.mux.Unlock()
	if s.slugs == nil {
		s.slugs = make(map[string]bool)
	}
	s.slugs[slug] = true
}

// parse is a filepath.WalkFunc that will load all the pages and articles into a context object.
func (s *Site) parse(path string, fileInfo os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	if fileInfo.Mode().IsDir() || !file.IsMarkdown(path) {
		return nil
	}

	file, err := file.New(path)
	if err != nil {
		return err
	}

	if _, present := s.slugs[file.Slug]; present {
		return fmt.Errorf("The slug '%s' already exists!", file.Slug)
	}
	s.AddSlug(file.Slug)

	// Add the pages or articles to the proper array on the site
	if file.IsPage {
		s.Context.Pages = append(s.Context.Pages, *file)
		return nil
	}

	// If it's not a page, it's an article
	s.Context.Articles = append(s.Context.Articles, *file)

	// Just supported on Articles
	if len(file.Tags) > 0 {
		s.Context.AppendUniqueTags(file.Tags)
	}
	s.Context.AppendUniqueCategory(file.Category)
	return nil
}

func (s *Site) Load() error {
	if err := filepath.Walk(s.source, s.parse); err != nil {
		return err
	}
	// Sort the articles after we got them
	sort.Sort(s.Context)
	return nil
}

func (s Site) Write() error {
	for _, article := range s.Context.Articles {
		fmt.Printf("%+v", article)
	}
	return nil
}
