package generator

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	sql.DB
}

type UniqueSlugError error

// Fill the DB with the articles an pages found
func (db *DB) Fill(root string) error {
	return filepath.Walk(root, db.parseAndSave)
}

// Append the paths to an array in case that they are markdown files.
// If there are pages (file inside the folder "pages") it's going to be
// saved with the value isPage = 1
func (db *DB) parseAndSave(path string, fileInfo os.FileInfo, err error) error {
	if err != nil {
		return err
	}

	slugsPresence := make(map[string]bool)

	if !fileInfo.Mode().IsDir() && strings.HasSuffix(path, ".md") {
		file := ParsedFile{}
		if err := file.load(path); err != nil {
			return err
		}

		if _, present := slugsPresence[file.Slug]; present {
			errorMessage := fmt.Sprintf("The slug '%s' already exists!", file.Slug)
			err = errors.New(errorMessage)
			return UniqueSlugError(err)
		}
		slugsPresence[file.Slug] = true

		if strings.HasPrefix(path, "pages/") || strings.Index(path, "/pages/") > 0 {
			file.isPage = true
		}

		if err := file.save(db); err != nil {
			return err
		}
	}

	return nil
}

// Create minimal DB struct.
// It's going to return a DB and it's your work to close it, we can not defer the close call.
func NewDB() (*DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	query := `
 	CREATE TABLE files (author text, title text, slug text, content text, category text, tags text, date text, status text, summary text, is_page integer);
	`
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Commit()

	if _, err = tx.Exec(query); err != nil {
		return nil, err
	}
	return &DB{*db}, nil
}
