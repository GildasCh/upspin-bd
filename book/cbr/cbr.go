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
	io.Reader
	pages int
}

func NewCBR(f io.Reader) (*CBR, error) {
	r, err := rar.NewReader(f, "")
	if err != nil {
		return nil, errors.Wrap(err, "could not create rar reader")
	}
	// defer (&rar.ReadCloser{Reader: *r}).Close()

	pages := 0
	for {
		_, err = r.Next()
		if err != nil {
			break
		}
		pages++
	}

	return &CBR{Reader: f, pages: pages}, nil
}

func NewCBRFromUpspin(pathName upspin.PathName,
	open func(name upspin.PathName) (upspin.File, error),
	lookup func(name upspin.PathName, followFinal bool) (*upspin.DirEntry, error)) (*CBR, bool, error) {
	f, err := open(pathName)
	if err != nil {
		fmt.Println(err)
		return nil, false, err
	}

	cb, err := NewCBR(f)
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

	r, err := rar.NewReader(c.Reader, "")
	if err != nil {
		return nil, true, errors.Wrap(err, "could not create rar reader")
	}
	// defer (&rar.ReadCloser{Reader: *r}).Close()

	n := 0
	for {
		_, err := r.Next()
		if err != nil {
			return nil, true,
				errors.Wrapf(err, "could not read rar file up to page %d", i)
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
