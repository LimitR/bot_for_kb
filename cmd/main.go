package main

import (
	"context"
	"go_crm_bot/internal/cache"
	"go_crm_bot/internal/core"
	"go_crm_bot/internal/tasks"
	"go_crm_bot/internal/team"
	"go_crm_bot/internal/users"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

func main() {
	secret, _ := godotenv.Read(".env")

	bot, err := tgbotapi.NewBotAPI(secret["TOKEN_BOT"])
	if err != nil {
		log.Panic(err)
	}

	bot.Request(tgbotapi.NewSetMyCommands(
		tgbotapi.BotCommand{
			Command:     "/get_team",
			Description: "Список команды",
		},
		tgbotapi.BotCommand{
			Command:     "/create_team",
			Description: "Создать команду",
		},
	))

	bot.Debug = false

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	update := bot.GetUpdatesChan(u)

	rdb := redis.NewClient(&redis.Options{
		Addr:     secret["REDIS_ADDR"],
		Password: "",
		DB:       0,
	})
	cache := cache.NewCache(rdb)
	task := tasks.NewBlankTask(rdb)
	users := users.NewUser(rdb)
	t := team.NewTeam(rdb)
	flag := core.NewFlag(rdb)
	parser := core.NewParser(flag, t, users, cache, &task, bot)

	parser.Parse(update)
}
