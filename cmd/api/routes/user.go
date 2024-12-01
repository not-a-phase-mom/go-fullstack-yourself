package routes

import (
	"github.com/gin-gonic/gin"
)

func (router *Router) RegisterUserRoutes(e *gin.Engine) {
	e.GET("/users/:id", router.handleGetUser)
	e.PUT("/users/:id", router.handleUpdateUser)
}

func (router *Router) handleGetUser(c *gin.Context) {
	// ...
}

func (router *Router) handleUpdateUser(c *gin.Context) {
	// ...
}
