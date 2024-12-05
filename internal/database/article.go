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
	Status      string    `json:"status"`
	Author      User      `json:"author"`
}

type ArticleCreation struct {
	Title    string `form:"title"`
	Content  string `form:"content"`
	Excerpt  string `form:"excerpt"`
	AuthorId string `form:"author_id"`
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
	query := `INSERT INTO articles (title, slug, content, excerpt, published_at, status, author_id) VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err = tx.QueryRow(context.Background(), query, article.Title, articleSlug, article.Content, article.Excerpt, publishedAt, "draft", article.AuthorId).Scan(&articleID)
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
	query := `SELECT a.id, a.title, a.slug, a.content, a.excerpt, a.published_at, a.status, t.id, t.name, u.id, u.email, u.role, u.name, u.password FROM articles a
	JOIN article_tags at ON a.id = at.article_id
	JOIN tags t ON at.tag_id = t.id
	JOIN users u ON a.author_id = u.id
	WHERE a.slug = $1`

	rows, err := a.DB.Query(context.Background(), query, slug)
	if err != nil {
		return Article{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		err := rows.Scan(&article.Id, &article.Title, &article.Slug, &article.Content, &article.Excerpt, &article.PublishedAt, &article.Status, &tag.Id, &tag.Name, &article.Author.Id, &article.Author.Email, &article.Author.Role, &article.Author.Name, &article.Author.Password)
		if err != nil {
			return Article{}, err
		}
		article.Tags = append(article.Tags, tag)
	}

	return article, nil
}

func (a *ArticleModel) ById(id string) (Article, error) {
	var article Article
	query := `SELECT a.id, a.title, a.slug, a.content, a.excerpt, a.published_at, a.status, t.id, t.name, u.id, u.email, u.role, u.name, u.password FROM articles a
	JOIN article_tags at ON a.id = at.article_id
	JOIN tags t ON at.tag_id = t.id
	JOIN users u ON a.author_id = u.id
	WHERE a.id = $1`

	rows, err := a.DB.Query(context.Background(), query, id)
	if err != nil {
		return Article{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		err := rows.Scan(&article.Id, &article.Title, &article.Slug, &article.Content, &article.Excerpt, &article.PublishedAt, &article.Status, &tag.Id, &tag.Name, &article.Author.Id, &article.Author.Email, &article.Author.Role, &article.Author.Name, &article.Author.Password)
		if err != nil {
			return Article{}, err
		}
		article.Tags = append(article.Tags, tag)
	}

	return article, nil
}

func (a *ArticleModel) All(status string) ([]Article, error) {
	articles := make(map[string]Article)
	var query string
	var rows pgx.Rows
	var err error
	if status != "" {
		query = `SELECT a.id, a.title, a.slug, a.content, a.excerpt, a.published_at, a.status, t.id, t.name, u.id, u.email, u.role, u.name, u.password FROM articles a
		JOIN article_tags at ON a.id = at.article_id
		JOIN tags t ON at.tag_id = t.id
		JOIN users u ON a.author_id = u.id
		WHERE a.status = $1`
		rows, err = a.DB.Query(context.Background(), query, status)
	} else {
		query = `SELECT a.id, a.title, a.slug, a.content, a.excerpt, a.published_at, a.status, t.id, t.name, u.id, u.email, u.role, u.name, u.password FROM articles a
		JOIN article_tags at ON a.id = at.article_id
		JOIN tags t ON at.tag_id = t.id
		JOIN users u ON a.author_id = u.id`
		rows, err = a.DB.Query(context.Background(), query)
	}
	if err != nil {
		fmt.Println("Router [ArticleModel.All] (query): " + err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var article Article
		var tag Tag
		err := rows.Scan(&article.Id, &article.Title, &article.Slug, &article.Content, &article.Excerpt, &article.PublishedAt, &article.Status, &tag.Id, &tag.Name, &article.Author.Id, &article.Author.Email, &article.Author.Role, &article.Author.Name, &article.Author.Password)
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

func (a *ArticleModel) Update(id string, article ArticleCreation, status string) error {
	query := "UPDATE articles SET title = $1, content = $2, excerpt = $3, status = $4 WHERE id = $5"
	_, err := a.DB.Exec(context.Background(), query, article.Title, article.Content, article.Excerpt, status, id)
	return err
}

func (a *ArticleModel) Publish(id string) (Article, error) {
	query := "UPDATE articles SET status = $1, published_at = $2 WHERE id = $3"
	_, err := a.DB.Exec(context.Background(), query, "published", time.Now(), id)
	if err != nil {
		return Article{}, err
	}

	var article Article
	query = `SELECT a.id, a.title, a.slug, a.content, a.excerpt, a.published_at, a.status, t.id, t.name, u.id, u.email, u.role, u.name, u.password FROM articles a
	JOIN article_tags at ON a.id = at.article_id
	JOIN tags t ON at.tag_id = t.id
	JOIN users u ON a.author_id = u.id
	WHERE a.id = $1`

	rows, err := a.DB.Query(context.Background(), query, id)
	if err != nil {
		return Article{}, err
	}
	defer rows.Close()

	for rows.Next() {
		var tag Tag
		err := rows.Scan(&article.Id, &article.Title, &article.Slug, &article.Content, &article.Excerpt, &article.PublishedAt, &article.Status, &tag.Id, &tag.Name, &article.Author.Id, &article.Author.Email, &article.Author.Role, &article.Author.Name, &article.Author.Password)
		if err != nil {
			return Article{}, err
		}
		article.Tags = append(article.Tags, tag)
	}

	return article, nil
}

func (a *ArticleModel) UnPublish(id string) (Article, error) {
	query := "UPDATE articles SET status = $1, published_at = $2 WHERE id = $3"
	_, err := a.DB.Exec(context.Background(), query, "draft", time.Now(), id)
	if err != nil {
		return Article{}, err
	}

	var article Article
	query = `SELECT a.id, a.title, a.slug, a.content, a.excerpt, a.published_at, a.status, t.id, t.name, u.id, u.email, u.role, u.name, u.password FROM articles a
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
		err := rows.Scan(&article.Id, &article.Title, &article.Slug, &article.Content, &article.Excerpt, &article.PublishedAt, &article.Status, &tag.Id, &tag.Name, &article.Author.Id, &article.Author.Email, &article.Author.Role, &article.Author.Name, &article.Author.Password)
		if err != nil {
			return Article{}, err
		}
		article.Tags = append(article.Tags, tag)
	}

	return article, nil
}

func (a *ArticleModel) UpdateTags(articleID string, tags []TagCreation) error {
	tx, err := a.DB.Begin(context.Background())
	if err != nil {
		return err
	}

	// Delete existing tags
	_, err = tx.Exec(context.Background(), "DELETE FROM article_tags WHERE article_id = $1", articleID)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}

	// Add new tags
	for _, tag := range tags {
		var tagID string
		tagQuery := "INSERT INTO tags (name) VALUES ($1) ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name RETURNING id"
		err = tx.QueryRow(context.Background(), tagQuery, tag.Name).Scan(&tagID)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}

		articleTagQuery := "INSERT INTO article_tags (article_id, tag_id) VALUES ($1, $2)"
		_, err = tx.Exec(context.Background(), articleTagQuery, articleID, tagID)
		if err != nil {
			tx.Rollback(context.Background())
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func (a Article) TagsNames() []string {
	var names []string
	for _, tag := range a.Tags {
		names = append(names, tag.Name)
	}
	return names
}
