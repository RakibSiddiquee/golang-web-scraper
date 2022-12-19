package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gocolly/colly"
)

func main() {

	db, err := sql.Open("mysql", "root:@/dinq")
	if err != nil {
		log.Fatalf("Error to connect: %s", err)
	}
	defer db.Close()

	c := colly.NewCollector(
		colly.Async(true),
	)
	cc := colly.NewCollector()
	c.SetRequestTimeout(25 * time.Second)
	//cc.SetRequestTimeout(25 * time.Second)

	c.OnXML("//loc", func(e *colly.XMLElement) {
		cc.Visit(e.Text)

	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	// Set error handler
	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	cc.OnHTML("html", func(h *colly.HTMLElement) {
		title := h.ChildText("h1[class=no-margin]")
		img := strings.Split(h.ChildAttr("meta[property=\"og:image\"]", "content"), "imgPath=")[1]
		//img := h.ChildAttr("div[class=featured-image] > img", "src")
		details := h.ChildText("div[class=content-details]")
		bgPath := "https://cdn.jagonews24.com/media/imgAllNew/BG/" + img
		smPath := "https://cdn.jagonews24.com/media/imgAllNew/SM/" + img
		xsPath := "https://cdn.jagonews24.com/media/imgAllNew/XS/" + img

		rand.Seed(time.Now().UnixNano())
		catId := rand.Intn(10-1) + 1

		if title != "" {
			stmt, e := db.Prepare("INSERT INTO bn_contents(cat_id, country_id, upozilla_id, content_heading, content_details, img_bg_path, img_sm_path, img_xs_path, uploader_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)")
			ErrorCheck(e)

			//execute
			res, e := stmt.Exec(catId, 19, 1, title, details, bgPath, smPath, xsPath, 1, time.Now(), time.Now())
			ErrorCheck(e)

			id, e := res.LastInsertId()
			ErrorCheck(e)

			fmt.Println("Insert id", id)
		}

		//fmt.Println(details)
	})

	cc.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting link", r.URL)
	})

	c.Visit("https://www.jagonews24.com/sitemaps/news-sitemap.xml")

	c.Wait()
}

func ErrorCheck(err error) {
	if err != nil {
		panic(err.Error())
	}
}
