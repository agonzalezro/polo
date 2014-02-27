package generator

import (
	"fmt"
	"log"
)

// getQueryInterface is going to cast the array of string to an array of
// interfaces adding the default parameter isPage in the position 0
func getQueryInterface(isPage bool, args ...string) []interface{} {
	// This funny casting is needed to call the sql .Query()
	sqlArgs := make([]interface{}, len(args)+1)
	sqlArgs[0] = interface{}(0)
	if isPage {
		sqlArgs[0] = interface{}(1)
	}
	for i, v := range args {
		sqlArgs[i+1] = interface{}(v)
	}
	return sqlArgs
}

func (site Site) Query(isPage bool, where string, args ...string) (files []*ParsedFile) {
	// In case that a where clausule needs to be added, add the AND at the beginning
	if where != "" {
		where = fmt.Sprintf("AND %s", where)
	}
	query := fmt.Sprintf(`
        SELECT title, slug, content, category, tags, date, summary
        FROM files
        WHERE is_page = ?
        AND status != 'draft'
        %s
        ORDER BY datetime(date) DESC
    `, where)

	sqlArgs := getQueryInterface(isPage, args...)
	rows, err := site.db.connection.Query(query, sqlArgs...)
	if err != nil {
		log.Panicf("Error querying '%s'\n%v", query, err)
	}

	for rows.Next() {
		file := &ParsedFile{isPage: isPage}
		rows.Scan(&file.Title, &file.Slug, &file.Content, &file.Category, &file.tags, &file.Date, &file.summary)
		files = append(files, file)
	}

	return files
}

func (site Site) QueryArticles(where string, args ...string) []*ParsedFile {
	return site.Query(false, where, args...)
}

func (site Site) QueryPages(where string, args ...string) []*ParsedFile {
	return site.Query(true, where, args...)
}
