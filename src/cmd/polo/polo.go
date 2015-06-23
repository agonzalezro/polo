package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"config"
	"site"

	flags "github.com/jessevdk/go-flags"
	fsnotify "gopkg.in/fsnotify.v1"
)

var opts struct {
	StartDaemon bool   `short:"d" long:"daemon" description:"start a simple HTTP server watching for markdown changes."`
	Config      string `short:"c" long:"config" default:"config.json" description:"the settings file."`
	ServerPort  int    `short:"p" long:"port" default:"8080" description:"port where to run the server."`
}

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

func main() {
	parser := flags.NewParser(&opts, flags.Default)
	parser.Usage = "[OPTIONS] sourcedir outputdir" + parser.Usage

	args, err := parser.Parse()
	if err != nil {
		panic(err)
	}

	if len(args) != 2 {
		parser.WriteHelp(os.Stderr)
		os.Exit(1)
	}

	sourcedir := args[0]
	outdir := args[1]

	config, err := config.New(opts.Config)
	if err != nil {
		log.Panic(err)
	}
	if err := writeSite(*config, sourcedir, outdir); err != nil {
		log.Fatal(err)
	}

	if opts.StartDaemon {
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

		addr := fmt.Sprintf(":%d", opts.ServerPort)
		log.Printf("Static server created on address %s\n", addr)
		log.Fatal(
			http.ListenAndServe(addr, http.FileServer(http.Dir(outdir))))
	}
}
