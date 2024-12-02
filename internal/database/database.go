package database

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

type Database struct {
	Db      *pgx.Conn
	User    *UserModel
	Article *ArticleModel
	Tag     *TagModel
}

func InitDatabase(connString string) (Database, error) {
	// Connect to the database
	conn, err := pgx.Connect(context.Background(), connString)
	if err != nil {
		return Database{}, fmt.Errorf("unable to connect to the database: %w", err)
	}

	// Test the connection
	if err := conn.Ping(context.Background()); err != nil {
		return Database{}, fmt.Errorf("failed to ping database: %w", err)
	}

	_, errr := conn.Exec(context.Background(), "CREATE EXTENSION IF NOT EXISTS pgcrypto")
	if errr != nil {
		return Database{}, fmt.Errorf("failed to create extension: %w", err)
	}

	// TODO:
	// err = DropAllTables(conn, []string{"users", "articles", "tags", "article_tags"})

	// create a user table
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS users (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		email VARCHAR(255) UNIQUE NOT NULL,
		name VARCHAR(255) NOT NULL,
		role VARCHAR(255) NOT NULL DEFAULT 'user',
		password VARCHAR(255) NOT NULL
	)`)
	if err != nil {
		return Database{}, fmt.Errorf("failed to create table: %w", err)
	}

	// create an article table
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS articles (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		title VARCHAR(255) NOT NULL,
		slug VARCHAR(255) UNIQUE NOT NULL,
		excerpt TEXT NOT NULL,
		content TEXT NOT NULL,
		published_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		status VARCHAR(50) NOT NULL DEFAULT 'draft',
		author_id UUID REFERENCES users(id)
	)`)
	if err != nil {
		return Database{}, fmt.Errorf("failed to create table: %w", err)
	}

	// create a tag table
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS tags (
		id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
		name VARCHAR(255) UNIQUE NOT NULL
	)`)
	if err != nil {
		return Database{}, fmt.Errorf("failed to create table: %w", err)
	}

	// create a table to store the relationship between articles and tags
	_, err = conn.Exec(context.Background(), `CREATE TABLE IF NOT EXISTS article_tags (
		article_id UUID REFERENCES articles(id),
		tag_id UUID REFERENCES tags(id),
		PRIMARY KEY (article_id, tag_id)
	)`)
	if err != nil {
		return Database{}, fmt.Errorf("failed to create table: %w", err)
	}

	// Assign the connection to the global variable
	db := Database{
		Db:      conn,
		User:    &UserModel{DB: conn},
		Article: &ArticleModel{DB: conn},
		Tag:     &TagModel{DB: conn},
	}

	// TODO: Populate the database with some initial data
	// err = db.populateDatabase()
	// if err != nil {
	// 	return Database{}, fmt.Errorf("failed to populate database: %w", err)
	// }

	log.Println("Database connection initialized successfully.")
	return db, nil
}

func (d *Database) Close() {
	d.Db.Close(context.Background())
	log.Println("Database connection closed.")
}

func (db *Database) populateDatabase() error {
	// articles
	articleFirstTags := []TagCreation{{
		Name: "debugging",
	},
		{
			Name: "javascript",
		},
		{
			Name: "neverGiveUp",
		},
	}
	articleSecondTags := []TagCreation{
		{
			Name: "foodFusion",
		},
		{
			Name: "polishBrazilianCuisine",
		},
		{
			Name: "pierogiFeijoada",
		},
	}

	articles := []ArticleCreation{
		{
			Title:   "Breaking: Local Dev Finds Missing Semicolon After 3-Day Debugging Marathon",
			Excerpt: "In an unprecedented display of perseverance, local developer \"Semicolon Sam\" has finally located the missing semicolon that's been wreaking havoc on their codebase. Witnesses report seeing tears of joy and hearing victorious screams of \"It compiles!\" echoing through the office",
			Content: "",
		},
		{
			Title:   "Polish-Brazilian Fusion: The Ultimate Pierogi de Feijoada Recipe",
			Excerpt: "In a culinary experiment that's raising eyebrows and tantalizing taste buds, our resident Polish-Brazilian chef has created a dish that bridges continents: the Pierogi de Feijoada. This unexpected fusion combines the comforting dough of Polish pierogis with the rich, hearty filling of Brazilian feijoada",
			Content: "# Pierogi de Feijoada: A Culinary Coding Adventure\n\nIn a world where fusion cuisine meets coding creativity, we bring you the ultimate crossover: Pierogi de Feijoada. This unexpected combination of Polish dumplings and Brazilian stew is taking the programming world by storm.\n\n## Ingredients\n\n- 2 cups of legacy code\n- 1 pound of refactored beans\n- A pinch of Polish syntax\n- A dash of Brazilian flair\n- Bugs to taste\n\n## Instructions\n\n1. Begin by boiling your legacy code until it's soft and pliable.\n2. While the code is boiling, prepare your refactored beans. Ensure all methods are properly encapsulated.\n3. Roll out the legacy code and cut it into circles. These will form the base of your pierogis.\n4. Place a spoonful of refactored beans onto each circle of code.\n5. Fold the code over and pinch the edges to seal. Be careful not to introduce any new bugs at this stage.\n6. Boil the pierogis in a pot of clean, well-commented code for 5 minutes.\n7. Serve hot, garnished with a sprinkle of Polish syntax and a side of Brazilian flair.\n\n## Tasting Notes\n\nThe resulting dish should be a perfect blend of structured Polish design and vibrant Brazilian functionality. If you taste any bugs, simply catch and handle them appropriately.\n\n```python\ntry:\n    eat(pierogi_de_feijoada)\nexcept BugInFoodError as e:\n    print(f\"Oops! Found a bug: {e}\")\n    debug_meal(pierogi_de_feijoada)\n```\n\nRemember, cooking is a lot like coding. Sometimes you need to debug your recipes, refactor your techniques, and always be ready to handle unexpected exceptions in the kitchen!\n\nBom apetite and smacznego!",
		},
	}

	// list of articles and tags grouped
	var articleTags = []struct {
		article ArticleCreation
		tags    []TagCreation
	}{
		{
			article: articles[0],
			tags:    articleFirstTags,
		},
		{
			article: articles[1],
			tags:    articleSecondTags,
		},
	}

	log.Println("Populating database with initial data...")

	// create the user
	userId, err := db.User.Create(User{
		Email:    "bartek@paczesny.pl",
		Name:     "Bartek Paczesny",
		Role:     "admin",
		Password: "password",
	})
	if err != nil {
		log.Fatalf("failed to create user: %v", err)
	}

	for _, article := range articleTags {
		log.Printf("Creating article: %s\n", article.article.Title)

		article.article.AuthorId = userId

		_, err := db.Article.Create(article.article, article.tags)
		if err != nil {
			log.Fatalf("failed to create article: %v", err)
		}
	}

	return nil
}

func DropAllTables(db *pgx.Conn, tables []string) error {
	for _, table := range tables {
		_, err := db.Exec(context.Background(), fmt.Sprintf("DROP TABLE IF EXISTS %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to drop table %s: %w", table, err)
		}
	}

	return nil
}
