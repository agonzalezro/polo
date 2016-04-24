package context

import (
	"sync"
	"testing"

	"github.com/agonzalezro/polo/config"
	"github.com/agonzalezro/polo/file"
	"github.com/stretchr/testify/assert"
)

func TestNumberOfPages(t *testing.T) {
	assert := assert.New(t)

	config := config.Config{
		PaginationSize: 1,
	}
	c := Context{
		Config:           config,
		numberOfPagesMux: &sync.Mutex{},
		Articles:         []file.ParsedFile{file.ParsedFile{}, file.ParsedFile{}},
	}
	assert.Equal(2, c.NumberOfPages())

	// Same context, old value
	c.Articles = []file.ParsedFile{file.ParsedFile{}}
	assert.Equal(2, c.NumberOfPages())

	// New context, old value
	c = Context{
		Config:           config,
		numberOfPagesMux: &sync.Mutex{},
		Articles:         []file.ParsedFile{},
	}
	assert.Equal(2, c.NumberOfPages())
}
