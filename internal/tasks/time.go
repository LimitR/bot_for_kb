package tasks

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetTimeButton(id string) []tgbotapi.InlineKeyboardButton {
	return []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("+10 м.", "taskTime_10_"+id),
		tgbotapi.NewInlineKeyboardButtonData("+30 м.", "taskTime_30_"+id),
		tgbotapi.NewInlineKeyboardButtonData("+1 ч.", "taskTime_60_"+id),
		tgbotapi.NewInlineKeyboardButtonData("+24 ч.", "taskTime_1440_"+id),
		tgbotapi.NewInlineKeyboardButtonData("Задать время", "taskTime_custom_"+id),
	}
}
