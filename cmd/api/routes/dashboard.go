package routes

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/not-a-phase-mom/go-fullstack-yourself/cmd/web/templates/component"
	"github.com/not-a-phase-mom/go-fullstack-yourself/cmd/web/templates/pages"
	"github.com/not-a-phase-mom/go-fullstack-yourself/internal/database"
	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/parser"
)

func (router *Router) RegisterDashboardRoutes(e *gin.Engine) {
	e.GET("/dashboard", router.handleDashboard)
	e.GET("/dashboard/articles", router.handleDashboardArticles)
	e.GET("/dashboard/articles/new", router.handleNewArticle)
	e.POST("/dashboard/articles/new", router.handleCreateArticle)
	e.POST("/dashboard/articles/upload", router.handleUploadArticle)
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
		c.Redirect(http.StatusSeeOther, "/login")
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
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	user, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		c.Redirect(http.StatusSeeOther, "/login")
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
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	user, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	emtpyArticle := database.ArticleCreation{}

	pages.NewArticlePage(&user, &emtpyArticle, "").Render(c.Request.Context(), c.Writer)
}

// Handler to create a new article
func (router *Router) handleCreateArticle(c *gin.Context) {
	// ...authentication logic...
	c.Writer.Header().Set("Content-Type", "text/html")
	token, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	_, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		c.Redirect(http.StatusSeeOther, "/login")
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
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	user, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		c.Redirect(http.StatusSeeOther, "/login")
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
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	_, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	var articleData database.ArticleCreation
	if err := c.ShouldBind(&articleData); err != nil {
		// Handle error
		pages.ErrorPage(http.StatusUnprocessableEntity, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

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

func (router *Router) handleUploadArticle(c *gin.Context) {
    c.Writer.Header().Set("Content-Type", "text/html")

    token, err := c.Cookie("token")
    if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
        return
    }

    user, errorObject := router.AuthenticateUser(token, true)
    if errorObject.Status != 0 {
        c.Redirect(http.StatusSeeOther, "/login")
        return
    }

    // Initialize the article data
    var articleData database.ArticleCreation

    // Set up Goldmark with metadata extension
    markdown := goldmark.New(
        goldmark.WithExtensions(
            meta.Meta,
        ),
    )

    // Retrieve the uploaded file from the form
    file, err := c.FormFile("file-upload")
    if err != nil {
        c.Status(http.StatusUnprocessableEntity)
        component := component.ErrorMessage("Error uploading file: " + err.Error())
        component.Render(c.Request.Context(), c.Writer)
        return
    }

    // Open the uploaded file
    fileContent, err := file.Open()
    if err != nil {
        c.Status(http.StatusUnprocessableEntity)
        component := component.ErrorMessage("Error opening file: " + err.Error())
        component.Render(c.Request.Context(), c.Writer)
        return
    }
    defer fileContent.Close()

    // Read the content of the file into a buffer
    buf := new(bytes.Buffer)
    _, err = buf.ReadFrom(fileContent)
    if err != nil {
        c.Status(http.StatusUnprocessableEntity)
        component := component.ErrorMessage("Error reading file: " + err.Error())
        component.Render(c.Request.Context(), c.Writer)
        return
    }

    // Parse the markdown content with metadata
    ctx := parser.NewContext()
    err = markdown.Convert(buf.Bytes(), io.Discard, parser.WithContext(ctx))
    if err != nil {
        c.Status(http.StatusUnprocessableEntity)
        component := component.ErrorMessage("Error parsing markdown: " + err.Error())
        component.Render(c.Request.Context(), c.Writer)
        return
    }

    // Extract metadata
    metaData := meta.Get(ctx)
	fmt.Printf("Metadata: %+v\n", metaData)

    // Initialize tag creations slice
    var tagCreations []database.TagCreation

    // Safely extract and process tags
    if tagsInterface, ok := metaData["tags"]; ok {
        switch tags := tagsInterface.(type) {
        case []interface{}:
            // Convert []interface{} to []string
            for _, tagInterface := range tags {
                if tagStr, ok := tagInterface.(string); ok {
                    tagCreations = append(tagCreations, database.TagCreation{Name: strings.TrimSpace(tagStr)})
                }
            }
        case []string:
            // Directly use []string
            for _, tag := range tags {
                tagCreations = append(tagCreations, database.TagCreation{Name: strings.TrimSpace(tag)})
            }
        default:
            // Handle unexpected type
            c.Status(http.StatusBadRequest)
            component := component.ErrorMessage("Error: 'tags' metadata should be a list of strings.")
            component.Render(c.Request.Context(), c.Writer)
            return
        }
    }

    // Safely extract 'title' from metadata
    if titleInterface, ok := metaData["title"]; ok {
        if title, ok := titleInterface.(string); ok {
            articleData.Title = title
        } else {
            c.Status(http.StatusBadRequest)
            component := component.ErrorMessage("Error: 'title' metadata should be a string.")
            component.Render(c.Request.Context(), c.Writer)
            return
        }
    } else {
        c.Status(http.StatusBadRequest)
        component := component.ErrorMessage("Error: 'title' metadata is required.")
        component.Render(c.Request.Context(), c.Writer)
        return
    }

    // Safely extract 'excerpt' from metadata
    if excerptInterface, ok := metaData["excerpt"]; ok {
        if excerpt, ok := excerptInterface.(string); ok {
            articleData.Excerpt = excerpt
        } else {
            c.Status(http.StatusBadRequest)
            component := component.ErrorMessage("Error: 'excerpt' metadata should be a string.")
            component.Render(c.Request.Context(), c.Writer)
            return
        }
    } else {
        c.Status(http.StatusBadRequest)
        component := component.ErrorMessage("Error: 'excerpt' metadata is required.")
        component.Render(c.Request.Context(), c.Writer)
        return
    }

    // Extract content without the front matter
    content := strings.TrimSpace(string(buf.Bytes()))
    if strings.HasPrefix(content, "---") {
        // Remove the front matter from content
        parts := strings.SplitN(content, "---", 3)
        if len(parts) == 3 {
            content = parts[2]
        }
    }
    articleData.Content = content
    articleData.Tags = tagCreations
    articleData.AuthorId = user.Id

	fmt.Printf("Article Data: %+v\n", articleData)

    // Set HTMX headers to swap the new article form
    c.Header("HX-Retarget", "body")
    c.Header("HX-Reswap", "innerHTML")

    // Render the new article page with the filled form
    pages.NewArticlePage(&user, &articleData, "").Render(c.Request.Context(), c.Writer)
}

// Handler to publish an article
func (router *Router) handlePublishArticle(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")
	id := c.Param("id")

	token, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	_, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		c.Redirect(http.StatusSeeOther, "/login")
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
		c.Redirect(http.StatusSeeOther, "/login")
		return
	}

	_, errorObject := router.AuthenticateUser(token, true)
	if errorObject.Status != 0 {
		c.Redirect(http.StatusSeeOther, "/login")
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
