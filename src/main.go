package main

import (
	"crawler/src/bot"
	"crawler/src/config"
	"crawler/src/crawler"

	"github.com/go-redis/redis/v7"
)

func main() {
	cfg := config.NewConfig()

	var b bot.Bot
	if cfg.Profile == "cloud" {
		telegramBot, err := bot.NewTelegramBot(cfg.TelegramBotConfig.Token)
		if err != nil {
			panic(err)
		}
		b = telegramBot
	} else {
		b = bot.NewMockBot()
	}

	client := NewRedisClient(cfg.RedisConfig)

	const houseSearch = "https://www.immobilienscout24.de/Suche/shape/haus-kaufen?constructionphasetypes=completed&shape=c256a0h7YHhyQHJ8R29UbGpGb3xIYVVzZkJybENrX0RsbEJfeUxjd0B1ZlF9ZEZ3X0NldUt5dkBrbUFgXXl7QXp6QnlmQHJAa29BZl9GbXRHdnVAZWlCcHBBa05qdEpxd0Fic0h2b0B_e0Z4Z0JyeUJ4ZEFnakB_Z0NkdUNqaUNieUE.&price=-650000.0&constructionyear=1990-&ground=320.0-&enteredFrom=result_list"
	const grundstuecksSuche = "https://www.immobilienscout24.de/Suche/shape/grundstueck-kaufen?plotarea=350.0-&shape=c256a0h7YHhyQHJ8R29UbGpGb3xIYVVzZkJybENrX0RsbEJfeUxjd0B1ZlF9ZEZ3X0NldUt5dkBrbUFgXXl7QXp6QnlmQHJAa29BZl9GbXRHdnVAZWlCcHBBa05qdEpxd0Fic0h2b0B_e0Z4Z0JyeUJ4ZEFnakB_Z0NkdUNqaUNieUE.&price=100000.0-320000.0&pricetype=buy&sitedevelopmenttypes=developedpartially,developed&sorting=2&enteredFrom=result_list#/"
	const zvgSearch = "https://www.zvg24.net/zwangsversteigerung/walldorf-baden?search_form%5Bcategories%5D%5B%5D=1&search_form%5Bcategories%5D%5B%5D=3&search_form%5Blimit%5D=2&search_form%5Blocation_id%5D=12279&search_form%5Bmax_range%5D=500000.00&search_form%5Bmin_range%5D=0.00&search_form%5Border%5D=1&search_form%5Bradius%5D=25&search_form%5Bsearch%5D=Walldorf%2C+Baden"
	runner := crawler.NewCrawlRunner([]crawler.Crawler{
		crawler.NewImmoscoutCrawler(houseSearch),
		crawler.NewImmoscoutCrawler(grundstuecksSuche),
		crawler.NewZvg24Crawler(zvgSearch),
	})

	results, _ := runner.StartPolling(cfg.CrawlConfig.PollInterval)
	for result := range results {
		if isNewResult(client, result) {
			b.SendMessage(cfg.TelegramBotConfig.ImmoChannelID, result.Url)
		}
	}
}

func isNewResult(client *redis.Client, result *crawler.CrawlResult) bool {
	numAdded := client.SAdd("crawled", result.Url)
	return numAdded.Val() == 1
}

func NewRedisClient(cfg config.RedisConfig) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.Hostname + ":6379",
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}
	return client
}
