package users

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/redis/go-redis/v9"
)

type User struct {
	rdb      *redis.Client
	Name     string `json:"name"`
	Id       int64  `json:"id"`
	Login    string `json:"login"`
	Position string `json:"position"`
}

func NewUser(rdb *redis.Client) *User {
	return &User{
		rdb: rdb,
	}
}

func (u *User) CreateUser(id int64, login, name, position string) {
	data, _ := json.Marshal(User{
		Name:     name,
		Position: position,
		Login:    login,
		Id:       id,
	})
	u.rdb.Set(context.Background(), strconv.Itoa(int(id)), string(data), 0)
}

func (u *User) GetUser(id int64) User {
	var res User
	data, _ := u.rdb.Get(context.Background(), strconv.Itoa(int(id))).Result()
	if data != "" {
		json.Unmarshal([]byte(data), &res)
		return res
	}
	return res
}
