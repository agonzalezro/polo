package parser

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/agonzalezro/go-polo/utils"
)

func parseMetadata(scanner *bufio.Scanner) (metadata map[string]string) {
	metadata = make(map[string]string)

	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "---") {
			scanner.Scan() // Read the last ---
			return metadata
		}

		metadataSplited := strings.Split(line, ":")
		key := strings.ToLower(metadataSplited[0])
		// It's possible that : is on the value of the metadata too (example: a date)
		metadata[key] = strings.Trim(strings.Join(metadataSplited[1:], ":"), " ")
	}

	log.Fatal("The metadata section was not properly closed!")
	return
}

func parseFile(filePath string) (parsedFile ParsedFile) {
	isFirstLine := true

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if line == "---" {
			parsedFile.Metadata = parseMetadata(scanner)
			continue
		}

		if isFirstLine {
			parsedFile.Title = line
			// TODO: check if there is a slug key on the metadata and don't assign it in that case
			parsedFile.Slug = utils.Slugify(line)
			scanner.Scan() // We don't want the title underlining

			isFirstLine = false
		} else {
			//bytesWithNewLine := append(line, []byte("\n")...)
			parsedFile.Content += line + "\n"
			//append(parsedFile.Content, bytesWithNewLine...)
		}

	}

	return parsedFile
}

func ParseFiles(articleFilePaths []string) []ParsedFile {
	var (
		slugsPresence map[string]bool
		parsedFiles   []ParsedFile
	)
	slugsPresence = make(map[string]bool)

	for _, filePath := range articleFilePaths {
		parsedFile := parseFile(filePath)
		parsedFiles = append(parsedFiles, parsedFile)

		if _, present := slugsPresence[parsedFile.Slug]; present {
			log.Fatalf("The slug '%s' already exists!", parsedFile.Slug)
		}
		slugsPresence[parsedFile.Slug] = true
	}
	return parsedFiles
}

func ParseConfig(configFile string) Config {
	file, err := os.Open(configFile)
	if err != nil {
		log.Panic(err)
	}

	decoder := json.NewDecoder(file)
	config := &Config{}

	err = decoder.Decode(&config)
	if err != nil {
		log.Panic(err)
	}
	return *config
}
