package main

import (
	"io/fs"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/not-a-phase-mom/go-fullstack-yourself/cmd/api/routes"
	"github.com/not-a-phase-mom/go-fullstack-yourself/cmd/web"
	"github.com/not-a-phase-mom/go-fullstack-yourself/internal/database"
	"github.com/not-a-phase-mom/go-fullstack-yourself/internal/redis"
)

// load the environment variables
func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	PORT := os.Getenv("PORT")

	POSTGRES_DATABASE := os.Getenv("POSTGRES_DATABASE")
	POSTGRES_USER := os.Getenv("POSTGRES_USER")
	POSTGRES_PASSWORD := os.Getenv("POSTGRES_PASSWORD")

	REDIS_ADDR := "localhost:6379"
	REDIS_PASSWORD := os.Getenv("REDIS_PASSWORD")
	REDIS_DB := 0

	// create the connection string
	connString := "postgres://" + POSTGRES_USER + ":" + POSTGRES_PASSWORD + "@localhost/" + POSTGRES_DATABASE + "?sslmode=disable"

	// connect to the database
	db, err := database.InitDatabase(connString)
	defer db.Close()
	if err != nil {
		log.Fatalf("failed to connect to the database: %v", err)
	}

	// connect to redis
	r, err := redis.InitRedis(REDIS_ADDR, REDIS_PASSWORD, REDIS_DB)
	defer r.Close()
	if err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}

	e := gin.Default()

	// set up the logger
	e.Use(gin.Logger())

	staticFiles, _ := fs.Sub(web.Files, "assets")
	log.Println(staticFiles)
	e.StaticFS("/assets", http.FS(staticFiles))

	router := routes.Router{
		Db: &db,
		R:  r,
	}

	router.RegisterAuthRoutes(e)
	router.RegisterUserRoutes(e)
	router.RegisterIndexRoutes(e)
	router.RegisterArticleRoutes(e)

	e.Run(":" + PORT)
}
