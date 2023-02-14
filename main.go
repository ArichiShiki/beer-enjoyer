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
			command := update.Message.Command()
			switch command {
			case "cur", "next":
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
			case "enjoy":
				anim := tgbotapi.NewAnimation(update.Message.Chat.ID, tgbotapi.FilePath("./resources/enjoy.gif.mp4"))
				if _, err := bot.Send(anim); err != nil {
					log.Panic(err)
				}
			}
		}
	}
}
