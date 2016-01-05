package site

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"

	"github.com/agonzalezro/polo/config"
	"github.com/agonzalezro/polo/context"
	"github.com/agonzalezro/polo/file"
)

const (
	atomPath           = "feeds/all.atom.xml"
	archivePath        = "archive.html"
	articlesPrefixPath = "" // TODO: perhaps allow configuration for this?
	pagesPrefixPath    = "pages/"

	categoryPathFormater = "category/%s.html"
	tagPathFormater      = "tag/%s.html"

	indexTemplate    = "index"
	atomTemplate     = "atom"
	articleTemplate  = "article"
	pageTemplate     = "page"
	archiveTemplate  = "archive"
	categoryTemplate = "category"
	tagTemplate      = "tag"
)

type Siteable interface {
	Load() error
	Write() error
}

type Site struct {
	port              int
	source, output    string
	templatesBasePath string

	slugs map[string]bool
	mux   *sync.Mutex

	Config  config.Config
	Context *context.Context
}

func New(source, output, configPath, templatesBasePath string) (*Site, error) {
	config, err := config.New(configPath)
	if err != nil {
		return nil, err
	}

	s := Site{
		source:            source,
		output:            output,
		templatesBasePath: templatesBasePath,
		Config:            *config,
		Context:           context.New(*config),
		mux:               &sync.Mutex{},
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
	var wg sync.WaitGroup

	errCh := make(chan error, 1024) // TODO: 1024 should be more than enough for the errors

	s.writeIndexes(&wg, errCh)
	s.writeFeeds(&wg, errCh)
	s.writeArticles(&wg, errCh)
	s.writePages(&wg, errCh)

	if s.Config.ShowArchive {
		s.writeArchive(&wg, errCh)
	}
	if s.Config.ShowCategories {
		s.writeCategories(&wg, errCh)
	}
	if s.Config.ShowTags {
		s.writeTags(&wg, errCh)
	}

	wg.Wait()

	for {
		select {
		case err := <-errCh:
			// I could append errors here but it isn't worthy
			return err
		default:
			return nil
		}
	}
}

func (s Site) writef(relativePath string, templateName string, c context.Context) error {
	tpl, err := s.getTemplate(templateName)
	if err != nil {
		return err
	}

	// Ensure absolute path exists
	err = s.mkdirP(relativePath)
	if err != nil {
		return err
	}

	f, err := os.Create(path.Join(s.output, relativePath))
	if err != nil {
		return err
	}

	return tpl.ExecuteTemplate(f, "base", c)
}

func (s Site) writeIndexes(wg *sync.WaitGroup, errCh chan<- error) {
	for i := 1; i <= s.Context.NumberOfPages(); i++ {
		wg.Add(1)
		go func(page int) {
			defer wg.Done()

			indexFile := "index.html"
			if page > 1 {
				indexFile = fmt.Sprintf("index%d.html", page)
			}

			c := s.Context.Copy()
			c.Articles = c.FilterByPage(page)

			c.CurrentPage = page

			if err := s.writef(indexFile, indexTemplate, *c); err != nil {
				errCh <- err
			}
		}(i)
	}
}

func (s Site) writeFeeds(wg *sync.WaitGroup, errCh chan<- error) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		limit := len(s.Context.Articles)
		if limit > 10 {
			limit = 10
		}

		c := s.Context.Copy()
		c.Articles = c.Articles[0:limit]

		// TODO: just ATOM feeds, not RSS unless somebody needs it and he/she is willing to implement it :)
		if err := s.writef(atomPath, atomTemplate, *c); err != nil {
			errCh <- err
		}
	}()
}

func (s Site) writeArticles(wg *sync.WaitGroup, errCh chan<- error) {
	for _, article := range s.Context.Articles {
		wg.Add(1)
		go func(article file.ParsedFile) {
			defer wg.Done()

			c := s.Context.Copy()
			c.Article = article

			p := path.Join(articlesPrefixPath, article.Slug)
			if err := s.writef(p, articleTemplate, *c); err != nil {
				errCh <- err
			}
		}(article)
	}
}

func (s Site) writePages(wg *sync.WaitGroup, errCh chan<- error) {
	for _, page := range s.Context.Pages {
		wg.Add(1)
		go func(page file.ParsedFile) {
			defer wg.Done()

			c := s.Context.Copy()
			c.Page = page

			p := path.Join(pagesPrefixPath, page.Slug)
			if err := s.writef(p, pageTemplate, *c); err != nil {
				errCh <- err
			}
		}(page)
	}
}

func (s Site) writeArchive(wg *sync.WaitGroup, errCh chan<- error) {
	wg.Add(1)
	defer wg.Done()

	if err := s.writef(archivePath, archiveTemplate, *s.Context); err != nil {
		errCh <- err
	}
}

func (s Site) writeCategories(wg *sync.WaitGroup, errCh chan<- error) {
	for _, category := range s.Context.Categories {
		wg.Add(1)
		go func(category string) {
			defer wg.Done()

			c := s.Context.Copy()
			c.Articles = c.FilterByCategory(category)
			c.Category = category

			p := fmt.Sprintf(categoryPathFormater, category)
			if err := s.writef(p, categoryTemplate, *c); err != nil {
				errCh <- err
			}
		}(category)
	}
}

func (s Site) writeTags(wg *sync.WaitGroup, errCh chan<- error) {
	for _, tag := range s.Context.Tags {
		wg.Add(1)
		go func(tag string) {
			defer wg.Done()

			c := s.Context.Copy()
			c.Articles = c.FilterByTag(tag)
			c.Tag = tag

			p := fmt.Sprintf(tagPathFormater, tag)
			if err := s.writef(p, tagTemplate, *c); err != nil {
				errCh <- err
			}
		}(tag)
	}
}

// mkdirP assures that the full path exists. It mimics `mkdir -p`
func (s *Site) mkdirP(elem ...string) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	// TODO: I am not pretty sure that this is the best way to do this, but it works (TM)
	toJoin := make([]string, 1, 1)
	toJoin[0] = s.output
	elem = append(toJoin, elem...)
	absolutePath := path.Join(elem...)

	dir := filepath.Dir(absolutePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0777); err != nil {
			return err
		}
	}

	return nil
}
