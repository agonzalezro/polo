package generator

import (
	"bufio"
	"strings"
	"testing"
)

func TestMetadataGeneration(t *testing.T) {
	content := `---
date: expected_date
tags: tag1, tag2
---

This should be out of the metadata
`

	pf := ParsedFile{}
	pf.scanner = bufio.NewScanner(strings.NewReader(content))
	// Read the first line (because it's done in other function)
	pf.scanner.Scan()

	pf.parseMetadata()
	if pf.Date != "expected_date" {
		t.Errorf("Date is not the expected: %s", pf.Date)
	}
	if pf.tags != "tag1, tag2" {
		t.Errorf("Tags not expected: %s", pf.tags)
	}
	pf.scanner.Scan()
	if pf.scanner.Text() != "This should be out of the metadata" {
		t.Errorf("This line should be here")
	}
}
