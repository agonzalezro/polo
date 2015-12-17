package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/agonzalezro/polo/config"
	"github.com/agonzalezro/polo/site"
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

	// TODO: this need to be clean up
	directoryExists := func(dir string) bool {
		if fi, err := os.Stat(dir); err == nil {
			return fi.IsDir()
		}
		return false
	}
	sourcedir := args[0]
	if !directoryExists(sourcedir) {
		fmt.Fprintf(os.Stderr, "The sourcedir must be an existent directory!")
		os.Exit(1)
	}
	outdir := args[1]
	if !directoryExists(outdir) {
		if err := os.Mkdir(outdir, os.ModePerm); err != nil {
			fmt.Fprintf(os.Stderr, "The outdir didn't exists and couldn't be created!")
		}
		time.Sleep(1 * time.Second) // Ugliest thing ever but it doesn't create the dir at time?
	}

	config, err := config.New(opts.Config)
	if err != nil {
		log.Println(err)
		log.Println("This usually happens because the JSON file is not well formed.")
		os.Exit(1)
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

		addr := fmt.Sprintf("localhost:%d", opts.ServerPort)
		log.Printf("Static server running on http://%s\n", addr)
		log.Fatal(
			http.ListenAndServe(addr, http.FileServer(http.Dir(outdir))))
	}
}
