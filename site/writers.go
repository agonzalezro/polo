package site

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path"
	"path/filepath"
	"strings"
	"sync"

	"github.com/agonzalezro/polo/context"
	"github.com/agonzalezro/polo/file"

	// TODO: Perhaps it worth moving the template rendering to the template
	// package
	assets "github.com/agonzalezro/polo/templates"
)

var templates map[string]*template.Template

// parseFiles is a wrapper similar to template.ParseFiles that is going to
// load the templates from the disk, and if they can not be found from the
// go-bindata file.
func parseFiles(filenames ...string) (*template.Template, error) {
	tpl := template.New(filenames[0])
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			b, err = assets.Asset(filename)
			if err != nil {
				log.Printf("Template: %s not found. Not in HD neihter in bindata!", filename)
				return nil, err
			}
		}
		tpl.Parse(string(b))
	}
	return tpl, nil
}

// loadTemplates is an ugly function but I need it to run the test without the
// template files. If I don't call .Write() I don't need the template files.
func loadTemplates() {
	templates = make(map[string]*template.Template)
	toRender := make(map[string][]string)

	toRender["atom"] = []string{"templates/atom.xml"}

	getAllTemplateInheritance := func(templates []string) []string {
		return append([]string{
			"templates/base/base.html",
			"templates/base/header.html",
			"templates/base/footer.html",
			"templates/base/analytics.html"}, templates...)
	}

	toRender["archives"] = getAllTemplateInheritance([]string{"templates/archives.html"})
	toRender["article"] = getAllTemplateInheritance([]string{"templates/article/article.html", "templates/article/disqus.html", "templates/article/sharethis.html"})
	toRender["category"] = getAllTemplateInheritance([]string{"templates/category.html"})
	toRender["index"] = getAllTemplateInheritance([]string{"templates/index.html"})
	toRender["page"] = getAllTemplateInheritance([]string{"templates/page.html"})
	toRender["tag"] = getAllTemplateInheritance([]string{"templates/tag.html"})

	for name, values := range toRender {
		templates[name] = template.Must(parseFiles(values...))
	}
}

