package generator

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"os"

	"github.com/agonzalezro/polo/templates"
)

var (
	archiveTemplate,
	articleTemplate,
	atomTemplate,
	baseTemplate,
	categoryTemplate,
	indexTemplate,
	pageTemplate,
	tagTemplate *template.Template
)

// parsedFiles is a wrapper similar to template.ParseFiles that is going to
// load the templates from the disk, and if they can not be found from the
// go-bindata file.
func parseFiles(filenames ...string) (*template.Template, error) {
	tpl := template.New(filenames[0])
	for _, filename := range filenames {
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			b, err = templates.Asset(filename)
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
	archiveTemplate = template.Must(parseFiles("templates/archives.html", "templates/base.html"))
	articleTemplate = template.Must(parseFiles("templates/article.html", "templates/base.html"))
	atomTemplate = template.Must(parseFiles("templates/atom.xml"))
	categoryTemplate = template.Must(parseFiles("templates/category.html", "templates/base.html"))
	indexTemplate = template.Must(parseFiles("templates/index.html", "templates/base.html"))
	pageTemplate = template.Must(parseFiles("templates/page.html", "templates/base.html"))
	tagTemplate = template.Must(parseFiles("templates/tag.html", "templates/base.html"))
}

// Dump all the site content to the disk
func (site Site) Write() {
	loadTemplates()

	site.writeIndexes()
	site.writeFeeds()
	site.writeArticles()
	site.writePages()
	if site.Config.ShowArchive {
		site.writeArchive()
	}
	if site.Config.ShowCategories {
		site.writeCategories()
	}
	if site.Config.ShowTags {
		site.writeTags()
	}
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

func (site Site) writeIndexes() {
	site.NumberOfPages = site.getNumberOfPages()

	for site.PageNumber = 1; site.PageNumber <= site.NumberOfPages; site.PageNumber++ {
		indexFile := fmt.Sprintf("%s/index%d.html", site.outputPath, site.PageNumber)
		if site.PageNumber == 1 {
			indexFile = fmt.Sprintf("%s/index.html", site.outputPath)
		}

		file, err := os.Create(indexFile)
		if err != nil {
			log.Panicf("Error creating index file for page '%d': %v", site.PageNumber, err)
		}

		if err := indexTemplate.ExecuteTemplate(file, "base", site); err != nil {
			log.Panicf("Error rendering the index file for page '%d': %v", site.PageNumber, err)
		}
	}
}

func (site Site) writeParsedFiles(rootPath string, files []*ParsedFile) {
	if rootPath != "" {
		if _, err := os.Stat(rootPath); os.IsNotExist(err) {
			os.Mkdir(rootPath, 0777)
		}
	}

	for _, parsedFile := range files {
		filePath := fmt.Sprintf("%s/%s.html", rootPath, parsedFile.Slug)

		var template *template.Template
		if files[0].isPage {
			template = pageTemplate
		} else {
			template = articleTemplate
		}

		file, err := os.Create(filePath)
		if err != nil {
			log.Panicf("Error creating the file: %s\n%v", filePath, err)
		}

		if files[0].isPage {
			site.Page = *parsedFile
		} else {
			site.Article = *parsedFile
		}
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			log.Panicf("Error rendering the template for the file: %s\n%v", filePath, err)
		}
	}
}

func (site Site) writeArticles() {
	site.writeParsedFiles(site.outputPath, site.Articles())
}

func (site Site) writePages() {
	pagesPath := fmt.Sprintf("%s/pages", site.outputPath)
	site.writeParsedFiles(pagesPath, site.Pages())
}

func (site Site) writeArchive() {
	archivesPath := fmt.Sprintf("%s/archives.html", site.outputPath)
	file, err := os.Create(archivesPath)
	if err != nil {
		log.Panicf("Error creating archive file: %v", err)
	}

	if err := archiveTemplate.ExecuteTemplate(file, "base", site); err != nil {
		log.Panicf("Error rendering the template for the archives: %v", err)
	}
}

func (site Site) writeCategories() {
	// First of all create the tags/ folder if it doesn't exist
	categoriesPath := fmt.Sprintf("%s/category", site.outputPath)
	if _, err := os.Stat(categoriesPath); os.IsNotExist(err) {
		os.Mkdir(categoriesPath, 0777)
	}

	for _, category := range site.Categories() {
		categoryFile := fmt.Sprintf("%s/%s.html", categoriesPath, category)

		file, err := os.Create(categoryFile)
		if err != nil {
			log.Panicf("Error creating the category '%s' file: %v", category, err)
		}

		site.Category = category
		if err := categoryTemplate.ExecuteTemplate(file, "base", site); err != nil {
			log.Panicf("Error rendering the template for the category '%s': %v", category, err)
		}
	}
}

func (site Site) writeTags() {
	// First of all create the tags/ folder if it doesn't exist
	tagsPath := fmt.Sprintf("%s/tag", site.outputPath)
	if _, err := os.Stat(tagsPath); os.IsNotExist(err) {
		os.Mkdir(tagsPath, 0777)
	}

	for _, tag := range site.Tags() {
		tagFile := fmt.Sprintf("%s/%s.html", tagsPath, tag)

		file, err := os.Create(tagFile)
		if err != nil {
			log.Panicf("Error creating the tag '%s' file: %v", tag, err)
		}

		site.Tag = tag
		if err := tagTemplate.ExecuteTemplate(file, "base", site); err != nil {
			log.Panicf("Error rendering the template for the tag '%s': %v", tag, err)
		}
	}
}

func (site Site) writeFeeds() {
	feedsPath := fmt.Sprintf("%s/feeds", site.outputPath)
	if _, err := os.Stat(feedsPath); os.IsNotExist(err) {
		os.Mkdir(feedsPath, 0777)
	}

	site.writeAtomFeed(feedsPath)
	site.writeRSSFeed(feedsPath)
}

func (site Site) writeAtomFeed(feedsPath string) {
	path := fmt.Sprintf("%s/all.atom.xml", feedsPath)

	file, err := os.Create(path)
	if err != nil {
		log.Panicf("Error creating the atom file: %v", err)
	}

	articles := site.Articles()
	limit := len(articles)
	if limit > 10 {
		limit = 10
	}
	site.FeedArticles = articles[:limit] // TODO: do it inside the function
	if err := atomTemplate.Execute(file, site); err != nil {
		log.Panicf("Error rendering the template for the atom file: %v", err)
	}
}

func (site Site) writeRSSFeed(feedsPath string) {
	// TODO (agonzalezro): to be implemented if somebody needs it
	return
}
