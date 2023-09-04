package core

import (
	"context"
	"go_crm_bot/internal/cache"
	"go_crm_bot/internal/tasks"
	"go_crm_bot/internal/team"
	"go_crm_bot/internal/users"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Parser struct {
	flag  *Flag
	team  *team.Team
	users *users.User
	cache *cache.Cache
	task  *tasks.Task
	tg    *tgbotapi.BotAPI
}

func NewParser(flag *Flag, team *team.Team, users *users.User, cache *cache.Cache, task *tasks.Task, tg *tgbotapi.BotAPI) *Parser {
	return &Parser{
		flag:  flag,
		team:  team,
		users: users,
		cache: cache,
		task:  task,
		tg:    tg,
	}
}

func (p *Parser) Parse(listener tgbotapi.UpdatesChannel) {
	for event := range listener {
		if event.Message != nil {
			switch p.flag.GetFlag(context.Background(), event.Message.Chat.ID) {
			case CreateTeam:
				p.createTeamParse(event)
				break
			case CreateUser:
				user := p.users.GetUser(event.Message.Chat.ID)
				if user.Name == "" {
					data := strings.Split(event.Message.Text, " ")
					p.users.CreateUser(event.Message.Chat.ID, event.Message.From.UserName, data[0], data[1])
					p.flag.SetFlag(context.Background(), event.Message.Chat.ID, Default)
					msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Аккаунт успешно создан!")
					p.tg.Send(msg)
				} else {
					p.flag.SetFlag(context.Background(), event.Message.Chat.ID, Default)
					msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Вы уже создали аккаунт")
					p.tg.Send(msg)
				}
				break
			case CreateTaskName:
				userStringId := strconv.Itoa(int(event.Message.Chat.ID))
				uuidTask := p.cache.Get(userStringId, "taskCreate")
				msgTaskId := p.cache.Get(userStringId, "msgId")
				task := p.task.Static_GetTaskById(uuidTask)
				task.NameTask = event.Message.Text
				task.Save()
				msgIdInt, _ := strconv.Atoi(msgTaskId)
				user := p.users.GetUser(task.ExecutorId)
				master := p.users.GetUser(task.Master)

				msg := tgbotapi.NewEditMessageText(event.Message.Chat.ID, msgIdInt, task.CreateMsgTask(user.Login, user.Name, master.Name))
				msg.ParseMode = "Markdown"
				msgRM := tgbotapi.NewEditMessageReplyMarkup(event.Message.Chat.ID, msgIdInt, tasks.GetReplyMarkup(uuidTask))

				p.tg.Send(msg)
				p.tg.Send(msgRM)

				p.flag.SetFlag(context.Background(), event.Message.Chat.ID, Default)
			case CreateTaskDescription:
				userStringId := strconv.Itoa(int(event.Message.Chat.ID))
				uuidTask := p.cache.Get(userStringId, "taskCreate")
				msgTaskId := p.cache.Get(userStringId, "msgId")
				task := p.task.Static_GetTaskById(uuidTask)
				task.Description = event.Message.Text
				task.Save()
				msgIdInt, _ := strconv.Atoi(msgTaskId)
				user := p.users.GetUser(task.ExecutorId)
				master := p.users.GetUser(task.Master)

				msg := tgbotapi.NewEditMessageText(event.Message.Chat.ID, msgIdInt, task.CreateMsgTask(user.Login, user.Name, master.Name))
				msg.ParseMode = "Markdown"
				msgRM := tgbotapi.NewEditMessageReplyMarkup(event.Message.Chat.ID, msgIdInt, tasks.GetReplyMarkup(uuidTask))

				p.tg.Send(msg)
				p.tg.Send(msgRM)

				p.flag.SetFlag(context.Background(), event.Message.Chat.ID, Default)
			default:
				switch event.Message.Text {
				case "/start":
					msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Введите ваше имя и должность (Пример: 'Иван Кассир')")
					p.flag.SetFlag(context.Background(), event.Message.Chat.ID, CreateUser)
					p.tg.Send(msg)
				case "/create_team":
					m := p.team.CreateTeam(event.Message.Chat.ID)
					p.tg.Send(m)
					p.flag.SetFlag(context.Background(), event.Message.Chat.ID, CreateTeam)
					break
				case "/get_team":
					msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Вся ваша команда:")
					users := p.team.GetUsersTeam(event.Message.Chat.ID)
					var buttons [][]tgbotapi.InlineKeyboardButton
					for _, user := range users {
						us := strings.Split(user, "_")
						buttons = append(buttons, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData(us[1], "users_list_"+us[0])})
					}
					msg.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttons...)
					p.tg.Send(msg)
					break
				}
			}
		}
		if event.CallbackQuery != nil {
			p.parseCallbackQuery(event.CallbackQuery.Data, event.CallbackQuery.From.ID, event.CallbackQuery.Message.MessageID)
		}
	}
}

