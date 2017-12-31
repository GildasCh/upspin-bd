package cbz

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewCBZ(t *testing.T) {
	f, _ := os.Open("test.cbz")
	fi, _ := f.Stat()
	c, err := NewCBZ(f, fi.Size())

	require.NoError(t, err)
	assert.NotNil(t, c.Reader)

	pages := ""
	for _, f := range c.pages {
		pages += f.FileHeader.Name + "\n"
	}

	expected := `mifune-free/01-Rashomon_poster_2.jpg
mifune-free/11629986985_267f712523_b.jpg
mifune-free/19971482_20b7f0fc5d_b.jpg
mifune-free/Shubun_poster_Toshiro_Mifune.jpg
mifune-free/toshiro_Mifune_character_in_filming_of_Hell_in_the_Pacific_4.jpg
mifune-free/Toshiro_Mifune_wearing_bandana.jpg
`

	assert.Equal(t, expected, pages)
}

func TestPage(t *testing.T) {
	f, _ := os.Open("test.cbz")
	fi, _ := f.Stat()
	c, err := NewCBZ(f, fi.Size())
	require.NoError(t, err)

	r, ok, err := c.Page(3)
	require.NoError(t, err)
	require.True(t, ok)

	b, err := ioutil.ReadAll(r)
	require.NoError(t, err)

	// sha1 sum of Shubun_poster_Toshiro_Mifune.jpg
	expected := "885cb5a39ef40fe0177e8571b10e0c4abf1e2c9e"

	sum := sha1.Sum(b)
	assert.Equal(t, expected, hex.EncodeToString(sum[:]))
}
