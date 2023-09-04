package team

import (
	"context"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (t *Team) CreateTeam(id int64) tgbotapi.MessageConfig {
	m := tgbotapi.NewMessage(id, "Перешлите сообщения от пользователей, которых вы хотите видеть в своей команде")
	m.ReplyMarkup = tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton("Сохранить"),
			tgbotapi.NewKeyboardButton("Начать с начала"),
		),
	)
	return m
}

func (t *Team) AddUserTeam(idMaster, idUser int64, nameUser string) {
	elements := t.GetUsersTeam(idMaster)
	if !findElement(elements, strconv.Itoa(int(idUser))+"_"+nameUser) {
		t.rdb.LPush(context.Background(), "team_"+strconv.Itoa(int(idMaster)), strconv.Itoa(int(idUser))+"_"+nameUser).Result()
	}
}

func (t *Team) GetUsersTeam(idMaster int64) []string {
	result, _ := t.rdb.LRange(context.Background(), "team_"+strconv.Itoa(int(idMaster)), 0, -1).Result()
	return result
}

func (t *Team) DeleteTeam(idMaster int64) {
	t.rdb.Del(context.Background(), "team_"+strconv.Itoa(int(idMaster)))
}

func findElement(elements []string, element string) bool {
	res := false
	for _, el := range elements {
		if el == element {
			res = true
		}
	}
	return res
}
