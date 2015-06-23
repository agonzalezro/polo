package main

import (
	"flag"
	"fmt"
	"gopkg.in/fsnotify.v1"
	"log"
	"net/http"
	"os"
	"path"
	"path/filepath"

	"config"
	"site"
)

var (
	configFile string
	daemon     bool
	port       int
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

func writeSite(config config.Config, sourcedir string, outdir string) error {
	s := site.New(config, outdir)
	if err := s.Populate(sourcedir); err != nil {
		return err
	}
	if err := s.Write(); err != nil {
		return err
	}
	return nil
}

func init() {
	flag.StringVar(&configFile, "config", "config.json", "the settings file to create your site.")
	flag.BoolVar(&daemon, "daemon", false, "create a simple HTTP server after the blog is created to see the result")
	flag.IntVar(&port, "port", 8080, "port where to run the server")

	flag.Usage = func() {
		fmt.Fprintf(
			os.Stderr, "Usage: %s [options] sourcedir outdir\n\n",
			path.Base(os.Args[0]))
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()
	if flag.NArg() != 2 {
		flag.Usage()
		os.Exit(1)
	}

	sourcedir := flag.Arg(0)
	outdir := flag.Arg(1)

	config, err := config.New(configFile)
	if err != nil {
		log.Panic(err)
	}
	if err := writeSite(*config, sourcedir, outdir); err != nil {
		log.Fatal(err)
	}

	if daemon {
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
						if err := writeSite(*config, sourcedir, outdir); err != nil {
							log.Fatal(err)
						}
					}
				case err := <-watcher.Errors:
					log.Fatal(err)
				}
			}
		}()

		paths, err := getAllSubdirectories(sourcedir)
		if err != nil {
			log.Fatal(err)
		}
		for _, path := range paths {
			watcher.Add(path)
		}

		addr := fmt.Sprintf(":%d", port)
		log.Printf("Static server created on address %s\n", addr)
		log.Fatal(
			http.ListenAndServe(addr, http.FileServer(http.Dir(outdir))))
	}
}
