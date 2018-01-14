package book

import (
	"strings"

	"github.com/gildasch/upspin-bd/book/cbr"
	"github.com/gildasch/upspin-bd/book/cbz"
	"github.com/gildasch/upspin-bd/book/dir"
	"upspin.io/upspin"
)

type Book interface {
	Page(i int) ([]byte, bool, error)
	Pages() int
}

func NewFromUpspin(path string, client upspin.Client, useCache bool) (b Book, ok bool, err error) {
	if useCache {
		if b, ok := cache[path]; ok {
			return b, true, nil
		}
	}

	if strings.HasSuffix(strings.ToLower(path), ".cbz") {
		// CBZ
		pathName := upspin.PathName(strings.TrimPrefix(path, "/"))
		b, ok, err = cbz.NewCBZFromUpspin(pathName, client.Open, client.Lookup)
	} else if strings.HasSuffix(strings.ToLower(path), ".cbr") {
		// CBR
		pathName := upspin.PathName(strings.TrimPrefix(path, "/"))
		b, ok, err = cbr.NewCBRFromUpspin(pathName, client.Open, client.Lookup)
	} else {
		// Directory
		pattern := extractPattern(path)
		b, ok, err = dir.NewDirFromUpspin(pattern, client.Glob, client.Open)
	}

	cacheUpdate(path, b, ok, err)
	return
}

func extractPattern(in string) string {
	pattern := strings.TrimPrefix(in, "/")
	if !strings.Contains(pattern, "*") {
		pattern = strings.TrimSuffix(pattern, "/") + "/*"
	}
	return pattern
}
