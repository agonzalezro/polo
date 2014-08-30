package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/agonzalezro/polo/generator"
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

	db, err := generator.NewDB()
	if err != nil {
		log.Panic(err)
	}
	if err := db.Fill(*inputPath); err != nil {
		log.Panic(err)
	}

	config, err := generator.ParseConfigFile(*configFile)
	if err != nil {
		log.Panic(err)
	}

	site := generator.NewSite(*db, *config, *outputPath)
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
