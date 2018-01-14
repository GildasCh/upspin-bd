package cbr

import (
	"fmt"
	"io"
	"io/ioutil"

	rar "github.com/nwaples/rardecode"
	"github.com/pkg/errors"
	"upspin.io/upspin"
)

type CBR struct {
	Reader func() (io.Reader, error)
	pages  int
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

	pages := 0
	for {
		fh, err := r.Next()
		if err != nil {
			break
		}
		if fh.IsDir {
			continue
		}
		pages++
	}

	return &CBR{Reader: f, pages: pages}, nil
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
	if i < 0 || i >= c.pages {
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

	n := 0
	for {
		fh, err := r.Next()
		if err != nil {
			return nil, true,
				errors.Wrapf(err, "could not read rar file up to page %d", i)
		}
		if fh.IsDir {
			continue
		}
		if n == i {
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
	return c.pages
}
