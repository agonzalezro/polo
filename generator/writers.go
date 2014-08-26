package generator

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"

	// TODO: Perhaps it worth moving the template rendering to the template
	// package
	assets "github.com/agonzalezro/polo/templates"
)

type ErrorCreate error
type ErrorExecuteTemplate error

var templates map[string]*template.Template

// parsedFiles is a wrapper similar to template.ParseFiles that is going to
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

	getAllTemplateInheritance := func(template string) []string {
		alwaysIncludeTemplates := []string{
			"templates/base/base.html",
			"templates/base/header.html",
			"templates/base/footer.html",
			"templates/base/analytics.html"}
		return append(alwaysIncludeTemplates, template)
	}

	toRender["archives"] = getAllTemplateInheritance("templates/archives.html")
	toRender["article"] = getAllTemplateInheritance("templates/article.html")
	toRender["category"] = getAllTemplateInheritance("templates/category.html")
	toRender["index"] = getAllTemplateInheritance("templates/index.html")
	toRender["page"] = getAllTemplateInheritance("templates/page.html")
	toRender["tag"] = getAllTemplateInheritance("templates/tag.html")

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
func logCreation(noun string, quantity int) {
	var pluralize string
	if quantity > 1 {
		pluralize = "s"
	}
	log.Printf("%4d %s%s created", quantity, strings.Title(noun), pluralize)
}

// Dump all the site content to the disk
func (site Site) Write() (err error) {
	loadTemplates()

	logs := make(map[string]int)

	chErrors := make(chan error, 7)
	subroutines := 4 // 4 are mandatory, the others depend of settings

	go func() {
		logs["index page"], err = site.writeIndexes()
		chErrors <- err
	}()

	go func() {
		logs["feed"], err = site.writeFeeds()
		chErrors <- err
	}()

	go func() {
		logs["article"], err = site.writeArticles()
		chErrors <- err
	}()

	go func() {
		logs["page"], err = site.writePages()
		chErrors <- err
	}()

	if site.Config.ShowArchive {
		subroutines++
		go func() {
			err := site.writeArchive()
			if err != nil {
				logs["archive"] = 1 // Not completely true, but... meh!
			}
			chErrors <- err
		}()
	}
	if site.Config.ShowCategories {
		subroutines++
		go func() {
			logs["category page"], err = site.writeCategories()
			chErrors <- err
		}()
	}
	if site.Config.ShowTags {
		subroutines++
		go func() {
			logs["tag page"], err = site.writeTags()
			chErrors <- err
		}()
	}

	for i := 0; i < subroutines; i++ {
		err := <-chErrors
		if err != nil {
			return err
		}
	}

	for noun, quantity := range logs {
		logCreation(noun, quantity)
	}

	return nil
}

func (site Site) getNumberOfPages() int {
	if site.NumberOfPages == 0 {
		site.NumberOfPages = len(site.Articles()) / site.Config.PaginationSize
	}
	// If it continue being 0, we don't have pages
	if site.NumberOfPages == 0 {
		site.NumberOfPages = -1
	}
	return site.NumberOfPages
}

func (site Site) writeIndexes() (int, error) {
	site.NumberOfPages = site.getNumberOfPages()

	for site.PageNumber = 1; site.PageNumber <= site.NumberOfPages; site.PageNumber++ {
		indexFile := "index.html"
		if site.PageNumber > 1 {
			indexFile = fmt.Sprintf("index%d.html", site.PageNumber)
		}

		file, err := site.createAbsolutePath(indexFile)
		if err != nil {
			return site.PageNumber, err
		}

		if err := templates["index"].ExecuteTemplate(file, "base", site); err != nil {
			err = errors.New(fmt.Sprintf("Error rendering the index file for page '%d': %v", site.PageNumber, err))
			return site.PageNumber, ErrorExecuteTemplate(err)
		}
	}

	return site.PageNumber, nil
}

func (site Site) writeParsedFiles(pathAppender string, files []*ParsedFile) (int, error) {
	var (
		i          int
		parsedFile *ParsedFile
	)

	for i, parsedFile = range files {
		var template *template.Template
		if files[0].isPage {
			template = templates["page"]
		} else {
			template = templates["article"]
		}

		filePath := fmt.Sprintf("%s/%s.html", pathAppender, parsedFile.Slug)
		file, err := site.createAbsolutePath(filePath)
		if err != nil {
			return i + 1, err
		}

		if files[0].isPage {
			site.Page = *parsedFile
		} else {
			site.Article = *parsedFile
		}
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			err = errors.New(fmt.Sprintf("Error rendering the template for the file: %s\n%v", filePath, err))
			return i + 1, ErrorExecuteTemplate(err)
		}
	}

	return i + 1, nil
}

func (site Site) writeArticles() (int, error) {
	return site.writeParsedFiles("/", site.Articles())
}

func (site Site) writePages() (int, error) {
	return site.writeParsedFiles("pages", site.Pages())
}

func (site Site) writeArchive() error {
	file, err := site.createAbsolutePath("archives.html")
	if err != nil {
		return err
	}

	if err := templates["archives"].ExecuteTemplate(file, "base", site); err != nil {
		err = errors.New(fmt.Sprintf("Error rendering the template for the archives: %v", err))
		return ErrorExecuteTemplate(err)
	}

	return nil
}

func (site Site) writeCategories() (int, error) {
	var (
		i        int
		category string
	)

	for i, category = range site.Categories() {
		categoryFile := fmt.Sprintf("category/%s.html", category)
		file, err := site.createAbsolutePath(categoryFile)
		if err != nil {
			return i + 1, err
		}

		site.Category = category
		if err := templates["category"].ExecuteTemplate(file, "base", site); err != nil {
			err = errors.New(fmt.Sprintf("Error rendering the template for the category '%s': %v", category, err))
			return i + 1, ErrorExecuteTemplate(err)
		}
	}

	return i + 1, nil
}

func (site Site) writeTags() (int, error) {
	var (
		i   int
		tag string
	)

	for i, tag = range site.Tags() {
		tagFile := fmt.Sprintf("tag/%s.html", tag)
		file, err := site.createAbsolutePath(tagFile)
		if err != nil {
			return i + 1, err
		}

		site.Tag = tag
		if err := templates["tag"].ExecuteTemplate(file, "base", site); err != nil {
			err = errors.New(fmt.Sprintf("Error rendering the template for the tag '%s': %v", tag, err))
			return i + 1, ErrorExecuteTemplate(err)
		}
	}

	return i + 1, nil
}

func (site Site) writeFeeds() (int, error) {
	var i int

	if err := site.writeAtomFeed(); err != nil {
		return i + 1, err
	}

	if err := site.writeRSSFeed(); err != nil {
		return i + 1, err
	}
	i += 0 // Not implemented yet

	return i + 1, nil
}

func (site Site) writeAtomFeed() error {
	file, err := site.createAbsolutePath("feeds/all.atom.xml")
	if err != nil {
		return err
	}

	articles := site.Articles()
	limit := len(articles)
	if limit > 10 {
		limit = 10
	}
	site.FeedArticles = articles[:limit] // TODO: do it inside the function
	if err := templates["atom"].Execute(file, site); err != nil {
		err = errors.New(fmt.Sprintf("Error rendering the template for the atom file: %v", err))
		return ErrorExecuteTemplate(err)
	}

	return nil
}

func (site Site) writeRSSFeed() error {
	// TODO (agonzalezro): to be implemented if somebody needs it
	return nil
}