func (p *Parser) parseCallbackQuery(query string, userId int64, msgId int) {
	if strings.HasPrefix(query, "users_list_") {
		id_string := strings.Split(query, "users_list_")[1]
		id, _ := strconv.Atoi(id_string)
		user := p.users.GetUser(int64(id))
		ReplyMarkup := tgbotapi.NewInlineKeyboardMarkup(
			[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Создать новую задачу", "usersCreateTasks_"+id_string)},
			[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Список задач", "usersTasks_"+id_string)},
			[]tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Подать сигнал", "ping_"+id_string)},
		)
		if user.Name != "" {
			msg_text := "Пользователя зовут: " + user.Name + "\nДолжность: " + user.Position
			msg := tgbotapi.NewMessage(userId, msg_text)
			msg.ReplyMarkup = ReplyMarkup
			p.tg.Send(msg)
		} else {
			chat, _ := p.tg.GetChat(tgbotapi.ChatInfoConfig{ChatConfig: tgbotapi.ChatConfig{ChatID: int64(id)}})
			msg := tgbotapi.NewMessage(userId, "Данных о пользователе '"+chat.UserName+"' пока нет, возможно он ещё не создал аккаунт")
			msg.ReplyMarkup = ReplyMarkup
			p.tg.Send(msg)
		}
	}
	if strings.HasPrefix(query, "taskNameCreate_") {
		taskId := strings.Split(query, "taskNameCreate_")[1]
		userStringId := strconv.Itoa(int(userId))
		msgStringId := strconv.Itoa(msgId)

		p.cache.Set(userStringId, "taskCreate", taskId)
		p.cache.Set(userStringId, "msgId", msgStringId)

		p.flag.SetFlag(context.Background(), userId, CreateTaskName)

		msg := tgbotapi.NewMessage(userId, "Введите назвние задачи:")
		p.tg.Send(msg)
	}
	if strings.HasPrefix(query, "taskNameDescription_") {
		taskId := strings.Split(query, "taskNameDescription_")[1]
		userStringId := strconv.Itoa(int(userId))
		msgStringId := strconv.Itoa(msgId)

		p.cache.Set(userStringId, "taskCreate", taskId)
		p.cache.Set(userStringId, "msgId", msgStringId)

		p.flag.SetFlag(context.Background(), userId, CreateTaskDescription)

		msg := tgbotapi.NewMessage(userId, "Введите описание задачи:")
		p.tg.Send(msg)
	}
	if strings.HasPrefix(query, "usersCreateTasks_") {
		id_string := strings.Split(query, "usersCreateTasks_")[1]
		id, _ := strconv.Atoi(id_string)
		user := p.users.GetUser(int64(id))
		master := p.users.GetUser(userId)
		blTask := tasks.NewBlankTask(p.flag.rdb)
		task := blTask.Static_CreateTask(userId, int64(id), "***", "***", "")
		msg := tgbotapi.NewMessage(userId, task.CreateMsgTask(user.Login, user.Name, master.Name))
		msg.ParseMode = "Markdown"
		msg.ReplyMarkup = tasks.GetReplyMarkup(task.Uuid)

		p.tg.Send(msg)
	}
}

func (p *Parser) createTeamParse(event tgbotapi.Update) {
	switch event.Message.Text {
	case "Сохранить":
		newMessage := tgbotapi.NewMessage(event.Message.From.ID, "Список команды успешно сформирован!")
		p.tg.Send(newMessage)
		p.flag.SetFlag(context.Background(), event.Message.From.ID, Default)
		break
	case "Начать с начала":
		p.team.DeleteTeam(event.Message.From.ID)
	default:
		if event.Message.ForwardFrom != nil {
			p.team.AddUserTeam(event.Message.From.ID, event.Message.ForwardFrom.ID, event.Message.ForwardFrom.UserName)
		} else {
			newMessage := tgbotapi.NewMessage(event.Message.From.ID, "Пожалуйста, присылайте только сообщение от нужных пользователей.")
			p.tg.Send(newMessage)
		}
	}
}

func (p *Parser) parseCommandButton(command *tgbotapi.CallbackQuery) {}
