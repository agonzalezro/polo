package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
)

// Dump all the site content to the disk
func (site Site) Write() {
	site.writeIndex()
	site.writeFeeds()
	site.writeArticles()
	site.writePages()
	site.writeCategories()
	site.writeTags()
}

func (site Site) writeIndex() {
	indexFile := fmt.Sprintf("%s/index.html", site.outputPath)

	template := template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))

	file, err := os.Create(indexFile)
	if err != nil {
		log.Panicf("Error creating index file: %v", err)
	}

	if err := template.ExecuteTemplate(file, "base", site); err != nil {
		log.Panicf("Error rendering the template for the index: %v", err)
	}
}

func (site Site) writeCategories() {
	// First of all create the tags/ folder if it doesn't exist
	pagesPath := fmt.Sprintf("%s/category", site.outputPath)
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		os.Mkdir(pagesPath, 0777)
	}

	for _, category := range site.Categories() {
		categoryFile := fmt.Sprintf("%s/category/%s.html", site.outputPath, category)

		template := template.Must(template.ParseFiles("templates/category.html", "templates/base.html"))

		file, err := os.Create(categoryFile)
		if err != nil {
			log.Panicf("Error creating the category '%s' file: %v", category, err)
		}

		site.Category = category
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			log.Panicf("Error rendering the template for the category '%s': %v", category, err)
		}
	}
}

func (site Site) writeTags() {
	// First of all create the tags/ folder if it doesn't exist
	pagesPath := fmt.Sprintf("%s/tag", site.outputPath)
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		os.Mkdir(pagesPath, 0777)
	}

	for _, tag := range site.Tags() {
		tagFile := fmt.Sprintf("%s/tag/%s.html", site.outputPath, tag)

		template := template.Must(template.ParseFiles("templates/tag.html", "templates/base.html"))

		file, err := os.Create(tagFile)
		if err != nil {
			log.Panicf("Error creating the tag '%s' file: %v", tag, err)
		}

		site.Tag = tag
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			log.Panicf("Error rendering the template for the tag '%s': %v", tag, err)
		}
	}
}

func (site Site) writeFeeds() {
	feedsPath := fmt.Sprintf("%s/feeds", site.outputPath)
	if _, err := os.Stat(feedsPath); os.IsNotExist(err) {
		os.Mkdir(feedsPath, 0777)
	}

	// TODO (agonzalezro): write the atom and RSS feeds
	site.writeAtomFeed(feedsPath)
	site.writeRSSFeed(feedsPath)
}

func (site Site) writeAtomFeed(feedsPath string) {
	path := fmt.Sprintf("%s/all.atom.xml", feedsPath)
	template := template.Must(template.ParseFiles("templates/atom.xml"))

	file, err := os.Create(path)
	if err != nil {
		log.Panicf("Error creating the atom file: %v", err)
	}

	articles := site.Articles()
	limit := len(articles)
	if limit > 10 {
		limit = 10
	}
	site.PaginatedArticles = articles[:limit] // TODO: do it inside the function
	if err := template.Execute(file, site); err != nil {
		log.Panicf("Error rendering the template for the atom file: %v", err)
	}
}

func (site Site) writeArticles() {
	for _, article := range site.Articles() {
		filePath := fmt.Sprintf("%s/%s.html", site.outputPath, article.Slug)

		template := template.Must(template.ParseFiles("templates/article.html", "templates/base.html"))

		file, err := os.Create(filePath)
		if err != nil {
			log.Panicf("Error creating the file: %s\n%v", filePath, err)
		}
		site.Article = *article
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			log.Panicf("Error rendering template for the file: %s\n%v", filePath, err)
		}
	}

}

func (site Site) writePages() {
	// First of all create the pages/ folder if it doesn't exist
	pagesPath := fmt.Sprintf("%s/pages", site.outputPath)
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		os.Mkdir(pagesPath, 0777)
	}

	for _, page := range site.Pages() {
		filePath := fmt.Sprintf("%s/%s.html", pagesPath, page.Slug)

		template := template.Must(template.ParseFiles("templates/page.html", "templates/base.html"))

		file, err := os.Create(filePath)
		if err != nil {
			log.Panicf("Error creating the file: %s\n%v", filePath, err)
		}
		site.Page = *page
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			log.Panicf("Error rendering the template for the file: %s\n%v", filePath, err)
		}
	}
}

func (site Site) writeRSSFeed(feedsPath string) {
	// TODO (agonzalezro): to be implemented if somebody needs it
	return
}
