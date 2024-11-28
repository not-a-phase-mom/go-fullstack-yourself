package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (router *Router) RegisterIndexRoutes(e *gin.Engine) {
	e.GET("/", router.handleIndex)
}

func (router *Router) handleIndex(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"user": nil,
		})
		return
	}

	id, err := router.R.GetValue(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	user, err := router.Db.GetUserById(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.HTML(http.StatusOK, "index.html", gin.H{
		"user": user,
	})
}
