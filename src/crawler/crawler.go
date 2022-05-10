package crawler

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/rs/zerolog/log"
	"net/http"
	"time"
)

type Crawler interface {
	Name() string
	Crawl(chan *CrawlResult) error
}

type CrawlRunner struct {
	crawlers []Crawler
}

func NewCrawlRunner(crawlers []Crawler) *CrawlRunner {
	return &CrawlRunner{crawlers: crawlers}
}

func (cr *CrawlRunner) StartPolling(duration time.Duration) (<-chan *CrawlResult, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	results := make(chan *CrawlResult)

	for _, cr := range cr.crawlers {
		ticker := time.NewTicker(duration)

		go func(c Crawler, ticker *time.Ticker) {
			for {
				select {
				case <-ticker.C:
					err := c.Crawl(results)
					if err != nil {
						log.Err(err)
					}
				case <-ctx.Done():
					return
				}
			}
		}(cr, ticker)
	}

	return results, cancel
}

type CrawlResult struct {
	Url     string
	details map[string]string
}

func NewCrawlResult(url string, details map[string]string) *CrawlResult {
	if details == nil {
		details = map[string]string{}
	}
	return &CrawlResult{Url: url, details: details}
}

func NewDocumentFromHttpGet(url string) (*goquery.Document, error) {
	res, err := http.Get(url)
	if err != nil {
		log.Err(err)
		return nil, err
	}
	defer res.Body.Close()

	l := log.Info()
	if res.StatusCode >= 300 {
		l = log.Error()
	}
	l.Str("searchUrl", url).Int("status_code", res.StatusCode).Msg("crawling")

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Err(err)
		return nil, err
	}

	return doc, nil
}
