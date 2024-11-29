package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/not-a-phase-mom/go-fullstack-yourself/templates/pages"
)

func (router *Router) RegisterIndexRoutes(e *gin.Engine) {
	e.GET("/", router.handleIndex)
}

func (router *Router) handleIndex(c *gin.Context) {
	token, err := c.Cookie("token")
	emptyUserIndex := pages.IndexPage(nil)

	c.Writer.Header().Set("Content-Type", "text/html")

	if err != nil {
		emptyUserIndex.Render(c.Request.Context(), c.Writer)
		return
	}

	id, err := router.R.GetValue(token)
	if err != nil {
		pages.ErrorPage(err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	user, err := router.Db.GetUserById(id)
	if err != nil {
		pages.ErrorPage(err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	component := pages.IndexPage(&user)
	component.Render(c.Request.Context(), c.Writer)
}
