package main

//go:generate go-bindata -o ../../templates/assets.go -pkg=assets -ignore=.DS_Store -ignore=assets.go -prefix=../.. ../../templates/...

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/Sirupsen/logrus"
	config "github.com/agonzalezro/polo/config"
	"github.com/agonzalezro/polo/site"
	kingpin "gopkg.in/alecthomas/kingpin.v2"
	fsnotify "gopkg.in/fsnotify.v1"
)

var (
	app = kingpin.New("polo", `Static site generator "compatible" with Jekyll & Pelican content.`)

	startDaemon = app.Flag("start-daemon", "Start a simple HTTP server watching for markdown changes.").Short('d').Bool()
	port        = app.Flag("port", "Port where to run the server.").Default("8080").Short('p').Int()
	configPath  = app.Flag("config", "The settings file.").Short('c').Default("config.json").String()

	templatesBasePath = app.Flag("templates-base-path", fmt.Sprintf("Where the '%s/' folder resides (in case it exists).", site.TemplatesRelativePath)).Default(".").ExistingDir()

	verbose = app.Flag("verbose", "Verbose logging.").Short('v').Bool()

	source = app.Arg("source", "Folder where the content resides.").Required().ExistingDir()
	output = app.Arg("output", "Where to store the published files.").Required().String()
)

func main() {
	app.HelpFlag.Short('h')
	kingpin.MustParse(app.Parse(os.Args[1:]))

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	if !dirExists(*output) {
		if err := os.Mkdir(*output, os.ModePerm); err != nil {
			app.FatalUsage("The output folder couldn't be created!")
		}
	}

	s, err := site.New(*source, *output, *configPath, *templatesBasePath)
	if err != nil {
		switch err.(type) {
		case config.ErrorParsingConfigFile:
			app.FatalUsage("Malformed JSON config file: ", err)
		default:
			log.Fatal(err)
		}
	}

	if err := s.Write(); err != nil {
		log.Fatal(err)
	}

	if *startDaemon {
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
						log.Info("Rewriting the site")
						if err := s.Write(); err != nil {
							log.Fatal(err)
						}
					}
				case err := <-watcher.Errors:
					log.Fatal(err) // TODO: perhaps return err
				}
			}
		}()

		for path := range subdirectories(*source) {
			watcher.Add(path)
		}

		addr := fmt.Sprintf(":%d", *port)
		log.Info("Static server running on ", addr)
		log.Fatal(http.ListenAndServe(addr, http.FileServer(http.Dir(*output))))
	}
}
