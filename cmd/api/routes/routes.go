package routes

import (
	"github.com/not-a-phase-mom/go-fullstack-yourself/internal/database"
	"github.com/not-a-phase-mom/go-fullstack-yourself/internal/redis"
)

type Router struct {
	Db *database.Database
	R  *redis.Redis
}
