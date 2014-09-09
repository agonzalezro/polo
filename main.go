package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/agonzalezro/polo/config"
	"github.com/agonzalezro/polo/site"
)

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

	site := site.New(*config, *outputPath)
	if err := site.Populate(*inputPath); err != nil {
		log.Panic(err)
	}
	if err := site.Write(); err != nil {
		log.Panic(err)
	}

	if *daemon {
		addr := fmt.Sprintf(":%d", *port)
		log.Printf("Static server created on address %s\n", addr)
		log.Fatal(
			http.ListenAndServe(addr, http.FileServer(http.Dir(*outputPath))))
	}
}
