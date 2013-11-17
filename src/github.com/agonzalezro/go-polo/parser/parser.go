package parser

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	//"github.com/agonzalezro/go-polo/utils"
	"github.com/russross/blackfriday"
)

func getMetadata(filePath string) (metadata map[string]string, content []byte) {
	var (
		line            string
		isThereMetadata bool = false
	)
	metadata = make(map[string]string)

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)

	ok := scanner.Scan()
	if ok {
		line = scanner.Text()
		if line == "---" {
			isThereMetadata = true
		} else {
			// TODO: this is crap but I don't know how to seek to go back to
			// the beginning of the file
			metadata["title"] = line
			content, _ := ioutil.ReadFile(filePath)
			return metadata, content
		}
	}

	if isThereMetadata {
		for scanner.Scan() {
			line = scanner.Text()

			// We have finished reading the metadata
			if strings.HasPrefix(line, "---") {
				scanner.Scan() // Read one more line for the \n
				break
			}

			metadataLine := strings.Split(line, ":")
			key := strings.ToLower(metadataLine[0])
			metadata[key] = metadataLine[1]
		}
	}

	var isFirstLine bool = true
	for scanner.Scan() {
		bytes := scanner.Bytes()

		if isFirstLine {
			metadata["title"] = string(bytes)
			isFirstLine = false
		}

		bytesWithNewLine := append(bytes, []byte("\n")...)
		content = append(content, bytesWithNewLine...)
	}

	return
}

func fileToHtml(filePath string) (map[string]string, []byte) {
	metadata, content := getMetadata(filePath)

	fmt.Printf("%+v", metadata)

	html := blackfriday.MarkdownCommon(content)
	return nil, html

	//fmt.Printf("%s\n", utils.Slugify(title))
}

func FilesToHtml(articleFilePaths []string, output string) {
	for _, filePath := range articleFilePaths {
		_, html := fileToHtml(filePath)
		fmt.Printf("%s", html)
	}
}
