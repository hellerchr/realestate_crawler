package config

import (
	"os"
	"strconv"
	"time"
)

const (
	Production  string = "prod"
	Development string = "dev"
)

type Config struct {
	Profile           string
	RedisConfig       RedisConfig
	CrawlConfig       CrawlConfig
	TelegramBotConfig TelegramBotConfig
	ImmoScoutConfig   ImmoScoutConfig
	ImmonetConfig     ImmonetConfig
}

type ImmonetConfig struct {
	SearchUrl string
}

func newImmonetConfig() ImmonetConfig {
	return ImmonetConfig{
		SearchUrl: mustGetEnv("IMMONET_SEARCH_URL"),
	}
}

type TelegramBotConfig struct {
	Token         string
	ImmoChannelID int64
}

type CrawlConfig struct {
	PollInterval time.Duration
}
type ImmoScoutConfig struct {
	SearchUrl string
}

func NewConfig() *Config {
	return &Config{
		Profile:           newProfile(),
		CrawlConfig:       newCrawlConfig(),
		TelegramBotConfig: newTelegramBotConfig(),
		ImmoScoutConfig:   newImmoScoutConfig(),
		ImmonetConfig:     newImmonetConfig(),
	}
}

func newImmoScoutConfig() ImmoScoutConfig {
	return ImmoScoutConfig{
		SearchUrl: mustGetEnv("IMMOSCOUT_SEARCH_URL"),
	}
}

func newTelegramBotConfig() TelegramBotConfig {
	return TelegramBotConfig{
		Token:         mustGetEnv("BOT_TOKEN"),
		ImmoChannelID: mustGetEnvAsInt64("BOT_IMMO_CHANNEL_ID"),
	}
}

func newCrawlConfig() CrawlConfig {
	return CrawlConfig{
		PollInterval: time.Duration(mustGetEnvAsInt("CRAWLER_POLL_INTERVAL_SECONDS")) * time.Second,
	}
}

type RedisConfig struct {
	Hostname string
}

func newProfile() string {
	return mustGetEnv("PROFILE")
}

func mustGetEnvAsInt(key string) int {
	val := mustGetEnv(key)

	num, err := strconv.Atoi(val)
	if err != nil {
		panic("no integer key: " + key)
	}

	return num
}

func mustGetEnvAsInt64(key string) int64 {
	val := mustGetEnv(key)

	num, err := strconv.ParseInt(val, 10, 64)
	if err != nil {
		panic("no integer key: " + key)
	}

	return num
}

func mustGetEnv(key string) string {
	value, ok := os.LookupEnv(key)
	if !ok {
		panic("no env for: " + key)
	}

	return value
}
