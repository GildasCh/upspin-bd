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

	router.Static("/static", "./static")

	router.LoadHTMLFiles("templates/index.html")
	router.GET("/read/*path", func(c *gin.Context) {
		cb, ok, err := loadCBZ(client, c.Param("path"))
		if !ok {
			c.Status(http.StatusBadRequest)
			return
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"resource": "/load" + c.Param("path"),
			"pages":    cb.Pages(),
		})
	})

	router.GET("/load/*path", func(c *gin.Context) {
		cb, ok, err := loadCBZ(client, c.Param("path"))
		if !ok {
			c.Status(http.StatusBadRequest)
			return
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			return
		}

		pageString := c.Query("page")
		page, err := strconv.Atoi(pageString)
		if err != nil {
			fmt.Println(err)
			page = 0
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

func loadCBZ(client upspin.Client, path string) (*cbz.CBZ, bool, error) {
	pathName := upspin.PathName(strings.TrimPrefix(path, "/"))

	f, err := client.Open(pathName)
	if err != nil {
		fmt.Println(err)
		return nil, false, err
	}

	de, err := client.Lookup(pathName, true)
	if err != nil {
		fmt.Println(err)
		return nil, false, err
	}
	size := int64(0)
	for _, db := range de.Blocks {
		size += db.Size
	}

	cb, err := cbz.NewCBZ(f, size)
	if err != nil {
		fmt.Println(err)
		return nil, true, err
	}

	return cb, true, nil
}
