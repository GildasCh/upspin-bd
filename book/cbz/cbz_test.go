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
	f, _ := os.Open("mifune.zip")
	c, err := NewCBZ(f)

	require.NoError(t, err)
	assert.NotNil(t, c.Reader)

	pages := ""
	for _, f := range c.pages {
		pages += f.FileHeader.Name + "\n"
	}

	expected := `mifune/3bf8c36bdeaf0785f43f83cc1987c57c--toshiro-mifune-red-beard.jpg
mifune/58af6ed8a5a1a.image.jpg
mifune/6-toshiro-mifune.w1200.h630.jpg
mifune/85c2bae500bf8ae72d485aaeb9864014--toshiro-mifune-sexy-men.jpg
mifune/9b74b5a06015f4f892d609c0e8cf6b03--rhys-davies-toshiro-mifune.jpg
mifune/9e06faf13a655750ca93ed4b2b74987e--toshiro-mifune-man-style.jpg
mifune/d4bdae14ac78527ae2e88c40df3ae8bf.jpg
mifune/Image-Mifune-Senses-08.jpg
mifune/mifune.jpg
mifune/mifune_4_large.jpg
mifune/Shubun_poster_Toshiro_Mifune.jpg
mifune/toshiro-mifune-star-wars.jpg
mifune/Toshiro_Mifune_EiganFan1952.jpg
`

	assert.Equal(t, expected, pages)
}

func TestPage(t *testing.T) {
	f, _ := os.Open("mifune.zip")
	c, err := NewCBZ(f)
	require.NoError(t, err)

	r, ok, err := c.Page(7)
	require.NoError(t, err)
	require.True(t, ok)

	b, err := ioutil.ReadAll(r)
	require.NoError(t, err)

	// sha1 sum of Image-Mifune-Senses-08.jpg
	expected := "2802672c394b2dd7a5d71a5d268905c827dc7416"

	sum := sha1.Sum(b)
	assert.Equal(t, expected, hex.EncodeToString(sum[:]))
}
