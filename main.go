package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"strings"
	"time"
)

type item struct {
	Title   string `json:"title"`
	BgImg   string `json:"bgimg"`
	SmImg   string `json:"smimg"`
	XsImg   string `json:"xsimg"`
	Details string `json:"details"`
}

func main() {
	//knownUrls := []string{}
	items := []item{}

	c := colly.NewCollector(
		colly.Async(true),
	)
	cc := colly.NewCollector()
	c.SetRequestTimeout(25 * time.Second)
	//cc.SetRequestTimeout(25 * time.Second)

	c.OnXML("//loc", func(e *colly.XMLElement) {
		cc.Visit(e.Text)
		//if e.Response.StatusCode == 200 {
		//	knownUrls = append(knownUrls, e.Text)
		//}

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
		if title != "" {
			item := item{
				Title:   title,
				BgImg:   "https://cdn.jagonews24.com/media/imgAllNew/BG/" + img,
				SmImg:   "https://cdn.jagonews24.com/media/imgAllNew/SM/" + img,
				XsImg:   "https://cdn.jagonews24.com/media/imgAllNew/XS/" + img,
				Details: details,
			}
			items = append(items, item)
		}

		//fmt.Println(details)
	})

	cc.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting link", r.URL)
	})

	c.Visit("https://www.jagonews24.com/sitemaps/news-sitemap.xml")

	//titles := []string{}
	//for _, link := range knownUrls {
	//	c.OnHTML("title", func(h *colly.HTMLElement) {
	//		fmt.Println(h.Text)
	//		titles = append(titles, h.Text)
	//	})
	//
	//	c.Visit(link)
	//}

	c.Wait()
	//cc.Wait()

	// Convert results to JSON data if the scraping job has finished
	jsonData, err := json.MarshalIndent(items, "", "  ")
	if err != nil {
		panic(err)
	}

	// Dump json to the standard output (can be redirected to a file)
	fmt.Println(string(jsonData))

	//fmt.Println(titles)
	//fmt.Println(knownUrls)
}
