package tasks

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetComplexityButton(id string) []tgbotapi.InlineKeyboardButton {
	return []tgbotapi.InlineKeyboardButton{
		tgbotapi.NewInlineKeyboardButtonData("Незначительно", "taskComplexity_1_"+id),
		tgbotapi.NewInlineKeyboardButtonData("Важно", "taskComplexity_2_"+id),
		tgbotapi.NewInlineKeyboardButtonData("Срочно", "taskComplexity_3_"+id),
	}
}
