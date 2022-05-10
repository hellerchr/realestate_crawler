package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

type ImmonetCrawler struct {
	searchUrl string
}

func NewImmonetCrawler(searchUrl string) *ImmonetCrawler {
	return &ImmonetCrawler{searchUrl: searchUrl}
}

func (ic ImmonetCrawler) Name() string {
	return "immonet"
}

func (ic ImmonetCrawler) Crawl(c chan *CrawlResult) error {
	document, err := NewDocumentFromHttpGet(ic.searchUrl)
	if err != nil {
		return err
	}

	document.Find(".listitem > a").Each(func(i int, s *goquery.Selection) {
		href, exists := s.Attr("href")
		url := "https://www.immowelt.de" + href
		if exists {
			c <- NewCrawlResult(url, nil)
		} else {
			log.Error().Msg("crawl error: tag not found")
		}
	})

	return nil
}
