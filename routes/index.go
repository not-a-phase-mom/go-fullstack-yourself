package routes

import (
	"github.com/gin-gonic/gin"
)

func (router *Router) RegisterIndexRoutes(e *gin.Engine) {
	e.GET("/", router.handleIndex)
}

func (router *Router) handleIndex(c *gin.Context) {
	token, err := c.Cookie("token")
	//...
}
