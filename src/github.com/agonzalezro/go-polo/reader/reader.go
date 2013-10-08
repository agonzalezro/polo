package reader

import (
	"os"
	"path/filepath"
	"strings"
)

var (
	articleFilePaths []string
	pageFilePaths    []string
)

// Append the paths to an array in case that they are RST files.
// If there are pages (file inside the folder "pages") it's going to be
// included in one array, otherwise it's going to be included in other.
func visit(path string, fileInfo os.FileInfo, err error) error {
	if !fileInfo.Mode().IsDir() && strings.HasSuffix(path, ".rst") {
		if strings.HasPrefix(path, "pages/") || strings.Index(path, "/pages/") > 0 {
			pageFilePaths = append(pageFilePaths, path)
		} else {
			articleFilePaths = append(articleFilePaths, path)
		}
	}
	return nil
}

// Return two arrays, one with the pages and other with the articles
func GetPagesAndArticles(root string) ([]string, []string) {
	filepath.Walk(root, visit)
	return pageFilePaths, articleFilePaths
}
