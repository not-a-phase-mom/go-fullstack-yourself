package routes

import (
	"bytes"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/not-a-phase-mom/go-fullstack-yourself/cmd/web/templates/pages"
	"github.com/not-a-phase-mom/go-fullstack-yourself/internal/database"

	"github.com/yuin/goldmark"
)

func (router *Router) RegisterArticleRoutes(e *gin.Engine) {
	e.GET("/articles/:slug", router.handleGetArticle)
	e.GET("/articles", router.handleGetAllArticles)
}

func (router *Router) handleGetArticle(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")
	slug := c.Param("slug")

	// Initialize user as nil
	var user *database.UserModel = nil

	// Get the token from the cookie
	token, err := c.Cookie("token")

	// If token exists, try to retrieve user
	if token != "" {
		userId, err := router.R.GetValue(token)
		if err == nil {
			// Try to fetch user, if successful update the user pointer
			result, err := router.Db.User.ById(userId)
			if err == nil {
				var userModel database.UserModel
				userModel.User = result
				user = &userModel
			}
		}
	}

	// Fetch the article
	article, err := router.Db.Article.BySlug(slug)
	if err != nil || article.Id == "" {
		// If article not found, render error page
		pages.ErrorPage(http.StatusNotFound, "Article not found").Render(c.Request.Context(), c.Writer)
		return
	}

	// Prepare user pointer for template
	var userPtr *database.User = nil
	if user != nil {
		userPtr = &user.User
	}

	source := []byte(article.Content)

	// parse the article content from markdown to HTML
	var buf bytes.Buffer
	if err := goldmark.Convert(source, &buf); err != nil {
		panic(err)
	}

	article.Content = buf.String()

	// Render the article page
	component := pages.ArticlePage(userPtr, &article)
	component.Render(c.Request.Context(), c.Writer)
}

func (router *Router) handleGetAllArticles(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")
	token, err := c.Cookie("token")

	// Initialize user as nil
	var user *database.UserModel = nil

	// If token exists, try to retrieve user
	if token != "" {
		userId, err := router.R.GetValue(token)
		if err == nil {
			// Try to fetch user, if successful update the user pointer
			result, err := router.Db.User.ById(userId)
			if err == nil {
				var userModel = database.UserModel{User: result}
				user = &userModel
			}
		}
	}

	// Fetch all articles
	articles, err := router.Db.Article.All("published")
	if err != nil {
		pages.ErrorPage(http.StatusInternalServerError, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	fmt.Printf("Articles: %v\n", articles[0].PublishedAt)

	// Render the articles page
	// Pass nil if user is nil
	var userPtr *database.User = nil
	if user != nil {
		userPtr = &user.User
	}
	component := pages.ArticlesPage(userPtr, articles)
	component.Render(c.Request.Context(), c.Writer)
}
