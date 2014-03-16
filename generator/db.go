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

// Append the paths to an array in case that they are markdown files.
// If there are pages (file inside the folder "pages") it's going to be
// saved with the value isPage = 1
func (db *DB) parseAndSave(path string, fileInfo os.FileInfo, err error) error {
	if err != nil {
		log.Panicf("Error walking through the path: %s\n%v", path, err)
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
		}

		if err := file.save(db); err != nil {
			log.Panicf("Error saving the parsed file: %s\n%v", path, err)
		}
	}

	return nil
}

// Create minimal DB struct.
// It's going to return a DB and it's your work to close it, we can not defer the close call.
func GetDB() *DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		log.Panicf("Impossible to open DB in memory!\n%v", err)
	}

	query := `
	CREATE table files (author text, title text, slug text, content text, category text, tags text, date text, status text, summary text, is_page integer);
	`
	if _, err = db.Exec(query); err != nil {
		log.Panicf("Error creating the DB: %v", err)
		return nil
	}
	return &DB{*db}
}
