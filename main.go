package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

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
		// if (update.Message != nil && update.Message.Text != "") {
		// 	t := []rune(update.Message.Text)
		// 	fmt.Printf("%v %v", len(t), t)
		// }
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
			case "ck":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "выберите год")

				msg.ReplyMarkup = newCalendarKeyboard(yearKeyboard, time.Now())

				if sendMsg, err := bot.Send(msg); err != nil {
					log.Panic(err)
				} else {
					calendarCache[sendMsg.MessageID] = &calendarInfo{state: yearKeyboard}
				}
			}
		} else if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			messageID := update.CallbackQuery.Message.MessageID
			var requests []tgbotapi.Chattable

			if info, ok := calendarCache[messageID]; ok {
				switch (info.state) {
				case yearKeyboard:
					info.state = monthKeyboard
					info.year, _ = strconv.Atoi(update.CallbackQuery.Data) // TODO catch error
					t, _ := time.Parse("2006-01-02", fmt.Sprintf("%04d-01-01", info.year))
					requests = []tgbotapi.Chattable{
						tgbotapi.NewEditMessageText(chatID, messageID, "выберите месяц"),
						tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, newCalendarKeyboard(monthKeyboard, t)),
					}
				case monthKeyboard:
					info.state = dayKeyboard
					info.month, _ = strconv.Atoi(update.CallbackQuery.Data) // TODO catch error
					t, _ := time.Parse("2006-01-02", fmt.Sprintf("%04d-%02d-01", info.year, info.month))
					requests = []tgbotapi.Chattable{
						tgbotapi.NewEditMessageText(chatID, messageID, "выберите день"),
						tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, newCalendarKeyboard(dayKeyboard, t)),
					}
				case dayKeyboard:
					if update.CallbackQuery.Data != " " {
						info.day, _ = strconv.Atoi(update.CallbackQuery.Data) // TODO catch error
						requests = []tgbotapi.Chattable{
							NewEditMessageDeleteInlineKeyboard(chatID, messageID),
							tgbotapi.NewEditMessageText(chatID, messageID, fmt.Sprintf("%02d/%02d/%04d", info.day, info.month, info.year)),
						}
						delete(calendarCache, messageID)
					}
				default:
					continue
				}
				//deleteRequest := tgbotapi.NewDeleteMessage(update.CallbackQuery.Message.Chat.ID, update.CallbackQuery.Message.MessageID)
				//msg := tgbotapi.NewMessage(update.CallbackQuery.Message.Chat.ID, fmt.Sprintf("Вы выбрали: %v", update.CallbackQuery.Data))
				for _, r := range requests {
					if _, err := bot.Request(r); err != nil {
						panic(err)
					}
				}
			}
		}
	}
}
