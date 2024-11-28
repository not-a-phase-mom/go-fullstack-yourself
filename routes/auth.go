package routes

import (
	"github.com/gin-gonic/gin"
)

func (router *Router) RegisterAuthRoutes(e *gin.Engine) {
	e.GET("/login", router.handleLogin)
	e.POST("/login", router.handleLoginPost)
	e.GET("/register", router.handleRegister)
	e.POST("/register", router.handleRegisterPost)
	e.GET("/logout", router.handleLogout)
}

func (router *Router) handleLogin(c *gin.Context) {
	// ...
}

func (router *Router) handleLoginPost(c *gin.Context) {
	// ...
}

func (router *Router) handleRegister(c *gin.Context) {
	// ...
}

func (router *Router) handleRegisterPost(c *gin.Context) {
	// ...
}

func (router *Router) handleLogout(c *gin.Context) {
	// ...
}
