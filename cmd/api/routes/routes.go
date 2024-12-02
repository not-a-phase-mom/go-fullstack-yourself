package routes

import (
	"fmt"
	"net/http"

	"github.com/not-a-phase-mom/go-fullstack-yourself/internal/database"
	"github.com/not-a-phase-mom/go-fullstack-yourself/internal/redis"
)

type Router struct {
	Db *database.Database
	R  *redis.Redis
}


type ErrorObject struct {
	Status int
	Error string
}

func (r *Router) AuthenticateUser(token string, shouldBeAdmin bool) (database.User, ErrorObject) {
	id, err := r.R.GetValue(token)
	if id == "" {
		return database.User{}, ErrorObject{http.StatusUnauthorized, "You are not authorized to view this page"}
	}

	user, err := r.Db.User.ById(id)
	if err != nil {
		fmt.Println("User fetching: " + err.Error())
		return database.User{}, ErrorObject{http.StatusUnauthorized, err.Error()}
	}

	if user.Role != "admin" && shouldBeAdmin {
		return database.User{}, ErrorObject{http.StatusUnauthorized, "You are not authorized to view this page"}
	}

	return user, ErrorObject{}
}