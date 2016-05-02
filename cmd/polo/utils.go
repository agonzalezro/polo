package main

import (
	"os"
	"path/filepath"

	log "github.com/Sirupsen/logrus"
)

func subdirectories(parentPath string) chan string {
	paths := make(chan string, 1)

	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Warning(err)
		}
		if info.IsDir() {
			paths <- path
		}
		return nil
	}

	go func() {
		filepath.Walk(parentPath, walkFn)
		close(paths)
	}()

	return paths
}

func dirExists(p string) bool {
	if fi, err := os.Stat(p); err == nil {
		return fi.IsDir()
	}
	return false
}
