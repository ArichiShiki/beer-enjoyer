package main

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// NewEditMessageDeleteInlineKeyboard allows you to edit the inline
// keyboard markup.
func NewEditMessageDeleteInlineKeyboard(chatID int64, messageID int) tgbotapi.EditMessageReplyMarkupConfig {
	return tgbotapi.EditMessageReplyMarkupConfig{
		BaseEdit: tgbotapi.BaseEdit{
			ChatID:    chatID,
			MessageID: messageID,
		},
	}
}
