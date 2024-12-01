package database

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
)

type ArticleModel struct {
	DB      *pgx.Conn
	Article Article
}

type Article struct {
	Id          string    `json:"id"`
	Slug        string    `json:"slug"`
	Title       string    `json:"title"`
	Content     string    `json:"content"`
	Excerpt     string    `json:"excerpt"`
	Tags        []Tag     `json:"tags"`
	PublishedAt time.Time `json:"published_at"`
}

type ArticleCreation struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Excerpt string `json:"excerpt"`
}

func (a *ArticleModel) Create(article ArticleCreation, tags []TagCreation) (string, error) {
	tx, err := a.DB.Begin(context.Background())
	if err != nil {
		return "", err
	}

	// build a slug from the title. get rid of the excape characters and replace spacec with dashes. use the url package
	articleSlug := strings.ReplaceAll(url.PathEscape(article.Title), "%20", "-")
	publishedAt := time.Now().Format("2006-01-02 15:04:05")

	var articleID string
	query := "INSERT INTO articles (title, slug, content, excerpt, published_at) VALUES ($1, $2, $3, $4, $5) RETURNING id"
	err = tx.QueryRow(context.Background(), query, article.Title, articleSlug, article.Content, article.Excerpt, publishedAt).Scan(&articleID)
	if err != nil {
		tx.Rollback(context.Background())
		return "", err
	}

	for _, tag := range tags {
		var tagID string
		tagQuery := "INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id"
		err = tx.QueryRow(context.Background(), tagQuery, tag.Name).Scan(&tagID)
		if err != nil {
			tx.Rollback(context.Background())
			return "", err
		}

		articleTagQuery := "INSERT INTO article_tags (article_id, tag_id) VALUES ($1, $2)"
		_, err = tx.Exec(context.Background(), articleTagQuery, articleID, tagID)
		if err != nil {
			tx.Rollback(context.Background())
			return "", err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return "", err
	}

	return articleID, nil
}

func (a *ArticleModel) BySlug(slug string) (Article, error) {
	var article Article
	query := `SELECT a.id, a.title, a.slug, a.content, a.excerpt, a.published_at, t.id, t.name FROM articles a
	JOIN article_tags at ON a.id = at.article_id
	JOIN tags t ON at.tag_id = t.id
	WHERE a.slug = $1`

	rows, err := a.DB.Query(context.Background(), query, slug)
	if err != nil {
		return Article{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		err := rows.Scan(&article.Id, &article.Title, &article.Slug, &article.Content, &article.Excerpt, &article.PublishedAt, &tag.Id, &tag.Name)
		if err != nil {
			return Article{}, err
		}
		article.Tags = append(article.Tags, tag)
	}

	return article, nil
}

func (a *ArticleModel) ById(id string) (Article, error) {
	var article Article
	query := `SELECT a.id, a.title, a.slug, a.content, a.excerpt, a.published_at, t.id, t.name FROM articles a
	JOIN article_tags at ON a.id = at.article_id
	JOIN tags t ON at.tag_id = t.id
	WHERE a.id = $1`

	rows, err := a.DB.Query(context.Background(), query, id)
	if err != nil {
		return Article{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		err := rows.Scan(&article.Id, &article.Title, &article.Slug, &article.Content, &article.Excerpt, &article.PublishedAt, &tag.Id, &tag.Name)
		if err != nil {
			return Article{}, err
		}
		article.Tags = append(article.Tags, tag)
	}

	return article, nil
}

func (a *ArticleModel) All() ([]Article, error) {
	articles := make(map[string]Article)
	query := `SELECT a.id, a.title, a.slug, a.content, a.excerpt, a.published_at, t.id, t.name FROM articles a
	JOIN article_tags at ON a.id = at.article_id
	JOIN tags t ON at.tag_id = t.id`

	rows, err := a.DB.Query(context.Background(), query)
	if err != nil {
		fmt.Println("Router [ArticleModel.All] (query): " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article Article
		var tag Tag
		err := rows.Scan(&article.Id, &article.Title, &article.Slug, &article.Content, &article.Excerpt, &article.PublishedAt, &tag.Id, &tag.Name)
		if err != nil {
			return nil, err
		}
		// Check if article already exists in the map
		existingArticle, ok := articles[article.Id]
		if ok {
			existingArticle.Tags = append(existingArticle.Tags, tag)
			articles[article.Id] = existingArticle
		} else {
			article.Tags = []Tag{tag}
			articles[article.Id] = article
		}
	}
	// ... convert map to slice of articles
	var result []Article
	for _, article := range articles {
		result = append(result, article)
	}
	return result, nil
}
