package book

import (
	"io"
	"strings"

	"github.com/gildasch/upspin-bd/book/cbz"
	"github.com/gildasch/upspin-bd/book/dir"
	"upspin.io/upspin"
)

type Book interface {
	Page(i int) (io.ReadCloser, bool, error)
	Pages() int
}

func NewFromUpspin(path string, client upspin.Client) (Book, bool, error) {
	// CBZ
	if strings.HasSuffix(strings.ToLower(path), ".cbz") {
		pathName := upspin.PathName(strings.TrimPrefix(path, "/"))
		return cbz.NewCBZFromUpspin(pathName, client.Open, client.Lookup)
	}

	// Directory
	pattern := extractPattern(path)
	return dir.NewDirFromUpspin(pattern, client.Glob, client.Open)
}

func extractPattern(in string) string {
	pattern := strings.TrimPrefix(in, "/")
	if !strings.Contains(pattern, "*") {
		pattern = strings.TrimSuffix(pattern, "/") + "/*"
	}
	return pattern
}
