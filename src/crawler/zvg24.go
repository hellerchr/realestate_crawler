package crawler

import (
	"github.com/PuerkitoBio/goquery"
)

type Zvg24Crawler struct {
	searchUrl string
}

func NewZvg24Crawler(searchUrl string) *Zvg24Crawler {
	return &Zvg24Crawler{searchUrl: searchUrl}
}

func (z Zvg24Crawler) Name() string {
	return "zvg24"
}

func (z Zvg24Crawler) Crawl(c chan *CrawlResult) error {
	document, err := NewDocumentFromHttpGet(z.searchUrl)
	if err != nil {
		return err
	}

	document.Find("article .content a").Each(func(i int, selection *goquery.Selection) {
		link, exists := selection.Attr("href")
		if exists {
			c <- NewCrawlResult(link, nil)
		}
	})
	return nil
}
