// TODO: THIS IS THE SH*$#ST PIECE OF CODE THAT I MADE IN AGES: take some time
// to refactor ASAP
package generator

import (
	"log"
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
	Article, Page     ParsedFile
	Tag               string
	Category          string
	PaginatedArticles []*ParsedFile
}

func NewSite(db DB, config Config, outputPath string) *Site {
	updated := time.Now().Format(time.RFC3339)
	return &Site{db: db, Config: config, outputPath: outputPath, Updated: updated}
}

func (site Site) Tags() (tags []string) {
	var storedTags string
	seenList := make(map[string]bool)

	query := "SELECT tags FROM files WHERE is_page = 0 AND status != 'draft'"
	rows, err := site.db.connection.Query(query)
	if err != nil {
		log.Panicf("Error querying for tags: %v", err)
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
	return tags
}

func (site Site) Categories() (categories []string) {
	var category string

	query := `SELECT DISTINCT category FROM files WHERE is_page = 0 AND status != 'draft' AND category != ""`
	rows, err := site.db.connection.Query(query)
	if err != nil {
		log.Panic("Error query for categories: %v", err)
	}

	for rows.Next() {
		rows.Scan(&category)
		categories = append(categories, category)
	}

	sort.Strings(categories)
	return categories
}

func (site Site) Articles() (articles []*ParsedFile) {
	return site.QueryArticles("")
}

func (site Site) ArticlesByTag(tag string) (articles []*ParsedFile) {
	// Concatenation from the hell
	return site.QueryArticles("tags LIKE \"%,\"||?||\",%\"", tag)
}

func (site Site) ArticlesByCategory(category string) (articles []*ParsedFile) {
	return site.QueryArticles("category = ?", category)
}

func (site Site) Pages() (pages []*ParsedFile) {
	return site.QueryPages("")
}
