package dir

import (
	"fmt"
	"io"

	"github.com/pkg/errors"
	"upspin.io/upspin"
)

type Dir struct {
	pages []func() (io.ReadCloser, error)
}

func NewDirFromUpspin(pattern string,
	glob func(pattern string) ([]*upspin.DirEntry, error),
	open func(name upspin.PathName) (upspin.File, error)) (*Dir, bool, error) {
	des, err := glob(pattern)
	if err != nil {
		fmt.Println(err)
		return nil, false, err
	}

	if len(des) <= 0 {
		return nil, false, errors.Errorf("no file matches pattern %q", pattern)
	}

	pages := []func() (io.ReadCloser, error){}

	for _, de := range des {
		name := de.Name
		pages = append(pages, func() (io.ReadCloser, error) {
			return open(name)
		})
	}

	return &Dir{pages: pages}, true, nil
}

func (d *Dir) Page(i int) (io.ReadCloser, bool, error) {
	if len(d.pages) <= i {
		return nil, false, nil
	}

	rc, err := d.pages[i]()
	if err != nil {
		return nil, true, errors.Wrap(err, "could not open file in dir")
	}

	return rc, true, nil
}

func (d *Dir) Pages() int {
	return len(d.pages)
}
