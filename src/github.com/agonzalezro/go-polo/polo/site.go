package polo

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
)

type Site struct {
	db DB

	Config     Config
	outputPath string

	Article, Page ParsedFile
}

func NewSite(db DB, config Config, outputPath string) *Site {
	return &Site{db: db, Config: config, outputPath: outputPath}
}

// Dump all the site content to the disk
func (site Site) Write() {
	site.writeIndex()
	site.writeArticles()
	site.writePages()
}

func (site Site) writeIndex() {
	indexFile := fmt.Sprintf("%s/index.html", site.outputPath)

	template := template.Must(template.ParseFiles("templates/index.html", "templates/base.html"))

	file, err := os.Create(indexFile)
	if err != nil {
		log.Panic(err)
	}

	if err := template.ExecuteTemplate(file, "base", site); err != nil {
		log.Panic(err)
	}
}

func (site Site) Tags() (tags []string) {
	// Not optimal, but it does the job
	var (
		seenList   map[string]bool
		storedTags string
	)

	query := "SELECT tags FROM files WHERE is_page = 0"
	rows, err := site.db.connection.Query(query)
	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		rows.Scan(&storedTags)
		for _, tag := range strings.Split(storedTags, ",") {
			if _, seen := seenList[tag]; !seen && tag != "" {
				tags = append(tags, strings.TrimSpace(tag))
			}
		}
	}

	return tags
}

func (site Site) Articles() (articles []*ParsedFile) {
	query := `
		SELECT title, slug, content, tags, date
		FROM files
		WHERE is_page = 0
		`
	rows, err := site.db.connection.Query(query)
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		article := ParsedFile{isPage: false}
		rows.Scan(&article.Title, &article.Slug, &article.Content, &article.tags, &article.Date)
		articles = append(articles, &article)
	}
	return articles
}

func (site Site) writeArticles() {
	for _, article := range site.Articles() {
		filePath := fmt.Sprintf("%s/%s.html", site.outputPath, article.Slug)

		template := template.Must(template.ParseFiles("templates/article.html", "templates/base.html"))

		file, err := os.Create(filePath)
		if err != nil {
			log.Panic(err)
		}
		site.Article = *article
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			log.Panic(err)
		}
	}

}

func (site Site) Pages() (pages []*ParsedFile) {
	query := `
		SELECT title, slug, content, tags, date
		FROM files
		WHERE is_page =1
		`
	rows, err := site.db.connection.Query(query)
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		page := ParsedFile{isPage: true}
		rows.Scan(&page.Title, &page.Slug, &page.Content, &page.tags, &page.Date)
		pages = append(pages, &page)
	}
	return pages
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
			log.Panic(err)
		}
		site.Page = *page
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			log.Panic(err)
		}
	}
}
