package crawler

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
)

type ImmoscoutCrawler struct {
	searchUrl string
}

func NewImmoscoutCrawler(searchUrl string) *ImmoscoutCrawler {
	return &ImmoscoutCrawler{searchUrl: searchUrl}
}

func (i *ImmoscoutCrawler) Name() string {
	return "immoscout24.de"
}

func (ic *ImmoscoutCrawler) Crawl(c chan *CrawlResult) error {
	document, err := NewDocumentFromHttpGet(ic.searchUrl)
	if err != nil {
		return err
	}

	document.Find(".result-list-entry__brand-title-container").Each(func(i int, s *goquery.Selection) {
		expose, exists := s.Attr("href")
		if exists {
			url := "https://www.immobilienscout24.de" + expose
			c <- NewCrawlResult(url, nil)
		} else {
			log.Error().Msg("crawl error: tag not found")
		}
	})

	return nil
}
