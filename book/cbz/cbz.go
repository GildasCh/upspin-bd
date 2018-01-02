package cbz

import (
	"archive/zip"
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"
	"sync"

	"github.com/pkg/errors"
	"upspin.io/upspin"
)

type CBZ struct {
	*zip.Reader
	pages []*zip.File
	mutex *sync.Mutex
}

func NewCBZ(f io.ReaderAt, size int64) (*CBZ, error) {
	r, err := zip.NewReader(f, size)
	if err != nil {
		return nil, errors.Wrap(err, "could not create zip reader")
	}

	return &CBZ{Reader: r, pages: pages(r), mutex: &sync.Mutex{}}, nil
}

func NewCBZFromUpspin(pathName upspin.PathName,
	open func(name upspin.PathName) (upspin.File, error),
	lookup func(name upspin.PathName, followFinal bool) (*upspin.DirEntry, error)) (*CBZ, bool, error) {
	f, err := open(pathName)
	if err != nil {
		fmt.Println(err)
		return nil, false, err
	}

	de, err := lookup(pathName, true)
	if err != nil {
		fmt.Println(err)
		return nil, false, err
	}
	size := int64(0)
	for _, db := range de.Blocks {
		size += db.Size
	}

	cb, err := NewCBZ(f, size)
	if err != nil {
		fmt.Println(err)
		return nil, true, err
	}

	return cb, true, nil
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

func (c *CBZ) Page(i int) ([]byte, bool, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if len(c.pages) <= i {
		return nil, false, nil
	}

	rc, err := c.pages[i].Open()
	if err != nil {
		return nil, true, errors.Wrap(err, "could not open file in cbz")
	}
	defer rc.Close()

	bytes, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, true, errors.Wrap(err, "could not read whole file in cbz")
	}

	return bytes, true, nil
}

func (c *CBZ) Pages() int {
	return len(c.pages)
}
