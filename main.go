package main

import (
	"log"
	"os"

	"arichi/beerEnjoyer/calendar"
	"arichi/beerEnjoyer/poll"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type calendarID struct {
	chatID int64
	messID int
	userID int64
}

var calendarCache = make(map[calendarID]*calendar.CalendarInfo)

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
				var p poll.Poll
				switch update.Message.Command() {
				case "cur":
					p = poll.CreatePoll(poll.CurWeekPoll)
				case "next":
					p = poll.CreatePoll(poll.NextWeekPoll)
				}

				poll := tgbotapi.NewPoll(update.Message.Chat.ID, p.Question, p.Options...)

				poll.IsAnonymous = false
				poll.AllowsMultipleAnswers = true

				if _, err := bot.Send(poll); err != nil {
					log.Panic(err)
				}
			case "enjoy":
				anim := tgbotapi.NewAnimation(update.Message.Chat.ID, tgbotapi.FilePath("./resources/enjoy.gif.mp4"))
				if _, err := bot.Send(anim); err != nil {
					log.Panic(err)
				}
			case "ck":
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "выберите год")

				calendarInfo, _ := calendar.NewCalendar(calendar.DayKeyboard)

				msg.ReplyMarkup = calendarInfo.Keyboard

				if sendMsg, err := bot.Send(msg); err != nil {
					log.Panic(err)
				} else {
					id := calendarID {
						chatID: sendMsg.Chat.ID,
						messID: sendMsg.MessageID,
						userID: update.Message.From.ID,
					}
					calendarCache[id] = calendarInfo
				}
			}
		} else if update.CallbackQuery != nil {
			chatID := update.CallbackQuery.Message.Chat.ID
			messageID := update.CallbackQuery.Message.MessageID
			calID := calendarID {
				chatID: chatID,
				messID: messageID,
				userID: update.CallbackQuery.From.ID,
			}
			var requests []tgbotapi.Chattable

			if info, ok := calendarCache[calID]; ok && update.CallbackQuery.Data != calendar.EmptyData {
				info.Update(update.CallbackQuery.Data)
				switch (info.State) {
				case calendar.YearKeyboard:
					requests = []tgbotapi.Chattable{
						tgbotapi.NewEditMessageText(chatID, messageID, "выберите год"),
						tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, info.Keyboard),
					}
				case calendar.MonthKeyboard:
					requests = []tgbotapi.Chattable{
						tgbotapi.NewEditMessageText(chatID, messageID, "выберите месяц"),
						tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, info.Keyboard),
					}
				case calendar.DayKeyboard:
					requests = []tgbotapi.Chattable{
						tgbotapi.NewEditMessageText(chatID, messageID, "выберите день"),
						tgbotapi.NewEditMessageReplyMarkup(chatID, messageID, info.Keyboard),
					}
				case calendar.Finish:
					requests = []tgbotapi.Chattable{
						NewEditMessageDeleteInlineKeyboard(chatID, messageID),
						tgbotapi.NewEditMessageText(chatID, messageID, info.String()),
					}
					delete(calendarCache, calID)
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
