package main

import (
	"flag"
	"fmt"
	"gopkg.in/fsnotify.v1"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/agonzalezro/polo/config"
	"github.com/agonzalezro/polo/site"
)

func getAllSubdirectories(parentPath string) (paths []string, err error) {
	walkFn := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			paths = append(paths, path)
		}
		return nil
	}
	err = filepath.Walk(parentPath, walkFn)
	return paths, err
}

func writeSite(config config.Config, inputPath string, outputPath string) error {
	s := site.New(config, outputPath)
	if err := s.Populate(inputPath); err != nil {
		return err
	}
	if err := s.Write(); err != nil {
		return err
	}
	return nil
}

func main() {
	var (
		inputPath = flag.String("input", ".",
			"path to your articles source files.")
		outputPath = flag.String("output", ".",
			"path where you want to creat the html files.")
		configFile = flag.String("config", "config.json",
			"the settings file to create your site.")
		daemon = flag.Bool("daemon", false,
			"create a simple HTTP server after the blog is created to see the result")
		port = flag.Int("port", 8080,
			"port where to run the server")
	)
	flag.Parse()

	config, err := config.New(*configFile)
	if err != nil {
		log.Panic(err)
	}
	if err := writeSite(*config, *inputPath, *outputPath); err != nil {
		log.Fatal(err)
	}

	if *daemon {
		watcher, err := fsnotify.NewWatcher()
		if err != nil {
			log.Fatal(err)
		}
		defer watcher.Close()

		go func() {
			for {
				select {
				case event := <-watcher.Events:
					if event.Op != fsnotify.Chmod {
						log.Println("Rewriting the site")
						if err := writeSite(*config, *inputPath, *outputPath); err != nil {
							log.Fatal(err)
						}
					}
				case err := <-watcher.Errors:
					log.Fatal(err)
				}
			}
		}()

		paths, err := getAllSubdirectories(*inputPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, path := range paths {
			watcher.Add(path)
		}

		addr := fmt.Sprintf(":%d", *port)
		log.Printf("Static server created on address %s\n", addr)
		log.Fatal(
			http.ListenAndServe(addr, http.FileServer(http.Dir(*outputPath))))
	}
}
