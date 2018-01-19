package main

import (
	"flag"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gildasch/upspin-bd/book"
	"github.com/gin-gonic/gin"
	"upspin.io/client"
	"upspin.io/config"
	_ "upspin.io/transports"
)

func main() {
	confPathPtr := flag.String("config", "~/upspin/config", "path to the upspin configuration file")
	baseURLPtr := flag.String("baseURL", "", "the base URL of the service")
	flag.Parse()

	cfg, err := config.FromFile(*confPathPtr)
	if err != nil {
		fmt.Println(err)
		return
	}

	baseURL := *baseURLPtr

	client := client.New(cfg)

	router := gin.Default()

	router.Static("/static", "./static")

	router.LoadHTMLFiles("templates/index.html", "templates/list.html")

	router.GET("/list/*path", func(c *gin.Context) {
		books, dirs, err := book.List(c.Param("path"), client, true)
		if err != nil {
			c.Status(http.StatusInternalServerError)
			fmt.Printf("error listing %q: %v\n", c.Param("path"), err)
			return
		}

		type bookAndThumb struct {
			Name  string
			Thumb string
		}

		bookAndThumbs := []bookAndThumb{}
		for _, b := range books {
			bookAndThumbs = append(bookAndThumbs,
				bookAndThumb{Name: b, Thumb: "/load/" + b + "?page=0"})
		}

		c.HTML(http.StatusOK, "list.html", gin.H{
			"books":   bookAndThumbs,
			"dirs":    dirs,
			"baseURL": baseURL,
		})
	})

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
			"baseURL":  baseURL,
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
