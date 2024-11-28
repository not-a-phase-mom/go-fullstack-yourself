package routes

import (
	"github.com/not-a-phase-mom/go-fullstack-yourself/services/database"
	"github.com/not-a-phase-mom/go-fullstack-yourself/services/redis"
)

type Router struct {
	Db *database.Database
	R  *redis.Redis
}
