package routes

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/not-a-phase-mom/go-fullstack-yourself/cmd/web/templates/pages"
	"github.com/not-a-phase-mom/go-fullstack-yourself/internal/database"
)

func (router *Router) RegisterDashboardRoutes(e *gin.Engine) {
	e.GET("/dashboard", router.handleDashboard)
	e.GET("/dashboard/articles", router.handleDashboardArticles)
	e.GET("/dashboard/articles/new", router.handleNewArticle)
	e.POST("/dashboard/articles/new", router.handleCreateArticle)
	e.GET("/dashboard/articles/edit/:id", router.handleEditArticle)
	e.POST("/dashboard/articles/edit/:id", router.handleUpdateArticle)
	e.POST("/dashboard/articles/publish/:id", router.handlePublishArticle)
	e.POST("/dashboard/articles/unpublish/:id", router.handleUnPublishArticle)
}

func (router *Router) handleDashboard(c *gin.Context) {
	// ...authentication logic...
	c.Writer.Header().Set("Content-Type", "text/html")
	token, err := c.Cookie("token")
	if err != nil {
		pages.ErrorPage(http.StatusNotFound, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	user, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		pages.ErrorPage(errorObject.Status, errorObject.Error).Render(c.Request.Context(), c.Writer)
		return
	}

	component := pages.DashboardPage(&user)
	component.Render(c.Request.Context(), c.Writer)

}

// Handler to view all articles (published and drafts)
func (router *Router) handleDashboardArticles(c *gin.Context) {
	// ...authentication logic...
	c.Writer.Header().Set("Content-Type", "text/html")
	token, err := c.Cookie("token")
	if err != nil {
		pages.ErrorPage(http.StatusNotFound, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	user, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		pages.ErrorPage(errorObject.Status, errorObject.Error).Render(c.Request.Context(), c.Writer)
		return
	}

	articles, err := router.Db.Article.All("")
	if err != nil {
		// Handle error
		return
	}
	// Render dashboard articles page
	pages.DashboardArticlesPage(&user, articles).Render(c.Request.Context(), c.Writer)
}

// Handler to display new article form
func (router *Router) handleNewArticle(c *gin.Context) {
	// ...authentication logic...
	c.Writer.Header().Set("Content-Type", "text/html")
	token, err := c.Cookie("token")
	if err != nil {
		pages.ErrorPage(http.StatusNotFound, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	user, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		pages.ErrorPage(errorObject.Status, errorObject.Error).Render(c.Request.Context(), c.Writer)
		return
	}

	pages.NewArticlePage(&user).Render(c.Request.Context(), c.Writer)
}

// Handler to create a new article
func (router *Router) handleCreateArticle(c *gin.Context) {
	// ...authentication logic...
	c.Writer.Header().Set("Content-Type", "text/html")
	token, err := c.Cookie("token")
	if err != nil {
		pages.ErrorPage(http.StatusNotFound, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	_, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		pages.ErrorPage(errorObject.Status, errorObject.Error).Render(c.Request.Context(), c.Writer)
		return
	}

	var article database.ArticleCreation
	if err := c.ShouldBind(&article); err != nil {
		// Handle error
		pages.ErrorPage(http.StatusUnprocessableEntity, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	tagsInput := c.PostForm("tags")
	tagNames := strings.Split(tagsInput, ",")
	var tagCreations []database.TagCreation
	for _, tagName := range tagNames {
		tagCreations = append(tagCreations, database.TagCreation{Name: strings.TrimSpace(tagName)})
	}
	_, errr := router.Db.Article.Create(article, tagCreations)
	if errr != nil {
		// Handle error
		pages.ErrorPage(http.StatusUnprocessableEntity, errr.Error()).Render(c.Request.Context(), c.Writer)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dashboard/articles")
}

// Handler to display edit article form
func (router *Router) handleEditArticle(c *gin.Context) {
	// ...authentication logic...
	c.Writer.Header().Set("Content-Type", "text/html")
	token, err := c.Cookie("token")
	if err != nil {
		pages.ErrorPage(http.StatusNotFound, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	user, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		pages.ErrorPage(errorObject.Status, errorObject.Error).Render(c.Request.Context(), c.Writer)
		return
	}

	articleId := c.Param("id")
	article, err := router.Db.Article.ById(articleId)
	if err != nil {
		// Handle error
		return
	}
	pages.EditArticlePage(&user, article).Render(c.Request.Context(), c.Writer)
}

// Handler to update an article
func (router *Router) handleUpdateArticle(c *gin.Context) {
	// ...authentication logic...
	c.Writer.Header().Set("Content-Type", "text/html")
	id := c.Param("id")

	token, err := c.Cookie("token")
	if err != nil {
		pages.ErrorPage(http.StatusNotFound, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	_, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		pages.ErrorPage(errorObject.Status, errorObject.Error).Render(c.Request.Context(), c.Writer)
		return
	}

	var articleData database.ArticleCreation
	if err := c.ShouldBind(&articleData); err != nil {
		// Handle error
		pages.ErrorPage(http.StatusUnprocessableEntity, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	fmt.Println(articleData)

	status := c.PostForm("status")
	errr := router.Db.Article.Update(id, articleData, status)
	if errr != nil {
		// Handle error
		pages.ErrorPage(http.StatusUnprocessableEntity, errr.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	tagsInput := c.PostForm("tags")
	tagNames := strings.Split(tagsInput, ",")
	var tagCreations []database.TagCreation
	for _, tagName := range tagNames {
		tagCreations = append(tagCreations, database.TagCreation{Name: strings.TrimSpace(tagName)})
	}
	err = router.Db.Article.UpdateTags(id, tagCreations)
	if err != nil {
		// Handle error
		pages.ErrorPage(http.StatusUnprocessableEntity, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}
	c.Redirect(http.StatusSeeOther, "/dashboard/articles")
}

// Handler to publish an article
func (router *Router) handlePublishArticle(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")
	id := c.Param("id")

	token, err := c.Cookie("token")
	if err != nil {
		pages.ErrorPage(http.StatusNotFound, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	_, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		pages.ErrorPage(errorObject.Status, errorObject.Error).Render(c.Request.Context(), c.Writer)
		return
	}

	updatedArticle, err := router.Db.Article.Publish(id)
	if err != nil {
		// Handle error
		pages.ErrorPage(http.StatusUnprocessableEntity, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	component := pages.DashboardArticlesRow(updatedArticle)
	c.Writer.Header().Set("HX-Retarget", "#article-"+updatedArticle.Id)
	component.Render(c.Request.Context(), c.Writer)
}

// Handler to unpublish an article
func (router *Router) handleUnPublishArticle(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")
	id := c.Param("id")

	token, err := c.Cookie("token")
	if err != nil {
		pages.ErrorPage(http.StatusNotFound, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	_, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		pages.ErrorPage(errorObject.Status, errorObject.Error).Render(c.Request.Context(), c.Writer)
		return
	}

	updatedArticle, err := router.Db.Article.UnPublish(id)
	if err != nil {
		// Handle error
		pages.ErrorPage(http.StatusUnprocessableEntity, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	component := pages.DashboardArticlesRow(updatedArticle)
	c.Writer.Header().Set("HX-Retarget", "#article-"+updatedArticle.Id)
	component.Render(c.Request.Context(), c.Writer)
}
