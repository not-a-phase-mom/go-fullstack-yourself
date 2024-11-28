package routes

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/not-a-phase-mom/go-fullstack-yourself/services/database"
	"github.com/not-a-phase-mom/go-fullstack-yourself/services/redis"
)

func (router *Router) RegisterAuthRoutes(e *gin.Engine) {
	e.GET("/login", router.handleLogin)
	e.POST("/login", router.handleLoginPost)
	e.GET("/register", router.handleRegister)
	e.POST("/register", router.handleRegisterPost)
	e.GET("/logout", router.handleLogout)
}

func (router *Router) handleLogin(c *gin.Context) {
	// check if the user is already logged in
	token, err := c.Cookie("token")
	if err == nil {
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

		if user.Id != "" {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
	}

	c.HTML(http.StatusOK, "login.html", nil)
}

func (router *Router) handleLoginPost(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	user, err := router.Db.GetUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if user.Password != password {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "invalid credentials",
		})
		return
	}

	// generate a random token for the user id
	token := redis.GenerateToken(user.Id)

	// store the token in redis
	err = router.R.SetValue(token, user.Id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// send an html page index.html and pass a set-cookie header with the token
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	c.Redirect(http.StatusSeeOther, "/")
}

func (router *Router) handleRegister(c *gin.Context) {
	// check if the user is already logged in
	token, err := c.Cookie("token")

	if err == nil {
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

		if user.Id != "" {
			c.Redirect(http.StatusSeeOther, "/")
			return
		}
	}

	c.HTML(http.StatusOK, "register.html", nil)
}

func (router *Router) handleRegisterPost(c *gin.Context) {
	email := c.PostForm("email")
	name := c.PostForm("name")
	password := c.PostForm("password")

	user := database.User{
		Email:    email,
		Name:     name,
		Password: password,
	}

	id, err := router.Db.CreateUser(user)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// convert id to string
	idStr := fmt.Sprintf("%d", id)
	// generate a random token for the user id
	token := redis.GenerateToken(idStr)

	// store the token in redis
	err = router.R.SetValue(token, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// send an html page index.html and pass a set-cookie header with the token
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)
	//convert the user id to a string
	c.Redirect(http.StatusSeeOther, "/")
}

func (router *Router) handleLogout(c *gin.Context) {
	token, err := c.Cookie("token")
	if err != nil {
		c.Redirect(http.StatusSeeOther, "/")
		return
	}

	err = router.R.DeleteValue(token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.SetCookie("token", "", -1, "/", "localhost", false, true)

	c.Redirect(http.StatusSeeOther, "/")
}
