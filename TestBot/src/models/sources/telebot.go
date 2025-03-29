package sources

import (
	"time"

	"gopkg.in/telebot.v4"
)

func TelebotInitializer() (*telebot.Bot, error) {
	pref := telebot.Settings{
        Token:  "7636991008:AAG3KlhZkWvje7mt_ZV8QxW8l0D_2jA3qww",
        Poller: &telebot.LongPoller{Timeout: 10 * time.Second},
    }

	bot, err := telebot.NewBot(pref)

	if err != nil {
        return nil, err
    }

	commands := []telebot.Command{
        {Text: "start", Description: "Запустить бота"},
        {Text: "about", Description: "Узнать о боте"},
    }

	if err := bot.SetCommands(commands); err != nil {
        return nil, err
    }

	return bot, nil
}
