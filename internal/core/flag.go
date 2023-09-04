package core

import (
	"context"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	CreateTeam            string = "create_team"
	CreateUser            string = "create_user"
	CreateTask            string = "create_task"
	CreateTaskName        string = "create_task_name"
	CreateTaskDescription string = "create_task_description"
	Default               string = "none"
)

type Flag struct {
	rdb *redis.Client
}

func NewFlag(rdb *redis.Client) *Flag {
	return &Flag{
		rdb: rdb,
	}
}

func (f *Flag) SetFlag(ctx context.Context, id int64, flag string) {
	f.rdb.Set(ctx, "flag_"+strconv.Itoa(int(id)), flag, 5*time.Minute)
}

func (f *Flag) GetFlag(ctx context.Context, id int64) string {
	result, _ := f.rdb.Get(ctx, "flag_"+strconv.Itoa(int(id))).Result()
	return result
}
