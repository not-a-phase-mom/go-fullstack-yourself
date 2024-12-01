package routes

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/not-a-phase-mom/go-fullstack-yourself/cmd/web/templates/component"
	"github.com/not-a-phase-mom/go-fullstack-yourself/cmd/web/templates/pages"
	"github.com/not-a-phase-mom/go-fullstack-yourself/internal/database"
	"github.com/not-a-phase-mom/go-fullstack-yourself/internal/redis"
)

func (router *Router) RegisterAuthRoutes(e *gin.Engine) {
	e.GET("/login", router.handleLogin)
	e.POST("/login", router.handleLoginPost)
	e.GET("/register", router.handleRegister)
	e.POST("/register", router.handleRegisterPost)
	e.GET("/logout", router.handleLogout)
}

func (router *Router) handleLogin(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")

	// check if the user is already logged in
	token, err := c.Cookie("token")
	if err == nil {
		id, err := router.R.GetValue(token)
		if err != nil {
			pages.ErrorPage(http.StatusNotFound, err.Error()).Render(c.Request.Context(), c.Writer)
			return
		}

		user, err := router.Db.User.ById(id)
		if err != nil {
			pages.ErrorPage(http.StatusInternalServerError, err.Error()).Render(c.Request.Context(), c.Writer)
			return
		}

		if user.Id != "" {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
	}

	component := pages.LoginPage("")
	component.Render(c.Request.Context(), c.Writer)
}

func (router *Router) handleLoginPost(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")

	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := router.Db.User.ByEmail(email)
	if err != nil {
		c.Status(http.StatusUnprocessableEntity)
		if strings.Contains(err.Error(), "no rows in result set") {
			component := component.ErrorMessage("No user found with that email.")
			component.Render(c.Request.Context(), c.Writer)
			return
		} else {
			component := component.ErrorMessage(err.Error())
			component.Render(c.Request.Context(), c.Writer)
			return
		}
	}

	if user.Password != password {
		c.Status(http.StatusUnprocessableEntity)
		component := component.ErrorMessage("Invalid credentials.")
		component.Render(c.Request.Context(), c.Writer)
		return
	}

	// generate a random token for the user id
	token := redis.GenerateToken(user.Id)

	// store the token in redis
	err = router.R.SetValue(token, user.Id)
	if err != nil {
		pages.ErrorPage(http.StatusUnprocessableEntity, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	// send an html page index.html and pass a set-cookie header with the token
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.Writer.Header().Set("HX-Redirect", "/")
	c.Status(http.StatusOK)
	return
}

func (router *Router) handleRegister(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")

	// check if the user is already logged in
	token, err := c.Cookie("token")

	if err == nil {
		id, err := router.R.GetValue(token)
		if err != nil {
			pages.ErrorPage(http.StatusUnauthorized, err.Error()).Render(c.Request.Context(), c.Writer)
			return
		}

		user, err := router.Db.User.ById(id)
		if err != nil {
			pages.ErrorPage(http.StatusNotFound, err.Error()).Render(c.Request.Context(), c.Writer)
			return
		}

		if user.Id != "" {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
	}

	component := pages.RegisterPage("")
	component.Render(c.Request.Context(), c.Writer)
}

func (router *Router) handleRegisterPost(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/html")

	email := c.PostForm("email")
	name := c.PostForm("name")
	password := c.PostForm("password")

	user := database.User{
		Email:    email,
		Name:     name,
		Password: password,
	}

	id, err := router.Db.User.Create(user)
	if err != nil {
		// check the sql error code if it is 23505
		c.Status(http.StatusUnprocessableEntity)
		component.ErrorMessage(err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	// generate a random token for the user id
	token := redis.GenerateToken(id)

	// store the token in redis
	err = router.R.SetValue(token, id)
	if err != nil {
		pages.ErrorPage(http.StatusUnprocessableEntity, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	// send an html page index.html and pass a set-cookie header with the token
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	//convert the user id to a string
	c.Writer.Header().Set("HX-Redirect", "/")
	c.Status(http.StatusOK)
	return
}

func (router *Router) handleLogout(c *gin.Context) {

	c.Writer.Header().Set("Content-Type", "text/html")

	token, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	err = router.R.DeleteValue(token)
	if err != nil {
		pages.ErrorPage(http.StatusInternalServerError, err.Error()).Render(c.Request.Context(), c.Writer)
		return
	}

	c.SetCookie("token", "", -1, "/", "localhost", false, true)

	c.Redirect(http.StatusSeeOther, "/")
}
