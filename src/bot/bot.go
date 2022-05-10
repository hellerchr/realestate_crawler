package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/rs/zerolog/log"
)

type Bot interface {
	SendMessage(channel int64, message string)
}

type MockBot struct {
}

func NewMockBot() *MockBot {
	return &MockBot{}
}

func (m MockBot) SendMessage(channel int64, message string) {
	println(message)
}

type TelegramBot struct {
	token string
	bot   *tgbotapi.BotAPI
}

func (t TelegramBot) SendMessage(channel int64, message string) {
	msg := tgbotapi.NewMessage(channel, message)
	t.bot.Send(msg)
	log.Info().Int64("channel", channel).Str("message", message).Msg("Telegram Bot sending message")
}

func NewTelegramBot(token string) (*TelegramBot, error) {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		return nil, err
	}
	u := tgbotapi.NewUpdate(0)
	updates, err := bot.GetUpdatesChan(u)

	go func() {
		for {
			<-updates
		}
	}()

	return &TelegramBot{token: token, bot: bot}, nil
}
