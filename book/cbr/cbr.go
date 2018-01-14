package cbr

import (
	"fmt"
	"io"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/gildasch/upspin-bd/book/types"
	rar "github.com/nwaples/rardecode"
	"github.com/pkg/errors"
	"upspin.io/upspin"
)

type CBR struct {
	Reader func() (io.Reader, error)
	pages  []page
}

type page struct {
	name           string
	placeInArchive int
}

func NewCBR(f func() (io.Reader, error)) (*CBR, error) {
	fr, err := f()
	if err != nil {
		return nil, errors.Wrap(err, "could not open upspin file")
	}

	r, err := rar.NewReader(fr, "")
	if err != nil {
		return nil, errors.Wrap(err, "could not create rar reader")
	}
	// defer (&rar.ReadCloser{Reader: *r}).Close()

	i := 0
	pages := []page{}
	for {
		fh, err := r.Next()
		if err != nil {
			break
		}
		if !types.IsImage(fh.Name, fh.IsDir) {
			continue
		}
		pages = append(pages, page{fh.Name, i})
		i++
	}

	sort.Sort(byName(pages))

	return &CBR{Reader: f, pages: pages}, nil
}

type byName []page

func (p byName) Len() int {
	return len(p)
}

func (p byName) Less(i, j int) bool {
	is, js := strings.ToLower(p[i].name), strings.ToLower(p[j].name)

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

func NewCBRFromUpspin(pathName upspin.PathName,
	open func(name upspin.PathName) (upspin.File, error),
	lookup func(name upspin.PathName, followFinal bool) (*upspin.DirEntry, error)) (*CBR, bool, error) {
	freader := func() (io.Reader, error) {
		return open(pathName)
	}

	cb, err := NewCBR(freader)
	if err != nil {
		fmt.Println(err)
		return nil, true, err
	}

	return cb, true, nil
}

func (c *CBR) Page(i int) ([]byte, bool, error) {
	if i < 0 || i >= c.Pages() {
		return nil, false, nil
	}

	fr, err := c.Reader()
	if err != nil {
		return nil, true, errors.Wrap(err, "could not open upspin file")
	}

	r, err := rar.NewReader(fr, "")
	if err != nil {
		return nil, true, errors.Wrap(err, "could not create rar reader")
	}
	// defer (&rar.ReadCloser{Reader: *r}).Close()

	stopAt := c.pages[i].placeInArchive

	n := 0
	for {
		fh, err := r.Next()
		if err != nil {
			return nil, true,
				errors.Wrapf(err, "could not read rar file up to page %d", i)
		}
		if !types.IsImage(fh.Name, fh.IsDir) {
			continue
		}
		if n == stopAt {
			break
		}
		n++
	}

	bytes, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, true, errors.Wrap(err, "could not read whole page in cbr")
	}

	return bytes, true, nil
}

func (c *CBR) Pages() int {
	return len(c.pages)
}
