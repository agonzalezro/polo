package generator

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type Site struct {
	db DB

	Config     Config
	outputPath string

	Updated string

	// Temporal stuff just for that page being rendered
	Article, Page ParsedFile
	Tag           string
	Category      string
	FeedArticles  []*ParsedFile
	PageNumber    int
	NumberOfPages int

	Cache map[string]interface{}
}

func NewSite(db DB, config Config, outputPath string) *Site {
	updated := time.Now().Format(time.RFC3339)
	site := &Site{db: db, Config: config, outputPath: outputPath, Updated: updated}
	site.Cache = make(map[string]interface{})
	return site
}

func (site Site) tags() (tags []string, err error) {
	if tags, ok := site.Cache["tags"]; ok {
		return tags.([]string), nil
	}

	var storedTags string
	seenList := make(map[string]bool)

	query := "SELECT tags FROM files WHERE is_page = 0 AND status != 'draft'"
	rows, err := site.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rows.Scan(&storedTags)
		for _, tag := range strings.Split(storedTags, ",") {
			tag = strings.TrimSpace(tag)
			if _, seen := seenList[tag]; !seen && tag != "" {
				seenList[tag] = true
				tags = append(tags, tag)
			}
		}
	}

	sort.Strings(tags)
	site.Cache["tags"] = tags
	return tags, nil
}

// ArrayOfPages is a dirty hack because we can not iterate over an integer on
// the template
func (site Site) ArrayOfPages() (pages []int) {
	nop, _ := site.getNumberOfPages() // If it errors here we would see the error before
	for i := 1; i < nop+1; i++ {
		pages = append(pages, i)
	}
	return pages
}

func (site Site) GetPreviousPageSlug(page int) (slug string) {
	switch page {
	case 1:
		slug = "#"
	case 2:
		slug = "/index.html"
	default:
		slug = fmt.Sprintf("/index%d.html", page-1)

	}
	return slug
}

func (site Site) GetNextPageSlug(page int) (slug string) {
	if page == site.NumberOfPages {
		return "#"
	}

	return fmt.Sprintf("/index%d.html", page+1)
}

func (site Site) categories() (categories []string, err error) {
	if categories, ok := site.Cache["categories"]; ok {
		return categories.([]string), nil
	}

	var category string

	query := `SELECT DISTINCT category FROM files WHERE is_page = 0 AND status != 'draft' AND category != ""`
	rows, err := site.db.Query(query)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		rows.Scan(&category)
		categories = append(categories, category)
	}

	sort.Strings(categories)

	site.Cache["categories"] = categories
	return categories, nil
}

// Tags is mean to be called from the templates, we don't want error handling
// there, so this function is just wrapping the call to .tags()
func (site Site) Tags() (tags []string) {
	tags, _ = site.tags()
	return
}

// Categories is just a wrapping to don't raise any error to the template
// rendering
func (site Site) Categories() (categories []string) {
	categories, _ = site.categories()
	return
}

func (site Site) articles() (articles []*ParsedFile, err error) {
	if articles, ok := site.Cache["articles"]; ok {
		return articles.([]*ParsedFile), nil
	}
	articles, err = site.QueryArticles("", -1)
	if err != nil {
		return nil, err
	}
	site.Cache["articles"] = articles
	return articles, nil
}

func (site Site) Articles() (articles []*ParsedFile) {
	articles, _ = site.articles()
	return
}

func (site Site) ArticlesByPage(page int) (articles []*ParsedFile) {
	articles, _ = site.QueryArticles("", page)
	return
}

func (site Site) ArticlesByTag(tag string) (articles []*ParsedFile) {
	// Concatenation from the hell
	articles, _ = site.QueryArticles("tags LIKE \"%,\"||?||\",%\"", -1, tag)
	return
}

func (site Site) ArticlesByCategory(category string) (articles []*ParsedFile) {
	articles, _ = site.QueryArticles("category = ?", -1, category)
	return
}

func (site Site) pages() (pages []*ParsedFile, err error) {
	pages, err = site.QueryPages()
	if err != nil {
		return nil, err
	}
	return pages, nil
}

func (site Site) Pages() (pages []*ParsedFile) {
	pages, _ = site.pages()
	return
}
