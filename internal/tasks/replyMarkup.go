package tasks

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetReplyMarkup(id string) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		GetComplexityButton(id),
		GetTimeButton(id),
		[]tgbotapi.InlineKeyboardButton{
			tgbotapi.NewInlineKeyboardButtonData("Задать название", "taskNameCreate_"+id),
			tgbotapi.NewInlineKeyboardButtonData("Задать описание", "taskNameDescription_"+id),
		},
		[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Сохранить", "ping_wad")},
	)
}
