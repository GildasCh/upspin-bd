package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gildasch/upspin-bd/book/cbz"
	"github.com/gin-gonic/gin"
	"upspin.io/client"
	"upspin.io/config"
	_ "upspin.io/transports"
	"upspin.io/upspin"
)

func main() {
	cfg, err := config.FromFile(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	client := client.New(cfg)

	router := gin.Default()

	router.GET("/read/*path", func(c *gin.Context) {
		fmt.Println("path:", c.Param("path"))
		path := upspin.PathName(strings.TrimPrefix(c.Param("path"), "/"))

		f, err := client.Open(path)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}

		de, err := client.Lookup(path, true)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusBadRequest)
			return
		}
		size := int64(0)
		for _, db := range de.Blocks {
			size += db.Size
		}

		pageString := c.Query("page")
		page, err := strconv.Atoi(pageString)
		if err != nil {
			fmt.Println(err)
			page = 0
		}

		cb, err := cbz.NewCBZ(f, size)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		rc, _, err := cb.Page(page)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}
		defer rc.Close()

		c.Stream(func(w io.Writer) bool {
			_, err := io.CopyN(w, rc, 1024*1024)
			return err == nil
		})
	})

	router.Run()
}
