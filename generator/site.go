package generator

import (
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"
)

type Site struct {
	db DB

	Config     Config
	outputPath string

	Article, Page ParsedFile

	PaginatedArticles []*ParsedFile

	Updated string
}

func NewSite(db DB, config Config, outputPath string) *Site {
	updated := time.Now().Format(time.RFC3339)
	return &Site{db: db, Config: config, outputPath: outputPath, Updated: updated}
}

// Dump all the site content to the disk
func (site Site) Write() {
	site.writeIndex()
	site.writeFeeds()
	site.writeArticles()
	site.writePages()
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

func (site Site) writeRSSFeed(feedsPath string) {
	// TODO (agonzalezro): to be implemented if somebody needs it
	return
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
		log.Panicf("Error querying for tags: %v", err)
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
		SELECT title, slug, content, tags, date, summary
		FROM files
		WHERE is_page = 0
		AND status != 'draft'
		ORDER BY datetime(date) DESC
		`
	rows, err := site.db.connection.Query(query)
	if err != nil {
		log.Panicf("Error querying for articles: %v", err)
	}
	for rows.Next() {
		article := ParsedFile{isPage: false}
		rows.Scan(&article.Title, &article.Slug, &article.Content, &article.tags, &article.Date, &article.summary)
		articles = append(articles, &article)
	}
	return articles
}

// TODO (agonzalezro): possibly duplicated, but the query params are different :(
func (site Site) ArticlesByTag(tag string) (articles []*ParsedFile) {
	// You can hit me for this ugly hack
	query := `
        SELECT title, slug, content, tags, date, summary
        FROM files
        WHERE is_page = 0
        AND status != 'draft'
        AND tags LIKE "%"||?||",%"
        OR tags LIKE "%, "||?||"%"
        OR tags LIKE "%,"||?||"%"
        OR tags = ?
        ORDER BY datetime(date) DESC
        `
	rows, err := site.db.connection.Query(query, tag, tag, tag, tag)
	if err != nil {
		log.Panicf("Error querying for articles with tag '%s': %v", tag, err)
	}
	for rows.Next() {
		article := ParsedFile{isPage: false}
		rows.Scan(&article.Title, &article.Slug, &article.Content, &article.tags, &article.Date, &article.summary)
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
			log.Panicf("Error creating the file: %s\n%v", filePath, err)
		}
		site.Article = *article
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			log.Panicf("Error rendering template for the file: %s\n%v", filePath, err)
		}
	}

}

func (site Site) Pages() (pages []*ParsedFile) {
	query := `
		SELECT title, slug, content, tags, date, summary
		FROM files
		WHERE is_page = 1
		AND status != 'draft'
		ORDER BY datetime(date) DESC
		`
	rows, err := site.db.connection.Query(query)
	if err != nil {
		log.Panicf("Error querying for the pages: %v", err)
	}
	for rows.Next() {
		page := ParsedFile{isPage: true}
		rows.Scan(&page.Title, &page.Slug, &page.Content, &page.tags, &page.Date, &page.summary)
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
			log.Panicf("Error creating the file: %s\n%v", filePath, err)
		}
		site.Page = *page
		if err := template.ExecuteTemplate(file, "base", site); err != nil {
			log.Panicf("Error rendering the template for the file: %s\n%v", filePath, err)
		}
	}
}
