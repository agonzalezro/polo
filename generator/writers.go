package generator

import (
	"errors"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
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

func (site Site) getAbsolutePath(elem ...string) string {
	// TODO: I am not pretty sure that this is the best way to do this
	s := make([]string, 1, 1)
	s[0] = site.outputPath
	elem = append(s, elem...)
	return path.Join(elem...)
}

// Dump all the site content to the disk
func (site Site) Write() (err error) {
	var i int
	log := func(noun string, qty int) {
		var pluralize string
		if qty > 1 {
			pluralize = "s"
		}
		log.Printf("%4d %s%s created", qty, strings.Title(noun), pluralize)
	}

	loadTemplates()

	i, err = site.writeIndexes()
	if err != nil {
		return err
	}
	log("index page", i)

	i, err = site.writeFeeds()
	if err != nil {
		return err
	}
	log("feed", i)

	i, err = site.writeArticles()
	if err != nil {
		return err
	}
	log("article", i)

	i, err = site.writePages()
	if err != nil {
		return err
	}
	log("page", i)

	if site.Config.ShowArchive {
		if err = site.writeArchive(); err != nil {
			return err
		}
		log("archive", 1)
	}
	if site.Config.ShowCategories {
		i, err = site.writeCategories()
		if err != nil {
			return err
		}
		log("category page", i)
	}
	if site.Config.ShowTags {
		i, err = site.writeTags()
		if err != nil {
			return err
		}
		log("tag page", i)
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
		indexFile := site.getAbsolutePath(fmt.Sprintf("index%d.html", site.PageNumber))
		if site.PageNumber == 1 {
			indexFile = site.getAbsolutePath("index.html")
		}

		file, err := os.Create(indexFile)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error creating index file for page '%d': %v", site.PageNumber, err))
			return site.PageNumber, ErrorCreate(err)
		}

		if err := templates["index"].ExecuteTemplate(file, "base", site); err != nil {
			err = errors.New(fmt.Sprintf("Error rendering the index file for page '%d': %v", site.PageNumber, err))
			return site.PageNumber, ErrorExecuteTemplate(err)
		}
	}

	return site.PageNumber, nil
}

func (site Site) writeParsedFiles(rootPath string, files []*ParsedFile) (int, error) {
	if rootPath != "" {
		if _, err := os.Stat(rootPath); os.IsNotExist(err) {
			os.Mkdir(rootPath, 0777)
		}
	}

	var (
		i          int
		parsedFile *ParsedFile
	)

	for i, parsedFile = range files {
		filePath := fmt.Sprintf("%s/%s.html", rootPath, parsedFile.Slug)

		var template *template.Template
		if files[0].isPage {
			template = templates["page"]
		} else {
			template = templates["article"]
		}

		file, err := os.Create(filePath)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error creating the file: %s\n%v", filePath, err))
			return i + 1, ErrorCreate(err)
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
	return site.writeParsedFiles(site.outputPath, site.Articles())
}

func (site Site) writePages() (int, error) {
	pagesPath := site.getAbsolutePath("pages")
	return site.writeParsedFiles(pagesPath, site.Pages())
}

func (site Site) writeArchive() error {
	archivesPath := site.getAbsolutePath("archives.html")
	file, err := os.Create(archivesPath)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error creating archive file: %v", err))
		return ErrorCreate(err)
	}

	if err := templates["archives"].ExecuteTemplate(file, "base", site); err != nil {
		err = errors.New(fmt.Sprintf("Error rendering the template for the archives: %v", err))
		return ErrorExecuteTemplate(err)
	}

	return nil
}

func (site Site) writeCategories() (int, error) {
	// First of all create the tags/ folder if it doesn't exist
	categoriesPath := site.getAbsolutePath("category")
	if _, err := os.Stat(categoriesPath); os.IsNotExist(err) {
		os.Mkdir(categoriesPath, 0777)
	}

	var (
		i        int
		category string
	)

	for i, category = range site.Categories() {
		categoryFile := fmt.Sprintf("%s/%s.html", categoriesPath, category)

		file, err := os.Create(categoryFile)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error creating the category '%s' file: %v", category, err))
			return i + 1, ErrorCreate(err)
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
	// First of all create the tags/ folder if it doesn't exist
	tagsPath := site.getAbsolutePath("tag")
	if _, err := os.Stat(tagsPath); os.IsNotExist(err) {
		os.Mkdir(tagsPath, 0777)
	}

	var (
		i   int
		tag string
	)

	for i, tag = range site.Tags() {
		tagFile := fmt.Sprintf("%s/%s.html", tagsPath, tag)

		file, err := os.Create(tagFile)
		if err != nil {
			err = errors.New(fmt.Sprintf("Error creating the tag '%s' file: %v", tag, err))
			return i + 1, ErrorCreate(err)
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

	feedsPath := site.getAbsolutePath("feeds")
	if _, err := os.Stat(feedsPath); os.IsNotExist(err) {
		os.Mkdir(feedsPath, 0777)
	}

	if err := site.writeAtomFeed(feedsPath); err != nil {
		return i + 1, err
	}

	if err := site.writeRSSFeed(feedsPath); err != nil {
		return i + 1, err
	}
	i += 0 // Not implemented yet

	return i + 1, nil
}

func (site Site) writeAtomFeed(feedsPath string) error {
	path := fmt.Sprintf("%s/all.atom.xml", feedsPath)

	file, err := os.Create(path)
	if err != nil {
		err = errors.New(fmt.Sprintf("Error creating the atom file: %v", err))
		return ErrorCreate(err)
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

func (site Site) writeRSSFeed(feedsPath string) error {
	// TODO (agonzalezro): to be implemented if somebody needs it
	return nil
}
