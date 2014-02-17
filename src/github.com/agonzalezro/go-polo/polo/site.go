package polo

import (
	"fmt"
	"html/template"
	"log"
	"os"
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

// TODO (agonzalezro): this needs to be merged with Pages() and used in writeXXXs
func (site Site) Articles() (articles []*ParsedFile) {
	query := `
		SELECT title, slug, content, tags, date, is_page
		FROM files
		WHERE is_page = 0
		`
	rows, err := site.db.connection.Query(query)
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		var title, slug, content, tags, date string
		var isPage int
		rows.Scan(&title, &slug, &content, &tags, &date, isPage)
		article := ParsedFile{Title: title, Slug: slug, Content: content, tags: tags, Date: date}
		articles = append(articles, &article)
	}
	return articles
}

// TODO (agonzalezro): use this function from write*
func (site Site) Pages() (pages []*ParsedFile) {
	query := `
		SELECT title, slug, content, tags, date, is_page
		FROM files
		WHERE is_page =1
		`
	rows, err := site.db.connection.Query(query)
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		var title, slug, content, tags, date string
		var isPage int
		rows.Scan(&title, &slug, &content, &tags, &date, isPage)
		page := ParsedFile{Title: title, Slug: slug, Content: content, tags: tags, Date: date}
		pages = append(pages, &page)
	}
	return pages
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
	template.ExecuteTemplate(file, "base", site)
}

// Write articles in different files
func (site Site) writeArticles() {
	query := `
    SELECT title, slug, content, tags, date, is_page
    FROM files
	WHERE is_page = 0
    `
	rows, err := site.db.connection.Query(query)
	if err != nil {
		log.Panic(err)
	}
	for rows.Next() {
		var title, slug, content, tags, date string
		var isPage int
		rows.Scan(&title, &slug, &content, &tags, &date, isPage)
		article := ParsedFile{Title: title, Slug: slug, Content: content, tags: tags, Date: date}
		filePath := fmt.Sprintf("%s/%s.html", site.outputPath, article.Slug)

		template := template.Must(template.ParseFiles("templates/article.html", "templates/base.html"))

		file, err := os.Create(filePath)
		if err != nil {
			log.Panic(err)
		}
		site.Article = article
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			log.Panic(err)
		}
	}

}

func (site Site) writePages() {
	query := `
	SELECT title, slug, content, tags, date, is_page
	FROM files
	WHERE is_page = 1
	`

	// First of all create the pages/ folder if it doesn't exist
	pagesPath := fmt.Sprintf("%s/pages", site.outputPath)
	if _, err := os.Stat(pagesPath); os.IsNotExist(err) {
		os.Mkdir(pagesPath, 0777)
	}

	rows, err := site.db.connection.Query(query)
	if err != nil {
		log.Panic(err)
	}

	for rows.Next() {
		var title, slug, content, tags, date string
		var isPage int
		rows.Scan(&title, &slug, &content, &tags, &date, &isPage)
		page := ParsedFile{Title: title, Slug: slug, Content: content, tags: tags, Date: date}

		filePath := fmt.Sprintf("%s/%s.html", pagesPath, page.Slug)

		template := template.Must(template.ParseFiles("templates/page.html", "templates/base.html"))

		file, err := os.Create(filePath)
		if err != nil {
			log.Panic(err)
		}
		site.Page = page
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			log.Panic(err)
		}
	}
}