// createAbsolutePath assures that the full dir tree exists and return the
// point to the file
func (site Site) createAbsolutePath(elem ...string) (file *os.File, err error) {
	// TODO: I am not pretty sure that this is the best way to do this
	s := make([]string, 1, 1)
	s[0] = site.outputPath
	elem = append(s, elem...)
	absolutePath := path.Join(elem...)

	dir := filepath.Dir(absolutePath)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.Mkdir(dir, 0777); err != nil {
			return nil, err
		}
	}

	file, err = os.Create(absolutePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

// Auxiliar function to quickly log the quantity of items created
func logCreation(noun string, quantity int) string {
	var pluralize string
	if quantity > 1 {
		pluralize = "s"
	}
	return fmt.Sprintf("%5d %s%s created", quantity, strings.Title(noun), pluralize)
}

// Dump all the site content to the disk
func (site Site) Write() (err error) {
	loadTemplates()

	var wg sync.WaitGroup
	chErrors := make(chan error, 12) // Enough buffer for all the subroutines

	wg.Add(4) // The following 4 subroutines are mandatory

	go func() {
		defer wg.Done()
		i, err := site.writeIndexes()
		fmt.Println(logCreation("index page", i))
		if err != nil {
			log.Panic(err)
		}
		chErrors <- err
	}()

	go func() {
		defer wg.Done()
		i, err := site.writeFeeds()
		fmt.Println(logCreation("feed", i))
		chErrors <- err
	}()

	go func() {
		defer wg.Done()
		i, err := site.writeArticles()
		fmt.Println(logCreation("article", i))
		chErrors <- err
	}()

	go func() {
		defer wg.Done()
		i, err := site.writePages()
		fmt.Println(logCreation("page", i))
		chErrors <- err
	}()

	if site.Config.ShowArchive {
		wg.Add(1)
		go func() {
			defer wg.Done()
			err := site.writeArchive()
			fmt.Println(logCreation("archive", 1)) // If error it will not be 1, but we don't care
			chErrors <- err
		}()
	}
	if site.Config.ShowCategories {
		wg.Add(1)
		go func() {
			defer wg.Done()
			i, err := site.writeCategories()
			fmt.Println(logCreation("category page", i))
			chErrors <- err
		}()
	}
	if site.Config.ShowTags {
		wg.Add(1)
		go func() {
			defer wg.Done()
			i, err := site.writeTags()
			fmt.Println(logCreation("tag page", i))
			chErrors <- err
		}()
	}

	wg.Wait()

LOOP:
	for {
		select {
		case err := <-chErrors:
			if err != nil {
				return err
			}
		default:
			break LOOP
		}
	}

	return nil
}

func (site Site) getNumberOfPages() int {
	return int(
		math.Ceil(
			float64(len(site.Articles)) / float64(site.Config.PaginationSize)))
}

func (site Site) writeIndexes() (int, error) {
	var pageNumber int

	for pageNumber = 1; pageNumber <= site.getNumberOfPages(); pageNumber++ {
		indexFile := "index.html"
		if pageNumber > 1 {
			indexFile = fmt.Sprintf("index%d.html", pageNumber)
		}

		file, err := site.createAbsolutePath(indexFile)
		if err != nil {
			return pageNumber, err
		}

		if err := templates["index"].ExecuteTemplate(file, "base", site.NewPaginatedContext(pageNumber)); err != nil {
			return pageNumber, err
		}
	}

	return pageNumber, nil
}

type contextCreator func(f file.ParsedFile) *context.Context

func (site Site) writeparsedFiles(
	files []*file.ParsedFile, template *template.Template, newContext contextCreator) (int, error) {

	var (
		parsedFile *file.ParsedFile
		wg         sync.WaitGroup
	)
	errChan := make(chan error, len(files))

	for _, parsedFile = range files {
		wg.Add(1)
		go func(parsedFile *file.ParsedFile) {
			defer wg.Done()

			file, err := site.createAbsolutePath(parsedFile.Slug)
			if err != nil {
				errChan <- err
				return
			}

			context := newContext(*parsedFile)
			if err := template.ExecuteTemplate(file, "base", context); err != nil {
				errChan <- err
			}
		}(parsedFile)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		if err != nil {
			return -1, err
		}
	default:
		break
	}

	return len(files), nil
}

func (site Site) writeArticles() (int, error) {
	return site.writeparsedFiles(site.Articles, templates["article"], site.NewArticleContext)
}

func (site Site) writePages() (int, error) {
	return site.writeparsedFiles(site.Pages, templates["page"], site.NewPageContext)
}

func (site Site) writeArchive() error {
	file, err := site.createAbsolutePath("archives.html")
	if err != nil {
		return err
	}

	if err := templates["archives"].ExecuteTemplate(file, "base", site.NewContext()); err != nil {
		return err
	}

	return nil
}

func (site Site) writeCategories() (int, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(site.Categories))

	for _, category := range site.Categories {
		wg.Add(1)
		go func(category string) {
			defer wg.Done()

			categoryFile := fmt.Sprintf("category/%s.html", category)
			file, err := site.createAbsolutePath(categoryFile)
			if err != nil {
				errChan <- err
				return
			}

			if err := templates["category"].ExecuteTemplate(file, "base", site.NewCategoryContext(category)); err != nil {
				errChan <- err
			}

			errChan <- nil
		}(category)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		if err != nil {
			return -1, err
		}
	default:
		break
	}

	return len(site.Categories), nil
}

func (site Site) writeTags() (int, error) {
	var wg sync.WaitGroup
	errChan := make(chan error, len(site.Tags))

	for _, tag := range site.Tags {
		wg.Add(1)
		go func(tag string) {
			defer wg.Done()

			tagFile := fmt.Sprintf("tag/%s.html", tag)
			file, err := site.createAbsolutePath(tagFile)
			if err != nil {
				errChan <- err
				return
			}

			if err := templates["tag"].ExecuteTemplate(file, "base", site.NewTagContext(tag)); err != nil {
				errChan <- err
			}

			errChan <- nil
		}(tag)
	}

	wg.Wait()

	select {
	case err := <-errChan:
		if err != nil {
			return -1, err
		}
	default:
		break
	}

	return len(site.Tags), nil
}

func (site Site) writeFeeds() (int, error) {
	var i int

	if err := site.writeAtomFeed(); err != nil {
		return i + 1, err
	}

	if err := site.writeRSSFeed(); err != nil {
		return i + 1, err
	}

	i++ // Not implemented yet, fake it
	return i, nil
}

func (site Site) writeAtomFeed() error {
	file, err := site.createAbsolutePath("feeds/all.atom.xml")
	if err != nil {
		return err
	}

	limit := len(site.Articles)
	if limit > 10 { // TODO: Move feed limit of news to the config
		limit = 10
	}

	if err := templates["atom"].Execute(file, site.NewAtomFeedContext(limit)); err != nil {
		return err
	}

	return nil
}

func (site Site) writeRSSFeed() error {
	// TODO (agonzalezro): to be implemented if somebody needs it
	return nil
}
