package team

import "github.com/redis/go-redis/v9"

type Team struct {
	rdb *redis.Client
}

func NewTeam(rdb *redis.Client) *Team {
	return &Team{
		rdb: rdb,
	}
}
