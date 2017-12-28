package cbz

import (
	"archive/zip"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

type CBZ struct {
	*zip.Reader
	pages []*zip.File
}

func NewCBZ(f *os.File) (*CBZ, error) {
	fi, err := f.Stat()
	if err != nil {
		return nil, errors.Wrap(err, "could not stat file")
	}

	r, err := zip.NewReader(f, fi.Size())
	if err != nil {
		return nil, errors.Wrap(err, "could not create zip reader")
	}

	return &CBZ{Reader: r, pages: pages(r)}, nil
}

func pages(r *zip.Reader) (pages []*zip.File) {
	for _, f := range r.File {
		if isDir(f) {
			continue
		}
		pages = append(pages, f)
	}

	sort.Sort(byName(pages))

	return
}

func isDir(f *zip.File) bool {
	return f.UncompressedSize == 0
}

type byName []*zip.File

func (p byName) Len() int {
	return len(p)
}

func (p byName) Less(i, j int) bool {
	is, js := strings.ToLower(p[i].Name), strings.ToLower(p[j].Name)

	k := 0
	for {
		if len(is) <= k {
			return true
		}
		if len(js) <= k {
			return false
		}
		if is[k] < js[k] {
			return true
		}
		if is[k] > js[k] {
			return false
		}
		k++
	}
}

func (p byName) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (c *CBZ) Page(i int) (io.ReadCloser, bool, error) {
	if len(c.pages) <= i {
		return nil, false, nil
	}

	rc, err := c.pages[i].Open()
	if err != nil {
		return nil, true, errors.Wrap(err, "could not open file in cbz")
	}

	return rc, true, nil
}
