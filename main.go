package main

import (
	"log"
	"os"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func main() {
	bot, err := tgbotapi.NewBotAPI(os.Getenv("TELEGRAM_APITOKEN"))
	if err != nil {
		log.Panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	// Loop through each update.
	for update := range updates {
		if update.Message != nil && update.Message.IsCommand() {
			var options []string
			question := "Mory:"
			
			switch update.Message.Command() {
			case "cur":
				question = curWeekQuestion
				options = getCurWeekPollParams()
			case "next":
				question = nextWeekQuestion
				options = getNextWeekPollParams()
			}

			if options != nil {
				poll := tgbotapi.NewPoll(update.Message.Chat.ID, question, options...)

				poll.IsAnonymous = false
				poll.AllowsMultipleAnswers = true

				if _, err := bot.Send(poll); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
