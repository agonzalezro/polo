package parser

type Site struct {
	Pages, Articles []ParsedFile
}

type ParsedFile struct {
	Metadata map[string]string
	Content  []byte
}
