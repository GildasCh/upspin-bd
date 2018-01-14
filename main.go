package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/gildasch/upspin-bd/book"
	"github.com/gin-gonic/gin"
	"upspin.io/client"
	"upspin.io/config"
	_ "upspin.io/transports"
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
		b, ok, err := book.NewFromUpspin(c.Param("path"), client, true)
		if !ok {
			c.Status(http.StatusBadRequest)
			fmt.Printf("%q not ok: %v\n", c.Param("path"), err)
			return
		}
		if err != nil {
			c.Status(http.StatusInternalServerError)
			fmt.Printf("error getting %q: %v\n", c.Param("path"), err)
			return
		}

		pages := []string{}
		for i := 0; i < b.Pages(); i++ {
			pages = append(pages,
				"/load"+c.Param("path")+"?page="+strconv.Itoa(i))
		}

		c.HTML(http.StatusOK, "index.html", gin.H{
			"resource": "/load" + c.Param("path"),
			"pages":    pages,
		})
	})

	router.GET("/load/*path", func(c *gin.Context) {
		b, ok, err := book.NewFromUpspin(c.Param("path"), client, true)
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

		bytes, _, err := b.Page(page)
		if err != nil {
			fmt.Println(err)
			c.Status(http.StatusInternalServerError)
			return
		}

		c.Data(http.StatusOK, "", bytes)
	})

	router.Run()
}
