package tasks

import (
	"context"
	"encoding/json"
	"go_crm_bot/pkg/utils"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

const timeFormat = "2006-01-02 15:04:05"

type Task struct {
	rdb         *redis.Client
	Uuid        string `json:"uuid"`
	Done        bool   `json:"done"`       // Выполнена ли задача
	ExecutorId  int64  `json:"executorId"` // На кого стоит задача
	Urgency     string `json:"urgency"`
	NameTask    string `json:"nameTask"`    // Название задачи
	Description string `json:"description"` // Описание задачи
	Execution   string `json:"execution"`
	Master      int64  `json:"master"` // Кто поставил задачу
	Result      string `json:"result"` // Итог по задаче
}

func NewBlankTask(rdb *redis.Client) Task {
	return Task{
		rdb: rdb,
	}
}

func (t *Task) Static_CreateTask(master int64, executorId int64, nameTask string, description string, execution string) Task {
	_id, _ := uuid.NewUUID()
	id := _id.String()
	task := Task{
		Uuid:        id,
		Done:        false,
		ExecutorId:  executorId,
		Urgency:     time.Now().Format(timeFormat),
		NameTask:    nameTask,
		Description: description,
		Execution:   execution,
		Master:      master,
	}
	taskJson, _ := json.Marshal(task)
	t.rdb.Set(context.Background(), "task_"+id, string(taskJson), 0)
	go func() {
		t.rdb.LPush(context.Background(), "myTask_"+strconv.Itoa(int(executorId)), id, 0)
		t.rdb.LPush(context.Background(), "taskWorkers_"+strconv.Itoa(int(master)), id, 0)
	}()
	return task
}

func (t *Task) Static_GetAllTasks() []Task {
	tasks := make([]Task, 0, 10)
	keys, _, _ := t.rdb.Scan(context.Background(), 0, "task_*", 0).Result()
	utils.ForEach(keys, func(_ int, element string) {
		var task Task
		json.Unmarshal([]byte(element), &task)
		task.rdb = t.rdb
		tasks = append(tasks, task)
	})
	return tasks
}

func (t *Task) Static_GetTaskById(id string) Task {
	var task Task
	element, _ := t.rdb.Get(context.Background(), "task_"+id).Result()
	json.Unmarshal([]byte(element), &task)
	task.rdb = t.rdb
	return task
}
func (t *Task) Static_GetAllTaskByUserId(id string) {}

func (t *Task) Save() {
	taskJson, _ := json.Marshal(t)
	if t.Done {
		t.rdb.Set(context.Background(), "task_"+t.Uuid, string(taskJson), 24*time.Hour)
	} else {
		t.rdb.Set(context.Background(), "task_"+t.Uuid, string(taskJson), 0)
	}
}

func (t *Task) DeleteTask(uuid string) {}

func (t *Task) CreateMsgTask(loginUser, name, nameMaster string) string {
	msg := "*Название задачи:* " + t.NameTask + "\n"
	msg += "*Описание задачи:* " + t.Description + "\n"
	msg += "*Поставил задачу:* " + nameMaster + "\n"
	msg += "*Исполнитель:* [" + name + "](https://t.me/" + loginUser + ")" + "\n"
	msg += "*Приоритет задачи:* " + t.Execution + "\n"
	msg += "*Срок выполнения:* " + t.Urgency + "\n"
	return msg
}
