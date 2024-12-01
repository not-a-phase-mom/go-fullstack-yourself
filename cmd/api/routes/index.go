package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/not-a-phase-mom/go-fullstack-yourself/cmd/web/templates/pages"
)

func (router *Router) RegisterIndexRoutes(e *gin.Engine) {
	e.GET("/", router.handleIndex)
}

func (router *Router) handleIndex(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")
	token, err := c.Cookie("token")

	articles, err := router.Db.Article.All()
	if err != nil {
		fmt.Println("Article fetching: " + err.Error())
		pages.ErrorPage(http.StatusInternalServerError, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	id, err := router.R.GetValue(token)

	if id == "" {
		pages.IndexPage(nil, articles, articles).Render(c.Request.Context(), c.Writer)
		return
	}

	user, err := router.Db.User.ById(id)
	if err != nil {
		fmt.Println("User fetching: " + err.Error())
		pages.ErrorPage(http.StatusUnauthorized, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	fmt.Printf("Articles: %v\n", articles)
	component := pages.IndexPage(&user, articles, articles)
	component.Render(c.Request.Context(), c.Writer)
}
