package generator

import (
	"database/sql"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	connection sql.DB
}

// Fill the DB with the articles an pages found
func (db *DB) Fill(root string) {
	filepath.Walk(root, db.parseAndSave)
}

var (
	articles []ParsedFile
	pages    []ParsedFile
)

// Append the paths to an array in case that they are RST files.
// If there are pages (file inside the folder "pages") it's going to be
// included in one array, otherwise it's going to be included in other.
func (db *DB) parseAndSave(path string, fileInfo os.FileInfo, err error) error {
	if err != nil {
		log.Panic(err)
	}

	slugsPresence := make(map[string]bool)

	if !fileInfo.Mode().IsDir() && strings.HasSuffix(path, ".md") {
		file := ParsedFile{}
		file.load(path)

		if _, present := slugsPresence[file.Slug]; present {
			log.Fatalf("The slug '%s' already exists!", file.Slug)
		}
		slugsPresence[file.Slug] = true

		if strings.HasPrefix(path, "pages/") || strings.Index(path, "/pages/") > 0 {
			file.isPage = true
			pages = append(pages, file)
		} else {
			articles = append(articles, file)
		}

		if err := file.save(db); err != nil {
			log.Panic(err)
		}
	}

	return nil
}

// Create minimal DB struct.
// It's going to return a DB and it's your work to close it, we can not defer the close call.
func GetDB() *DB {
	db, err := sql.Open("sqlite3", "/tmp/db.sqlite")
	if err != nil {
		log.Panic("Impossible to open DB in memory!")
	}

	query := `
	CREATE table files (title text, slug text, content text, tags text, date text, is_page integer);
	`
	if _, err = db.Exec(query); err != nil {
		log.Panic("%q: %s", err, query)
		return nil
	}
	return &DB{*db}
}
