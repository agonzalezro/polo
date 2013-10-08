package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/agonzalezro/go-polo/utils"
)

func FilesToHtml(articleFilePaths []string, output string) {
	for _, element := range articleFilePaths {
		file, err := os.Open(element)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Scan()
		title := scanner.Text()

		fmt.Printf("%s\n", utils.Slugify(title))
	}
}
